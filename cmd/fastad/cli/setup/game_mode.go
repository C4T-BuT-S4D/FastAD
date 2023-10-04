package setup

import (
	"fmt"
	"strings"

	gspb "github.com/c4t-but-s4d/fastad/pkg/proto/data/game_state"
	"gopkg.in/yaml.v3"
)

type GameMode gspb.GameMode

func (g *GameMode) UnmarshalYAML(value *yaml.Node) error {
	var s string
	if err := value.Decode(&s); err != nil {
		return fmt.Errorf("decoding game mode as string: %w", err)
	}

	enumName := fmt.Sprintf("GAME_MODE_%s", strings.ToUpper(s))
	enumValue, ok := gspb.GameMode_value[enumName]
	if !ok || enumValue == 0 {
		return fmt.Errorf("unknown game mode: %s", s)
	}

	*g = GameMode(enumValue)
	return nil
}
