package common

import (
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"

	checkerpb "github.com/c4t-but-s4d/fastad/pkg/proto/checker"
)

type CheckerAction checkerpb.Action

//goland:noinspection GoMixedReceiverTypes
func (a *CheckerAction) UnmarshalYAML(value *yaml.Node) error {
	var s string
	if err := value.Decode(&s); err != nil {
		return fmt.Errorf("decoding checker action as string: %w", err)
	}

	enumName := fmt.Sprintf("ACTION_%s", strings.ToUpper(s))
	enumValue, ok := checkerpb.Action_value[enumName]
	if !ok || enumValue == 0 {
		return fmt.Errorf("unknown checker action: %s", s)
	}

	*a = CheckerAction(enumValue)
	return nil
}

//goland:noinspection GoMixedReceiverTypes
func (a CheckerAction) MarshalYAML() (any, error) {
	enumName, ok := checkerpb.Action_name[int32(a)]
	if !ok {
		return nil, fmt.Errorf("unknown checker action: %d", a)
	}

	return strings.ToLower(strings.TrimPrefix(enumName, "ACTION_")), nil
}
