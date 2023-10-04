package checkers

import (
	"github.com/c4t-but-s4d/fastad/internal/clients/services"
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
	teamsClient    *teams.Client
	servicesClient *services.Client
}

func NewActivityState(
	teamsClient *teams.Client,
	servicesClient *services.Client,
) *ActivityState {
	return &ActivityState{
		teamsClient:    teamsClient,
		servicesClient: servicesClient,
	}
}
