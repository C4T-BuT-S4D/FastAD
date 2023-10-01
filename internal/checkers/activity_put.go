package checkers

import (
	"context"
	"strings"

	"github.com/c4t-but-s4d/fastad/internal/models"
	checkerpb "github.com/c4t-but-s4d/fastad/pkg/proto/checker"
	"github.com/sirupsen/logrus"
)

type PutActivityParameters struct {
	GameState *models.GameState
	Team      *models.Team
	Service   *models.Service
}

type PutActivityResult struct {
	Verdict *models.CheckerVerdict
}

func (s *ActivityState) PutActivityDefinition(ctx context.Context, params *PutActivityParameters) (*PutActivityResult, error) {
	logger := logrus.WithFields(logrus.Fields{
		"team":    params.Team.Name,
		"service": params.Service.Name,
		"action":  checkerpb.Action_ACTION_PUT,
	})

	flag := &models.Flag{
		Flag:    strings.Repeat("A", 31) + "=",
		Private: "some-flag-id",
	}

	logger.Info("starting")
	verdict := RunPutAction(ctx, params, flag)
	logger.Info("finished: %v", verdict)

	return &PutActivityResult{Verdict: verdict}, nil

}
