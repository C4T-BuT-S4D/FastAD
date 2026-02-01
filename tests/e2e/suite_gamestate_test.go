//go:build e2e

package e2e

import (
	"io"
	"net/http"
	"time"

	"google.golang.org/protobuf/encoding/protojson"

	gspb "github.com/c4t-but-s4d/fastad/pkg/proto/data/game_state"
)

type GameStateSuite struct {
	BaseSuite
}

func (s *GameStateSuite) TestGetGameStateJSON() {
	resp, err := http.Get(baseURL + "/api/game")
	s.Require().NoError(err)
	defer resp.Body.Close()

	s.Assert().Equal(http.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	s.Require().NoError(err)

	var gs gspb.GameState
	err = protojson.Unmarshal(body, &gs)
	s.Require().NoError(err)

	s.Assert().Equal(10.0, gs.GetHardness(), "Hardness should match config")
	s.Assert().False(gs.GetInflation(), "Inflation should match config")
	s.Assert().Equal(gspb.GameStatus_GAME_STATUS_RUNNING, gs.GetStatus(), "Game should be running")
	s.Assert().Equal(uint64(5), gs.GetFlagLifetimeRounds(), "Flag lifetime should match config")
	s.Assert().NotNil(gs.GetStartTime(), "Start time should be set")
	s.Assert().NotNil(gs.GetRoundDuration(), "Round duration should be set")

	expectedDuration := 10 * time.Second
	s.Assert().Equal(expectedDuration.Seconds(), float64(gs.GetRoundDuration().GetSeconds()),
		"Round duration should be 10 seconds as configured")
}

func (s *GameStateSuite) TestGetGameStateProtoFormat() {
	resp, err := http.Get(baseURL + "/api/game?format=proto")
	s.Require().NoError(err)
	defer resp.Body.Close()

	s.Assert().Equal(http.StatusOK, resp.StatusCode)
	s.Assert().Equal("application/octet-stream", resp.Header.Get("Content-Type"))

	body, err := io.ReadAll(resp.Body)
	s.Require().NoError(err)
	s.Assert().NotEmpty(body, "Proto response should not be empty")
}

func (s *GameStateSuite) TestRoundsIncrement() {
	var initialRound uint64

	s.WaitForCondition(func() bool {
		initialRound = s.GetCurrentRound()
		return initialRound > 0
	}, 30*time.Second, time.Second, "Game should start with round > 0")

	s.T().Logf("Initial round: %d", initialRound)

	currentRound := s.WaitForRoundIncrement(initialRound, 30*time.Second)
	s.T().Logf("Round incremented from %d to %d", initialRound, currentRound)

	s.Assert().Greater(currentRound, initialRound, "Round should have incremented")
}
