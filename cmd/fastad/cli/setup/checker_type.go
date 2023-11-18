package setup

import (
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"

	checkerpb "github.com/c4t-but-s4d/fastad/pkg/proto/checker"
)

type CheckerType checkerpb.Type

func (t *CheckerType) UnmarshalYAML(value *yaml.Node) error {
	var s string
	if err := value.Decode(&s); err != nil {
		return fmt.Errorf("decoding checker type as string: %w", err)
	}

	if s == "" {
		*t = CheckerType(checkerpb.Type_TYPE_LEGACY)
		return nil
	}

	enumName := fmt.Sprintf("TYPE_%s", strings.ToUpper(s))
	enumValue, ok := checkerpb.Type_value[enumName]
	if !ok || enumValue == 0 {
		return fmt.Errorf("unknown checker type: %s", s)
	}

	*t = CheckerType(enumValue)
	return nil
}
