// cmd/agents.go

package cmd

import (
	"github.com/spf13/cobra"

	"github.com/biztos/greenhead/ghd/runner"
)

// AgentsCmd represents the "agents" command set.
var AgentsCmd = &cobra.Command{
	Use:   "agents [list|check|run]",
	Short: "Work with agents.",
	Long: `The agents commands help manage the configured agents.

Additional functionality will be added later (one hopes).`,
}

// AgentsListCmd represents the "agents list" subcommand.
var AgentsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all configured agents.",
	Long: `Shows the basic information of all configured agents.

Agents are instantiated to check for configuration errors.

Note that each will have a unique identifier (a ULID) when running.

If no agents are configured, shows the available named agents, if any.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		r, err := runner.NewRunner(Config)
		if err != nil {
			return err
		}
		return r.ListAgents(Stdout)
	},
}

// TODO: decide whether an AgentsShowCmd as in the tools one makes any sense.
// It would be nice to be able to confirm that an agent has the tools you
// expect, for example.  Anything else we would want?
// What about just doing "agents list --tools" and it lists all tools below
// the other agent info?

// AgentsCheckCmd represents the "agents check" subcommand.
var AgentsCheckCmd = &cobra.Command{
	Use:   "check",
	Short: "Check the configured agents.",
	Long: `Runs each configured agent's Check function and fails on error.

At least one agent must be configured.

Output is logged, and on success only the message "OK" is printed.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		r, err := runner.NewRunner(Config)
		if err != nil {
			return err
		}
		return r.CheckAgents(Stdout)

	},
}

var agentsRunJson = false

// AgentsRunCmd represents the "agents run" subcommand.
var AgentsRunCmd = &cobra.Command{
	Use:   "run [--json] $MY_PROMPT",
	Short: "Run a completion with the configured agents.",
	Long: `Runs a single completion with each configured agent, sequentially.

The completion may include tool calls, which will be executed.

If --json is specified, the output will be in JSON format.

If the prompt begins with '@' then it will be read from a file, e.g. @foo.txt.

In order to pipe output, e.g. to jq, it is advisable to use the --silent flag
and not use the --stream flag.  This is the recommended way to capture output
for multiple agents in a single run.
`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		r, err := runner.NewRunner(Config)
		if err != nil {
			return err
		}
		return r.RunAgents(Stdout, args[0], agentsRunJson)
	},
}

// AgentsColorCmd represents the "agents color" subcommand.
var AgentsColorCmd = &cobra.Command{
	Use:   "color [NAME...]",
	Short: "Print configured agent (and named) colors.",
	Long: `For any configured agents, and any names passed as arguments,
prints a short summary in that color and its pair color.

For named colors, a foreground/background pair can be specified as "fg/bg":

	yellow/white

`,
	RunE: func(cmd *cobra.Command, args []string) error {

		r, err := runner.NewRunner(Config)
		if err != nil {
			return err
		}
		return r.PrintColors(Stdout, args...)
	},
}

func init() {

	// Flags:
	AgentsRunCmd.Flags().BoolVar(&agentsRunJson, "json", false,
		"Print the whole completion as JSON.")

	// Registration:
	AgentsCmd.AddCommand(AgentsListCmd)
	AgentsCmd.AddCommand(AgentsCheckCmd)
	AgentsCmd.AddCommand(AgentsRunCmd)
	AgentsCmd.AddCommand(AgentsColorCmd)
	RootCmd.AddCommand(AgentsCmd)
}
