package models

import (
	"time"

	"github.com/uptrace/bun"
)

type (
	ServiceAttackData map[string][]string          // map[team.address][]flag.public
	AttackDataPayload map[string]ServiceAttackData // map[service.name]map[team.address][]flag.public
)

type AttackDataSnapshot struct {
	bun.BaseModel `bun:"attack_data_snapshots,alias:ads"`

	ID        int               `bun:"id,pk,autoincrement"`
	CreatedAt time.Time         `bun:"created_at,notnull"`
	Round     uint64            `bun:"round,notnull"`
	Payload   AttackDataPayload `bun:"payload,type:jsonb,notnull"`
}
