package models

import (
	"github.com/uptrace/bun"

	checkerpb "github.com/c4t-but-s4d/fastad/pkg/proto/checker"
)

type CheckerExecution struct {
	bun.BaseModel `bun:"checker_execution,alias:ce"`

	ID int `bun:"id,pk,autoincrement"`

	TSSHistoryEntryID int `bun:"tss_history_entry_id,notnull"`

	Status  checkerpb.Status `bun:"status,notnull"`
	Public  string           `bun:"public"`
	Private string           `bun:"private"`
	Command string           `bun:"command"`
}
