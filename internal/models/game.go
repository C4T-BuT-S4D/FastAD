package models

import (
	"time"
)

type GameSettings struct {
	FlagLifetimeRounds int
	RoundTime          time.Duration

	StartTime time.Time
	EndTime   time.Time

	CheckersBasePath string
}
