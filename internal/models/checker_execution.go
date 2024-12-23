package models

import (
	"time"

	"github.com/uptrace/bun"

	checkerpb "github.com/c4t-but-s4d/fastad/pkg/proto/checker"
)

type CheckerExecution struct {
	bun.BaseModel `bun:"checker_executions,alias:ce"`

	ID int `bun:"id,pk,autoincrement"`

	ExecutionID string `bun:"execution_id,notnull,unique"`
	TeamID      int    `bun:"team_id,notnull"`
	ServiceID   int    `bun:"service_id,notnull"`

	Action checkerpb.Action `bun:"action,notnull"`
	Status checkerpb.Status `bun:"status,notnull"`

	Public  string `bun:"public"`
	Private string `bun:"private"`
	Command string `bun:"command"`

	CreatedAt time.Time `bun:"created_at,nullzero,notnull,default:current_timestamp"`
}
