package models

import (
	"time"

	"github.com/uptrace/bun"
)

type GameState struct {
	bun.BaseModel `bun:"game_state,alias:gs"`

	StartTime time.Time  `bun:"start_time,notnull"`
	EndTime   *time.Time `bun:"end_time"`

	Paused bool `bun:"paused"`

	FlagLifetimeRounds int           `bun:"flag_lifetime_rounds,notnull"`
	RoundDuration      time.Duration `bun:"round_duration,notnull"`

	// TODO: game_mode.

	RunningRound      int       `bun:"running_round"`
	RunningRoundStart time.Time `bun:"running_round_start"`
}
