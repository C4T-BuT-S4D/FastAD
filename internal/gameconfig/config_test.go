package gameconfig_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/c4t-but-s4d/fastad/internal/gameconfig"
)

func TestGame_Validate(t *testing.T) {
	tests := []struct {
		name    string
		game    gameconfig.Game
		wantErr string
	}{
		{
			name:    "missing start_time",
			game:    gameconfig.Game{Hardness: 10},
			wantErr: "start_time required",
		},
		{
			name: "end_time before start_time",
			game: gameconfig.Game{
				StartTime: time.Date(2025, 1, 1, 12, 0, 0, 0, time.UTC),
				EndTime:   ptrTime(time.Date(2025, 1, 1, 10, 0, 0, 0, time.UTC)),
				Hardness:  10,
			},
			wantErr: "end_time is before start_time",
		},
		{
			name: "zero hardness",
			game: gameconfig.Game{
				StartTime: time.Date(2025, 1, 1, 12, 0, 0, 0, time.UTC),
				Hardness:  0,
			},
			wantErr: "hardness must be positive",
		},
		{
			name: "negative hardness",
			game: gameconfig.Game{
				StartTime: time.Date(2025, 1, 1, 12, 0, 0, 0, time.UTC),
				Hardness:  -5,
			},
			wantErr: "hardness must be positive",
		},
		{
			name: "valid game with defaults",
			game: gameconfig.Game{
				StartTime: time.Date(2025, 1, 1, 12, 0, 0, 0, time.UTC),
				Hardness:  10,
			},
			wantErr: "",
		},
		{
			name: "valid game with end_time",
			game: gameconfig.Game{
				StartTime: time.Date(2025, 1, 1, 12, 0, 0, 0, time.UTC),
				EndTime:   ptrTime(time.Date(2025, 1, 1, 18, 0, 0, 0, time.UTC)),
				Hardness:  10,
			},
			wantErr: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.game.Validate()
			if tt.wantErr != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestGame_Validate_DefaultCheckersBasePath(t *testing.T) {
	g := gameconfig.Game{
		StartTime:        time.Date(2025, 1, 1, 12, 0, 0, 0, time.UTC),
		Hardness:         10,
		CheckersBasePath: "",
	}

	err := g.Validate()
	require.NoError(t, err)
	assert.Equal(t, "checkers", g.CheckersBasePath)
}

func TestTeam_Validate(t *testing.T) {
	tests := []struct {
		name    string
		team    gameconfig.Team
		wantErr string
	}{
		{
			name:    "missing name",
			team:    gameconfig.Team{Address: "10.0.0.1"},
			wantErr: "name required",
		},
		{
			name:    "missing address",
			team:    gameconfig.Team{Name: "Team A"},
			wantErr: "address required",
		},
		{
			name: "valid team",
			team: gameconfig.Team{
				Name:    "Team A",
				Address: "10.0.0.1",
			},
			wantErr: "",
		},
		{
			name: "valid team with labels",
			team: gameconfig.Team{
				Name:    "Team A",
				Address: "10.0.0.1",
				Labels:  map[string]string{"country": "US"},
			},
			wantErr: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.team.Validate()
			if tt.wantErr != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestChecker_Validate(t *testing.T) {
	tests := []struct {
		name    string
		checker gameconfig.Checker
		wantErr string
	}{
		{
			name:    "missing path",
			checker: gameconfig.Checker{DefaultTimeout: time.Second * 30},
			wantErr: "path required",
		},
		{
			name:    "missing default_timeout",
			checker: gameconfig.Checker{Path: "checker.py"},
			wantErr: "default_timeout required",
		},
		{
			name: "valid checker",
			checker: gameconfig.Checker{
				Path:           "checker.py",
				DefaultTimeout: time.Second * 30,
			},
			wantErr: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.checker.Validate()
			if tt.wantErr != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestService_Validate(t *testing.T) {
	tests := []struct {
		name    string
		service gameconfig.Service
		wantErr string
	}{
		{
			name:    "missing name",
			service: gameconfig.Service{DefaultScore: 2500},
			wantErr: "name required",
		},
		{
			name:    "missing default_score",
			service: gameconfig.Service{Name: "test"},
			wantErr: "default_score required",
		},
		{
			name: "missing checker",
			service: gameconfig.Service{
				Name:         "test",
				DefaultScore: 2500,
			},
			wantErr: "path required",
		},
		{
			name: "valid service",
			service: gameconfig.Service{
				Name:         "test",
				DefaultScore: 2500,
				Checker: &gameconfig.Checker{
					Path:           "checker.py",
					DefaultTimeout: time.Second * 30,
				},
			},
			wantErr: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.service.Validate()
			if tt.wantErr != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestFastAD_Validate(t *testing.T) {
	t.Run("defaults are set", func(t *testing.T) {
		f := gameconfig.FastAD{}
		err := f.Validate()
		require.NoError(t, err)

		assert.Equal(t, 8080, f.ListenPort)
		assert.Equal(t, "info", f.LogLevel)
		assert.NotEmpty(t, f.IntercomToken)
	})

	t.Run("preserves existing values", func(t *testing.T) {
		f := gameconfig.FastAD{
			ListenPort:    9090,
			LogLevel:      "debug",
			IntercomToken: "my-token",
		}
		err := f.Validate()
		require.NoError(t, err)

		assert.Equal(t, 9090, f.ListenPort)
		assert.Equal(t, "debug", f.LogLevel)
		assert.Equal(t, "my-token", f.IntercomToken)
	})
}

func TestAdmin_Validate(t *testing.T) {
	tests := []struct {
		name    string
		admin   gameconfig.Admin
		wantErr string
	}{
		{
			name:    "missing username",
			admin:   gameconfig.Admin{Password: "secret"},
			wantErr: "username required",
		},
		{
			name:    "missing password",
			admin:   gameconfig.Admin{Username: "admin"},
			wantErr: "password required",
		},
		{
			name: "valid admin",
			admin: gameconfig.Admin{
				Username: "admin",
				Password: "secret",
			},
			wantErr: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.admin.Validate()
			if tt.wantErr != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestGameConfig_Validate(t *testing.T) {
	validGame := &gameconfig.Game{
		StartTime: time.Date(2025, 1, 1, 12, 0, 0, 0, time.UTC),
		Hardness:  10,
	}

	validTeam := &gameconfig.Team{
		Name:    "Team A",
		Address: "10.0.0.1",
	}

	validService := &gameconfig.Service{
		Name:         "test",
		DefaultScore: 2500,
		Checker: &gameconfig.Checker{
			Path:           "checker.py",
			DefaultTimeout: time.Second * 30,
		},
	}

	t.Run("missing game", func(t *testing.T) {
		cfg := &gameconfig.GameConfig{}
		err := cfg.Validate()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "game required")
	})

	t.Run("invalid game", func(t *testing.T) {
		cfg := &gameconfig.GameConfig{
			Game: &gameconfig.Game{},
		}
		err := cfg.Validate()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "game:")
	})

	t.Run("nil team in list", func(t *testing.T) {
		cfg := &gameconfig.GameConfig{
			Game:  validGame,
			Teams: []*gameconfig.Team{nil},
		}
		err := cfg.Validate()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "team 0: nil")
	})

	t.Run("invalid team", func(t *testing.T) {
		cfg := &gameconfig.GameConfig{
			Game:  validGame,
			Teams: []*gameconfig.Team{{Name: ""}},
		}
		err := cfg.Validate()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "team 0:")
	})

	t.Run("nil service in list", func(t *testing.T) {
		cfg := &gameconfig.GameConfig{
			Game:     validGame,
			Teams:    []*gameconfig.Team{validTeam},
			Services: []*gameconfig.Service{nil},
		}
		err := cfg.Validate()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "service 0: nil")
	})

	t.Run("invalid service", func(t *testing.T) {
		cfg := &gameconfig.GameConfig{
			Game:     validGame,
			Teams:    []*gameconfig.Team{validTeam},
			Services: []*gameconfig.Service{{Name: ""}},
		}
		err := cfg.Validate()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "service 0:")
	})

	t.Run("auto-generates admin credentials", func(t *testing.T) {
		cfg := &gameconfig.GameConfig{
			Game:     validGame,
			Teams:    []*gameconfig.Team{validTeam},
			Services: []*gameconfig.Service{validService},
		}
		err := cfg.Validate()
		require.NoError(t, err)
		assert.Equal(t, "fastad", cfg.Admin.Username)
		assert.NotEmpty(t, cfg.Admin.Password)
		assert.NotContains(t, cfg.Admin.Password, "-")
	})

	t.Run("valid config", func(t *testing.T) {
		cfg := &gameconfig.GameConfig{
			Game: validGame,
			Admin: &gameconfig.Admin{
				Username: "admin",
				Password: "secret",
			},
			Teams:    []*gameconfig.Team{validTeam},
			Services: []*gameconfig.Service{validService},
		}
		err := cfg.Validate()
		require.NoError(t, err)
	})
}

func ptrTime(t time.Time) *time.Time {
	return &t
}
