//go:build e2e && e2e_simple

package e2e

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

func TestE2ESimplePreset(t *testing.T) {
	sharedState := NewSharedTestState(t, "simple")
	defer sharedState.TearDown()

	t.Run("TeamsSuite", func(t *testing.T) {
		s := &TeamsSuite{}
		s.sharedState = sharedState
		suite.Run(t, s)
	})

	t.Run("ServicesSuite", func(t *testing.T) {
		s := &ServicesSuite{}
		s.sharedState = sharedState
		suite.Run(t, s)
	})

	t.Run("GameStateSuite", func(t *testing.T) {
		s := &GameStateSuite{}
		s.sharedState = sharedState
		suite.Run(t, s)
	})

	t.Run("ScoreboardSuite", func(t *testing.T) {
		s := &ScoreboardSuite{}
		s.sharedState = sharedState
		suite.Run(t, s)
	})

	t.Run("FlagsSuite", func(t *testing.T) {
		s := &FlagsSuite{}
		s.sharedState = sharedState
		suite.Run(t, s)
	})

	t.Run("CentrifugeSuite", func(t *testing.T) {
		s := &CentrifugeSuite{}
		s.sharedState = sharedState
		suite.Run(t, s)
	})
}
