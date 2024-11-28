package checkers

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"strings"

	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/c4t-but-s4d/fastad/internal/models"
	checkerpb "github.com/c4t-but-s4d/fastad/pkg/proto/checker"
	gspb "github.com/c4t-but-s4d/fastad/pkg/proto/data/game_state"
)

type PrepareRoundActivityParameters struct {
	GameState *models.GameState
	Teams     []*models.Team
	Services  []*models.Service
}

type FQFlagInfo struct {
	Team    *models.Team
	Service *models.Service
	Flag    *models.Flag
}

type PrepareRoundActivityResult struct {
	Flags []*FQFlagInfo
}

func (s *ActivityState) PrepareRoundActivityDefinition(ctx context.Context, params *PrepareRoundActivityParameters) (*PrepareRoundActivityResult, error) {
	logger := logrus.WithFields(logrus.Fields{
		"action": "PrepareRound",
	})

	logger.Info("starting")
	flags, err := s.prepareRoundPutState(ctx, params, logger)
	if err != nil {
		return nil, fmt.Errorf("preparing round put state: %w", err)
	}
	logger.Infof("finished")

	return &PrepareRoundActivityResult{Flags: flags}, nil
}

func (s *ActivityState) prepareRoundPutState(
	ctx context.Context,
	params *PrepareRoundActivityParameters,
	logger *logrus.Entry,
) ([]*FQFlagInfo, error) {
	logger.Infof("preparing flags for %d teams and %d services", len(params.Teams), len(params.Services))

	flags := make([]*FQFlagInfo, 0, len(params.Teams)*len(params.Services))
	flagModels := make([]*models.Flag, 0, len(params.Teams)*len(params.Services))
	for _, team := range params.Teams {
		for _, service := range params.Services {
			for range service.GetRunCount(checkerpb.Action_ACTION_PUT) {
				flag := &models.Flag{
					Flag:      generateFlag(service),
					TeamID:    team.ID,
					ServiceID: service.ID,
					CreatedAt: params.GameState.RunningRoundStart,
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

	logger.Infof(
		"inserting %d flags, bumping round to %v",
		len(flagModels),
		params.GameState.RunningRound+1,
	)

	if err := s.checkersController.AddFlags(ctx, flagModels); err != nil {
		return nil, fmt.Errorf("adding flags: %w", err)
	}

	logger.Infof("inserted %d flags: %+v", len(flagModels), flagModels[0])

	if _, err := s.gameStateClient.UpdateRound(ctx, &gspb.UpdateRoundRequest{
		RunningRound:      params.GameState.RunningRound,
		RunningRoundStart: timestamppb.New(params.GameState.RunningRoundStart),
	}); err != nil {
		return nil, fmt.Errorf("updating round: %w", err)
	}

	return flags, nil
}

func generateFlag(service *models.Service) string {
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

	return strings.ToUpper(service.Name[:1]) + string(result[:]) + "="
}
