package checkers

import (
	"context"
	"fmt"

	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/log"

	"github.com/c4t-but-s4d/fastad/internal/models"
	gspb "github.com/c4t-but-s4d/fastad/pkg/proto/data/game_state"
	servicespb "github.com/c4t-but-s4d/fastad/pkg/proto/data/services"
	teamspb "github.com/c4t-but-s4d/fastad/pkg/proto/data/teams"
)

const PickGetFlagActivityName = "PickGetFlag"

type PickGetFlagActivity struct {
	checkersController *Controller
}

func NewPickGetFlagActivity(checkersController *Controller) *PickGetFlagActivity {
	return &PickGetFlagActivity{checkersController: checkersController}
}

type PickGetFlagActivityParameters struct {
	GameState *gspb.GameState
	Team      *teamspb.Team
	Service   *servicespb.Service
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
		int(params.Team.Id),
		int(params.Service.Id),
		params.GameState.RunningRound,
		params.GameState.FlagLifetimeRounds,
	)
	if err != nil {
		return nil, fmt.Errorf("picking flag: %w", err)
	}

	logger.Info("picked flag", "flag", flag)
	return &PickGetFlagActivityResult{Flag: flag}, nil
}
