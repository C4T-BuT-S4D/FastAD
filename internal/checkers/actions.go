package checkers

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os/exec"
	"path/filepath"
	"syscall"
	"time"

	checkerpb "github.com/c4t-but-s4d/fastad/pkg/proto/checker"
)

// TODO: add process reaper to avoid zombies.

const (
	checkAction = "check"
	putAction   = "put"
	getAction   = "get"
)

const checkerKillDelay = time.Second * 3

func RunCheckAction(
	ctx context.Context,
	params *CheckActivityParameters,
) *Verdict {
	checkerPath := filepath.Join("checkers", params.Service.CheckerPath)
	return RunAction(
		ctx,
		checkerPath,
		checkerpb.Action_ACTION_CHECK,
		[]string{checkAction, params.Team.Address},
		params.Service.CheckerTimeout(checkerpb.Action_ACTION_CHECK),
	)
}

func RunPutAction(
	ctx context.Context,
	params *PutActivityParameters,
) *Verdict {
	checkerPath := filepath.Join("checkers", params.FlagInfo.Service.CheckerPath)
	return RunAction(
		ctx,
		checkerPath,
		checkerpb.Action_ACTION_PUT,
		[]string{
			putAction,
			params.FlagInfo.Team.Address,
			params.FlagInfo.Flag.Private,
			params.FlagInfo.Flag.Flag,
			"1", // TODO: vulns.
		},
		params.FlagInfo.Service.CheckerTimeout(checkerpb.Action_ACTION_PUT),
	)
}

func RunGetAction(
	ctx context.Context,
	params *GetActivityParameters,
) *Verdict {
	checkerPath := filepath.Join("checkers", params.Service.CheckerPath)
	return RunAction(
		ctx,
		checkerPath,
		checkerpb.Action_ACTION_GET,
		[]string{
			getAction,
			params.Team.Address,
			params.Flag.Private,
			params.Flag.Flag,
			"1", // TODO: vulns.
		},
		params.Service.CheckerTimeout(checkerpb.Action_ACTION_GET),
	)
}

func RunAction(
	ctx context.Context,
	checkerPath string,
	action checkerpb.Action,
	args []string,
	softTimeout time.Duration,
) *Verdict {
	ctx, cancel := context.WithTimeout(ctx, softTimeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, checkerPath, args...)
	cmd.Cancel = func() error {
		if err := cmd.Process.Signal(syscall.SIGTERM); err != nil {
			return fmt.Errorf("sending SIGTERM to process: %w", err)
		}
		return nil
	}
	cmd.WaitDelay = checkerKillDelay

	// TODO: limit buffer size.
	stdout := bytes.Buffer{}
	stderr := bytes.Buffer{}

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	verdict := &Verdict{
		Action:  action,
		Command: cmd.String(),
	}

	err := cmd.Run()
	switch {
	case err == nil:
		verdict.Status = checkerpb.Status_STATUS_UP
		verdict.Public = stdout.String()
		verdict.Private = stderr.String()

	case errors.Is(err, context.DeadlineExceeded):
		verdict.Status = checkerpb.Status_STATUS_DOWN
		verdict.Public = "timeout"
		// TODO: truncate.
		verdict.Private = fmt.Sprintf("err: %v\nstdout: %s\nstderr: %s", err, stdout.String(), stderr.String())

	default:
		verdict.Status = checkerpb.Status_STATUS_CHECK_FAILED
		verdict.Public = "internal error"
		// TODO: truncate.
		verdict.Private = fmt.Sprintf("err: %v\nstdout: %s\nstderr: %s", err, stdout.String(), stderr.String())
	}
	return verdict
}
