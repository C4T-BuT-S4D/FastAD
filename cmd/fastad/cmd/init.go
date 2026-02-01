package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"

	"github.com/c4t-but-s4d/fastad/internal/gameconfig"
)

func NewInitCmd(cc *Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize game infrastructure (generate compose.yml, .env, etc.)",
		RunE: func(_ *cobra.Command, _ []string) error {
			root, err := gameconfig.GetFastADRoot()
			if err != nil {
				return fmt.Errorf("getting fastad root: %w", err)
			}

			configPath := cc.Viper.GetString("game_config")
			if configPath == "" {
				configPath = filepath.Join(root, "fastad.yaml")
			}

			zap.S().Infof("reading game config from %s", configPath)

			content, err := os.ReadFile(configPath)
			if os.IsNotExist(err) {
				examplePath := filepath.Join(root, "fastad.yaml.example")
				return fmt.Errorf(
					"config file not found: %s\n\nTo get started:\n  cp %s %s\n  # edit %s with your game settings\n  fastad init",
					configPath, examplePath, configPath, configPath,
				)
			}
			if err != nil {
				return fmt.Errorf("reading game config: %w", err)
			}

			var cfg *gameconfig.GameConfig
			if err := yaml.Unmarshal(content, &cfg); err != nil {
				return fmt.Errorf("unmarshalling game config: %w", err)
			}

			zap.S().Infof("parsed game config: %+v", cfg)

			if err := cfg.Validate(); err != nil {
				return fmt.Errorf("validating game config: %w", err)
			}

			zap.S().Info("game config validated")

			// Check if already initialized
			generatedDir := filepath.Join(root, gameconfig.GeneratedDir)
			generatedConfigPath := filepath.Join(generatedDir, gameconfig.GeneratedGameConfig)
			if _, err := os.Stat(generatedConfigPath); err == nil && !cc.Viper.GetBool("force") {
				return fmt.Errorf(
					"game already initialized (config at %s exists). Use --force to reinitialize",
					generatedConfigPath,
				)
			}

			switch cc.Viper.GetString("preset") {
			case "simple":
				zap.L().Info("initializing simple preset")
				if err := PresetSimple(root, cfg); err != nil {
					return fmt.Errorf("initializing simple preset: %w", err)
				}
				zap.L().Info("simple preset initialized")
			case "production":
				zap.L().Info("initializing production preset")
				if err := PresetProduction(root, cfg); err != nil {
					return fmt.Errorf("initializing production preset: %w", err)
				}
				zap.L().Info("production preset initialized")
			default:
				return fmt.Errorf("unsupported preset: %s (supported: simple, production)", cc.Viper.GetString("preset"))
			}

			// Save generated config
			if err := os.MkdirAll(generatedDir, 0o755); err != nil {
				return fmt.Errorf("creating generated dir: %w", err)
			}
			configContent, err := yaml.Marshal(cfg)
			if err != nil {
				return fmt.Errorf("marshalling game config: %w", err)
			}
			if err := os.WriteFile(generatedConfigPath, configContent, 0o644); err != nil {
				return fmt.Errorf("writing generated game config: %w", err)
			}

			zap.L().Info("initialization complete",
				zap.String("generated_config", generatedConfigPath),
				zap.String("compose_file", filepath.Join(root, "compose.yml")),
			)
			zap.L().Info("run 'fastad run' to start the game")

			return nil
		},
	}

	cmd.Flags().StringP("game-config", "c", "", "path to game config yaml file (defaults to fastad.yaml in fastad root)")
	cmd.Flags().String("preset", "simple", "preset to use ('simple' or 'production')")
	cmd.Flags().Bool("force", false, "overwrite existing generated files")

	return cmd
}
