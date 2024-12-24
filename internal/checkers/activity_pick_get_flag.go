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
	GameState *models.GameState
	Team      *models.Team
	Service   *models.Service
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
		"team", params.Team.Name,
		"service", params.Service.Name,
		"activity", PickGetFlagActivityName,
	)

	logger.Info("picking flag")

	flag, err := a.checkersController.PickFlag(
		ctx,
		params.Team.ID,
		params.Service.ID,
		params.GameState.RunningRound,
		params.GameState.FlagLifetimeRounds,
	)
	if err != nil {
		return nil, fmt.Errorf("picking flag: %w", err)
	}

	logger.Info("picked flag", "flag", flag)
	return &PickGetFlagActivityResult{Flag: flag}, nil
}
