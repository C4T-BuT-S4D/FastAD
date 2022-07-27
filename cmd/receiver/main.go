package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"FastAD/internal/pinger"
	"FastAD/internal/receiver"
	"FastAD/internal/wrapper"
	pingerpb "FastAD/pkg/proto/pinger"
	receiverpb "FastAD/pkg/proto/receiver"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	initLogger()

	logrus.Info("Starting flag receiver")

	grpcServer := grpc.NewServer()
	receiverpb.RegisterReceiverServiceServer(grpcServer, receiver.New())
	pingerpb.RegisterPingerServiceServer(grpcServer, pinger.New())
	reflection.Register(grpcServer)

	httpServer := wrapper.NewServer(grpcServer)
	httpServer.Addr = "0.0.0.0:8002"

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

func initLogger() {
	mainFormatter := &logrus.TextFormatter{}
	mainFormatter.FullTimestamp = true
	mainFormatter.ForceColors = true
	mainFormatter.PadLevelText = true
	mainFormatter.TimestampFormat = "2006-01-02 15:04:05"
	logrus.SetFormatter(mainFormatter)
}
