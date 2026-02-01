package viperext_test

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/c4t-but-s4d/fastad/pkg/viperext"
)

func TestRunBindCommandFlags(t *testing.T) {
	t.Run("binds regular flags", func(t *testing.T) {
		v := viper.New()
		cmd := &cobra.Command{Use: "test"}
		cmd.Flags().String("my-flag", "default", "test flag")

		err := cmd.ParseFlags([]string{"--my-flag", "custom-value"})
		require.NoError(t, err)

		viperext.RunBindCommandFlags(v, cmd)

		assert.Equal(t, "custom-value", v.GetString("my_flag"))
	})

	t.Run("binds persistent flags", func(t *testing.T) {
		v := viper.New()
		cmd := &cobra.Command{Use: "test"}
		cmd.PersistentFlags().String("persistent-flag", "default", "test flag")

		err := cmd.ParseFlags([]string{"--persistent-flag", "persistent-value"})
		require.NoError(t, err)

		viperext.RunBindCommandFlags(v, cmd)

		assert.Equal(t, "persistent-value", v.GetString("persistent_flag"))
	})

	t.Run("normalizes flag names with dashes to underscores", func(t *testing.T) {
		v := viper.New()
		cmd := &cobra.Command{Use: "test"}
		cmd.Flags().String("kebab-case-flag", "default", "test flag")

		err := cmd.ParseFlags([]string{"--kebab-case-flag", "value"})
		require.NoError(t, err)

		viperext.RunBindCommandFlags(v, cmd)

		assert.Equal(t, "value", v.GetString("kebab_case_flag"))
	})

	t.Run("binds multiple flags", func(t *testing.T) {
		v := viper.New()
		cmd := &cobra.Command{Use: "test"}
		cmd.Flags().String("flag-one", "", "first flag")
		cmd.Flags().Int("flag-two", 0, "second flag")
		cmd.Flags().Bool("flag-three", false, "third flag")

		err := cmd.ParseFlags([]string{
			"--flag-one", "value1",
			"--flag-two", "42",
			"--flag-three",
		})
		require.NoError(t, err)

		viperext.RunBindCommandFlags(v, cmd)

		assert.Equal(t, "value1", v.GetString("flag_one"))
		assert.Equal(t, 42, v.GetInt("flag_two"))
		assert.True(t, v.GetBool("flag_three"))
	})

	t.Run("uses default values when flag not provided", func(t *testing.T) {
		v := viper.New()
		cmd := &cobra.Command{Use: "test"}
		cmd.Flags().String("my-flag", "default-value", "test flag")

		viperext.RunBindCommandFlags(v, cmd)

		assert.Equal(t, "default-value", v.GetString("my_flag"))
	})
}

func TestNewViper(t *testing.T) {
	t.Run("creates viper with env prefix", func(t *testing.T) {
		t.Setenv("TEST_PREFIX_MY_VAR", "env-value")

		v, err := viperext.NewViper("TEST_PREFIX")
		require.NoError(t, err)

		assert.Equal(t, "env-value", v.GetString("my_var"))
	})

	t.Run("replaces dashes with underscores in env names", func(t *testing.T) {
		t.Setenv("TEST_PREFIX_MY_DASHED_VAR", "dashed-value")

		v, err := viperext.NewViper("TEST_PREFIX")
		require.NoError(t, err)

		assert.Equal(t, "dashed-value", v.GetString("my-dashed-var"))
	})

	t.Run("replaces dots with underscores in env names", func(t *testing.T) {
		t.Setenv("TEST_PREFIX_MY_DOTTED_VAR", "dotted-value")

		v, err := viperext.NewViper("TEST_PREFIX")
		require.NoError(t, err)

		assert.Equal(t, "dotted-value", v.GetString("my.dotted.var"))
	})
}
