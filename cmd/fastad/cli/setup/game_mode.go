package setup

import (
	"fmt"

	gspb "github.com/c4t-but-s4d/fastad/pkg/proto/data/game_state"
)

type GameMode string

const (
	GameModeClassic GameMode = "classic"
)

func (g GameMode) String() string {
	return string(g)
}

func (g GameMode) Validate() error {
	switch g {
	case GameModeClassic:
		return nil
	default:
		return fmt.Errorf("invalid game mode: %s", g)
	}
}

func (g GameMode) ToProto() gspb.GameMode {
	switch g {
	case GameModeClassic:
		return gspb.GameMode_GAME_MODE_CLASSIC
	default:
		return gspb.GameMode_GAME_MODE_UNSPECIFIED
	}
}
