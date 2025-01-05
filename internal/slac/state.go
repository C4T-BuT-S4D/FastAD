package slac

import (
	"github.com/samber/lo"

	"github.com/c4t-but-s4d/fastad/internal/models"
	checkerpb "github.com/c4t-but-s4d/fastad/pkg/proto/checker"
	slacpb "github.com/c4t-but-s4d/fastad/pkg/proto/slac"
)

const keepLastChecks = 3

type TeamServiceKey struct {
	TeamID    int
	ServiceID int
}

type State struct {
	TeamServiceStates map[TeamServiceKey]*slacpb.TeamServiceState
}

func NewState() *State {
	return &State{
		TeamServiceStates: make(map[TeamServiceKey]*slacpb.TeamServiceState),
	}
}

func (s *State) Apply(execution *models.CheckerExecution) {
	key := TeamServiceKey{
		TeamID:    execution.TeamID,
		ServiceID: execution.ServiceID,
	}

	tss, ok := s.TeamServiceStates[key]

	if !ok {
		tss = &slacpb.TeamServiceState{
			TeamId:    int64(execution.TeamID),
			ServiceId: int64(execution.ServiceID),
		}
		s.TeamServiceStates[key] = tss
	}

	tss.ChecksTotal++
	if execution.Status == checkerpb.Status_STATUS_UP {
		tss.ChecksPassed++
	}
	tss.CheckStatuses = append(tss.GetCheckStatuses(), &slacpb.TeamServiceState_CheckStatus{
		Status:  execution.Status,
		Message: execution.Public,
	})
	if len(tss.GetCheckStatuses()) > keepLastChecks {
		tss.CheckStatuses = tss.GetCheckStatuses()[len(tss.GetCheckStatuses())-keepLastChecks:]
	}
}

func (s *State) Clone() *State {
	return &State{
		TeamServiceStates: lo.MapValues(
			s.TeamServiceStates,
			func(value *slacpb.TeamServiceState, _ TeamServiceKey) *slacpb.TeamServiceState {
				return value.CloneVT()
			},
		),
	}
}

func (s *State) ToProto() *slacpb.State {
	return &slacpb.State{
		TeamServiceStates: lo.Values(s.TeamServiceStates),
	}
}
