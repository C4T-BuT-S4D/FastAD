package models

import (
	"time"

	"github.com/samber/lo"
	"github.com/uptrace/bun"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	gspb "github.com/c4t-but-s4d/fastad/pkg/proto/data/game_state"
)

type GameState struct {
	bun.BaseModel `bun:"game_state,alias:gs"`

	ID int `bun:"id,pk" json:"id"`

	StartTime   time.Time  `bun:"start_time,notnull" json:"start_time"`
	EndTime     *time.Time `bun:"end_time" json:"end_time,omitempty"`
	TotalRounds uint64     `bun:"total_rounds" json:"total_rounds,omitempty"`

	Paused bool `bun:"paused" json:"paused,omitempty"`

	FlagLifetimeRounds uint64        `bun:"flag_lifetime_rounds,notnull" json:"flag_lifetime_rounds"`
	RoundDuration      time.Duration `bun:"round_duration,notnull" json:"round_duration"`

	// TODO: game_mode.
	Hardness  float64 `bun:"hardness,notnull" json:"hardness,omitempty"`
	Inflation bool    `bun:"inflation" json:"inflation,omitempty"`

	RunningRound      uint64    `bun:"running_round" json:"running_round"`
	RunningRoundStart time.Time `bun:"running_round_start" json:"running_round_start"`
}

func (gs *GameState) ToProto() *gspb.GameState {
	res := &gspb.GameState{
		StartTime:   timestamppb.New(gs.StartTime),
		TotalRounds: gs.TotalRounds,
		Paused:      gs.Paused,

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

func NewGameStateFromProto(p *gspb.GameState) *GameState {
	res := &GameState{
		StartTime:   p.StartTime.AsTime(),
		TotalRounds: p.TotalRounds,
		Paused:      p.Paused,

		FlagLifetimeRounds: p.FlagLifetimeRounds,
		RoundDuration:      p.RoundDuration.AsDuration(),

		RunningRound:      p.RunningRound,
		RunningRoundStart: p.RunningRoundStart.AsTime(),

		Hardness:  p.Hardness,
		Inflation: p.Inflation,
	}
	if p.EndTime != nil {
		res.EndTime = lo.ToPtr(p.EndTime.AsTime())
	}
	return res
}
