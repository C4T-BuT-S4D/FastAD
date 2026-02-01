//go:build e2e

package e2e

import (
	"encoding/json"
	"io"
	"net/http"
	"sort"
	"time"

	checkerpb "github.com/c4t-but-s4d/fastad/pkg/proto/checker"
	scoreboardpb "github.com/c4t-but-s4d/fastad/pkg/proto/scoreboard"
)

type ScoreboardSuite struct {
	BaseSuite
}

func (s *ScoreboardSuite) TestGetScoreboard() {
	var states []*scoreboardpb.Scoreboard_TeamServiceState

	s.Require().Eventually(func() bool {
		sb := s.GetScoreboard()
		states = sb.GetTeamServiceStates()
		return len(states) == 3
	}, 30*time.Second, 2*time.Second, "Expected 3 team-service states (3 teams x 1 service)")

	for _, state := range states {
		s.Assert().Greater(state.GetTeamId(), int64(0))
		s.Assert().Greater(state.GetServiceId(), int64(0))
		s.Assert().GreaterOrEqual(state.GetPoints(), float64(0))
		s.Assert().GreaterOrEqual(state.GetChecksTotal(), int64(0))
		s.Assert().GreaterOrEqual(state.GetChecksPassed(), int64(0))
		s.Assert().GreaterOrEqual(state.GetFlagsStolen(), int64(0))
		s.Assert().GreaterOrEqual(state.GetFlagsLost(), int64(0))
	}
}

