package metrics

import (
	"context"
	"errors"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

func RunServer(runCtx, shutdownCtx context.Context, addr string) {
	logger := zap.L().With(
		zap.String("service", "metrics"),
		zap.String("listen_address", addr),
	)

	logger.Info("starting server")

	m := http.NewServeMux()
	m.Handle("/metrics", promhttp.Handler())
	s := &http.Server{
		Addr:    addr,
		Handler: m,
	}

	go func() {
		<-runCtx.Done()
		if err := s.Shutdown(shutdownCtx); err != nil {
			logger.Error("error shutting down server", zap.Error(err))
		}
	}()

	if err := s.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Error("running server", zap.Error(err))
	}
}
