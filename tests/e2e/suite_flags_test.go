//go:build e2e

package e2e

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"time"

	gspb "github.com/c4t-but-s4d/fastad/pkg/proto/data/game_state"
	receiverpb "github.com/c4t-but-s4d/fastad/pkg/proto/receiver"
)

type FlagsSuite struct {
	BaseSuite
}

func (s *FlagsSuite) TestSubmitInvalidFlag() {
	s.Require().NotEmpty(s.teamTokens(), "Team tokens required")

	resp := s.SubmitFlags(s.teamTokens()[0].Token, []string{"INVALID_FLAG_12345"})
	s.Require().Len(resp.GetResponses(), 1)

	s.Assert().Equal(receiverpb.FlagResponse_VERDICT_INVALID, resp.GetResponses()[0].GetVerdict())
	s.Assert().Equal("INVALID_FLAG_12345", resp.GetResponses()[0].GetFlag())
}

func (s *FlagsSuite) TestSubmitEmptyFlags() {
	s.Require().NotEmpty(s.teamTokens(), "Team tokens required")

	req, err := http.NewRequest(
		"POST",
		baseURL+"/flags",
		bytes.NewReader([]byte(`{"flags":[]}`)),
	)
	s.Require().NoError(err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Team-Token", s.teamTokens()[0].Token)

	resp, err := http.DefaultClient.Do(req)
	s.Require().NoError(err)
	defer resp.Body.Close()

	s.Assert().Equal(http.StatusBadRequest, resp.StatusCode)
}

func (s *FlagsSuite) TestSubmitTooManyFlags() {
	s.Require().NotEmpty(s.teamTokens(), "Team tokens required")

	// Generate 101 flags (exceeds maxFlagsInRequest = 100)
	flags := make([]string, 101)
	for i := range flags {
		flags[i] = "FLAG_" + string(rune('A'+i%26)) + string(rune('0'+i/26))
	}

	flagsJSON, err := json.Marshal(map[string][]string{"flags": flags})
	s.Require().NoError(err)

	req, err := http.NewRequest(
		"POST",
		baseURL+"/flags",
		bytes.NewReader(flagsJSON),
	)
	s.Require().NoError(err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Team-Token", s.teamTokens()[0].Token)

	resp, err := http.DefaultClient.Do(req)
	s.Require().NoError(err)
	defer resp.Body.Close()

	s.Assert().Equal(http.StatusBadRequest, resp.StatusCode,
		"Submitting >100 flags should return BadRequest")
}

func (s *FlagsSuite) TestSubmitWithoutToken() {
	req, err := http.NewRequest(
		"POST",
		baseURL+"/flags",
		bytes.NewReader([]byte(`{"flags":["test"]}`)),
	)
	s.Require().NoError(err)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	s.Require().NoError(err)
	defer resp.Body.Close()

	s.Assert().Equal(http.StatusBadRequest, resp.StatusCode)
}

func (s *FlagsSuite) TestSubmitWithInvalidToken() {
	req, err := http.NewRequest(
		"POST",
		baseURL+"/flags",
		bytes.NewReader([]byte(`{"flags":["test"]}`)),
	)
	s.Require().NoError(err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Team-Token", "invalid-token-12345")

	resp, err := http.DefaultClient.Do(req)
	s.Require().NoError(err)
	defer resp.Body.Close()

	// Invalid token returns 404 (team not found)
	s.Assert().Equal(http.StatusNotFound, resp.StatusCode)
}

func (s *FlagsSuite) TestSubmitMultipleInvalidFlags() {
	s.Require().NotEmpty(s.teamTokens(), "Team tokens required")

	flags := []string{"FLAG1", "FLAG2", "FLAG3"}
	resp := s.SubmitFlags(s.teamTokens()[0].Token, flags)

	s.Assert().Len(resp.GetResponses(), 3)
	for i, r := range resp.GetResponses() {
		s.Assert().Equal(receiverpb.FlagResponse_VERDICT_INVALID, r.GetVerdict(), "Flag %d should be invalid", i)
	}
}

func (s *FlagsSuite) TestSubmitDuplicateFlagsInSameRequest() {
	s.Require().NotEmpty(s.teamTokens(), "Team tokens required")

	flags := []string{"SAME_FLAG", "SAME_FLAG"}
	resp := s.SubmitFlags(s.teamTokens()[0].Token, flags)

	s.Assert().Len(resp.GetResponses(), 1, "API should deduplicate identical flags")
	s.Assert().Equal(receiverpb.FlagResponse_VERDICT_INVALID, resp.GetResponses()[0].GetVerdict())
}

func (s *FlagsSuite) TestSubmitValidFlag() {
	s.Require().GreaterOrEqual(len(s.teamTokens()), 2, "Need at least 2 team tokens")
	s.Require().NotEmpty(s.services(), "Services required")
	s.Require().NotNil(s.db(), "Database connection required")

	serviceID := int(s.services()[0].GetId())
	victimTeamID := int(s.teamTokens()[1].ID)
	currentRound := s.GetCurrentRound()

	flag := s.InsertTestFlag(victimTeamID, serviceID, currentRound, true)

	resp := s.SubmitFlags(s.teamTokens()[0].Token, []string{flag.Flag})
	s.Require().Len(resp.GetResponses(), 1)

	flagResp := resp.GetResponses()[0]
	s.Assert().Equal(flag.Flag, flagResp.GetFlag())
	s.Assert().Equal(receiverpb.FlagResponse_VERDICT_ACCEPTED, flagResp.GetVerdict(),
		"Expected VERDICT_ACCEPTED, got %v: %s", flagResp.GetVerdict(), flagResp.GetMessage())
	s.Assert().Greater(flagResp.GetAttackerDelta(), float64(0), "Attacker should gain points")
}

func (s *FlagsSuite) TestSubmitOwnFlag() {
	s.Require().NotEmpty(s.teamTokens(), "Team tokens required")
	s.Require().NotEmpty(s.services(), "Services required")
	s.Require().NotNil(s.db(), "Database connection required")

	serviceID := int(s.services()[0].GetId())
	attackerTeamID := int(s.teamTokens()[0].ID)
	currentRound := s.GetCurrentRound()

	flag := s.InsertTestFlag(attackerTeamID, serviceID, currentRound, true)

	resp := s.SubmitFlags(s.teamTokens()[0].Token, []string{flag.Flag})
	s.Require().Len(resp.GetResponses(), 1)

	flagResp := resp.GetResponses()[0]
	s.Assert().Equal(receiverpb.FlagResponse_VERDICT_OWN, flagResp.GetVerdict(),
		"Expected VERDICT_OWN for own flag, got %v: %s", flagResp.GetVerdict(), flagResp.GetMessage())
}

func (s *FlagsSuite) TestSubmitOldFlag() {
	s.Require().GreaterOrEqual(len(s.teamTokens()), 2, "Need at least 2 team tokens")
	s.Require().NotEmpty(s.services(), "Services required")
	s.Require().NotNil(s.db(), "Database connection required")

	gs := s.GetGameState()
	flagLifetime := gs.GetFlagLifetimeRounds()
	currentRound := s.GetCurrentRound()

	if currentRound <= flagLifetime {
		s.T().Skipf("Current round %d <= lifetime %d, cannot test old flags yet", currentRound, flagLifetime)
	}

	serviceID := int(s.services()[0].GetId())
	victimTeamID := int(s.teamTokens()[1].ID)

	oldRound := uint64(0)
	oldFlag := s.InsertTestFlag(victimTeamID, serviceID, oldRound, true)

	resp := s.SubmitFlags(s.teamTokens()[0].Token, []string{oldFlag.Flag})
	s.Require().Len(resp.GetResponses(), 1)

	flagResp := resp.GetResponses()[0]
	s.Assert().Equal(receiverpb.FlagResponse_VERDICT_OLD, flagResp.GetVerdict(),
		"Expected VERDICT_OLD for expired flag from round %d (current: %d, lifetime: %d), got %v: %s",
		oldRound, currentRound, flagLifetime, flagResp.GetVerdict(), flagResp.GetMessage())
}

func (s *FlagsSuite) TestSubmitDuplicateValidFlag() {
	s.Require().GreaterOrEqual(len(s.teamTokens()), 2, "Need at least 2 team tokens")
	s.Require().NotEmpty(s.services(), "Services required")
	s.Require().NotNil(s.db(), "Database connection required")

	serviceID := int(s.services()[0].GetId())
	victimTeamID := int(s.teamTokens()[1].ID)
	currentRound := s.GetCurrentRound()

	flag := s.InsertTestFlag(victimTeamID, serviceID, currentRound, true)

	// First submission
	resp1 := s.SubmitFlags(s.teamTokens()[0].Token, []string{flag.Flag})
	s.Require().Len(resp1.GetResponses(), 1)
	s.Require().Equal(receiverpb.FlagResponse_VERDICT_ACCEPTED, resp1.GetResponses()[0].GetVerdict(),
		"First submission should be accepted")

	// Second submission (duplicate)
	resp2 := s.SubmitFlags(s.teamTokens()[0].Token, []string{flag.Flag})
	s.Require().Len(resp2.GetResponses(), 1)

	s.Assert().Equal(receiverpb.FlagResponse_VERDICT_DUPLICATE, resp2.GetResponses()[0].GetVerdict(),
		"Re-submitting same flag should return DUPLICATE")
}

func (s *FlagsSuite) TestScoreboardUpdatesOnFlagSubmission() {
	s.Require().GreaterOrEqual(len(s.teamTokens()), 2, "Need at least 2 team tokens")
	s.Require().NotEmpty(s.services(), "Services required")
	s.Require().NotNil(s.db(), "Database connection required")

	currentRound := s.GetCurrentRound()
	serviceID := int64(s.services()[0].GetId())
	attackerID := int64(s.teamTokens()[0].ID)
	victimID := int64(s.teamTokens()[1].ID)

	// Get scoreboard state before
	sbBefore := s.GetScoreboard()
	attackerStateBefore := s.GetTeamServiceState(sbBefore, attackerID, serviceID)
	victimStateBefore := s.GetTeamServiceState(sbBefore, victimID, serviceID)

	var attackerStolenBefore, victimLostBefore int64
	var attackerPointsBefore float64
	if attackerStateBefore != nil {
		attackerStolenBefore = attackerStateBefore.GetFlagsStolen()
		attackerPointsBefore = attackerStateBefore.GetPoints()
	}
	if victimStateBefore != nil {
		victimLostBefore = victimStateBefore.GetFlagsLost()
	}

	// Insert and submit flag
	flag := s.InsertTestFlag(int(victimID), int(serviceID), currentRound, true)
	resp := s.SubmitFlags(s.teamTokens()[0].Token, []string{flag.Flag})
	s.Require().Len(resp.GetResponses(), 1)
	s.Require().Equal(receiverpb.FlagResponse_VERDICT_ACCEPTED, resp.GetResponses()[0].GetVerdict())

	// Wait for scoreboard update
	s.Require().Eventually(func() bool {
		sbAfter := s.GetScoreboard()
		attackerStateAfter := s.GetTeamServiceState(sbAfter, attackerID, serviceID)
		if attackerStateAfter == nil {
			return false
		}
		return attackerStateAfter.GetFlagsStolen() > attackerStolenBefore
	}, 10*time.Second, 500*time.Millisecond, "Attacker's FlagsStolen should increase")

	// Verify scoreboard changes
	sbAfter := s.GetScoreboard()
	attackerStateAfter := s.GetTeamServiceState(sbAfter, attackerID, serviceID)
	victimStateAfter := s.GetTeamServiceState(sbAfter, victimID, serviceID)

	s.Require().NotNil(attackerStateAfter, "Attacker should have scoreboard entry")
	s.Require().NotNil(victimStateAfter, "Victim should have scoreboard entry")

	s.Assert().Greater(attackerStateAfter.GetFlagsStolen(), attackerStolenBefore,
		"Attacker FlagsStolen should increase: before=%d, after=%d",
		attackerStolenBefore, attackerStateAfter.GetFlagsStolen())

	s.Assert().Greater(victimStateAfter.GetFlagsLost(), victimLostBefore,
		"Victim FlagsLost should increase: before=%d, after=%d",
		victimLostBefore, victimStateAfter.GetFlagsLost())

	s.Assert().Greater(attackerStateAfter.GetPoints(), attackerPointsBefore,
		"Attacker points should increase: before=%.2f, after=%.2f",
		attackerPointsBefore, attackerStateAfter.GetPoints())
}

func (s *FlagsSuite) TestFlagNotPutFinished() {
	s.Require().GreaterOrEqual(len(s.teamTokens()), 2, "Need at least 2 team tokens")
	s.Require().NotEmpty(s.services(), "Services required")
	s.Require().NotNil(s.db(), "Database connection required")

	serviceID := int(s.services()[0].GetId())
	victimTeamID := int(s.teamTokens()[1].ID)
	currentRound := s.GetCurrentRound()

	// Insert flag with put_finished = false
	flag := s.InsertTestFlag(victimTeamID, serviceID, currentRound, false)

	resp := s.SubmitFlags(s.teamTokens()[0].Token, []string{flag.Flag})
	s.Require().Len(resp.GetResponses(), 1)

	flagResp := resp.GetResponses()[0]
	s.Assert().Equal(receiverpb.FlagResponse_VERDICT_FLAG_NOT_READY, flagResp.GetVerdict(),
		"Flag with put_finished=false should be FLAG_NOT_READY, got %v: %s",
		flagResp.GetVerdict(), flagResp.GetMessage())
}

func (s *FlagsSuite) TestAttackDataEndpoint() {
	resp, err := http.Get(baseURL + "/attack_data")
	s.Require().NoError(err)
	defer resp.Body.Close()

	s.Assert().Equal(http.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	s.Require().NoError(err)

	s.Assert().True(len(body) >= 0, "Attack data endpoint should be accessible")
}

func (s *FlagsSuite) TestFlagSubmissionWhenGamePaused() {
	s.Require().GreaterOrEqual(len(s.teamTokens()), 2, "Need at least 2 team tokens")
	s.Require().NotEmpty(s.services(), "Services required")

	s.SetGameStatus(gspb.GameStatus_GAME_STATUS_PAUSED)
	defer s.SetGameStatus(gspb.GameStatus_GAME_STATUS_RUNNING)

	s.Require().Eventually(func() bool {
		return s.GetGameState().GetStatus() == gspb.GameStatus_GAME_STATUS_PAUSED
	}, 5*time.Second, 100*time.Millisecond, "Game should be paused")

	serviceID := int(s.services()[0].GetId())
	victimTeamID := int(s.teamTokens()[1].ID)
	currentRound := s.GetCurrentRound()

	flag := s.InsertTestFlag(victimTeamID, serviceID, currentRound, true)

	req, err := http.NewRequest(
		"POST",
		baseURL+"/flags",
		bytes.NewReader([]byte(`{"flags":["`+flag.Flag+`"]}`)),
	)
	s.Require().NoError(err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Team-Token", s.teamTokens()[0].Token)

	resp, err := http.DefaultClient.Do(req)
	s.Require().NoError(err)
	defer resp.Body.Close()

	s.Assert().Equal(http.StatusServiceUnavailable, resp.StatusCode,
		"Flag submission should fail when game is paused")
}

func (s *FlagsSuite) TestFlagSubmissionWhenGameFinished() {
	s.Require().GreaterOrEqual(len(s.teamTokens()), 2, "Need at least 2 team tokens")
	s.Require().NotEmpty(s.services(), "Services required")

	s.SetGameStatus(gspb.GameStatus_GAME_STATUS_FINISHED)
	defer s.SetGameStatus(gspb.GameStatus_GAME_STATUS_RUNNING)

	s.Require().Eventually(func() bool {
		return s.GetGameState().GetStatus() == gspb.GameStatus_GAME_STATUS_FINISHED
	}, 5*time.Second, 100*time.Millisecond, "Game should be finished")

	serviceID := int(s.services()[0].GetId())
	victimTeamID := int(s.teamTokens()[1].ID)
	currentRound := s.GetCurrentRound()

	flag := s.InsertTestFlag(victimTeamID, serviceID, currentRound, true)

	req, err := http.NewRequest(
		"POST",
		baseURL+"/flags",
		bytes.NewReader([]byte(`{"flags":["`+flag.Flag+`"]}`)),
	)
	s.Require().NoError(err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Team-Token", s.teamTokens()[0].Token)

	resp, err := http.DefaultClient.Do(req)
	s.Require().NoError(err)
	defer resp.Body.Close()

	s.Assert().Equal(http.StatusPreconditionFailed, resp.StatusCode,
		"Flag submission should fail when game is finished")
}

func (s *FlagsSuite) TestFlagSubmissionForDisabledService() {
	s.Require().GreaterOrEqual(len(s.teamTokens()), 2, "Need at least 2 team tokens")
	s.Require().NotEmpty(s.services(), "Services required")

	serviceID := int(s.services()[0].GetId())
	victimTeamID := int(s.teamTokens()[1].ID)
	currentRound := s.GetCurrentRound()

	s.SetServiceDisabled(serviceID, true)
	defer s.SetServiceDisabled(serviceID, false)

	s.Require().Eventually(func() bool {
		resp, err := http.Get(baseURL + "/api/services")
		if err != nil {
			return false
		}
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)
		return bytes.Contains(body, []byte(`"disabled":true`))
	}, 5*time.Second, 100*time.Millisecond, "Service should be disabled")

	flag := s.InsertTestFlag(victimTeamID, serviceID, currentRound, true)

	resp := s.SubmitFlags(s.teamTokens()[0].Token, []string{flag.Flag})
	s.Require().Len(resp.GetResponses(), 1)

	flagResp := resp.GetResponses()[0]
	s.Assert().Equal(receiverpb.FlagResponse_VERDICT_INVALID, flagResp.GetVerdict(),
		"Flag submission should return INVALID for disabled service, got %v: %s",
		flagResp.GetVerdict(), flagResp.GetMessage())
	s.Assert().Contains(flagResp.GetMessage(), "disabled",
		"Message should mention service is disabled")
}
