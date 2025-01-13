package tokens

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/urfave/cli/v2"

	"github.com/c4t-but-s4d/fastad/cmd/fastad/cli/common"
	"github.com/c4t-but-s4d/fastad/pkg/clients/teams"
	"github.com/c4t-but-s4d/fastad/pkg/grpcext"
	teamspb "github.com/c4t-but-s4d/fastad/pkg/proto/data/teams"
)

func NewCommand(_ *common.CommandContext) *cli.Command {
	return &cli.Command{
		Name:  "tokens",
		Usage: "print team tokens",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "format",
				Aliases: []string{"f"},
				Usage:   "output format (json, text)",
				Value:   "text",
			},
		},
		Action: func(c *cli.Context) error {
			cfg, err := common.ReadGeneratedConfig()
			if err != nil {
				return fmt.Errorf("reading generated config: %w", err)
			}

			apiAddress := fmt.Sprintf("127.0.0.1:%d", cfg.FastAD.ListenPort)
			apiConn, err := grpcext.Dial(
				apiAddress,
				"fastad-setup",
				grpcext.AuthDialOptions(cfg.FastAD.IntercomToken)...,
			)
			if err != nil {
				return fmt.Errorf("dialing data service: %w", err)
			}

			teamsClient := teams.NewClient(teamspb.NewTeamsServiceClient(apiConn))

			teamsList, err := teamsClient.List(c.Context)
			if err != nil {
				return fmt.Errorf("listing teams: %w", err)
			}

			switch c.String("format") {
			case "json":
				type teamToken struct {
					ID    int64  `json:"id"`
					Name  string `json:"name"`
					Token string `json:"token"`
				}
				data := make([]teamToken, 0, len(teamsList))
				for _, team := range teamsList {
					data = append(data, teamToken{
						ID:    team.GetId(),
						Name:  team.GetName(),
						Token: team.GetToken(),
					})
				}
				if err := json.NewEncoder(os.Stdout).Encode(data); err != nil {
					return fmt.Errorf("encoding JSON: %w", err)
				}
			case "text":
				for _, team := range teamsList {
					fmt.Printf("%s:%s\n", team.GetName(), team.GetToken())
				}
			default:
				return fmt.Errorf("unsupported format: %s", c.String("format"))
			}

			return nil
		},
	}
}
