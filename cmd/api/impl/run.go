package impl

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"

	"github.com/c4t-but-s4d/fastad/internal/api"
	"github.com/c4t-but-s4d/fastad/pkg/clients/gamestate"
	"github.com/c4t-but-s4d/fastad/pkg/clients/services"
	"github.com/c4t-but-s4d/fastad/pkg/clients/teams"
	"github.com/c4t-but-s4d/fastad/pkg/grpcext"
	"github.com/c4t-but-s4d/fastad/pkg/httpext"
	gspb "github.com/c4t-but-s4d/fastad/pkg/proto/data/game_state"
	servicespb "github.com/c4t-but-s4d/fastad/pkg/proto/data/services"
	teamspb "github.com/c4t-but-s4d/fastad/pkg/proto/data/teams"
	receiverpb "github.com/c4t-but-s4d/fastad/pkg/proto/receiver"
	scoreboardpb "github.com/c4t-but-s4d/fastad/pkg/proto/scoreboard"
)

func Run(runCtx, shutdownCtx context.Context, cfg *api.Config) error {
	dataServiceConn, err := grpcext.Dial(cfg.DataService.Address, cfg.UserAgent)
	if err != nil {
		return fmt.Errorf("connecting to data service: %w", err)
	}

	teamsClient := teams.NewClient(teamspb.NewTeamsServiceClient(dataServiceConn))
	servicesClient := services.NewClient(servicespb.NewServicesServiceClient(dataServiceConn))
	gameStateClient := gamestate.NewClient(gspb.NewGameStateServiceClient(dataServiceConn))

	receiverConn, err := grpcext.Dial(cfg.ReceiverAddress, cfg.UserAgent)
	if err != nil {
		return fmt.Errorf("connecting to receiver service: %w", err)
	}
	receiverClient := receiverpb.NewReceiverServiceClient(receiverConn)

	scoreboardConn, err := grpcext.Dial(cfg.ScoreboardAddress, cfg.UserAgent)
	if err != nil {
		return fmt.Errorf("connecting to scoreboard service: %w", err)
	}
	scoreboardClient := scoreboardpb.NewScoreboardServiceClient(scoreboardConn)

	apiService := api.NewService(
		teamsClient,
		servicesClient,
		gameStateClient,
		receiverClient,
		scoreboardClient,
	)

	e := echo.New()
	e.Use(
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
		middleware.RequestID(),
	)

	e.HTTPErrorHandler = httpext.ErrorHandler()
	apiService.RegisterRoutes(e)

	go func() {
		<-runCtx.Done()

		zap.L().Info("shutting down api server")
		if err := e.Shutdown(shutdownCtx); err != nil {
			zap.L().Error("error shutting down api server", zap.Error(err))
		}
	}()

	if err := e.Start(cfg.ListenAddress); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("running server: %w", err)
	}

	return nil
}
