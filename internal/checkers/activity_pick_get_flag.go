package checkers

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"

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
	logger := logrus.WithFields(logrus.Fields{
		"activity": "PickGetFlag",
	})

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

	logger.Infof("picked flag %d", flag.ID)
	return &PickGetFlagActivityResult{Flag: flag}, nil
}
