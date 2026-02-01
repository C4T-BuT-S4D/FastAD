package cmd_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/c4t-but-s4d/fastad/cmd/fastad/cmd"
)

func TestDefaultTemporalPostgresConfig(t *testing.T) {
	cfg := cmd.DefaultTemporalPostgresConfig()

	assert.Equal(t, "temporal-postgres", cfg.Host)
	assert.Equal(t, "5432", cfg.Port)
	assert.Equal(t, "temporal", cfg.User)
	assert.Equal(t, "temporal", cfg.Password)
	assert.Equal(t, "temporal", cfg.Database)
	assert.Equal(t, "false", cfg.TLS)
	assert.Equal(t, "false", cfg.SkipDBCreate)
}

func TestParseTemporalDSN(t *testing.T) {
	tests := []struct {
		name    string
		dsn     string
		want    cmd.TemporalPostgresConfig
		wantErr string
	}{
		{
			name: "full dsn",
			dsn:  "postgres://user:pass@host:5433/mydb",
			want: cmd.TemporalPostgresConfig{
				User:         "user",
				Password:     "pass",
				Host:         "host",
				Port:         "5433",
				Database:     "mydb",
				TLS:          "false",
				SkipDBCreate: "true",
			},
		},
		{
			name: "dsn without password",
			dsn:  "postgres://user@host:5432/mydb",
			want: cmd.TemporalPostgresConfig{
				User:         "user",
				Password:     "",
				Host:         "host",
				Port:         "5432",
				Database:     "mydb",
				TLS:          "false",
				SkipDBCreate: "true",
			},
		},
		{
			name: "dsn with leading slash in path",
			dsn:  "postgres://user:pass@host:5432/database",
			want: cmd.TemporalPostgresConfig{
				User:         "user",
				Password:     "pass",
				Host:         "host",
				Port:         "5432",
				Database:     "database",
				TLS:          "false",
				SkipDBCreate: "true",
			},
		},
		{
			name:    "invalid scheme",
			dsn:     "mysql://user:pass@host:3306/db",
			wantErr: "scheme must be postgres",
		},
		{
			name:    "invalid url",
			dsn:     "://invalid",
			wantErr: "parsing temporal dsn",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := cmd.ParseTemporalDSN(tt.dsn)
			if tt.wantErr != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErr)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestComposeManipulator(t *testing.T) {
	tmpDir := t.TempDir()
	composePath := filepath.Join(tmpDir, "compose.yml")

	initialContent := `services:
  app:
    image: myapp
    depends_on:
      postgres:
        condition: service_healthy
    volumes:
      - ./data:/data
  postgres:
    image: postgres:17
  redis:
    image: redis:7
volumes:
  app-data:
  postgres-data:
`

	err := os.WriteFile(composePath, []byte(initialContent), 0o644)
	require.NoError(t, err)

	t.Run("LoadCompose", func(t *testing.T) {
		compose, err := cmd.LoadCompose(composePath)
		require.NoError(t, err)
		assert.Len(t, compose.Services(), 3)
		assert.Len(t, compose.Volumes(), 2)
	})

	t.Run("RemoveService", func(t *testing.T) {
		compose, err := cmd.LoadCompose(composePath)
		require.NoError(t, err)

		compose.RemoveService("postgres")
		assert.Len(t, compose.Services(), 2)
		_, exists := compose.Services()["postgres"]
		assert.False(t, exists)
	})

	t.Run("RemoveVolume", func(t *testing.T) {
		compose, err := cmd.LoadCompose(composePath)
		require.NoError(t, err)

		compose.RemoveVolume("postgres-data")
		assert.Len(t, compose.Volumes(), 1)
		_, exists := compose.Volumes()["postgres-data"]
		assert.False(t, exists)
	})

	t.Run("RemoveDependsOn", func(t *testing.T) {
		compose, err := cmd.LoadCompose(composePath)
		require.NoError(t, err)

		compose.RemoveDependsOn("app")
		appService := compose.Services()["app"].(map[string]any)
		_, hasDeps := appService["depends_on"]
		assert.False(t, hasDeps)
	})

	t.Run("SetDependsOn", func(t *testing.T) {
		compose, err := cmd.LoadCompose(composePath)
		require.NoError(t, err)

		newDeps := map[string]any{
			"redis": map[string]any{"condition": "service_started"},
		}
		compose.SetDependsOn("app", newDeps)

		appService := compose.Services()["app"].(map[string]any)
		deps := appService["depends_on"].(map[string]any)
		redisDep := deps["redis"].(map[string]any)
		assert.Equal(t, "service_started", redisDep["condition"])
	})

	t.Run("SetVolumes", func(t *testing.T) {
		compose, err := cmd.LoadCompose(composePath)
		require.NoError(t, err)

		compose.SetVolumes("app", []string{"/new/path:/container/path"})

		appService := compose.Services()["app"].(map[string]any)
		volumes := appService["volumes"].([]string)
		require.Len(t, volumes, 1)
		assert.Equal(t, "/new/path:/container/path", volumes[0])
	})

	t.Run("Write", func(t *testing.T) {
		compose, err := cmd.LoadCompose(composePath)
		require.NoError(t, err)

		compose.RemoveService("redis")

		outputPath := filepath.Join(tmpDir, "output.yml")
		err = compose.Write(outputPath)
		require.NoError(t, err)

		reloaded, err := cmd.LoadCompose(outputPath)
		require.NoError(t, err)
		assert.Len(t, reloaded.Services(), 2)
	})

	t.Run("LoadCompose_FileNotFound", func(t *testing.T) {
		_, err := cmd.LoadCompose("/nonexistent/path/compose.yml")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "reading compose file")
	})
}

func TestWriteEnvFile(t *testing.T) {
	tmpDir := t.TempDir()
	envPath := filepath.Join(tmpDir, ".env")

	content := []byte("KEY=value\nANOTHER=test\n")
	err := cmd.WriteEnvFile(envPath, content)
	require.NoError(t, err)

	read, err := os.ReadFile(envPath)
	require.NoError(t, err)
	assert.Equal(t, content, read)
}

func TestPresetComposePath(t *testing.T) {
	path := cmd.PresetComposePath("/root/project", "simple")
	assert.Equal(t, "/root/project/docker/presets/simple/compose.yml", path)
}
