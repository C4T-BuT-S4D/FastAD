package checkers

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/c4t-but-s4d/fastad/internal/models"
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
) *models.CheckerVerdict {
	cmd := exec.CommandContext(
		ctx,
		filepath.Join(
			params.GameSettings.CheckersBasePath,
			params.Service.CheckerPath,
		),
		checkAction,
		params.Team.Address,
	)
	return RunAction(ctx, cmd, params.Service.CheckerTimeout())
}

func RunPutAction(
	ctx context.Context,
	params *PutActivityParameters,
	flag *models.Flag,
) *models.CheckerVerdict {
	cmd := exec.CommandContext(
		ctx,
		filepath.Join(
			params.GameSettings.CheckersBasePath,
			params.Service.CheckerPath,
		),
		putAction,
		params.Team.Address,
		flag.Private,
		flag.Flag,
		"1",
	)
	return RunAction(ctx, cmd, params.Service.CheckerTimeout())
}

func RunGetAction(
	ctx context.Context,
	params *GetActivityParameters,
	flag *models.Flag,
) *models.CheckerVerdict {
	cmd := exec.CommandContext(
		ctx,
		filepath.Join(
			params.GameSettings.CheckersBasePath,
			params.Service.CheckerPath,
		),
		getAction,
		params.Team.Address,
		flag.Private,
		flag.Flag,
		"1",
	)
	return RunAction(ctx, cmd, params.Service.CheckerTimeout())
}

func RunAction(ctx context.Context, cmd *exec.Cmd, softTimeout time.Duration) *models.CheckerVerdict {
	// TODO: limit buffer size.
	stdout := bytes.Buffer{}
	stderr := bytes.Buffer{}

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	commandContext, commandCancel := context.WithTimeout(ctx, softTimeout)
	defer commandCancel()

	verdict := &models.CheckerVerdict{
		Action:  checkerpb.Action_ACTION_CHECK,
		Command: cmd.String(),
	}

	err := RunCommandGracefully(commandContext, cmd, checkerKillDelay)
	switch {
	case err == nil:
		verdict.Status = checkerpb.Status_STATUS_UP
		verdict.Public = stdout.String()
		verdict.Private = stderr.String()

	case errors.Is(err, TimeoutError):
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
