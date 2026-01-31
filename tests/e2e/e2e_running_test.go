//go:build e2e_running

package e2e

import (
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const defaultBaseURL = "http://localhost:8080"

func TestE2EWithRunningGame(t *testing.T) {
	baseURL := os.Getenv("E2E_BASE_URL")
	if baseURL == "" {
		baseURL = defaultBaseURL
	}

	t.Run("HealthCheck", func(t *testing.T) {
		resp, err := http.Get(baseURL + "/healthcheck")
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("TeamsAPI", func(t *testing.T) {
		resp, err := http.Get(baseURL + "/api/teams")
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		body, _ := io.ReadAll(resp.Body)
		t.Logf("Teams: %s", string(body))
	})

	t.Run("ServicesAPI", func(t *testing.T) {
		resp, err := http.Get(baseURL + "/api/services")
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		body, _ := io.ReadAll(resp.Body)
		t.Logf("Services: %s", string(body))
	})

	t.Run("GameStateAPI", func(t *testing.T) {
		resp, err := http.Get(baseURL + "/api/game")
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		body, _ := io.ReadAll(resp.Body)
		t.Logf("Game State: %s", string(body))
	})

	t.Run("ScoreboardAPI", func(t *testing.T) {
		resp, err := http.Get(baseURL + "/api/scoreboard")
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		body, _ := io.ReadAll(resp.Body)
		t.Logf("Scoreboard: %s", string(body))
	})
}
