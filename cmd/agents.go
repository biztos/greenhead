// cmd/agents.go

package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/biztos/greenhead/runner"
)

// AgentsCmd represents the "agents" command set.
var AgentsCmd = &cobra.Command{
	Use:   "agents [list|check]",
	Short: "Work with agents (without running them).",
	Long: `The agents command helps manage the configured agents.

Additional functionality will be added later (one hopes).`,
}

// AgentsListCmd represents the "agents list" subcommand.
var AgentsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all configured agents.",
	Long: `Shows the basic information of all configured agents.

Agents are instantiated to check for configuration errors.

Note that each will have a unique identifier (a ULID) when running.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runner.ListAgents(Config, os.Stdout)
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

Agents are instantiated to check for configuration errors.

Output is logged, and on success only the message "OK" is printed.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := runner.CheckAgents(Config, context.Background()); err != nil {
			return err
		}
		fmt.Fprintln(Stdout, "OK")
		return nil
	},
}

func init() {

	// Registration:
	AgentsCmd.AddCommand(AgentsListCmd)
	AgentsCmd.AddCommand(AgentsCheckCmd)
	RootCmd.AddCommand(AgentsCmd)
}
