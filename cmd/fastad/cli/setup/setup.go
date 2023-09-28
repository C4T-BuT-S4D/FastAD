package setup

import (
	"fmt"
	"os"

	"github.com/c4t-but-s4d/fastad/cmd/fastad/cli/common"
	"github.com/c4t-but-s4d/fastad/internal/clients/teams"
	"github.com/c4t-but-s4d/fastad/pkg/grpctools"
	teamspb "github.com/c4t-but-s4d/fastad/pkg/proto/data/teams"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v3"
)

func NewSetupCommand(*common.CommandContext) *cli.Command {
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
				Value: "localhost:8080",
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

			teamsToCreate := lo.Map(cfg.Teams, func(t Team, _ int) *teamspb.Team {
				return t.ToProto()
			})

			createdTeams, err := teamsClient.CreateBatch(c.Context, teamsToCreate)
			if err != nil {
				return fmt.Errorf("creating teams: %w", err)
			}

			logrus.Infof("created teams: %+v", createdTeams)

			return nil
		},
	}
}
