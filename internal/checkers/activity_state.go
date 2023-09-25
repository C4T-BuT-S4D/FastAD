package checkers

import (
	"github.com/c4t-but-s4d/fastad/internal/clients/teams"
)

const (
	ActivityFetchDataName = "ActivityFetchData"
	ActivityCheckName     = "ActivityCheck"
	ActivityPutName       = "ActivityPut"
	ActivityGetName       = "ActivityGet"
	ActivityLastName      = "ActivityLast"
)

type ActivityState struct {
	teamsClient *teams.Client
}

func NewActivityState(teamsClient *teams.Client) *ActivityState {
	return &ActivityState{
		teamsClient: teamsClient,
	}
}
