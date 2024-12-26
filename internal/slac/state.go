package slac

import (
	"github.com/samber/lo"

	"github.com/c4t-but-s4d/fastad/internal/models"
	checkerpb "github.com/c4t-but-s4d/fastad/pkg/proto/checker"
	slacpb "github.com/c4t-but-s4d/fastad/pkg/proto/slac"
)

type TeamServiceState struct {
	TeamID       int
	ServiceID    int
	ChecksTotal  int
	ChecksPassed int
	Status       checkerpb.Status
}

func (s *TeamServiceState) Clone() *TeamServiceState {
	return &TeamServiceState{
		TeamID:       s.TeamID,
		ServiceID:    s.ServiceID,
		ChecksTotal:  s.ChecksTotal,
		ChecksPassed: s.ChecksPassed,
		Status:       s.Status,
	}
}

func (s *TeamServiceState) ToProto() *slacpb.TeamServiceState {
	return &slacpb.TeamServiceState{
		TeamId:       int64(s.TeamID),
		ServiceId:    int64(s.ServiceID),
		ChecksTotal:  int64(s.ChecksTotal),
		ChecksPassed: int64(s.ChecksPassed),
		Status:       s.Status,
	}
}

type TeamServiceKey struct {
	TeamID    int
	ServiceID int
}

type State struct {
	TeamServiceStates map[TeamServiceKey]*TeamServiceState
}

func NewState() *State {
	return &State{
		TeamServiceStates: make(map[TeamServiceKey]*TeamServiceState),
	}
}

func (s *State) Apply(execution *models.CheckerExecution) {
	key := TeamServiceKey{
		TeamID:    execution.TeamID,
		ServiceID: execution.ServiceID,
	}

	tss, ok := s.TeamServiceStates[key]

	if !ok {
		tss = &TeamServiceState{
			TeamID:    execution.TeamID,
			ServiceID: execution.ServiceID,
		}
		s.TeamServiceStates[key] = tss
	}

	tss.Status = execution.Status
	tss.ChecksTotal++
	if execution.Status == checkerpb.Status_STATUS_UP {
		tss.ChecksPassed++
	}
}

func (s *State) Clone() *State {
	return &State{
		TeamServiceStates: lo.MapValues(
			s.TeamServiceStates,
			func(value *TeamServiceState, _ TeamServiceKey) *TeamServiceState {
				return value.Clone()
			},
		),
	}
}

func (s *State) ToProto() *slacpb.State {
	return &slacpb.State{
		TeamServiceStates: lo.MapToSlice(
			s.TeamServiceStates,
			func(_ TeamServiceKey, value *TeamServiceState) *slacpb.TeamServiceState {
				return value.ToProto()
			},
		),
	}
}
