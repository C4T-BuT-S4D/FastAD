package cmd

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

type TemporalPostgresConfig struct {
	User         string
	Password     string
	Database     string
	Host         string
	Port         string
	TLS          string
	SkipDBCreate string
}

func DefaultTemporalPostgresConfig() TemporalPostgresConfig {
	return TemporalPostgresConfig{
		Host:         "temporal-postgres",
		Port:         "5432",
		User:         "temporal",
		Password:     "temporal",
		Database:     "temporal",
		TLS:          "false",
		SkipDBCreate: "false",
	}
}

func ParseTemporalDSN(dsn string) (TemporalPostgresConfig, error) {
	parsedDSN, err := url.Parse(dsn)
	if err != nil {
		return TemporalPostgresConfig{}, fmt.Errorf("parsing temporal dsn: %w", err)
	}
	if parsedDSN.Scheme != "postgres" {
		return TemporalPostgresConfig{}, errors.New("temporal dsn scheme must be postgres")
	}

	cfg := TemporalPostgresConfig{
		User:         parsedDSN.User.Username(),
		Host:         parsedDSN.Hostname(),
		Port:         parsedDSN.Port(),
		Database:     strings.TrimPrefix(parsedDSN.Path, "/"),
		TLS:          "false",
		SkipDBCreate: "true",
	}
	cfg.Password, _ = parsedDSN.User.Password()

	return cfg, nil
}

// ComposeFile represents a docker-compose.yml file.
type ComposeFile struct {
	Services map[string]any `yaml:"services,omitempty"`
	Volumes  map[string]any `yaml:"volumes,omitempty"`
	Networks map[string]any `yaml:"networks,omitempty"`
}

// ComposeManipulator provides simple compose file manipulation.
type ComposeManipulator struct {
	data *ComposeFile
}

func LoadCompose(path string) (*ComposeManipulator, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading compose file: %w", err)
	}

	var data ComposeFile
	if err := yaml.Unmarshal(content, &data); err != nil {
		return nil, fmt.Errorf("parsing compose file: %w", err)
	}

	if data.Services == nil {
		data.Services = make(map[string]any)
	}
	if data.Volumes == nil {
		data.Volumes = make(map[string]any)
	}

	return &ComposeManipulator{data: &data}, nil
}

func (c *ComposeManipulator) Services() map[string]any {
	return c.data.Services
}

func (c *ComposeManipulator) Volumes() map[string]any {
	return c.data.Volumes
}

func (c *ComposeManipulator) RemoveService(name string) {
	delete(c.data.Services, name)
}

func (c *ComposeManipulator) RemoveVolume(name string) {
	delete(c.data.Volumes, name)
}

func (c *ComposeManipulator) getService(name string) map[string]any {
	svc, ok := c.data.Services[name]
	if !ok {
		return nil
	}
	svcMap, ok := svc.(map[string]any)
	if !ok {
		return nil
	}
	return svcMap
}

func (c *ComposeManipulator) RemoveDependsOn(serviceName string) {
	svc := c.getService(serviceName)
	if svc == nil {
		return
	}
	delete(svc, "depends_on")
}

func (c *ComposeManipulator) SetDependsOn(serviceName string, deps map[string]any) {
	svc := c.getService(serviceName)
	if svc == nil {
		return
	}
	svc["depends_on"] = deps
}

func (c *ComposeManipulator) SetVolumes(serviceName string, volumes []string) {
	svc := c.getService(serviceName)
	if svc == nil {
		return
	}
	svc["volumes"] = volumes
}

func (c *ComposeManipulator) Write(path string) error {
	content, err := yaml.Marshal(c.data)
	if err != nil {
		return fmt.Errorf("marshalling compose file: %w", err)
	}

	if err := os.WriteFile(path, content, 0o644); err != nil {
		return fmt.Errorf("writing compose file: %w", err)
	}

	zap.L().Info("wrote compose file", zap.String("path", path))
	return nil
}

func WriteEnvFile(path string, content []byte) error {
	if err := os.WriteFile(path, content, 0o644); err != nil {
		return fmt.Errorf("writing env file: %w", err)
	}

	zap.L().Info("wrote .env file", zap.String("path", path))
	return nil
}

func PresetComposePath(root, preset string) string {
	return filepath.Join(root, "docker", "presets", preset, "compose.yml")
}
