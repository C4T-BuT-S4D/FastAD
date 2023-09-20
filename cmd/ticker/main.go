package main

import (
	"context"
	"os/signal"
	"syscall"
	"time"

	"github.com/c4t-but-s4d/fastad/internal/logging"
	"github.com/c4t-but-s4d/fastad/internal/ticker"
	"github.com/sirupsen/logrus"
	"go.temporal.io/sdk/client"
)

func main() {
	logging.Init()

	temporalClient, err := client.Dial(client.Options{
		HostPort: "localhost:7233",
		Logger: logging.NewTemporalAdapter(
			logrus.WithFields(logrus.Fields{
				"component": "ticker",
			}),
		),
	})
	if err != nil {
		logrus.Fatalf("dialing temporal: %v", err)
	}
	defer temporalClient.Close()

	t := ticker.NewTicker(time.Second*10, temporalClient)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	t.Run(ctx)
}
