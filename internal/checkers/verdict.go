package checkers

import (
	"fmt"

	checkerpb "github.com/c4t-but-s4d/fastad/pkg/proto/checker"
)

type Verdict struct {
	Action  checkerpb.Action
	Status  checkerpb.Status
	Public  string
	Private string
	Command string
}

func (v *Verdict) String() string {
	return fmt.Sprintf("%v %v", v.Action, v.Status)
}

func (v *Verdict) IsUp() bool {
	return v.Status == checkerpb.Status_STATUS_UP
}
