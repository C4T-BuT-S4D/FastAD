package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/c4t-but-s4d/fastad/internal/gameconfig"
)

func NewResetCmd(_ *Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "reset",
		Short: "Reset the game (stops containers, removes volumes and data)",
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx := cmd.Context()
			root, err := gameconfig.GetFastADRoot()
			if err != nil {
				return fmt.Errorf("getting fastad root: %w", err)
			}

			composePath := filepath.Join(root, "compose.yml")
			if _, err := os.Stat(composePath); os.IsNotExist(err) {
				zap.L().Info("compose.yml not found, nothing to reset")
				return nil
			}

			zap.L().Info("stopping and removing containers")
			stopCmd := exec.CommandContext(ctx, "docker", "compose", "-f", composePath, "down", "-v", "--remove-orphans")
			stopCmd.Dir = root
			stopCmd.Stdout = os.Stdout
			stopCmd.Stderr = os.Stderr
			if err := stopCmd.Run(); err != nil {
				return fmt.Errorf("stopping containers: %w", err)
			}

			dataDir := filepath.Join(root, "data")
			if _, err := os.Stat(dataDir); err == nil {
				zap.L().Info("removing data directory", zap.String("path", dataDir))
				if err := os.RemoveAll(dataDir); err != nil {
					zap.L().Warn("failed to remove data directory", zap.Error(err))
				}
			}

			zap.L().Info("reset complete")
			return nil
		},
	}

	return cmd
}
