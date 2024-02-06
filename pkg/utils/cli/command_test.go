package cli

import (
	"os"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testCmd struct {
	String      string   `name:"string" short:"s" usage:"string flag" env:"STRING" default:"string"`
	Int         int      `name:"int" short:"i" usage:"int flag" env:"INT" default:"1"`
	Bool        bool     `name:"bool" short:"b" usage:"bool flag" env:"BOOL" default:"true"`
	StringSlice []string `name:"string-slice" short:"S" usage:"string slice flag"`
}

func (c *testCmd) Run(cmd *cobra.Command, args []string) error {
	return nil
}

func TestCommand(t *testing.T) {
	var c testCmd
	cmd := Command(&c, &cobra.Command{
		Short: "test",
	})
	require.NoError(t, cmd.Execute())
	assert.Equal(t, "string", c.String)
	assert.Equal(t, 1, c.Int)
	assert.Equal(t, true, c.Bool)
	assert.Equal(t, []string{}, c.StringSlice)
}

func TestCommandEnv(t *testing.T) {
	require.NoError(t, os.Setenv("STRING", "env-string"))
	require.NoError(t, os.Setenv("INT", "2"))
	require.NoError(t, os.Setenv("BOOL", "false"))
	require.NoError(t, os.Setenv("STRING_SLICE", "env-string1,env-string2"))
	var c testCmd
	cmd := Command(&c, &cobra.Command{
		Short: "test",
	})
	require.NoError(t, cmd.Execute())
	assert.Equal(t, "env-string", c.String)
	assert.Equal(t, 2, c.Int)
	assert.Equal(t, false, c.Bool)
	assert.Equal(t, []string{}, c.StringSlice)
}
