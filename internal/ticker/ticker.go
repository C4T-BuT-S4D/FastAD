package ticker

import (
	"context"
	"fmt"
	"time"

	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
	"go.temporal.io/sdk/client"

	"github.com/c4t-but-s4d/fastad/internal/checkers"
	"github.com/c4t-but-s4d/fastad/internal/models"
)

type Ticker struct {
	checkInterval time.Duration

	temporalClient client.Client
}

func NewTicker(checkInterval time.Duration, temporalClient client.Client) *Ticker {
	return &Ticker{
		checkInterval: checkInterval,

		temporalClient: temporalClient,
	}
}

func (t *Ticker) Run(ctx context.Context) {
	t.RunOnce(ctx)

	ticker := time.NewTicker(t.checkInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			t.RunOnce(ctx)
		}
	}
}

func (t *Ticker) RunOnce(ctx context.Context) {
	// TODO: check interval/game mode/etc changed.

	if err := t.scheduleTasks(ctx); err != nil {
		logrus.Errorf("scheduling tasks: %v", err)
	} else {
		logrus.Info("scheduled tasks")
	}
}

func (t *Ticker) scheduleTasks(ctx context.Context) error {
	workflowRun, err := t.temporalClient.ExecuteWorkflow(
		context.Background(),
		client.StartWorkflowOptions{
			TaskQueue: "checkers",
		},
		"WorkflowDefinition",
		checkers.WorkflowParameters{
			GameState: &models.GameState{
				StartTime: time.Now(),
				EndTime:   lo.ToPtr(time.Now()),

				FlagLifetimeRounds: 10,
				RoundDuration:      time.Second * 10,

				RunningRound: 1,
			},
		},
	)
	if err != nil {
		return fmt.Errorf("executing workflow: %v", err)
	}

	if err := workflowRun.Get(ctx, nil); err != nil {
		return fmt.Errorf("waiting for workflow: %v", err)
	}

	return nil
}
