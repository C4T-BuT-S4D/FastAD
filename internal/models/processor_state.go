package models

import (
	"time"

	"github.com/uptrace/bun"
)

type ProcessorState struct {
	bun.BaseModel `bun:"processor_states,alias:ps"`

	ProcessorID string    `bun:"processor_id,pk"`
	EntityID    int       `bun:"entity_id,index"`
	ProcessedAt time.Time `bun:"processed_at"`
}
