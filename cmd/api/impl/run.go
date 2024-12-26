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

	"github.com/c4t-but-s4d/fastad/internal/api"
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
	scoreboardpb "github.com/c4t-but-s4d/fastad/pkg/proto/scoreboard"
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
		return fmt.Errorf("connecting to receiver service: %w", err)
	}
	receiverClient := receiverpb.NewReceiverServiceClient(receiverConn)

	scoreboardConn, err := grpcext.Dial(cfg.ScoreboardAddress, cfg.Installation)
	if err != nil {
		return fmt.Errorf("connecting to scoreboard service: %w", err)
	}
	scoreboardClient := scoreboardpb.NewScoreboardServiceClient(scoreboardConn)

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

	apiService := api.NewService(
		cfg,
		node,
		teamsClient,
		servicesClient,
		gameStateClient,
		receiverClient,
		scoreboardClient,
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

	go func() {
		<-runCtx.Done()

		zap.L().Info("shutting down api server")
		if err := e.Shutdown(shutdownCtx); err != nil {
			zap.L().Error("error shutting down api server", zap.Error(err))
		}
		if err := node.Shutdown(shutdownCtx); err != nil {
			zap.L().Error("error shutting down centrifuge node", zap.Error(err))
		}
	}()

	if cfg.MetricsAddress != "" {
		go metrics.RunServer(runCtx, shutdownCtx, cfg.MetricsAddress)
	}

	if err := e.Start(cfg.ListenAddress); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("running server: %w", err)
	}

	return nil
}
