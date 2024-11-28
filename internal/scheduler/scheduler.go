package scheduler

import (
	"context"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"go.temporal.io/sdk/client"

	"github.com/c4t-but-s4d/fastad/internal/checkers"
)

type Scheduler struct {
	checkInterval time.Duration

	temporalClient client.Client
}

func New(checkInterval time.Duration, temporalClient client.Client) *Scheduler {
	return &Scheduler{
		checkInterval: checkInterval,

		temporalClient: temporalClient,
	}
}

func (t *Scheduler) Run(ctx context.Context) {
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

func (t *Scheduler) RunOnce(ctx context.Context) {
	// TODO: check interval/game mode/etc changed.

	if err := t.scheduleTasks(ctx); err != nil {
		logrus.Errorf("scheduling tasks: %v", err)
	} else {
		logrus.Info("scheduled tasks")
	}
}

func (t *Scheduler) scheduleTasks(ctx context.Context) error {
	workflowRun, err := t.temporalClient.ExecuteWorkflow(
		context.Background(),
		client.StartWorkflowOptions{
			TaskQueue: "checkers",
		},
		"RoundWorkflowDefinition",
		checkers.WorkflowParameters{},
	)
	if err != nil {
		return fmt.Errorf("executing workflow: %w", err)
	}

	if err := workflowRun.Get(ctx, nil); err != nil {
		return fmt.Errorf("waiting for workflow: %w", err)
	}

	return nil
}
