package checkers

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"strings"

	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/log"

	"github.com/c4t-but-s4d/fastad/internal/models"
	"github.com/c4t-but-s4d/fastad/pkg/clients/gamestate"
	"github.com/c4t-but-s4d/fastad/pkg/modelsutil"
	checkerpb "github.com/c4t-but-s4d/fastad/pkg/proto/checker"
	gspb "github.com/c4t-but-s4d/fastad/pkg/proto/data/game_state"
	servicespb "github.com/c4t-but-s4d/fastad/pkg/proto/data/services"
	teamspb "github.com/c4t-but-s4d/fastad/pkg/proto/data/teams"
)

const PrepareRoundActivityName = "PrepareRound"

type PrepareRoundActivity struct {
	checkersController *Controller
	gameStateClient    *gamestate.Client
}

func NewPrepareRoundActivity(
	checkersController *Controller,
	gameStateClient *gamestate.Client,
) *PrepareRoundActivity {
	return &PrepareRoundActivity{
		checkersController: checkersController,
		gameStateClient:    gameStateClient,
	}
}

type PrepareRoundActivityParameters struct {
	GameState *gspb.GameState
	Teams     []*teamspb.Team
	Services  []*servicespb.Service
}

type FQFlagInfo struct {
	Team    *teamspb.Team
	Service *servicespb.Service
	Flag    *models.Flag
}

type PrepareRoundActivityResult struct {
	Flags []*FQFlagInfo
}

func (a *PrepareRoundActivity) ActivityDefinition(ctx context.Context, params *PrepareRoundActivityParameters) (*PrepareRoundActivityResult, error) {
	logger := log.With(
		activity.GetLogger(ctx),
		"activity", PrepareRoundActivityName,
		"running_round", params.GameState.GetRunningRound(),
	)

	logger.Info("starting")
	flags, err := a.prepareRoundPutState(ctx, params, logger)
	if err != nil {
		return nil, fmt.Errorf("preparing round put state: %w", err)
	}
	logger.Info("finished")

	return &PrepareRoundActivityResult{Flags: flags}, nil
}

func (a *PrepareRoundActivity) prepareRoundPutState(
	ctx context.Context,
	params *PrepareRoundActivityParameters,
	logger log.Logger,
) ([]*FQFlagInfo, error) {
	logger.Info("updating round", "round", params.GameState.GetRunningRound())
	if _, err := a.gameStateClient.UpdateRound(ctx, &gspb.UpdateRoundRequest{
		RunningRound:      params.GameState.GetRunningRound(),
		RunningRoundStart: params.GameState.GetRunningRoundStart(),
	}); err != nil {
		return nil, fmt.Errorf("updating round: %w", err)
	}

	logger.Info(
		"preparing flags for teams and services",
		"teams", len(params.Teams),
		"services", len(params.Services),
	)

	flags := make([]*FQFlagInfo, 0, len(params.Teams)*len(params.Services))
	flagModels := make([]*models.Flag, 0, len(params.Teams)*len(params.Services))
	for _, team := range params.Teams {
		for _, service := range params.Services {
			for range modelsutil.ServiceCheckerRunCount(service, checkerpb.Action_ACTION_PUT) {
				flag := &models.Flag{
					Flag:      generateFlag(service),
					TeamID:    int(team.GetId()),
					ServiceID: int(service.GetId()),
					Round:     params.GameState.GetRunningRound(),
					CreatedAt: params.GameState.GetRunningRoundStart().AsTime(),
				}
				flags = append(flags, &FQFlagInfo{
					Team:    team,
					Service: service,
					Flag:    flag,
				})
				flagModels = append(flagModels, flag)
			}
		}
	}

	if len(flagModels) == 0 {
		logger.Warn("no flags to insert")
		return flags, nil
	}

	logger.Info("inserting flags", "flags", len(flagModels))

	if err := a.checkersController.AddFlags(ctx, flagModels); err != nil {
		return nil, fmt.Errorf("adding flags: %w", err)
	}

	logger.Info("inserted flags", "flags", len(flagModels))

	return flags, nil
}

func generateFlag(service *servicespb.Service) string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const length = 30
	var result [length]byte

	for i := range result {
		randomIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			panic(err)
		}
		result[i] = charset[randomIndex.Int64()]
	}

	return strings.ToUpper(service.GetName()[:1]) + string(result[:]) + "="
}
