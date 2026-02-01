package viperext_test

import (
	"context"
	"os"
	"testing"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/c4t-but-s4d/fastad/pkg/viperext"
)

func TestMain(m *testing.M) {
	cobra.EnableTraverseRunHooks = true
	os.Exit(m.Run())
}

func TestNestedSubcommandsFlagBinding(t *testing.T) {
	v := viper.New()

	var capturedSettings map[string]any
	var preRunOrder []string

	rootCmd := &cobra.Command{
		Use: "root",
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			preRunOrder = append(preRunOrder, "root")
			viperext.RunBindCommandFlags(v, cmd)
			return nil
		},
	}
	rootCmd.PersistentFlags().String("root-flag", "root-default", "Root persistent flag")
	rootCmd.PersistentFlags().Bool("verbose", false, "Verbose output")

	level1Cmd := &cobra.Command{
		Use: "level1",
		PersistentPreRunE: func(_ *cobra.Command, _ []string) error {
			preRunOrder = append(preRunOrder, "level1")
			return nil
		},
	}
	level1Cmd.PersistentFlags().String("level1-flag", "level1-default", "Level1 persistent flag")
	level1Cmd.PersistentFlags().Int("level1-count", 0, "Level1 count")

	level2Cmd := &cobra.Command{
		Use: "level2",
		PersistentPreRunE: func(_ *cobra.Command, _ []string) error {
			preRunOrder = append(preRunOrder, "level2")
			return nil
		},
	}
	level2Cmd.PersistentFlags().String("level2-flag", "level2-default", "Level2 persistent flag")

	level3Cmd := &cobra.Command{
		Use: "level3",
		PreRunE: func(_ *cobra.Command, _ []string) error {
			preRunOrder = append(preRunOrder, "level3-prerun")
			return nil
		},
		RunE: func(_ *cobra.Command, _ []string) error {
			preRunOrder = append(preRunOrder, "level3-run")
			capturedSettings = v.AllSettings()
			return nil
		},
	}
	level3Cmd.Flags().String("level3-flag", "level3-default", "Level3 local flag")
	level3Cmd.Flags().Bool("dry-run", false, "Dry run mode")

	level2Cmd.AddCommand(level3Cmd)
	level1Cmd.AddCommand(level2Cmd)
	rootCmd.AddCommand(level1Cmd)

	rootCmd.SetArgs([]string{
		"level1", "level2", "level3",
		"--root-flag", "custom-root",
		"--verbose",
		"--level1-flag", "custom-level1",
		"--level1-count", "42",
		"--level2-flag", "custom-level2",
		"--level3-flag", "custom-level3",
		"--dry-run",
	})

	err := rootCmd.ExecuteContext(context.Background())
	require.NoError(t, err)

	t.Run("all hooks executed in order", func(t *testing.T) {
		assert.Equal(t, []string{
			"root",
			"level1",
			"level2",
			"level3-prerun",
			"level3-run",
		}, preRunOrder)
	})

	t.Run("root persistent flags bound", func(t *testing.T) {
		assert.Equal(t, "custom-root", capturedSettings["root_flag"])
		assert.Equal(t, true, capturedSettings["verbose"])
	})

	t.Run("level1 persistent flags bound", func(t *testing.T) {
		assert.Equal(t, "custom-level1", capturedSettings["level1_flag"])
		assert.Equal(t, 42, capturedSettings["level1_count"])
	})

	t.Run("level2 persistent flags bound", func(t *testing.T) {
		assert.Equal(t, "custom-level2", capturedSettings["level2_flag"])
	})

	t.Run("level3 local flags bound", func(t *testing.T) {
		assert.Equal(t, "custom-level3", capturedSettings["level3_flag"])
		assert.Equal(t, true, capturedSettings["dry_run"])
	})

	t.Run("viper has all expected keys", func(t *testing.T) {
		expectedKeys := []string{
			"root_flag",
			"verbose",
			"level1_flag",
			"level1_count",
			"level2_flag",
			"level3_flag",
			"dry_run",
		}
		for _, key := range expectedKeys {
			_, exists := capturedSettings[key]
			assert.True(t, exists, "expected key %q to exist in viper settings", key)
		}
	})
}

