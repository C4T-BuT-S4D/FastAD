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

	ID int `bun:"id,pk"`

	StartTime   time.Time  `bun:"start_time,notnull"`
	EndTime     *time.Time `bun:"end_time"`
	TotalRounds uint64     `bun:"total_rounds"`

	Status gspb.GameStatus `bun:"status,notnull,default:1"`

	FlagLifetimeRounds uint64        `bun:"flag_lifetime_rounds,notnull"`
	RoundDuration      time.Duration `bun:"round_duration,notnull"`

	Hardness  float64 `bun:"hardness,notnull"`
	Inflation bool    `bun:"inflation"`

	RunningRound      uint64    `bun:"running_round"`
	RunningRoundStart time.Time `bun:"running_round_start"`
}

func NewGameStateFromProto(gs *gspb.GameState) *GameState {
	res := &GameState{
		StartTime:   gs.GetStartTime().AsTime(),
		TotalRounds: gs.GetTotalRounds(),

		Status: gs.GetStatus(),

		FlagLifetimeRounds: gs.GetFlagLifetimeRounds(),
		RoundDuration:      gs.GetRoundDuration().AsDuration(),

		RunningRound:      gs.GetRunningRound(),
		RunningRoundStart: gs.GetRunningRoundStart().AsTime(),

		Hardness:  gs.GetHardness(),
		Inflation: gs.GetInflation(),
	}
	if gs.GetEndTime() != nil {
		res.EndTime = lo.ToPtr(gs.GetEndTime().AsTime())
	}
	return res
}

func (gs *GameState) ToProto() *gspb.GameState {
	res := &gspb.GameState{
		StartTime:   timestamppb.New(gs.StartTime),
		TotalRounds: gs.TotalRounds,

		Status: gs.Status,

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
