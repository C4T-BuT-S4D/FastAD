package checkers

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/log"

	"github.com/c4t-but-s4d/fastad/internal/models"
	servicespb "github.com/c4t-but-s4d/fastad/pkg/proto/data/services"
	teamspb "github.com/c4t-but-s4d/fastad/pkg/proto/data/teams"
)

const SaveVerdictActivityName = "SaveVerdict"

type SaveVerdictActivity struct {
	checkersController *Controller
}

func NewSaveVerdictActivity(checkersController *Controller) *SaveVerdictActivity {
	return &SaveVerdictActivity{
		checkersController: checkersController,
	}
}

type SaveVerdictActivityParameters struct {
	Team    *teamspb.Team
	Service *servicespb.Service
	Verdict *Verdict
}

type SaveVerdictActivityResult struct{}

func (s *SaveVerdictActivity) ActivityDefinition(ctx context.Context, params *SaveVerdictActivityParameters) (*SaveRoundDataActivityResult, error) {
	logger := log.With(
		activity.GetLogger(ctx),
		"team", params.Team.Name,
		"service", params.Service.Name,
		"action", params.Verdict.Action,
		"activity", SaveVerdictActivityName,
	)

	logger.Info("saving verdict", "verdict", params.Verdict)
	execution := &models.CheckerExecution{
		ExecutionID: uuid.NewString(),
		TeamID:      int(params.Team.Id),
		ServiceID:   int(params.Service.Id),
		Action:      params.Verdict.Action,
		Status:      params.Verdict.Status,
		Public:      params.Verdict.Public,
		Private:     params.Verdict.Private,
		Command:     params.Verdict.Command,
		CreatedAt:   time.Now(),
	}
	if err := s.checkersController.AddCheckerExecutions(ctx, execution); err != nil {
		return nil, fmt.Errorf("adding checker executions: %w", err)
	}

	return &SaveRoundDataActivityResult{}, nil
}
