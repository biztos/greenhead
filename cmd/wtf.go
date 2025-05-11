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

		// OK, we start with a standard runner.
		r, err := runner.NewRunner(Config)
		if err != nil {
			return err
		}
		// Now we have whatever agents.
		fmt.Println(len(r.Agents), "agents in the runner")
		for _, a := range r.Agents {
			fmt.Println("loaded", a.ULID, a.Name)
			c, err := a.Clone()
			if err != nil {
				return err
			}
			fmt.Println("cloned", c.ULID, c.Name)
			c.Logger().Info("george cloney", "name", c.Name)
			c.Print(c)
		}

		return nil
	},
}

func init() {

	RootCmd.AddCommand(WtfCmd)

	// well at least it's clear how we test this...
	runner.NamedAgentConfigs["chatty"].Color = "red" // e.g. "break me"

}
