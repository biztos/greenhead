package cmd

import (
	"context"
	"fmt"

	"github.com/biztos/greenhead/runner"
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
		// PAIR CHAT as in: ghd pair run
		r, err := runner.NewRunner(Config)
		if err != nil {
			return err
		}
		if len(r.Agents) != 2 {
			return fmt.Errorf("exactly two agents must be configured")
		}

		// We will have this in an arg... may be tricky to find a good one,
		// but anyway this is the starting prompt which will be FROM a1 TO a2,
		// such that actually a2 does the first completion.
		prompt := `Parse this URL please: https://www.google.com/foo?bar+baz

Then generate some other random URLs that could be confused for this one by
an uncareful user.
`
		first := r.Agents[0]
		second := r.Agents[1]

		// Just for shits and giggles let's try a few rounds and see what
		// happens.
		ctx := context.Background()
		for count := 0; count < 3; count++ {
			res, err := second.RunCompletion(ctx, prompt)
			if err != nil {
				return err
			}
			prompt = res.Content
			if prompt == "" {
				return fmt.Errorf("received no-content completion")
			}
			res, err = first.RunCompletion(ctx, prompt)
			if err != nil {
				return err
			}
			prompt = res.Content
			if prompt == "" {
				return fmt.Errorf("received no-content completion")
			}
		}
		return nil

	},
}

func init() {
	RootCmd.AddCommand(WtfCmd)
}
