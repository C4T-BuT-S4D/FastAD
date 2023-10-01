package models

import (
	"time"

	gspb "github.com/c4t-but-s4d/fastad/pkg/proto/data/game_state"
	"github.com/samber/lo"
	"github.com/uptrace/bun"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type GameState struct {
	bun.BaseModel `bun:"game_state,alias:gs"`

	StartTime   time.Time  `bun:"start_time,notnull"`
	EndTime     *time.Time `bun:"end_time"`
	TotalRounds int        `bun:"total_rounds"`

	Paused bool `bun:"paused"`

	FlagLifetimeRounds int           `bun:"flag_lifetime_rounds,notnull"`
	RoundDuration      time.Duration `bun:"round_duration,notnull"`

	// TODO: game_mode.

	RunningRound      int       `bun:"running_round"`
	RunningRoundStart time.Time `bun:"running_round_start"`
}

func (gs *GameState) ToProto() *gspb.GameState {
	res := &gspb.GameState{
		StartTime:   timestamppb.New(gs.StartTime),
		TotalRounds: int32(gs.TotalRounds),
		Paused:      gs.Paused,

		FlagLifetimeRounds: int32(gs.FlagLifetimeRounds),
		RoundDuration:      durationpb.New(gs.RoundDuration),

		RunningRound:      int32(gs.RunningRound),
		RunningRoundStart: timestamppb.New(gs.RunningRoundStart),
	}
	if gs.EndTime != nil {
		res.EndTime = timestamppb.New(*gs.EndTime)
	}
	return res
}

func NewGameStateFromProto(p *gspb.GameState) *GameState {
	res := &GameState{
		StartTime:   p.StartTime.AsTime(),
		TotalRounds: int(p.TotalRounds),
		Paused:      p.Paused,

		FlagLifetimeRounds: int(p.FlagLifetimeRounds),
		RoundDuration:      p.RoundDuration.AsDuration(),

		RunningRound:      int(p.RunningRound),
		RunningRoundStart: p.RunningRoundStart.AsTime(),
	}
	if p.EndTime != nil {
		res.EndTime = lo.ToPtr(p.EndTime.AsTime())
	}
	return res
}
