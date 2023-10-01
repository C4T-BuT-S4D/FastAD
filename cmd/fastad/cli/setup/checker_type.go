package setup

import (
	"fmt"

	checkerpb "github.com/c4t-but-s4d/fastad/pkg/proto/checker"
)

type CheckerType string

const (
	CheckerTypeLegacy CheckerType = "legacy"
)

func (a CheckerType) String() string {
	return string(a)
}

func (a CheckerType) Validate() error {
	switch a {
	case CheckerTypeLegacy:
		return nil
	default:
		return fmt.Errorf("invalid checker type: %s", a)
	}
}

func (a CheckerType) ToProto() checkerpb.Type {
	switch a {
	case CheckerTypeLegacy:
		return checkerpb.Type_TYPE_LEGACY
	default:
		return checkerpb.Type_TYPE_UNSPECIFIED
	}
}
