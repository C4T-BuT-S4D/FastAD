package checkers

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/log"

	"github.com/c4t-but-s4d/fastad/internal/models"
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
	TeamID    int
	ServiceID int
	Verdict   *Verdict
}

type SaveVerdictActivityResult struct{}

func (s *SaveVerdictActivity) ActivityDefinition(ctx context.Context, params *SaveVerdictActivityParameters) (*SaveRoundDataActivityResult, error) {
	logger := log.With(
		activity.GetLogger(ctx),
		"team", params.TeamID,
		"service", params.ServiceID,
		"action", params.Verdict.Action,
		"activity", SaveVerdictActivityName,
	)

	logger.Info("saving verdict", "verdict", params.Verdict)
	execution := &models.CheckerExecution{
		ExecutionID: uuid.NewString(),
		TeamID:      params.TeamID,
		ServiceID:   params.ServiceID,
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
