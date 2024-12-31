package impl

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/centrifugal/centrifuge"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"github.com/c4t-but-s4d/fastad/internal/api"
	"github.com/c4t-but-s4d/fastad/internal/centutil"
	"github.com/c4t-but-s4d/fastad/pkg/clients/gamestate"
	"github.com/c4t-but-s4d/fastad/pkg/clients/services"
	"github.com/c4t-but-s4d/fastad/pkg/clients/teams"
	"github.com/c4t-but-s4d/fastad/pkg/grpcext"
	"github.com/c4t-but-s4d/fastad/pkg/httpext"
	"github.com/c4t-but-s4d/fastad/pkg/metrics"
	gspb "github.com/c4t-but-s4d/fastad/pkg/proto/data/game_state"
	servicespb "github.com/c4t-but-s4d/fastad/pkg/proto/data/services"
	teamspb "github.com/c4t-but-s4d/fastad/pkg/proto/data/teams"
	receiverpb "github.com/c4t-but-s4d/fastad/pkg/proto/receiver"
	slacpb "github.com/c4t-but-s4d/fastad/pkg/proto/slac"
)

func Run(runCtx, shutdownCtx context.Context, cfg *api.Config) error {
	dataServiceConn, err := grpcext.Dial(cfg.DataService.Address, cfg.Installation)
	if err != nil {
		return fmt.Errorf("connecting to data service: %w", err)
	}

	teamsClient := teams.NewClient(teamspb.NewTeamsServiceClient(dataServiceConn))
	servicesClient := services.NewClient(servicespb.NewServicesServiceClient(dataServiceConn))
	gameStateClient := gamestate.NewClient(gspb.NewGameStateServiceClient(dataServiceConn))

	receiverConn, err := grpcext.Dial(cfg.ReceiverAddress, cfg.Installation)
	if err != nil {
		return fmt.Errorf("connecting to receiver: %w", err)
	}
	receiverClient := receiverpb.NewReceiverServiceClient(receiverConn)

	slacConn, err := grpcext.Dial(cfg.SlacAddress, cfg.Installation)
	if err != nil {
		return fmt.Errorf("connecting to slac: %w", err)
	}
	slacClient := slacpb.NewSlacServiceClient(slacConn)

	centConfig := centrifuge.Config{
		Name:     cfg.Installation,
		LogLevel: centrifuge.LogLevelDebug,
		LogHandler: func(entry centrifuge.LogEntry) {
			zap.L().Debug(
				"centrifuge log",
				zap.String("message", entry.Message),
				zap.Any("entry", entry.Fields),
			)
		},
	}
	node, err := centrifuge.New(centConfig)
	if err != nil {
		return fmt.Errorf("creating centrifuge node: %w", err)
	}

	producer := centutil.NewNodeProducer(node, cfg.ScoreboardChannel)

	boardBuilder := api.NewBoardBuilder(
		teamsClient,
		servicesClient,
		receiverClient,
		slacClient,
		producer,
	)

	db := cfg.Postgres.BunDB()

	apiService := api.NewService(
		cfg,
		db,
		node,
		teamsClient,
		servicesClient,
		gameStateClient,
		receiverClient,
		slacClient,
		boardBuilder,
	)

	e := echo.New()
	e.Use(
		httpext.RequestIDMiddleware(),
		middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
			LogMethod:    true,
			LogURI:       true,
			LogStatus:    true,
			LogLatency:   true,
			LogRemoteIP:  true,
			LogRoutePath: true,
			LogValuesFunc: func(_ echo.Context, v middleware.RequestLoggerValues) error {
				zap.L().Info("request",
					zap.String("method", v.Method),
					zap.String("URI", v.URI),
					zap.Int("status", v.Status),
					zap.Duration("latency", v.Latency),
					zap.String("remote_ip", v.RemoteIP),
					zap.String("path", v.RoutePath),
				)

				return nil
			},
		}),
		middleware.Recover(),
		middleware.Gzip(),
		middleware.Secure(), // TODO: allow images from external sources for avatars.
		middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins:                             []string{"*"},
			AllowCredentials:                         true,
			UnsafeWildcardOriginWithAllowCredentials: true,
		}),
	)
	e.GET("/health", httpext.HealthHandler())

	e.HTTPErrorHandler = httpext.ErrorHandler()
	apiService.RegisterRoutes(e)
	apiService.RegisterNode()

	if err := node.Run(); err != nil {
		return fmt.Errorf("running centrifuge node: %w", err)
	}

	g, gctx := errgroup.WithContext(runCtx)
	g.Go(func() error {
		boardBuilder.Run(gctx)
		return nil
	})

	if cfg.MetricsAddress != "" {
		g.Go(func() error {
			metrics.RunServer(gctx, shutdownCtx, cfg.MetricsAddress)
			return nil
		})
	}

	g.Go(func() error {
		zap.L().Info("starting api server", zap.String("address", cfg.ListenAddress))
		if err := e.Start(cfg.ListenAddress); err != nil && !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("running server: %w", err)
		}
		zap.L().Info("api server stopped")
		return nil
	})

	g.Go(func() error {
		<-gctx.Done()

		zap.L().Info("shutting down api server")
		if err := e.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("shutting down server: %w", err)
		}
		if err := node.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("shutting down centrifuge node: %w", err)
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		return fmt.Errorf("waiting for goroutines: %w", err)
	}

	return nil
}
