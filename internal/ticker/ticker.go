package ticker

import (
	"context"
	"fmt"
	"time"

	"github.com/c4t-but-s4d/fastad/internal/checkers"
	"github.com/c4t-but-s4d/fastad/internal/models"
	"github.com/sirupsen/logrus"
	"go.temporal.io/sdk/client"
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
	workflowOptions := client.StartWorkflowOptions{
		TaskQueue: "checkers",
	}
	workflowRun, err := t.temporalClient.ExecuteWorkflow(context.Background(), workflowOptions, "WorkflowDefinition", checkers.WorkflowParameters{
		GameSettings: &models.GameSettings{
			FlagLifetimeRounds: 10,
			RoundTime:          10,
			StartTime:          time.Now(),
			EndTime:            time.Now(),
			CheckersBasePath:   "checkers",
		},
	})
	if err != nil {
		return fmt.Errorf("executing workflow: %v", err)
	}

	if err := workflowRun.Get(ctx, nil); err != nil {
		return fmt.Errorf("waiting for workflow: %v", err)
	}

	return nil
}
