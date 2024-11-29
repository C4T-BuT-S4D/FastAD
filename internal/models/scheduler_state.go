package models

import (
	"time"

	"github.com/uptrace/bun"
)

type SchedulerState struct {
	bun.BaseModel `bun:"scheduler_states,alias:ss"`

	SchedulerID     string    `bun:"scheduler_id,pk"`
	ExpectedNextRun time.Time `bun:"expected_next_run"`
}
