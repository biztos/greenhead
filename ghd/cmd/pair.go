package cmd

import (
	"github.com/spf13/cobra"

	"github.com/biztos/greenhead/ghd/runner"
)

// PairCmd represents the "pair" command set.
var PairCmd = &cobra.Command{
	Use:   "pair run 'hello world'",
	Short: "Run pairs of agents.",
	Long: `The pair commands run pairs of agents.

See the subcommands for details.`,
}

var pairDumpJson bool

// PairRunCmd represents the "pair run" subcommand.
var PairRunCmd = &cobra.Command{
	Use:   "run $PROMPT",
	Short: "Run the two configured agents in conversation.",
	Long: `Runs completions between two configured agents in order.

The prompt is sent to the first agent; then the first agent's response
becomes the prompt for the second agent.  The second agent's response is
the prompt for the first againt again, and so on.

The conversation will continue until the maximum number of allowed
completions is reached, or an error occurs.

From each agent's point of view, the other agent is the "user".

If the prompt begins with '@' then it will be read from a file, e.g. @foo.txt.

Note that output and limit controls are very important when running pairs:
if --max-completions or --max-toolchain is set to a high number (or zero) the
run can be very costly.  Likewise, if output is not monitored it may be lost.
`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		r, err := runner.NewRunner(Config)
		if err != nil {
			return err
		}
		return r.RunPair(args[0], Stdout)
	},
}

func init() {

	PairCmd.AddCommand(PairRunCmd)
	RootCmd.AddCommand(PairCmd)
}
