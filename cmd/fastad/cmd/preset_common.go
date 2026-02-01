package cmd

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/compose-spec/compose-go/v2/dotenv"
	"github.com/compose-spec/compose-go/v2/loader"
	"github.com/compose-spec/compose-go/v2/types"
	"go.uber.org/zap"
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

// ComposeManipulator wraps a compose-go Project for manipulation.
type ComposeManipulator struct {
	project *types.Project
}

func LoadCompose(path string) (*ComposeManipulator, error) {
	return LoadComposeWithEnv(path, nil)
}

func LoadComposeWithEnv(path string, env map[string]string) (*ComposeManipulator, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading compose file: %w", err)
	}

	workingDir := filepath.Dir(path)

	environment := make(map[string]string)
	for _, e := range os.Environ() {
		if k, v, ok := strings.Cut(e, "="); ok {
			environment[k] = v
		}
	}

	envFile := filepath.Join(workingDir, ".env")
	if envFromFile, err := dotenv.Read(envFile); err == nil {
		for k, v := range envFromFile {
			environment[k] = v
		}
	}

	for k, v := range env {
		environment[k] = v
	}

	configDetails := types.ConfigDetails{
		WorkingDir: workingDir,
		ConfigFiles: []types.ConfigFile{
			{
				Filename: path,
				Content:  content,
			},
		},
		Environment: environment,
	}

	project, err := loader.LoadWithContext(context.Background(), configDetails, func(opts *loader.Options) {
		opts.SkipValidation = true
		opts.SkipNormalization = true
		opts.ResolvePaths = false
	})
	if err != nil {
		return nil, fmt.Errorf("loading compose file: %w", err)
	}

	return &ComposeManipulator{project: project}, nil
}

// Services returns the services map for compatibility with existing tests.
func (c *ComposeManipulator) Services() types.Services {
	return c.project.Services
}

// Volumes returns the volumes map for compatibility with existing tests.
func (c *ComposeManipulator) Volumes() types.Volumes {
	return c.project.Volumes
}

func (c *ComposeManipulator) RemoveService(name string) {
	delete(c.project.Services, name)
}

func (c *ComposeManipulator) RemoveVolume(name string) {
	delete(c.project.Volumes, name)
}

func (c *ComposeManipulator) RemoveDependsOn(serviceName string) {
	if svc, exists := c.project.Services[serviceName]; exists {
		svc.DependsOn = nil
		c.project.Services[serviceName] = svc
	}
}

func (c *ComposeManipulator) SetDependsOn(serviceName string, deps map[string]any) {
	svc, exists := c.project.Services[serviceName]
	if !exists {
		return
	}

	dependsOn := make(types.DependsOnConfig)
	for depName, depValue := range deps {
		switch v := depValue.(type) {
		case map[string]any:
			condition := ""
			if cond, ok := v["condition"].(string); ok {
				condition = cond
			}
			dependsOn[depName] = types.ServiceDependency{
				Condition: condition,
				Required:  true,
			}
		case string:
			dependsOn[depName] = types.ServiceDependency{
				Condition: v,
				Required:  true,
			}
		default:
			dependsOn[depName] = types.ServiceDependency{
				Condition: types.ServiceConditionStarted,
				Required:  true,
			}
		}
	}

	svc.DependsOn = dependsOn
	c.project.Services[serviceName] = svc
}

func (c *ComposeManipulator) SetVolumes(serviceName string, volumes []string) {
	svc, exists := c.project.Services[serviceName]
	if !exists {
		return
	}

	volumeConfigs := make([]types.ServiceVolumeConfig, 0, len(volumes))
	for _, vol := range volumes {
		parts := strings.SplitN(vol, ":", 3)
		source := parts[0]

		volType := types.VolumeTypeBind
		if !strings.HasPrefix(source, "/") && !strings.HasPrefix(source, ".") {
			volType = types.VolumeTypeVolume
		}

		volumeConfig := types.ServiceVolumeConfig{
			Type:   volType,
			Source: source,
		}
		if len(parts) > 1 {
			volumeConfig.Target = parts[1]
		}
		if len(parts) > 2 && parts[2] == "ro" {
			volumeConfig.ReadOnly = true
		}
		volumeConfigs = append(volumeConfigs, volumeConfig)
	}

	svc.Volumes = volumeConfigs
	c.project.Services[serviceName] = svc
}

func (c *ComposeManipulator) Write(path string) error {
	c.cleanupBindMounts()

	raw, err := c.project.MarshalYAML()
	if err != nil {
		return fmt.Errorf("marshalling compose file: %w", err)
	}

	if err := os.WriteFile(path, raw, 0o644); err != nil {
		return fmt.Errorf("writing compose file: %w", err)
	}

	zap.L().Info("wrote compose file", zap.String("path", path))
	return nil
}

func (c *ComposeManipulator) cleanupBindMounts() {
	for name, svc := range c.project.Services {
		for i := range svc.Volumes {
			vol := &svc.Volumes[i]
			if vol.Type == types.VolumeTypeBind {
				if vol.Bind == nil || (vol.Bind.SELinux == "" && vol.Bind.Propagation == "" && vol.Bind.Recursive == "") {
					vol.Bind = nil
				}
			}
		}
		c.project.Services[name] = svc
	}
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
