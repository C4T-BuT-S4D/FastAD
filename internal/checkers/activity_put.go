package checkers

import (
	"context"

	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/log"

	checkerpb "github.com/c4t-but-s4d/fastad/pkg/proto/checker"
)

const PutActivityName = "Put"

type PutActivity struct{}

func NewPutActivity() *PutActivity {
	return &PutActivity{}
}

type PutActivityParameters struct {
	FlagInfo *FQFlagInfo
}

type PutActivityResult struct {
	FlagInfo *FQFlagInfo
	Verdict  *Verdict
}

func (a *PutActivity) ActivityDefinition(ctx context.Context, params *PutActivityParameters) (*PutActivityResult, error) {
	logger := log.With(
		activity.GetLogger(ctx),
		"team", params.FlagInfo.Team.Name,
		"service", params.FlagInfo.Service.Name,
		"action", checkerpb.Action_ACTION_PUT,
	)

	logger.Info("starting")
	verdict := RunPutAction(ctx, params)
	logger.Info("finished", "verdict", verdict)

	return &PutActivityResult{
		FlagInfo: params.FlagInfo,
		Verdict:  verdict,
	}, nil
}
