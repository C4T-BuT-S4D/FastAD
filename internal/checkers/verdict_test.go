package checkers_test

import (
	"testing"

	"github.com/c4t-but-s4d/fastad/internal/checkers"
	checkerpb "github.com/c4t-but-s4d/fastad/pkg/proto/checker"
)

func TestVerdict_String(t *testing.T) {
	tests := []struct {
		name     string
		verdict  checkers.Verdict
		expected string
	}{
		{
			name: "check up",
			verdict: checkers.Verdict{
				Action: checkerpb.Action_ACTION_CHECK,
				Status: checkerpb.Status_STATUS_UP,
			},
			expected: "ACTION_CHECK STATUS_UP",
		},
		{
			name: "put mumble",
			verdict: checkers.Verdict{
				Action: checkerpb.Action_ACTION_PUT,
				Status: checkerpb.Status_STATUS_MUMBLE,
			},
			expected: "ACTION_PUT STATUS_MUMBLE",
		},
		{
			name: "get corrupt",
			verdict: checkers.Verdict{
				Action: checkerpb.Action_ACTION_GET,
				Status: checkerpb.Status_STATUS_CORRUPT,
			},
			expected: "ACTION_GET STATUS_CORRUPT",
		},
		{
			name: "check down",
			verdict: checkers.Verdict{
				Action: checkerpb.Action_ACTION_CHECK,
				Status: checkerpb.Status_STATUS_DOWN,
			},
			expected: "ACTION_CHECK STATUS_DOWN",
		},
		{
			name: "check failed",
			verdict: checkers.Verdict{
				Action: checkerpb.Action_ACTION_CHECK,
				Status: checkerpb.Status_STATUS_CHECK_FAILED,
			},
			expected: "ACTION_CHECK STATUS_CHECK_FAILED",
		},
		{
			name: "unspecified",
			verdict: checkers.Verdict{
				Action: checkerpb.Action_ACTION_UNSPECIFIED,
				Status: checkerpb.Status_STATUS_UNSPECIFIED,
			},
			expected: "ACTION_UNSPECIFIED STATUS_UNSPECIFIED",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.verdict.String()
			if result != tt.expected {
				t.Errorf("Verdict.String() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestVerdict_IsUp(t *testing.T) {
	tests := []struct {
		name    string
		verdict checkers.Verdict
		want    bool
	}{
		{
			name: "status up returns true",
			verdict: checkers.Verdict{
				Action: checkerpb.Action_ACTION_CHECK,
				Status: checkerpb.Status_STATUS_UP,
			},
			want: true,
		},
		{
			name: "status mumble returns false",
			verdict: checkers.Verdict{
				Action: checkerpb.Action_ACTION_CHECK,
				Status: checkerpb.Status_STATUS_MUMBLE,
			},
			want: false,
		},
		{
			name: "status corrupt returns false",
			verdict: checkers.Verdict{
				Action: checkerpb.Action_ACTION_GET,
				Status: checkerpb.Status_STATUS_CORRUPT,
			},
			want: false,
		},
		{
			name: "status down returns false",
			verdict: checkers.Verdict{
				Action: checkerpb.Action_ACTION_CHECK,
				Status: checkerpb.Status_STATUS_DOWN,
			},
			want: false,
		},
		{
			name: "status check_failed returns false",
			verdict: checkers.Verdict{
				Action: checkerpb.Action_ACTION_CHECK,
				Status: checkerpb.Status_STATUS_CHECK_FAILED,
			},
			want: false,
		},
		{
			name: "status unspecified returns false",
			verdict: checkers.Verdict{
				Action: checkerpb.Action_ACTION_CHECK,
				Status: checkerpb.Status_STATUS_UNSPECIFIED,
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.verdict.IsUp(); got != tt.want {
				t.Errorf("Verdict.IsUp() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVerdict_Fields(t *testing.T) {
	// Test that all fields can be set and accessed
	verdict := checkers.Verdict{
		Action:  checkerpb.Action_ACTION_PUT,
		Status:  checkerpb.Status_STATUS_UP,
		Public:  "flag stored successfully",
		Private: "internal error details",
		Command: "put flag ABC123",
	}

	if verdict.Action != checkerpb.Action_ACTION_PUT {
		t.Errorf("Action = %v, want %v", verdict.Action, checkerpb.Action_ACTION_PUT)
	}
	if verdict.Status != checkerpb.Status_STATUS_UP {
		t.Errorf("Status = %v, want %v", verdict.Status, checkerpb.Status_STATUS_UP)
	}
	if verdict.Public != "flag stored successfully" {
		t.Errorf("Public = %q, want %q", verdict.Public, "flag stored successfully")
	}
	if verdict.Private != "internal error details" {
		t.Errorf("Private = %q, want %q", verdict.Private, "internal error details")
	}
	if verdict.Command != "put flag ABC123" {
		t.Errorf("Command = %q, want %q", verdict.Command, "put flag ABC123")
	}
}
