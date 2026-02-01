package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/c4t-but-s4d/fastad/internal/gameconfig"
	"github.com/c4t-but-s4d/fastad/pkg/apiwait"
	"github.com/c4t-but-s4d/fastad/pkg/clients/gamestate"
	"github.com/c4t-but-s4d/fastad/pkg/clients/services"
	"github.com/c4t-but-s4d/fastad/pkg/clients/teams"
	"github.com/c4t-but-s4d/fastad/pkg/grpcext"
	gspb "github.com/c4t-but-s4d/fastad/pkg/proto/data/game_state"
	servicespb "github.com/c4t-but-s4d/fastad/pkg/proto/data/services"
	teamspb "github.com/c4t-but-s4d/fastad/pkg/proto/data/teams"
)

func NewRunCmd(cc *Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "run",
		Short: "Start the game (runs docker compose and initializes game state)",
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx := cmd.Context()
			root, err := gameconfig.GetFastADRoot()
			if err != nil {
				return fmt.Errorf("getting fastad root: %w", err)
			}

			// Read the generated config (created by 'fastad init')
			cfg, err := gameconfig.ReadGeneratedConfig()
			if err != nil {
				return fmt.Errorf("reading generated config: %w (did you run 'fastad init' first?)", err)
			}

			composePath := filepath.Join(root, "compose.yml")
			if _, err := os.Stat(composePath); os.IsNotExist(err) {
				return fmt.Errorf("compose.yml not found at %s (did you run 'fastad init' first?)", composePath)
			}

			if !cc.Viper.GetBool("only_init") {
				// Stop existing services first
				zap.L().Info("stopping existing services")
				stopCmd := exec.CommandContext(ctx, "docker", "compose", "-f", composePath, "down", "-v", "--remove-orphans")
				stopCmd.Dir = root
				stopCmd.Stdout = os.Stdout
				stopCmd.Stderr = os.Stderr
				if err := stopCmd.Run(); err != nil {
					return fmt.Errorf("stopping existing services: %w", err)
				}

				// Start services
				args := []string{"compose", "-f", composePath, "up", "-d"}
				if !cc.Viper.GetBool("no_build") {
					args = append(args, "--build")
				}

				dockerCmd := exec.CommandContext(ctx, "docker", args...)
				dockerCmd.Dir = root
				dockerCmd.Stdout = os.Stdout
				dockerCmd.Stderr = os.Stderr

				zap.L().Info("starting services", zap.String("command", dockerCmd.String()))
				if err := dockerCmd.Run(); err != nil {
					return fmt.Errorf("running docker compose: %w", err)
				}
			} else {
				zap.L().Info("skipping starting the services (--only-init)")
			}

			// Wait for services to be ready
			zap.L().Info("waiting for services to start")
			httpAddr := fmt.Sprintf("127.0.0.1:%d", cfg.FastAD.ListenPort)
			if err := apiwait.HTTP(ctx, httpAddr); err != nil {
				return fmt.Errorf("waiting for API service: %w", err)
			}

			// Initialize game state
			apiAddress := fmt.Sprintf("127.0.0.1:%d", cfg.FastAD.ListenPort)
			zap.L().Info("initializing game", zap.String("api_address", apiAddress))

			apiConn, err := grpcext.Dial(
				apiAddress,
				"fastad-setup",
				grpcext.AuthDialOptions(cfg.FastAD.IntercomToken)...,
			)
			if err != nil {
				return fmt.Errorf("dialing data service: %w", err)
			}
			defer apiConn.Close()

			teamsClient := teams.NewClient(teamspb.NewTeamsServiceClient(apiConn), "fastad-cli")
			servicesClient := services.NewClient(servicespb.NewServicesServiceClient(apiConn), "fastad-cli")
			gameStateClient := gamestate.NewClient(gspb.NewGameStateServiceClient(apiConn), "fastad-cli")

			teamsToCreate := lo.Map(cfg.Teams, func(t *gameconfig.Team, _ int) *teamspb.Team {
				return t.ToProto()
			})

			servicesToCreate := lo.Map(cfg.Services, func(s *gameconfig.Service, _ int) *servicespb.Service {
				return s.ToProto()
			})

			createdTeams, err := teamsClient.CreateBatch(ctx, teamsToCreate)
			if err != nil {
				return fmt.Errorf("creating teams: %w", err)
			}
			zap.S().Infof("created teams: %v", createdTeams)

			createdServices, err := servicesClient.CreateBatch(ctx, servicesToCreate)
			if err != nil {
				return fmt.Errorf("creating services: %w", err)
			}
			zap.S().Infof("created services: %v", createdServices)

			// Try to get existing game state first
			existingGameState, err := gameStateClient.Get(ctx)
			if err == nil && existingGameState != nil {
				zap.S().Infof("game state already exists, skipping creation: %+v", existingGameState)
			} else {
				// Create new game state
				createdGameState, err := gameStateClient.Create(ctx, cfg.Game.ToCreateRequestProto())
				if err != nil {
					return fmt.Errorf("creating game state: %w", err)
				}
				zap.S().Infof("created game state: %+v", createdGameState)
			}

			zap.L().Info("game started successfully")

			return nil
		},
	}

	cmd.Flags().Bool("no-build", false, "skip building images (use --build by default)")
	cmd.Flags().Bool("only-init", false, "only initialize the game state (skip starting containers)")

	return cmd
}
