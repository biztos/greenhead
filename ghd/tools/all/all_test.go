package all_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/biztos/greenhead/ghd/registry"
	_ "github.com/biztos/greenhead/ghd/tools/all"
)

// TODO (maybe): programmatically prove we have every tool.
//
//	(requires some kind of tool spec file, annoying?)
//
// TODO (maybe): prove every tool is correctly namespaced.
//
//	(same, tool -> pkg is not known)
//
// For now just keep this list current when adding any tools.
func TestHaveAllTools(t *testing.T) {
	require := require.New(t)

	exp := []string{
		"demo_store",
		"demo_recall",
		"demo_sum",
		"tictactoe_new_game",
		"tictactoe_state",
		"tictactoe_move",
	}
	require.Equal(exp, registry.Names())
}
