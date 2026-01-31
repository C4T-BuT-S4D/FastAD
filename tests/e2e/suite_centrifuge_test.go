//go:build e2e

package e2e

import (
	"sync"
	"time"

	"github.com/centrifugal/centrifuge-go"
	"google.golang.org/protobuf/encoding/protojson"

	receiverpb "github.com/c4t-but-s4d/fastad/pkg/proto/receiver"
	scoreboardpb "github.com/c4t-but-s4d/fastad/pkg/proto/scoreboard"
)

const centrifugeURL = "ws://localhost:8080/centrifuge/websocket"

type CentrifugeSuite struct {
	BaseSuite
}

func (s *CentrifugeSuite) TestScoreboardSubscriptionAndUpdates() {
	client := centrifuge.NewJsonClient(centrifugeURL, centrifuge.Config{})

	var receivedUpdate bool
	var scoreboard *scoreboardpb.Scoreboard
	updateChan := make(chan struct{}, 1)

	sub, err := client.NewSubscription("scoreboard")
	s.Require().NoError(err, "Failed to create subscription")

	sub.OnPublication(func(e centrifuge.PublicationEvent) {
		var sb scoreboardpb.Scoreboard
		if err := protojson.Unmarshal(e.Data, &sb); err != nil {
			s.T().Logf("Failed to unmarshal scoreboard: %v", err)
			return
		}
		scoreboard = &sb
		receivedUpdate = true
		select {
		case updateChan <- struct{}{}:
		default:
		}
	})

	err = client.Connect()
	s.Require().NoError(err, "Failed to connect")

	err = sub.Subscribe()
	s.Require().NoError(err, "Failed to subscribe")

	select {
	case <-updateChan:
	case <-time.After(10 * time.Second):
		s.T().Fatal("Timeout waiting for scoreboard update")
	}

	s.Assert().True(receivedUpdate, "Should receive scoreboard update")
	s.Assert().NotNil(scoreboard, "Scoreboard should not be nil")
	s.Assert().NotEmpty(scoreboard.GetTeamServiceStates(), "Scoreboard should have team service states")

	for _, state := range scoreboard.GetTeamServiceStates() {
		s.Assert().Greater(state.GetTeamId(), int64(0), "TeamId should be positive")
		s.Assert().Greater(state.GetServiceId(), int64(0), "ServiceId should be positive")
	}

	client.Close()
}

func (s *CentrifugeSuite) TestAttackNotificationOnFlagSubmission() {
	s.Require().GreaterOrEqual(len(s.teamTokens), 2, "Need at least 2 team tokens")
	s.Require().NotEmpty(s.services, "Services required")
	s.Require().NotNil(s.db, "Database connection required")

	client := centrifuge.NewJsonClient(centrifugeURL, centrifuge.Config{})

	var receivedNotification bool
	var attackBatch *receiverpb.AttackNotification_Batch
	notificationChan := make(chan struct{}, 1)

	sub, err := client.NewSubscription("attacks")
	s.Require().NoError(err, "Failed to create subscription")

	sub.OnPublication(func(e centrifuge.PublicationEvent) {
		var batch receiverpb.AttackNotification_Batch
		if err := protojson.Unmarshal(e.Data, &batch); err != nil {
			s.T().Logf("Failed to unmarshal attack notification: %v", err)
			return
		}
		attackBatch = &batch
		receivedNotification = true
		select {
		case notificationChan <- struct{}{}:
		default:
		}
	})

	err = client.Connect()
	s.Require().NoError(err, "Failed to connect")

	err = sub.Subscribe()
	s.Require().NoError(err, "Failed to subscribe")

	time.Sleep(500 * time.Millisecond)

	serviceID := int(s.services[0].GetId())
	victimTeamID := int(s.teamTokens[1].ID)
	currentRound := s.GetCurrentRound()

	flag := s.InsertTestFlag(victimTeamID, serviceID, currentRound, true)

	resp := s.SubmitFlags(s.teamTokens[0].Token, []string{flag.Flag})
	s.Require().Len(resp.GetResponses(), 1)
	s.Require().Equal(receiverpb.FlagResponse_VERDICT_ACCEPTED, resp.GetResponses()[0].GetVerdict(),
		"Flag should be accepted")

	select {
	case <-notificationChan:
	case <-time.After(10 * time.Second):
		s.T().Fatal("Timeout waiting for attack notification")
	}

	s.Assert().True(receivedNotification, "Should receive attack notification")
	s.Assert().NotNil(attackBatch, "Attack batch should not be nil")
	s.Assert().NotEmpty(attackBatch.GetAttacks(), "Should have at least one attack notification")

	attack := attackBatch.GetAttacks()[0]
	s.Assert().Equal(int64(s.teamTokens[0].ID), attack.GetAttackerId(), "Attacker team ID should match")
	s.Assert().Equal(int64(victimTeamID), attack.GetVictimId(), "Victim team ID should match")
	s.Assert().Equal(int64(serviceID), attack.GetServiceId(), "Service ID should match")

	client.Close()
}

func (s *CentrifugeSuite) TestMultipleScoreboardUpdates() {
	client := centrifuge.NewJsonClient(centrifugeURL, centrifuge.Config{})

	var updates []*scoreboardpb.Scoreboard
	var mu sync.Mutex
	updateChan := make(chan struct{}, 10)

	sub, err := client.NewSubscription("scoreboard")
	s.Require().NoError(err, "Failed to create subscription")

	sub.OnPublication(func(e centrifuge.PublicationEvent) {
		var sb scoreboardpb.Scoreboard
		if err := protojson.Unmarshal(e.Data, &sb); err != nil {
			s.T().Logf("Failed to unmarshal scoreboard: %v", err)
			return
		}
		mu.Lock()
		updates = append(updates, &sb)
		mu.Unlock()
		select {
		case updateChan <- struct{}{}:
		default:
		}
	})

	err = client.Connect()
	s.Require().NoError(err, "Failed to connect")

	err = sub.Subscribe()
	s.Require().NoError(err, "Failed to subscribe")

	timeout := time.After(15 * time.Second)
	updateCount := 0

	for updateCount < 3 {
		select {
		case <-updateChan:
			updateCount++
			s.T().Logf("Received scoreboard update %d", updateCount)
		case <-timeout:
			s.T().Fatalf("Timeout waiting for scoreboard updates, received %d", updateCount)
		}
	}

	mu.Lock()
	defer mu.Unlock()

	s.Assert().GreaterOrEqual(len(updates), 3, "Should receive at least 3 scoreboard updates")

	for i, sb := range updates {
		s.Assert().NotEmpty(sb.GetTeamServiceStates(), "Update %d should have team service states", i)
	}

	client.Close()
}
