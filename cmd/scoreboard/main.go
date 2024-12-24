package main

import (
	"context"
	"errors"
	"net/http"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"go.uber.org/zap"

	"github.com/c4t-but-s4d/fastad/internal/baseconfig"
	"github.com/c4t-but-s4d/fastad/internal/logging"
	"github.com/c4t-but-s4d/fastad/internal/multiproto"
	"github.com/c4t-but-s4d/fastad/internal/scoreboard"
	"github.com/c4t-but-s4d/fastad/pkg/grpcext"
	scoreboardpb "github.com/c4t-but-s4d/fastad/pkg/proto/scoreboard"
)

func main() {
	defer logging.Init().Close()

	cfg := baseconfig.MustSetupAll(&scoreboard.Config{}, baseconfig.WithEnvPrefix("FASTAD_DATA_SERVICE"))

	db := cfg.Postgres.BunDB()

	service := scoreboard.NewService(db, cfg)

	grpcServer := grpcext.NewServer()
	scoreboardpb.RegisterScoreboardServiceServer(grpcServer, service)

	httpServer := &http.Server{
		Addr:              "0.0.0.0:8003",
		Handler:           multiproto.NewHandler(grpcServer),
		ReadHeaderTimeout: time.Second * 10,
	}

	runCtx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	if err := service.RestoreState(runCtx); err != nil {
		zap.L().Fatal("failed to restore state", zap.Error(err))
	}

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		service.Run(runCtx)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		zap.L().With(zap.String("listen_address", httpServer.Addr)).Info("Running http server")
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			zap.L().With(zap.Error(err)).Fatal("error running http server")
		}
	}()

	<-runCtx.Done()

	zap.L().Info("shutting down")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		zap.L().With(zap.Error(err)).Fatal("error shutting down server")
	}

	wg.Wait()
}
