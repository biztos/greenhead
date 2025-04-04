package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	_ "github.com/biztos/greenhead/tools/tictactoe"
)

const ExitCodeWTF = 999

// WtfCmd is for work in progress.
var WtfCmd = &cobra.Command{
	Hidden: true,
	Use:    "wtf",
	Short:  "WTF? WTF!",
	Long: `What.
The.
(Actual.)
F**K?!

Work-in-progress random experiments.  Use at your own risk!

TODO: equivalent of this, but in an external package.  Hard? Easy?

Maybe just hook into pre/postrun funcs?`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return fmt.Errorf("nothing doing here ATM, WTF")

	},
}

func init() {
	RootCmd.AddCommand(WtfCmd)
}
