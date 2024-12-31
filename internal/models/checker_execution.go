package models

import (
	"time"

	"github.com/uptrace/bun"
	"google.golang.org/protobuf/types/known/timestamppb"

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

	// Foreign keys.
	Team    *Team    `bun:"rel:belongs-to,join:team_id=id"`
	Service *Service `bun:"rel:belongs-to,join:service_id=id"`
}

func (e *CheckerExecution) ToProto() *checkerpb.Execution {
	return &checkerpb.Execution{
		TeamId:    int64(e.TeamID),
		ServiceId: int64(e.ServiceID),
		Action:    e.Action,
		Status:    e.Status,
		Public:    e.Public,
		Private:   e.Private,
		Command:   e.Command,
		CreatedAt: timestamppb.New(e.CreatedAt),
	}
}
