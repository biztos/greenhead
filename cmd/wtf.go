package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/biztos/greenhead/runner"
)

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

		fmt.Println("EHLO WTF")

		return nil
	},
}

func init() {

	RootCmd.AddCommand(WtfCmd)

	// well at least it's clear how we test this...
	runner.NamedAgentConfigs["chatty"].Color = "red" // e.g. "break me"

}
