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

	AttackerID int `bun:"attacker_id,notnull,unique:attacker_flag_id"`
	VictimID   int `bun:"victim_id,notnull"`
	FlagID     int `bun:"flag_id,notnull,unique:attacker_flag_id"`

	AttackerDelta float64 `bun:"attacker_delta,notnull"`
	VictimDelta   float64 `bun:"victim_delta,notnull"`

	RequestID string `bun:"request_id,notnull,type:uuid"`

	// Foreign keys.
	Service  *Service `bun:"rel:belongs-to,join:service_id=id"`
	Attacker *Team    `bun:"rel:belongs-to,join:attacker_id=id"`
	Victim   *Team    `bun:"rel:belongs-to,join:victim_id=id"`
	Flag     *Flag    `bun:"rel:belongs-to,join:flag_id=id"`
}
