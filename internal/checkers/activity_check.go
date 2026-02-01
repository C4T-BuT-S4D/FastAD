package checkers

import (
	"context"

	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/log"

	checkerpb "github.com/c4t-but-s4d/fastad/pkg/proto/checker"
	gspb "github.com/c4t-but-s4d/fastad/pkg/proto/data/game_state"
	servicespb "github.com/c4t-but-s4d/fastad/pkg/proto/data/services"
	teamspb "github.com/c4t-but-s4d/fastad/pkg/proto/data/teams"
)

const CheckActivityName = "Check"

type CheckActivity struct{}

func NewCheckActivity() *CheckActivity {
	return &CheckActivity{}
}

type CheckActivityParameters struct {
	GameState *gspb.GameState
	Team      *teamspb.Team
	Service   *servicespb.Service
}

type CheckActivityResult struct {
	Verdict *Verdict
}

func (*CheckActivity) ActivityDefinition(
	ctx context.Context,
	params *CheckActivityParameters,
) (*CheckActivityResult, error) {
	logger := log.With(
		activity.GetLogger(ctx),
		"team", params.Team.GetId(),
		"service", params.Service.GetId(),
		"action", checkerpb.Action_ACTION_CHECK,
		"activity", CheckActivityName,
	)

	logger.Info("starting")

	verdict := RunCheckAction(ctx, params)
	logger.Info("finished", "verdict", verdict)

	return &CheckActivityResult{Verdict: verdict}, nil
}