func (s *ScoreboardSuite) TestGetCTFTimeScoreboard() {
	resp, err := http.Get(baseURL + "/api/ctftime")
	s.Require().NoError(err)
	defer resp.Body.Close()

	s.Assert().Equal(http.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	s.Require().NoError(err)

	var ctftimeResp struct {
		Standings []struct {
			Pos   int     `json:"pos"`
			Team  string  `json:"team"`
			Score float64 `json:"score"`
		} `json:"standings"`
	}
	err = json.Unmarshal(body, &ctftimeResp)
	s.Require().NoError(err)

	s.Assert().NotNil(ctftimeResp.Standings)
	s.Assert().Len(ctftimeResp.Standings, 3, "Should have 3 teams in standings")

	for i, standing := range ctftimeResp.Standings {
		s.Assert().Equal(i+1, standing.Pos, "Position should be 1-indexed")
		s.Assert().NotEmpty(standing.Team, "Team name should not be empty")
	}
}

func (s *ScoreboardSuite) TestScoreboardTeamsAreSortedByScore() {
	s.Require().GreaterOrEqual(len(s.teamTokens()), 2, "Need at least 2 teams")
	s.Require().NotEmpty(s.services(), "Services required")

	serviceID := int(s.services()[0].GetId())
	victimTeamID := int(s.teamTokens()[1].ID)
	currentRound := s.GetCurrentRound()

	for i := 0; i < 3; i++ {
		flag := s.InsertTestFlag(victimTeamID, serviceID, currentRound, true)
		resp := s.SubmitFlags(s.teamTokens()[0].Token, []string{flag.Flag})
		s.Require().Len(resp.GetResponses(), 1)
	}

	attackerID := int64(s.teamTokens()[0].ID)
	var states []*scoreboardpb.Scoreboard_TeamServiceState
	s.Require().Eventually(func() bool {
		sb := s.GetScoreboard()
		states = sb.GetTeamServiceStates()
		for _, state := range states {
			if state.GetTeamId() == attackerID && state.GetFlagsStolen() >= 3 {
				return true
			}
		}
		return false
	}, 10*time.Second, 500*time.Millisecond, "Attacker should have stolen flags reflected in scoreboard")

	type teamScore struct {
		teamID int64
		points float64
	}
	var scores []teamScore
	for _, state := range states {
		scores = append(scores, teamScore{teamID: state.GetTeamId(), points: state.GetPoints()})
	}

	isSorted := sort.SliceIsSorted(scores, func(i, j int) bool {
		return scores[i].points >= scores[j].points
	})

	s.Assert().True(isSorted || len(scores) <= 1, "Teams should be sorted by score (descending)")
}

func (s *ScoreboardSuite) TestChecksAreBeingRecorded() {
	var totalChecks int64

	s.Require().Eventually(func() bool {
		sb := s.GetScoreboard()
		totalChecks = 0
		for _, state := range sb.GetTeamServiceStates() {
			totalChecks += state.GetChecksTotal()
		}
		return totalChecks > 0
	}, 60*time.Second, 2*time.Second, "Checks should be recorded")
}

func (s *ScoreboardSuite) TestScoreboardUpdatesOnRoundProgression() {
	sbBefore := s.GetScoreboard()
	var totalChecksBefore int64
	for _, state := range sbBefore.GetTeamServiceStates() {
		totalChecksBefore += state.GetChecksTotal()
	}

	initialRound := s.GetCurrentRound()
	s.WaitForRoundIncrement(initialRound, 30*time.Second)

	var totalChecksAfter int64
	s.Require().Eventually(func() bool {
		sbAfter := s.GetScoreboard()
		totalChecksAfter = 0
		for _, state := range sbAfter.GetTeamServiceStates() {
			totalChecksAfter += state.GetChecksTotal()
		}
		return totalChecksAfter > totalChecksBefore
	}, 30*time.Second, time.Second, "Total checks should increase after round progression")
}

func (s *ScoreboardSuite) TestSLACalculationWithUpStatus() {
	s.Require().NotEmpty(s.teams(), "Teams required")
	s.Require().NotEmpty(s.services(), "Services required")

	teamID := int(s.teams()[0].GetId())
	serviceID := int(s.services()[0].GetId())

	sbBefore := s.GetScoreboard()
	stateBefore := s.GetTeamServiceState(sbBefore, int64(teamID), int64(serviceID))

	var checksPassedBefore int64
	if stateBefore != nil {
		checksPassedBefore = stateBefore.GetChecksPassed()
	}

	s.InsertCheckerExecution(teamID, serviceID, checkerpb.Action_ACTION_CHECK, checkerpb.Status_STATUS_UP, "Service is up")

	var checksPassedAfter int64
	s.Require().Eventually(func() bool {
		sbAfter := s.GetScoreboard()
		stateAfter := s.GetTeamServiceState(sbAfter, int64(teamID), int64(serviceID))
		if stateAfter == nil {
			return false
		}
		checksPassedAfter = stateAfter.GetChecksPassed()
		return checksPassedAfter > checksPassedBefore
	}, 30*time.Second, time.Second, "ChecksPassed should increase after STATUS_UP execution")
}

func (s *ScoreboardSuite) TestSLACalculationWithDownStatus() {
	s.Require().NotEmpty(s.teams(), "Teams required")
	s.Require().NotEmpty(s.services(), "Services required")

	teamID := int(s.teams()[0].GetId())
	serviceID := int(s.services()[0].GetId())

	sbBefore := s.GetScoreboard()
	stateBefore := s.GetTeamServiceState(sbBefore, int64(teamID), int64(serviceID))

	var checksTotalBefore, checksPassedBefore int64
	if stateBefore != nil {
		checksTotalBefore = stateBefore.GetChecksTotal()
		checksPassedBefore = stateBefore.GetChecksPassed()
	}

	s.InsertCheckerExecution(teamID, serviceID, checkerpb.Action_ACTION_CHECK, checkerpb.Status_STATUS_DOWN, "Connection error")

	var checksTotalAfter, checksPassedAfter int64
	s.Require().Eventually(func() bool {
		sbAfter := s.GetScoreboard()
		stateAfter := s.GetTeamServiceState(sbAfter, int64(teamID), int64(serviceID))
		if stateAfter == nil {
			return false
		}
		checksTotalAfter = stateAfter.GetChecksTotal()
		checksPassedAfter = stateAfter.GetChecksPassed()
		return checksTotalAfter > checksTotalBefore
	}, 30*time.Second, time.Second, "ChecksTotal should increase after STATUS_DOWN execution")

	s.Assert().Equal(checksPassedBefore, checksPassedAfter,
		"ChecksPassed should NOT increase after STATUS_DOWN: before=%d, after=%d",
		checksPassedBefore, checksPassedAfter)
}

func (s *ScoreboardSuite) TestSLADecreasesWithFailedChecks() {
	s.Require().NotEmpty(s.teams(), "Teams required")
	s.Require().NotEmpty(s.services(), "Services required")

	teamID := int(s.teams()[0].GetId())
	serviceID := int(s.services()[0].GetId())

	s.InsertCheckerExecution(teamID, serviceID, checkerpb.Action_ACTION_CHECK, checkerpb.Status_STATUS_UP, "Up 1")
	s.InsertCheckerExecution(teamID, serviceID, checkerpb.Action_ACTION_CHECK, checkerpb.Status_STATUS_UP, "Up 2")

	var stateBeforeDown *scoreboardpb.Scoreboard_TeamServiceState
	s.Require().Eventually(func() bool {
		sbBeforeDown := s.GetScoreboard()
		stateBeforeDown = s.GetTeamServiceState(sbBeforeDown, int64(teamID), int64(serviceID))
		return stateBeforeDown != nil && stateBeforeDown.GetChecksPassed() >= 2
	}, 10*time.Second, 500*time.Millisecond, "Checks should be reflected in scoreboard")

	slaBeforeDown := float64(0)
	if stateBeforeDown.GetChecksTotal() > 0 {
		slaBeforeDown = float64(stateBeforeDown.GetChecksPassed()) / float64(stateBeforeDown.GetChecksTotal())
	}

	checksTotalBefore := stateBeforeDown.GetChecksTotal()
	s.InsertCheckerExecution(teamID, serviceID, checkerpb.Action_ACTION_CHECK, checkerpb.Status_STATUS_DOWN, "Down 1")
	s.InsertCheckerExecution(teamID, serviceID, checkerpb.Action_ACTION_CHECK, checkerpb.Status_STATUS_DOWN, "Down 2")
	s.InsertCheckerExecution(teamID, serviceID, checkerpb.Action_ACTION_CHECK, checkerpb.Status_STATUS_DOWN, "Down 3")

	var stateAfterDown *scoreboardpb.Scoreboard_TeamServiceState
	s.Require().Eventually(func() bool {
		sbAfterDown := s.GetScoreboard()
		stateAfterDown = s.GetTeamServiceState(sbAfterDown, int64(teamID), int64(serviceID))
		return stateAfterDown != nil && stateAfterDown.GetChecksTotal() >= checksTotalBefore+3
	}, 10*time.Second, 500*time.Millisecond, "Down checks should be reflected in scoreboard")

	slaAfterDown := float64(0)
	if stateAfterDown.GetChecksTotal() > 0 {
		slaAfterDown = float64(stateAfterDown.GetChecksPassed()) / float64(stateAfterDown.GetChecksTotal())
	}

	s.Assert().Less(slaAfterDown, slaBeforeDown,
		"SLA should decrease after failed checks: before=%.4f, after=%.4f",
		slaBeforeDown, slaAfterDown)
}

func (s *ScoreboardSuite) TestCheckStatusHistory() {
	s.Require().NotEmpty(s.teams(), "Teams required")
	s.Require().NotEmpty(s.services(), "Services required")

	var sb *scoreboardpb.Scoreboard
	s.Require().Eventually(func() bool {
		sb = s.GetScoreboard()
		if len(sb.GetTeamServiceStates()) == 0 {
			return false
		}
		for _, state := range sb.GetTeamServiceStates() {
			if len(state.GetCheckStatuses()) > 0 {
				return true
			}
		}
		return false
	}, 60*time.Second, 2*time.Second, "Should have check status history")

	for _, state := range sb.GetTeamServiceStates() {
		statuses := state.GetCheckStatuses()
		s.Assert().LessOrEqual(len(statuses), 3,
			"Should keep at most 3 check statuses, got %d", len(statuses))
	}
}
