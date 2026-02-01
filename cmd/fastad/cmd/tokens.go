package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/c4t-but-s4d/fastad/internal/gameconfig"
	"github.com/c4t-but-s4d/fastad/pkg/clients/teams"
	"github.com/c4t-but-s4d/fastad/pkg/grpcext"
	teamspb "github.com/c4t-but-s4d/fastad/pkg/proto/data/teams"
)

func NewTokensCmd(cc *Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tokens",
		Short: "Print team tokens",
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx := cmd.Context()
			cfg, err := gameconfig.ReadGeneratedConfig()
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
			defer apiConn.Close()

			teamsClient := teams.NewClient(teamspb.NewTeamsServiceClient(apiConn), "fastad-cli")

			teamsList, err := teamsClient.List(ctx)
			if err != nil {
				return fmt.Errorf("listing teams: %w", err)
			}

			switch cc.Viper.GetString("format") {
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
				return fmt.Errorf("unsupported format: %s", cc.Viper.GetString("format"))
			}

			return nil
		},
	}

	cmd.Flags().StringP("format", "f", "text", "output format (json, text)")

	return cmd
}
