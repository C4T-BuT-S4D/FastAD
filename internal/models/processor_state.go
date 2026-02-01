package models

import (
	"time"

	"github.com/uptrace/bun"
)

type SlacProcessedItem struct {
	bun.BaseModel `bun:"slac_processed_items,alias:spi"`

	ID                 int       `bun:"id,pk,autoincrement"`
	CheckerExecutionID int       `bun:"checker_execution_id,unique"`
	ProcessedAt        time.Time `bun:"processed_at"`

	CheckerExecution *CheckerExecution `bun:"rel:belongs-to,join:checker_execution_id=id"`
}
