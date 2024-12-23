package receiver

import (
	"fmt"
	"math"

	"github.com/samber/lo"

	"github.com/c4t-but-s4d/fastad/internal/models"
	receiverpb "github.com/c4t-but-s4d/fastad/pkg/proto/receiver"
)

type TeamState struct {
	TeamID      int     `json:"team_id"`
	Points      float64 `json:"points"`
	StolenFlags int     `json:"stolen_flags"`
	LostFlags   int     `json:"lost_flags"`
}

func (s *TeamState) ToProto() *receiverpb.State_Team {
	return &receiverpb.State_Team{
		Id:          int64(s.TeamID),
		Points:      s.Points,
		StolenFlags: int64(s.StolenFlags),
		LostFlags:   int64(s.LostFlags),
	}
}

func (s *TeamState) Clone() *TeamState {
	return &TeamState{
		Points:      s.Points,
		StolenFlags: s.StolenFlags,
		LostFlags:   s.LostFlags,
	}
}

type ServiceState struct {
	ServiceID    int                `json:"service_id"`
	DefaultScore float64            `json:"default_score"`
	TeamStates   map[int]*TeamState `json:"team_states"`
}

func newServiceState(service *models.Service) *ServiceState {
	return &ServiceState{
		ServiceID:    service.ID,
		DefaultScore: service.DefaultScore,
		TeamStates:   make(map[int]*TeamState),
	}
}

func (s *ServiceState) Clone() *ServiceState {
	return &ServiceState{
		ServiceID:    s.ServiceID,
		DefaultScore: s.DefaultScore,
		TeamStates: lo.MapEntries(s.TeamStates, func(key int, value *TeamState) (int, *TeamState) {
			return key, value.Clone()
		}),
	}
}

func (s *ServiceState) ProcessAttack(gs *models.GameState, attack *models.Attack) error {
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
		victimState.LostFlags--
	}
}

func (s *ServiceState) ToProto() *receiverpb.State_Service {
	return &receiverpb.State_Service{
		Id: int64(s.ServiceID),
		Teams: lo.MapToSlice(s.TeamStates, func(_ int, value *TeamState) *receiverpb.State_Team {
			return value.ToProto()
		}),
	}
}

func (s *ServiceState) getOrCreate(teamID int) *TeamState {
	res, ok := s.TeamStates[teamID]
	if !ok {
		res = &TeamState{Points: s.DefaultScore}
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
	return s.getOrCreate(service).ProcessAttack(gs, attack)
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
		Services: lo.MapToSlice(s.ServiceStates, func(_ int, value *ServiceState) *receiverpb.State_Service {
			return value.ToProto()
		}),
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
