package checkers

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
)

var TimeoutError = errors.New("timeout")

func RunCommandGracefully(ctx context.Context, cmd *exec.Cmd, killDelay time.Duration) error {
	err := cmd.Start()
	if err != nil {
		return fmt.Errorf("failed to start command: %w", err)
	}

	done := make(chan error)
	go func() {
		done <- cmd.Wait()
	}()

	select {
	case <-ctx.Done():
		// terminate process gracefully.
		if err := cmd.Process.Signal(syscall.SIGTERM); err != nil {
			logrus.Debugf("failed to send SIGTERM to process: %s", err)
		}

		t := time.NewTimer(killDelay)
		defer t.Stop()

		select {
		case err := <-done:
			// process exited before kill delay,
			// return timeout + original error (most likely terminated by signal).
			return errors.Join(TimeoutError, fmt.Errorf("waiting for process: %w", err))
		case <-t.C:
			// kill process.
		}

		<-done
		return errors.Join(TimeoutError, ctx.Err())
	case err := <-done:
		if err != nil {
			return fmt.Errorf("waiting for process: %w", err)
		}
		return nil
	}
}
