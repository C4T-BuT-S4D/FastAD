//go:build e2e && e2e_production

package e2e

import "testing"

func TestE2EProductionPreset(t *testing.T) {
	tc := setupTest(t, "production")
	defer tc.teardown()

	runAllTests(t, tc)
}
