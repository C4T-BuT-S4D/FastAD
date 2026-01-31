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

type ComposeManipulator struct {
	Services map[string]any `yaml:"services"`
	Volumes  map[string]any `yaml:"volumes"`
}

func LoadCompose(path string) (*ComposeManipulator, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading compose file: %w", err)
	}

	var compose ComposeManipulator
	if err := yaml.Unmarshal(content, &compose); err != nil {
		return nil, fmt.Errorf("unmarshalling compose file: %w", err)
	}

	return &compose, nil
}

func (c *ComposeManipulator) RemoveService(name string) {
	delete(c.Services, name)
}

func (c *ComposeManipulator) RemoveVolume(name string) {
	delete(c.Volumes, name)
}

func (c *ComposeManipulator) RemoveDependsOn(serviceName string) {
	if svc, ok := c.Services[serviceName].(map[string]any); ok {
		delete(svc, "depends_on")
	}
}

func (c *ComposeManipulator) SetDependsOn(serviceName string, deps map[string]any) {
	if svc, ok := c.Services[serviceName].(map[string]any); ok {
		svc["depends_on"] = deps
	}
}

func (c *ComposeManipulator) SetVolumes(serviceName string, volumes []string) {
	if svc, ok := c.Services[serviceName].(map[string]any); ok {
		svc["volumes"] = volumes
	}
}

func (c *ComposeManipulator) Write(path string) error {
	raw, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("marshalling compose file: %w", err)
	}

	if err := os.WriteFile(path, raw, 0o644); err != nil {
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
