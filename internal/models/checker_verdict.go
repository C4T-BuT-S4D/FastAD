package models

import (
	"fmt"

	checkerpb "github.com/c4t-but-s4d/fastad/pkg/proto/checker"
)

type CheckerVerdict struct {
	Action  checkerpb.Action
	Public  string
	Private string
	Status  checkerpb.Status
	Command string
}

func (v *CheckerVerdict) String() string {
	return fmt.Sprintf("%v %v", v.Action, v.Status)
}

func (v *CheckerVerdict) IsUp() bool {
	return v.Status == checkerpb.Status_STATUS_UP
}
