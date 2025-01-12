package run

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/samber/lo"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"

	"github.com/c4t-but-s4d/fastad/cmd/fastad/cli/common"
	"github.com/c4t-but-s4d/fastad/pkg/apiwait"
	"github.com/c4t-but-s4d/fastad/pkg/clients/gamestate"
	"github.com/c4t-but-s4d/fastad/pkg/clients/services"
	"github.com/c4t-but-s4d/fastad/pkg/clients/teams"
	"github.com/c4t-but-s4d/fastad/pkg/grpcext"
	gspb "github.com/c4t-but-s4d/fastad/pkg/proto/data/game_state"
	servicespb "github.com/c4t-but-s4d/fastad/pkg/proto/data/services"
	teamspb "github.com/c4t-but-s4d/fastad/pkg/proto/data/teams"
)

func NewRunCommand(_ *common.CommandContext) *cli.Command {
	return &cli.Command{
		Name:  "run",
		Usage: "Setup & run the game",
		Flags: []cli.Flag{
			&cli.PathFlag{
				Name:    "game-config",
				Aliases: []string{"c"},
				Usage:   "path to game config yaml file (defaults to fastad.yaml in fastad root)",
			},
			&cli.StringFlag{
				Name:  "preset",
				Usage: "preset to use (only 'simple' is supported)",
				Value: "simple",
			},
			&cli.BoolFlag{
				Name:  "only-init",
				Usage: "only initialize the game state for an already running game",
			},
		},
		Action: func(c *cli.Context) error {
			root, err := common.GetFastADRoot()
			if err != nil {
				return fmt.Errorf("getting fastad root: %w", err)
			}

			configPath := c.Path("game-config")
			if configPath == "" {
				configPath = filepath.Join(root, "fastad.yaml")
			}

			zap.S().Infof("reading game config from %s", configPath)

			content, err := os.ReadFile(configPath)
			if err != nil {
				return fmt.Errorf("reading game config: %w", err)
			}

			var cfg *GameConfig
			if err := yaml.Unmarshal(content, &cfg); err != nil {
				return fmt.Errorf("unmarshalling game config: %w", err)
			}

			zap.S().Infof("parsed game config: %+v", cfg)

			if err := cfg.Validate(); err != nil {
				return fmt.Errorf("validating game config: %w", err)
			}

			zap.S().Info("game config validated")

			if !c.Bool("only-init") {
				switch c.String("preset") {
				case "simple":
					zap.L().Info("starting simple preset")
					if err := StartPresetSimple(c.Context, root, cfg); err != nil {
						return fmt.Errorf("starting simple preset: %w", err)
					}
					zap.L().Info("simple preset started")
				default:
					return fmt.Errorf("unsupported preset: %s", c.String("preset"))
				}
			} else {
				zap.L().Info("skipping starting the services")
			}

			zap.L().Info("waiting for services to start")
			httpAddr := fmt.Sprintf("127.0.0.1:%d", cfg.FastAD.ListenPort)
			if err := apiwait.HTTP(c.Context, httpAddr); err != nil {
				return fmt.Errorf("waiting for API service: %w", err)
			}

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

			teamsClient := teams.NewClient(teamspb.NewTeamsServiceClient(apiConn))
			servicesClient := services.NewClient(servicespb.NewServicesServiceClient(apiConn))
			gameStateClient := gamestate.NewClient(gspb.NewGameStateServiceClient(apiConn))

			teamsToCreate := lo.Map(cfg.Teams, func(t *Team, _ int) *teamspb.Team {
				return t.ToProto()
			})

			servicesToCreate := lo.Map(cfg.Services, func(s *Service, _ int) *servicespb.Service {
				return s.ToProto()
			})

			createdTeams, err := teamsClient.CreateBatch(c.Context, teamsToCreate)
			if err != nil {
				return fmt.Errorf("creating teams: %w", err)
			}
			zap.S().Infof("created teams: %v", createdTeams)

			createdServices, err := servicesClient.CreateBatch(c.Context, servicesToCreate)
			if err != nil {
				return fmt.Errorf("creating services: %w", err)
			}
			zap.S().Infof("created services: %v", createdServices)

			createdGameService, err := gameStateClient.Create(c.Context, cfg.Game.ToCreateRequestProto())
			if err != nil {
				return fmt.Errorf("updating game state: %w", err)
			}
			zap.S().Infof("created game state: %+v", createdGameService)

			return nil
		},
	}
}
