package setup

import (
	"fmt"
	"os"

	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v3"

	"github.com/c4t-but-s4d/fastad/cmd/fastad/cli/common"
	"github.com/c4t-but-s4d/fastad/internal/clients/gamestate"
	"github.com/c4t-but-s4d/fastad/internal/clients/services"
	"github.com/c4t-but-s4d/fastad/internal/clients/teams"
	"github.com/c4t-but-s4d/fastad/pkg/grpctools"
	gspb "github.com/c4t-but-s4d/fastad/pkg/proto/data/game_state"
	servicespb "github.com/c4t-but-s4d/fastad/pkg/proto/data/services"
	teamspb "github.com/c4t-but-s4d/fastad/pkg/proto/data/teams"
)

func NewSetupCommand(_ *common.CommandContext) *cli.Command {
	return &cli.Command{
		Name: "setup",
		Flags: []cli.Flag{
			&cli.PathFlag{
				Name:    "game-config",
				Aliases: []string{"c"},
				Usage:   "path to game config yaml file",
				Value:   "game.yml",
			},
			&cli.StringFlag{
				Name:  "api-address",
				Usage: "fastad api address",
				Value: "127.0.0.1:1337",
			},
		},
		Action: func(c *cli.Context) error {
			logrus.Infof("reading game config from %s", c.Path("game-config"))

			content, err := os.ReadFile(c.Path("game-config"))
			if err != nil {
				return fmt.Errorf("reading game config: %w", err)
			}

			var cfg *GameConfig
			if err := yaml.Unmarshal(content, &cfg); err != nil {
				return fmt.Errorf("unmarshalling game config: %w", err)
			}

			logrus.Infof("parsed game config: %+v", cfg)

			if err := cfg.Validate(); err != nil {
				return fmt.Errorf("validating game config: %w", err)
			}

			logrus.Info("game config validated")

			apiConn, err := grpctools.Dial(c.String("api-address"), "fastad-setup")
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
			logrus.Infof("created teams: %+v", createdTeams)

			createdServices, err := servicesClient.CreateBatch(c.Context, servicesToCreate)
			if err != nil {
				return fmt.Errorf("creating services: %w", err)
			}
			logrus.Infof("created services: %+v", createdServices)

			createdGameService, err := gameStateClient.Update(c.Context, cfg.Game.ToUpdateRequestProto())
			if err != nil {
				return fmt.Errorf("updating game state: %w", err)
			}
			logrus.Infof("created game state: %+v", createdGameService)

			return nil
		},
	}
}
