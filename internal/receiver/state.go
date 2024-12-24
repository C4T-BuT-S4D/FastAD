package receiver

import (
	"fmt"
	"math"

	"github.com/samber/lo"

	"github.com/c4t-but-s4d/fastad/internal/models"
	receiverpb "github.com/c4t-but-s4d/fastad/pkg/proto/receiver"
)

type TeamServiceState struct {
	TeamID      int     `json:"team_id"`
	Points      float64 `json:"points"`
	StolenFlags int     `json:"stolen_flags"`
	LostFlags   int     `json:"lost_flags"`
}

func (s *TeamServiceState) ToProto(serviceID int) *receiverpb.State_TeamService {
	return &receiverpb.State_TeamService{
		TeamId:      int64(s.TeamID),
		ServiceId:   int64(serviceID),
		Points:      s.Points,
		StolenFlags: int64(s.StolenFlags),
		LostFlags:   int64(s.LostFlags),
	}
}

func (s *TeamServiceState) Clone() *TeamServiceState {
	return &TeamServiceState{
		Points:      s.Points,
		StolenFlags: s.StolenFlags,
		LostFlags:   s.LostFlags,
	}
}

type ServiceState struct {
	ServiceID    int                       `json:"service_id"`
	DefaultScore float64                   `json:"default_score"`
	TeamStates   map[int]*TeamServiceState `json:"team_states"`
}

func newServiceState(service *models.Service) *ServiceState {
	return &ServiceState{
		ServiceID:    service.ID,
		DefaultScore: service.DefaultScore,
		TeamStates:   make(map[int]*TeamServiceState),
	}
}

func (s *ServiceState) Clone() *ServiceState {
	return &ServiceState{
		ServiceID:    s.ServiceID,
		DefaultScore: s.DefaultScore,
		TeamStates: lo.MapEntries(s.TeamStates, func(key int, value *TeamServiceState) (int, *TeamServiceState) {
			return key, value.Clone()
		}),
	}
}

func (s *ServiceState) Apply(gs *models.GameState, attack *models.Attack) error {
	attackerState := s.getOrCreate(attack.AttackerID)
	victimState := s.getOrCreate(attack.VictimID)

	attackerScore := attackerState.Points
	victimScore := victimState.Points

	scale := 50 * math.Sqrt(gs.Hardness)
	norm := math.Log(math.Log(gs.Hardness)) / 12
	ratingDelta := math.Sqrt(attackerScore) - math.Sqrt(victimScore)
	ratingDeltaNorm := ratingDelta * norm
	attackerDelta := scale / (1 + math.Exp(ratingDeltaNorm))
	victimDelta := -min(victimScore, attackerDelta)
	if !gs.Inflation {
		attackerDelta = min(attackerDelta, -victimDelta)
	}
	attack.AttackerDelta = attackerDelta
	attack.VictimDelta = victimDelta

	attackerState.Points += attackerDelta
	attackerState.StolenFlags++

	victimState.Points += victimDelta
	victimState.LostFlags++

	return nil
}

func (s *ServiceState) ApplyRaw(attacks ...*models.Attack) {
	for _, attack := range attacks {
		attackerState := s.getOrCreate(attack.AttackerID)
		victimState := s.getOrCreate(attack.VictimID)

		attackerState.Points += attack.AttackerDelta
		attackerState.StolenFlags++

		victimState.Points += attack.VictimDelta
		victimState.LostFlags++
	}
}

func (s *ServiceState) ToProto() []*receiverpb.State_TeamService {
	return lo.MapToSlice(s.TeamStates, func(_ int, value *TeamServiceState) *receiverpb.State_TeamService {
		return value.ToProto(s.ServiceID)
	})
}

func (s *ServiceState) getOrCreate(teamID int) *TeamServiceState {
	res, ok := s.TeamStates[teamID]
	if !ok {
		res = &TeamServiceState{Points: s.DefaultScore}
		s.TeamStates[teamID] = res
	}
	return res
}

type State struct {
	ServiceStates map[int]*ServiceState
}

func NewState() *State {
	return &State{
		ServiceStates: make(map[int]*ServiceState),
	}
}

func (s *State) ProcessAttack(gs *models.GameState, service *models.Service, attack *models.Attack) error {
	return s.getOrCreate(service).Apply(gs, attack)
}

func (s *State) ApplyRaw(services map[int]*models.Service, attacks ...*models.Attack) error {
	attacksByService := lo.GroupBy(attacks, func(attack *models.Attack) int {
		return attack.ServiceID
	})

	for serviceID, serviceAttacks := range attacksByService {
		service, ok := services[serviceID]
		if !ok {
			return fmt.Errorf("missing service id %d", serviceID)
		}
		s.getOrCreate(service).ApplyRaw(serviceAttacks...)
	}

	return nil
}

func (s *State) Clone() *State {
	res := &State{
		ServiceStates: lo.MapEntries(s.ServiceStates, func(key int, value *ServiceState) (int, *ServiceState) {
			return key, value.Clone()
		}),
	}
	return res
}

func (s *State) ToProto() *receiverpb.State {
	return &receiverpb.State{
		TeamServices: lo.Flatten(lo.MapToSlice(s.ServiceStates, func(_ int, value *ServiceState) []*receiverpb.State_TeamService {
			return value.ToProto()
		})),
	}
}

func (s *State) getOrCreate(service *models.Service) *ServiceState {
	ss, ok := s.ServiceStates[service.ID]
	if !ok {
		ss = newServiceState(service)
		s.ServiceStates[service.ID] = ss
	}
	return ss
}
