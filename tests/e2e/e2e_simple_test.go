//go:build e2e && e2e_simple

package e2e

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

// TestE2ESimplePreset runs all E2E test suites using the "simple" preset.
// Each suite embeds BaseSuite which handles game initialization, teardown,
// and provides common helper methods.
func TestE2ESimplePreset(t *testing.T) {
	// Run suites sequentially - they share game state
	// Order: teams -> services -> gamestate -> scoreboard -> flags -> centrifuge

	t.Run("TeamsSuite", func(t *testing.T) {
		s := &TeamsSuite{}
		s.preset = "simple"
		suite.Run(t, s)
	})

	t.Run("ServicesSuite", func(t *testing.T) {
		s := &ServicesSuite{}
		s.preset = "simple"
		suite.Run(t, s)
	})

	t.Run("GameStateSuite", func(t *testing.T) {
		s := &GameStateSuite{}
		s.preset = "simple"
		suite.Run(t, s)
	})

	t.Run("ScoreboardSuite", func(t *testing.T) {
		s := &ScoreboardSuite{}
		s.preset = "simple"
		suite.Run(t, s)
	})

	t.Run("FlagsSuite", func(t *testing.T) {
		s := &FlagsSuite{}
		s.preset = "simple"
		suite.Run(t, s)
	})

	t.Run("CentrifugeSuite", func(t *testing.T) {
		s := &CentrifugeSuite{}
		s.preset = "simple"
		suite.Run(t, s)
	})
}
