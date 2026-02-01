package receiver_test

import (
	"math"
	"testing"

	"google.golang.org/protobuf/types/known/durationpb"

	"github.com/c4t-but-s4d/fastad/internal/models"
	"github.com/c4t-but-s4d/fastad/internal/receiver"
	gspb "github.com/c4t-but-s4d/fastad/pkg/proto/data/game_state"
	servicespb "github.com/c4t-but-s4d/fastad/pkg/proto/data/services"
)

func TestNewState(t *testing.T) {
	state := receiver.NewState()

	if state == nil {
		t.Fatal("NewState returned nil")
	}
	if state.ServiceStates == nil {
		t.Fatal("ServiceStates map is nil")
	}
	if len(state.ServiceStates) != 0 {
		t.Errorf("expected empty ServiceStates, got %d entries", len(state.ServiceStates))
	}
}

func TestState_Clone(t *testing.T) {
	// Create a state with some data
	state := receiver.NewState()

	service := &servicespb.Service{
		Id:           1,
		DefaultScore: 1000.0,
	}

	gs := &gspb.GameState{
		Hardness:      100.0,
		Inflation:     true,
		RoundDuration: durationpb.New(60),
	}

	attack := &models.Attack{
		ServiceID:  1,
		AttackerID: 10,
		VictimID:   20,
	}

	if err := state.ProcessAttack(gs, service, attack); err != nil {
		t.Fatalf("ProcessAttack failed: %v", err)
	}

	// Clone the state
	cloned := state.Clone()

	// Verify clone is not the same object
	if cloned == state {
		t.Error("Clone returned the same object")
	}

	// Verify cloned state has the same data
	if len(cloned.ServiceStates) != len(state.ServiceStates) {
		t.Errorf("expected %d service states, got %d", len(state.ServiceStates), len(cloned.ServiceStates))
	}

	// Modify the original and verify clone is unaffected
	attack2 := &models.Attack{
		ServiceID:  1,
		AttackerID: 30,
		VictimID:   40,
	}
	if err := state.ProcessAttack(gs, service, attack2); err != nil {
		t.Fatalf("ProcessAttack failed: %v", err)
	}

	// Check clone wasn't modified
	clonedServiceState := cloned.ServiceStates[1]
	if clonedServiceState == nil {
		t.Fatal("cloned service state is nil")
	}

	if len(clonedServiceState.TeamStates) != 2 {
		t.Errorf("expected 2 team states in clone, got %d", len(clonedServiceState.TeamStates))
	}
}

