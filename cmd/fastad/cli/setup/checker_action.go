package setup

import (
	"fmt"

	checkerpb "github.com/c4t-but-s4d/fastad/pkg/proto/checker"
)

type CheckerAction string

const (
	CheckerActionCheck CheckerAction = "check"
	CheckerActionPut   CheckerAction = "put"
	CheckerActionGet   CheckerAction = "get"
)

func (a CheckerAction) String() string {
	return string(a)
}

func (a CheckerAction) Validate() error {
	switch a {
	case CheckerActionCheck:
		return nil
	case CheckerActionPut:
		return nil
	case CheckerActionGet:
		return nil
	default:
		return fmt.Errorf("invalid checker action: %s", a)
	}
}

func (a CheckerAction) ToProto() checkerpb.Action {
	switch a {
	case CheckerActionCheck:
		return checkerpb.Action_ACTION_CHECK
	case CheckerActionPut:
		return checkerpb.Action_ACTION_PUT
	case CheckerActionGet:
		return checkerpb.Action_ACTION_GET
	default:
		return checkerpb.Action_ACTION_UNSPECIFIED
	}
}
