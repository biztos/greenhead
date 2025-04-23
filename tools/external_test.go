package tools_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/biztos/greenhead/tools"
)

func TestExternalToolConfigCopyOK(t *testing.T) {
	require := require.New(t)
	c := &tools.ExternalToolConfig{
		Name:        "list_files",
		Description: "The UNIX `ls` command, used to list files and directories.",
		Command:     "/bin/ls",
		Options: []tools.ExternalToolOption{
			{"-l", "long", "List files in the long format."},
			{"-A", "all", "List all files, including dotfiles."},
		},
		Args: []string{"[FILE...]"},
	}
	n := c.Copy()
	require.EqualValues(n, c)
	c.Name = "changed"
	c.Options = append(c.Options, tools.ExternalToolOption{})
	c.Args = append(c.Args, "more")
	require.Equal("list_files", n.Name, "Name not linked")
	require.Equal(2, len(n.Options), "Options not linked")
	require.Equal([]string{"[FILE...]"}, n.Args, "Args not linked")

}
