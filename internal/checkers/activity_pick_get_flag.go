package checkers

import (
	"context"
	"fmt"

	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/log"

	"github.com/c4t-but-s4d/fastad/internal/models"
)

const PickGetFlagActivityName = "PickGetFlag"

type PickGetFlagActivity struct {
	checkersController *Controller
}

func NewPickGetFlagActivity(checkersController *Controller) *PickGetFlagActivity {
	return &PickGetFlagActivity{checkersController: checkersController}
}

type PickGetFlagActivityParameters struct {
	TeamID             int
	ServiceID          int
	RunningRound       uint64
	FlagLifetimeRounds uint64
}

type PickGetFlagActivityResult struct {
	Flag *models.Flag
}

func (a *PickGetFlagActivity) ActivityDefinition(
	ctx context.Context,
	params *PickGetFlagActivityParameters,
) (*PickGetFlagActivityResult, error) {
	logger := log.With(
		activity.GetLogger(ctx),
		"team", params.TeamID,
		"service", params.ServiceID,
		"activity", PickGetFlagActivityName,
	)

	logger.Info("picking flag")

	flag, err := a.checkersController.PickFlag(
		ctx,
		params.TeamID,
		params.ServiceID,
		params.RunningRound,
		params.FlagLifetimeRounds,
	)
	if err != nil {
		return nil, fmt.Errorf("picking flag: %w", err)
	}

	logger.Info("picked flag", "flag", flag)
	return &PickGetFlagActivityResult{Flag: flag}, nil
}
