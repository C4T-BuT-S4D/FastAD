package models

import (
	"time"

	"github.com/uptrace/bun"
)

type Flag struct {
	bun.BaseModel `bun:"flags,alias:f"`

	ID      int    `bun:"id,pk,autoincrement"`
	Flag    string `bun:"flag,notnull,unique"`
	Public  string `bun:"public"`
	Private string `bun:"private"`

	Round uint64 `bun:"round,notnull"`

	TeamID    int `bun:"team_id,notnull"`
	ServiceID int `bun:"service_id,notnull"`

	// CreatedAt is set to round workflow start.
	CreatedAt time.Time `bun:"created_at,nullzero,notnull"`

	// Foreign keys.
	Team    *Team    `bun:"rel:belongs-to,join:team_id=id"`
	Service *Service `bun:"rel:belongs-to,join:service_id=id"`
}
