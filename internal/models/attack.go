package models

import (
	"time"

	"github.com/uptrace/bun"
)

type Attack struct {
	bun.BaseModel `bun:"attacks,alias:a"`

	ID        int       `bun:"id,pk,autoincrement"`
	CreatedAt time.Time `bun:"created_at,notnull"`

	ServiceID int `bun:"service_id,notnull"`

	// TODO: unique constraints, foreign keys, indices.
	AttackerID int `bun:"attacker_id,notnull,unique:attacker_flag_id"`
	VictimID   int `bun:"victim_id,notnull"`
	FlagID     int `bun:"flag_id,notnull,unique:attacker_flag_id"`

	AttackerDelta float64 `bun:"attacker_delta,notnull"`
	VictimDelta   float64 `bun:"victim_delta,notnull"`

	RequestID string `bun:"request_id,notnull,type:uuid"`
}
