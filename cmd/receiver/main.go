package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/c4t-but-s4d/fastad/internal/logging"
	"github.com/c4t-but-s4d/fastad/internal/multiproto"
	"github.com/c4t-but-s4d/fastad/internal/pinger"
	"github.com/c4t-but-s4d/fastad/internal/receiver"
	pingerpb "github.com/c4t-but-s4d/fastad/pkg/proto/pinger"
	receiverpb "github.com/c4t-but-s4d/fastad/pkg/proto/receiver"
)

func main() {
	logging.Init()

	logrus.Info("Starting flag receiver")

	grpcServer := grpc.NewServer()
	receiverpb.RegisterReceiverServiceServer(grpcServer, receiver.New())
	pingerpb.RegisterPingerServiceServer(grpcServer, pinger.New())
	reflection.Register(grpcServer)

	httpServer := &http.Server{
		Addr:    "0.0.0.0:8002",
		Handler: multiproto.NewHandler(grpcServer),
	}

	go func() {
		logrus.Infof("Running http server on %s", httpServer.Addr)
		if err := httpServer.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			logrus.Fatalf("Error running http server: %v", err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		logrus.Fatalf("Shutting down http server: %v", err)
	}
}
