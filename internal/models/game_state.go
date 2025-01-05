package models

import (
	"time"

	"github.com/uptrace/bun"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	gspb "github.com/c4t-but-s4d/fastad/pkg/proto/data/game_state"
)

type GameState struct {
	bun.BaseModel `bun:"game_state,alias:gs"`

	ID int `bun:"id,pk"`

	StartTime   time.Time  `bun:"start_time,notnull"`
	EndTime     *time.Time `bun:"end_time"`
	TotalRounds uint64     `bun:"total_rounds"`

	Paused   bool `bun:"paused"`
	Finished bool `bun:"finished"`

	FlagLifetimeRounds uint64        `bun:"flag_lifetime_rounds,notnull"`
	RoundDuration      time.Duration `bun:"round_duration,notnull"`

	// TODO: game_mode.
	Hardness  float64 `bun:"hardness,notnull"`
	Inflation bool    `bun:"inflation"`

	RunningRound      uint64    `bun:"running_round"`
	RunningRoundStart time.Time `bun:"running_round_start"`
}

func (gs *GameState) ToProto() *gspb.GameState {
	res := &gspb.GameState{
		StartTime:   timestamppb.New(gs.StartTime),
		TotalRounds: gs.TotalRounds,

		Paused:   gs.Paused,
		Finished: gs.Finished,

		FlagLifetimeRounds: gs.FlagLifetimeRounds,
		RoundDuration:      durationpb.New(gs.RoundDuration),

		RunningRound:      gs.RunningRound,
		RunningRoundStart: timestamppb.New(gs.RunningRoundStart),

		Hardness:  gs.Hardness,
		Inflation: gs.Inflation,
	}
	if gs.EndTime != nil {
		res.EndTime = timestamppb.New(*gs.EndTime)
	}
	return res
}