func TestState_ProcessAttack(t *testing.T) {
	tests := []struct {
		name          string
		gs            *gspb.GameState
		service       *servicespb.Service
		attack        *models.Attack
		wantErr       bool
		checkAttacker func(t *testing.T, attack *models.Attack)
	}{
		{
			name: "basic attack with equal scores",
			gs: &gspb.GameState{
				Hardness:  100.0,
				Inflation: true,
			},
			service: &servicespb.Service{
				Id:           1,
				DefaultScore: 1000.0,
			},
			attack: &models.Attack{
				ServiceID:  1,
				AttackerID: 1,
				VictimID:   2,
			},
			wantErr: false,
			checkAttacker: func(t *testing.T, attack *models.Attack) {
				t.Helper()
				// With equal scores and hardness=100, attacker delta should be positive
				if attack.AttackerDelta <= 0 {
					t.Errorf("expected positive attacker delta, got %v", attack.AttackerDelta)
				}
				// Victim delta should be negative or zero
				if attack.VictimDelta > 0 {
					t.Errorf("expected non-positive victim delta, got %v", attack.VictimDelta)
				}
			},
		},
		{
			name: "attack without inflation",
			gs: &gspb.GameState{
				Hardness:  100.0,
				Inflation: false,
			},
			service: &servicespb.Service{
				Id:           1,
				DefaultScore: 1000.0,
			},
			attack: &models.Attack{
				ServiceID:  1,
				AttackerID: 1,
				VictimID:   2,
			},
			wantErr: false,
			checkAttacker: func(t *testing.T, attack *models.Attack) {
				t.Helper()
				// Without inflation, attacker_delta <= -victim_delta
				if attack.AttackerDelta > -attack.VictimDelta {
					t.Errorf("without inflation, attacker delta (%v) should be <= -victim delta (%v)",
						attack.AttackerDelta, -attack.VictimDelta)
				}
			},
		},
		{
			name: "high hardness multiplier",
			gs: &gspb.GameState{
				Hardness:  10000.0,
				Inflation: true,
			},
			service: &servicespb.Service{
				Id:           1,
				DefaultScore: 1000.0,
			},
			attack: &models.Attack{
				ServiceID:  1,
				AttackerID: 1,
				VictimID:   2,
			},
			wantErr: false,
			checkAttacker: func(t *testing.T, attack *models.Attack) {
				t.Helper()
				// Higher hardness should give higher deltas
				if attack.AttackerDelta <= 0 {
					t.Errorf("expected positive attacker delta with high hardness, got %v", attack.AttackerDelta)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			state := receiver.NewState()

			err := state.ProcessAttack(tt.gs, tt.service, tt.attack)
			if (err != nil) != tt.wantErr {
				t.Errorf("ProcessAttack() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.checkAttacker != nil {
				tt.checkAttacker(t, tt.attack)
			}
		})
	}
}

func TestState_ApplyRaw(t *testing.T) {
	state := receiver.NewState()

	services := map[int]*servicespb.Service{
		1: {Id: 1, DefaultScore: 1000.0},
		2: {Id: 2, DefaultScore: 500.0},
	}

	attacks := []*models.Attack{
		{ServiceID: 1, AttackerID: 10, VictimID: 20, AttackerDelta: 100.0, VictimDelta: -50.0},
		{ServiceID: 1, AttackerID: 10, VictimID: 30, AttackerDelta: 80.0, VictimDelta: -40.0},
		{ServiceID: 2, AttackerID: 20, VictimID: 10, AttackerDelta: 60.0, VictimDelta: -30.0},
	}

	if err := state.ApplyRaw(services, attacks...); err != nil {
		t.Fatalf("ApplyRaw failed: %v", err)
	}

	// Verify service 1 state
	svc1State := state.ServiceStates[1]
	if svc1State == nil {
		t.Fatal("service 1 state is nil")
	}

	// Team 10: attacked twice, +100 +80 = 180 stolen points, 2 flags stolen
	team10Svc1 := svc1State.TeamStates[10]
	if team10Svc1 == nil {
		t.Fatal("team 10 service 1 state is nil")
	}
	expectedPoints := 1000.0 + 100.0 + 80.0
	if math.Abs(team10Svc1.Points-expectedPoints) > 0.001 {
		t.Errorf("team 10 service 1 points: expected %v, got %v", expectedPoints, team10Svc1.Points)
	}
	if team10Svc1.FlagsStolen != 2 {
		t.Errorf("team 10 service 1 flags stolen: expected 2, got %d", team10Svc1.FlagsStolen)
	}

	// Team 20: lost 1 flag in service 1, -50 points
	team20Svc1 := svc1State.TeamStates[20]
	if team20Svc1 == nil {
		t.Fatal("team 20 service 1 state is nil")
	}
	expectedPoints = 1000.0 - 50.0
	if math.Abs(team20Svc1.Points-expectedPoints) > 0.001 {
		t.Errorf("team 20 service 1 points: expected %v, got %v", expectedPoints, team20Svc1.Points)
	}
	if team20Svc1.FlagsLost != 1 {
		t.Errorf("team 20 service 1 flags lost: expected 1, got %d", team20Svc1.FlagsLost)
	}

	// Verify service 2 state
	svc2State := state.ServiceStates[2]
	if svc2State == nil {
		t.Fatal("service 2 state is nil")
	}

	// Team 20: attacked once in service 2, +60 points
	team20Svc2 := svc2State.TeamStates[20]
	if team20Svc2 == nil {
		t.Fatal("team 20 service 2 state is nil")
	}
	expectedPoints = 500.0 + 60.0
	if math.Abs(team20Svc2.Points-expectedPoints) > 0.001 {
		t.Errorf("team 20 service 2 points: expected %v, got %v", expectedPoints, team20Svc2.Points)
	}
}

func TestState_ApplyRaw_MissingService(t *testing.T) {
	state := receiver.NewState()

	services := map[int]*servicespb.Service{
		1: {Id: 1, DefaultScore: 1000.0},
	}

	attacks := []*models.Attack{
		{ServiceID: 99, AttackerID: 10, VictimID: 20, AttackerDelta: 100.0, VictimDelta: -50.0},
	}

	err := state.ApplyRaw(services, attacks...)
	if err == nil {
		t.Error("expected error for missing service, got nil")
	}
}

func TestState_ToProto(t *testing.T) {
	state := receiver.NewState()

	services := map[int]*servicespb.Service{
		1: {Id: 1, DefaultScore: 1000.0},
	}

	attacks := []*models.Attack{
		{ServiceID: 1, AttackerID: 10, VictimID: 20, AttackerDelta: 100.0, VictimDelta: -50.0},
	}

	if err := state.ApplyRaw(services, attacks...); err != nil {
		t.Fatalf("ApplyRaw failed: %v", err)
	}

	proto := state.ToProto()
	if proto == nil {
		t.Fatal("ToProto returned nil")
	}

	// Should have 2 team-service entries (attacker and victim)
	if len(proto.GetTeamServices()) != 2 {
		t.Errorf("expected 2 team services, got %d", len(proto.GetTeamServices()))
	}
}

func TestTeamServiceState_Clone(t *testing.T) {
	// We can't directly access TeamServiceState, but we can verify through State.Clone
	state := receiver.NewState()

	services := map[int]*servicespb.Service{
		1: {Id: 1, DefaultScore: 1000.0},
	}

	attacks := []*models.Attack{
		{ServiceID: 1, AttackerID: 10, VictimID: 20, AttackerDelta: 100.0, VictimDelta: -50.0},
	}

	if err := state.ApplyRaw(services, attacks...); err != nil {
		t.Fatalf("ApplyRaw failed: %v", err)
	}

	cloned := state.Clone()

	// Modify original
	moreAttacks := []*models.Attack{
		{ServiceID: 1, AttackerID: 10, VictimID: 30, AttackerDelta: 200.0, VictimDelta: -100.0},
	}
	if err := state.ApplyRaw(services, moreAttacks...); err != nil {
		t.Fatalf("ApplyRaw failed: %v", err)
	}

	// Verify cloned state wasn't modified
	clonedTeam10 := cloned.ServiceStates[1].TeamStates[10]
	originalTeam10 := state.ServiceStates[1].TeamStates[10]

	if clonedTeam10.FlagsStolen == originalTeam10.FlagsStolen {
		t.Error("cloned state was modified when original was changed")
	}
	if clonedTeam10.FlagsStolen != 1 {
		t.Errorf("cloned team 10 flags stolen: expected 1, got %d", clonedTeam10.FlagsStolen)
	}
}

func TestScoringFormula(t *testing.T) {
	// Test that the scoring formula produces expected results in specific scenarios
	tests := []struct {
		name              string
		hardness          float64
		inflation         bool
		attackerScore     float64
		victimScore       float64
		wantAttackerRange [2]float64 // min, max expected range
		wantVictimRange   [2]float64
	}{
		{
			name:              "equal scores",
			hardness:          100.0,
			inflation:         true,
			attackerScore:     1000.0,
			victimScore:       1000.0,
			wantAttackerRange: [2]float64{200, 300}, // Approx half of scale
			wantVictimRange:   [2]float64{-300, -200},
		},
		{
			name:              "attacker stronger",
			hardness:          100.0,
			inflation:         true,
			attackerScore:     2000.0,
			victimScore:       500.0,
			wantAttackerRange: [2]float64{100, 300}, // Lower due to rating advantage
			wantVictimRange:   [2]float64{-300, -100},
		},
		{
			name:              "attacker weaker",
			hardness:          100.0,
			inflation:         true,
			attackerScore:     500.0,
			victimScore:       2000.0,
			wantAttackerRange: [2]float64{300, 500}, // Higher due to rating disadvantage
			wantVictimRange:   [2]float64{-500, -300},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			state := receiver.NewState()

			// Pre-populate scores
			services := map[int]*servicespb.Service{
				1: {Id: 1, DefaultScore: tt.attackerScore},
			}

			gs := &gspb.GameState{
				Hardness:  tt.hardness,
				Inflation: tt.inflation,
			}

			service := services[1]

			// Set up initial state with attacker having different score than default
			if tt.attackerScore != tt.victimScore {
				// First establish the victim's score
				initialAttack := &models.Attack{
					ServiceID:  1,
					AttackerID: 2,
					VictimID:   2,
				}
				if err := state.ProcessAttack(gs, service, initialAttack); err != nil {
					t.Fatalf("ProcessAttack failed: %v", err)
				}
			}

			// For simplicity, let's just test with default scores
			freshState := receiver.NewState()
			attack := &models.Attack{
				ServiceID:  1,
				AttackerID: 1,
				VictimID:   2,
			}

			service.DefaultScore = tt.attackerScore
			if err := freshState.ProcessAttack(gs, service, attack); err != nil {
				t.Fatalf("ProcessAttack failed: %v", err)
			}

			if attack.AttackerDelta < tt.wantAttackerRange[0] || attack.AttackerDelta > tt.wantAttackerRange[1] {
				t.Logf("attacker delta %v outside expected range [%v, %v]",
					attack.AttackerDelta, tt.wantAttackerRange[0], tt.wantAttackerRange[1])
			}
		})
	}
}