func TestNestedSubcommandsWithDefaults(t *testing.T) {
	v := viper.New()

	var capturedSettings map[string]any

	rootCmd := &cobra.Command{
		Use: "root",
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			viperext.RunBindCommandFlags(v, cmd)
			return nil
		},
	}
	rootCmd.PersistentFlags().String("root-flag", "root-default", "Root flag")

	childCmd := &cobra.Command{
		Use: "child",
		RunE: func(_ *cobra.Command, _ []string) error {
			capturedSettings = v.AllSettings()
			return nil
		},
	}
	childCmd.Flags().String("child-flag", "child-default", "Child flag")

	rootCmd.AddCommand(childCmd)

	rootCmd.SetArgs([]string{"child"})

	err := rootCmd.ExecuteContext(context.Background())
	require.NoError(t, err)

	t.Run("default values are captured", func(t *testing.T) {
		assert.Equal(t, "root-default", capturedSettings["root_flag"])
		assert.Equal(t, "child-default", capturedSettings["child_flag"])
	})
}

func TestNestedSubcommandsWithEnvOverride(t *testing.T) {
	v := viper.New()
	v.SetEnvPrefix("TEST")
	v.AutomaticEnv()

	t.Setenv("TEST_ROOT_FLAG", "env-value")

	var capturedSettings map[string]any

	rootCmd := &cobra.Command{
		Use: "root",
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			viperext.RunBindCommandFlags(v, cmd)
			return nil
		},
	}
	rootCmd.PersistentFlags().String("root-flag", "default", "Root flag")

	childCmd := &cobra.Command{
		Use: "child",
		RunE: func(_ *cobra.Command, _ []string) error {
			capturedSettings = v.AllSettings()
			return nil
		},
	}

	rootCmd.AddCommand(childCmd)
	rootCmd.SetArgs([]string{"child"})

	err := rootCmd.ExecuteContext(context.Background())
	require.NoError(t, err)

	t.Run("env variable overrides default", func(t *testing.T) {
		assert.Equal(t, "env-value", capturedSettings["root_flag"])
	})
}

func TestIntermediateCommandExecution(t *testing.T) {
	v := viper.New()

	var capturedSettings map[string]any
	var executedCommand string

	rootCmd := &cobra.Command{
		Use: "root",
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			viperext.RunBindCommandFlags(v, cmd)
			return nil
		},
	}
	rootCmd.PersistentFlags().String("root-flag", "root-default", "Root flag")

	level1Cmd := &cobra.Command{
		Use: "level1",
		RunE: func(_ *cobra.Command, _ []string) error {
			executedCommand = "level1"
			capturedSettings = v.AllSettings()
			return nil
		},
	}
	level1Cmd.PersistentFlags().String("level1-flag", "level1-default", "Level1 flag")

	level2Cmd := &cobra.Command{
		Use: "level2",
		RunE: func(_ *cobra.Command, _ []string) error {
			executedCommand = "level2"
			capturedSettings = v.AllSettings()
			return nil
		},
	}
	level2Cmd.Flags().String("level2-flag", "level2-default", "Level2 flag")

	level1Cmd.AddCommand(level2Cmd)
	rootCmd.AddCommand(level1Cmd)

	t.Run("executing intermediate command binds its flags", func(t *testing.T) {
		v = viper.New()
		capturedSettings = nil

		rootCmd.SetArgs([]string{"level1", "--root-flag", "r", "--level1-flag", "l1"})
		err := rootCmd.ExecuteContext(context.Background())
		require.NoError(t, err)

		assert.Equal(t, "level1", executedCommand)
		assert.Equal(t, "r", capturedSettings["root_flag"])
		assert.Equal(t, "l1", capturedSettings["level1_flag"])
		_, hasLevel2 := capturedSettings["level2_flag"]
		assert.False(t, hasLevel2, "level2 flag should not be bound when executing level1")
	})
}
