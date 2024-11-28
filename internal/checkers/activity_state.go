package checkers

import (
	"github.com/c4t-but-s4d/fastad/internal/clients/gamestate"
	"github.com/c4t-but-s4d/fastad/internal/clients/services"
	"github.com/c4t-but-s4d/fastad/internal/clients/teams"
)

const (
	ActivityFetchDataName = "ActivityFetchData"
	ActivityCheckName     = "ActivityCheck"
	ActivityGetName       = "ActivityGet"
	ActivityLastName      = "ActivityLast"

	// Round-related activities.

	// ActivityPrepareRoundStateName is an activity name for preparing round state.
	ActivityPrepareRoundStateName = "ActivityPrepareRoundState"
	// ActivitySaveRoundStateName is an activity name for saving round state.
	ActivitySaveRoundStateName = "ActivitySaveRoundState"
	// ActivityPutName is an activity name for actually calling checker PUT action.
	ActivityPutName = "ActivityPut"
)

type ActivityState struct {
	teamsClient        *teams.Client
	servicesClient     *services.Client
	gameStateClient    *gamestate.Client
	checkersController *Controller
}

func NewActivityState(
	teamsClient *teams.Client,
	servicesClient *services.Client,
	gameStateClient *gamestate.Client,
	checkersCntroller *Controller,
) *ActivityState {
	return &ActivityState{
		teamsClient:        teamsClient,
		servicesClient:     servicesClient,
		gameStateClient:    gameStateClient,
		checkersController: checkersCntroller,
	}
}
