//go:build e2e

package e2e

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"google.golang.org/protobuf/encoding/protojson"

	checkerpb "github.com/c4t-but-s4d/fastad/pkg/proto/checker"
	teamspb "github.com/c4t-but-s4d/fastad/pkg/proto/data/teams"
)

type TeamsSuite struct {
	BaseSuite
}

func (s *TeamsSuite) TestListTeams() {
	resp, err := http.Get(baseURL + "/api/teams")
	s.Require().NoError(err)
	defer resp.Body.Close()

	s.Assert().Equal(http.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	s.Require().NoError(err)

	var teamsResp teamspb.Team_Batch
	err = protojson.Unmarshal(body, &teamsResp)
	s.Require().NoError(err)

	teams := teamsResp.GetTeams()
	s.Assert().Len(teams, 3, "Expected 3 teams")

	teamNames := make([]string, len(teams))
	for i, team := range teams {
		teamNames[i] = team.GetName()
		s.Assert().Empty(team.GetToken(), "Team token should be stripped from API response")
		s.Assert().NotEmpty(team.GetAddress(), "Team should have an address")
		s.Assert().Greater(team.GetId(), int64(0), "Team should have positive ID")
	}

	s.Assert().Contains(teamNames, "Team Alpha")
	s.Assert().Contains(teamNames, "Team Beta")
	s.Assert().Contains(teamNames, "Team Gamma")
}

func (s *TeamsSuite) TestUpdateTeamAvatar() {
	s.Require().NotEmpty(s.teamTokens, "Team tokens required")

	newAvatarURL := "https://example.com/avatar.png"
	body := fmt.Sprintf(`{"avatar_url": "%s"}`, newAvatarURL)

	req, err := http.NewRequest(http.MethodPut, baseURL+"/api/teams", bytes.NewReader([]byte(body)))
	s.Require().NoError(err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Team-Token", s.teamTokens[0].Token)

	resp, err := http.DefaultClient.Do(req)
	s.Require().NoError(err)
	defer resp.Body.Close()

	s.Assert().Equal(http.StatusOK, resp.StatusCode)

	respBody, err := io.ReadAll(resp.Body)
	s.Require().NoError(err)

	var team teamspb.Team
	err = protojson.Unmarshal(respBody, &team)
	s.Require().NoError(err)

	s.Assert().Equal(newAvatarURL, team.GetAvatarUrl())
}

func (s *TeamsSuite) TestUpdateTeamAvatarWithoutToken() {
	body := `{"avatar_url": "https://example.com/avatar.png"}`

	req, err := http.NewRequest(http.MethodPut, baseURL+"/api/teams", bytes.NewReader([]byte(body)))
	s.Require().NoError(err)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	s.Require().NoError(err)
	defer resp.Body.Close()

	s.Assert().Equal(http.StatusBadRequest, resp.StatusCode)
}

func (s *TeamsSuite) TestUpdateTeamAvatarWithInvalidToken() {
	body := `{"avatar_url": "https://example.com/avatar.png"}`

	req, err := http.NewRequest(http.MethodPut, baseURL+"/api/teams", bytes.NewReader([]byte(body)))
	s.Require().NoError(err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Team-Token", "invalid-token-12345")

	resp, err := http.DefaultClient.Do(req)
	s.Require().NoError(err)
	defer resp.Body.Close()

	s.Assert().Equal(http.StatusNotFound, resp.StatusCode)
}

func (s *TeamsSuite) TestGetTeamHistoryShowsCheckerResults() {
	s.Require().NotEmpty(s.teams, "Teams required")
	s.Require().NotEmpty(s.services, "Services required")

	teamID := int(s.teams[0].GetId())
	serviceID := int(s.services[0].GetId())

	s.InsertCheckerExecution(teamID, serviceID, checkerpb.Action_ACTION_CHECK, checkerpb.Status_STATUS_UP, "Test check")

	url := fmt.Sprintf("%s/api/teams/%d/history", baseURL, teamID)
	resp, err := http.Get(url)
	s.Require().NoError(err)
	defer resp.Body.Close()

	s.Assert().Equal(http.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	s.Require().NoError(err)

	var historyResp checkerpb.Execution_Batch
	err = protojson.Unmarshal(body, &historyResp)
	s.Require().NoError(err)

	s.Assert().NotEmpty(historyResp.GetExecutions(), "History should contain executions")

	for _, exec := range historyResp.GetExecutions() {
		s.Assert().Empty(exec.GetPrivate(), "Private field should be stripped")
	}
}

func (s *TeamsSuite) TestGetTeamHistoryWithServiceFilter() {
	s.Require().NotEmpty(s.teams, "Teams required")
	s.Require().NotEmpty(s.services, "Services required")

	teamID := int(s.teams[0].GetId())
	serviceID := int(s.services[0].GetId())

	s.InsertCheckerExecution(teamID, serviceID, checkerpb.Action_ACTION_CHECK, checkerpb.Status_STATUS_UP, "Filtered check")

	url := fmt.Sprintf("%s/api/teams/%d/history?service_id=%d", baseURL, teamID, serviceID)
	resp, err := http.Get(url)
	s.Require().NoError(err)
	defer resp.Body.Close()

	s.Assert().Equal(http.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	s.Require().NoError(err)

	var historyResp checkerpb.Execution_Batch
	err = protojson.Unmarshal(body, &historyResp)
	s.Require().NoError(err)

	for _, exec := range historyResp.GetExecutions() {
		s.Assert().Equal(int64(serviceID), exec.GetServiceId(),
			"All executions should be for the filtered service")
	}
}

func (s *TeamsSuite) TestGetTeamHistoryWithLimit() {
	s.Require().NotEmpty(s.teams, "Teams required")
	s.Require().NotEmpty(s.services, "Services required")

	teamID := int(s.teams[0].GetId())
	serviceID := int(s.services[0].GetId())

	for i := 0; i < 10; i++ {
		s.InsertCheckerExecution(teamID, serviceID, checkerpb.Action_ACTION_CHECK, checkerpb.Status_STATUS_UP, fmt.Sprintf("Check %d", i))
	}

	limit := 5
	url := fmt.Sprintf("%s/api/teams/%d/history?limit=%d", baseURL, teamID, limit)
	resp, err := http.Get(url)
	s.Require().NoError(err)
	defer resp.Body.Close()

	s.Assert().Equal(http.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	s.Require().NoError(err)

	var historyResp checkerpb.Execution_Batch
	err = protojson.Unmarshal(body, &historyResp)
	s.Require().NoError(err)

	s.Assert().LessOrEqual(len(historyResp.GetExecutions()), limit,
		"Should return at most %d executions", limit)
}

func (s *TeamsSuite) TestGetTeamHistoryForNonExistentTeam() {
	url := fmt.Sprintf("%s/api/teams/99999/history", baseURL)
	resp, err := http.Get(url)
	s.Require().NoError(err)
	defer resp.Body.Close()

	s.Assert().Equal(http.StatusOK, resp.StatusCode)
}
