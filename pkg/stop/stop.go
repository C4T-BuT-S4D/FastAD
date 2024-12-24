package stop

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

// SetupCtx returns two contexts: runCtx and shutdownCtx.
// RunCtx is used to run the program until the first ^C (interrupt signal) is received.
// ShutdownCtx is used to handle the shutdown process until the second ^C is received.
func SetupCtx() (context.Context, context.Context, context.CancelFunc) {
	runCtx, runCancel := context.WithCancel(context.Background())
	shutdownCtx, shutdownCancel := context.WithCancel(context.Background())

	ch := make(chan os.Signal, 2) //nolint:mnd // 2 is the number of signals we are interested in
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		select {
		case <-ch:
			runCancel()
		case <-runCtx.Done():
		}

		if shutdownCtx.Err() == nil {
			go func() {
				select {
				case <-ch:
					shutdownCancel()
				case <-shutdownCtx.Done():
				}
			}()
		}
	}()

	return runCtx, shutdownCtx, func() {
		shutdownCancel()
		runCancel()
		signal.Stop(ch)
	}
}
