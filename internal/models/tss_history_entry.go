package models

import (
	checkerpb "github.com/c4t-but-s4d/fastad/pkg/proto/checker"
	"github.com/uptrace/bun"
)

type TSSHistoryEntry struct {
	bun.BaseModel `bun:"tss_history_entry,alias:tss_he"`

	ID int `bun:"id,pk,autoincrement"`

	TeamID    int `bun:"team_id,unique:team_service_round,notnull"`
	ServiceID int `bun:"service_id,unique:team_service_round,notnull"`
	Round     int `bun:"round,unique:team_service_round,notnull"`

	Checks int `bun:"checks"`
	// TODO: add CONSTRAINT sla_valid CHECK ( checks >= 0 AND checks_passed >= 0 AND checks_passed <= checks ).
	ChecksPassed int `bun:"checks_passed"`

	Status checkerpb.Status `bun:"status,notnull"`
	Public string           `bun:"public"`

	// Relations.
	Team    *Team    `bun:"rel:belongs-to,join:team_id=id"`
	Service *Service `bun:"rel:belongs-to,join:service_id=id"`

	CheckerExecutions []CheckerExecution `bun:"rel:has-many,join:id=tss_history_entry_id"`
}
