package cmd

import (
	"fmt"
	"io"

	"github.com/spf13/cobra"

	"github.com/biztos/greenhead/agent"
	"github.com/biztos/greenhead/runner"
	"github.com/biztos/greenhead/utils"
)

const ExitCodeToolsError = 2

var toolsListNames bool

// ToolsCmd represents the "tools" command set.
var ToolsCmd = &cobra.Command{
	Use:   "tools [list|show NAME]",
	Short: "Show information about available tools.",
	Long: `The tools command helps manage registered tools.

Note that all tools must be specifically enabled by name for each agent in the
agent's configuration, and that some tools may be registered at runtime from
the configuration itself.

It is even possible for a tool to register more tools, though this must be
explicitly enabled.`,
}

// ToolsListCmd represents the "tools list" subcommand.
var ToolsListCmd = &cobra.Command{
	Use:   "list [--names]",
	Short: "List all registered tools which can be enabled for agents.",
	Long: `Lists all the registered tools, optionally only listing names.

For important caveats, see the parent command's help text.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		r, err := runner.NewRunner(Config)
		if err != nil {
			return err
		}
		return r.ListTools(toolsListNames, Stdout)
	},
}

// ToolsShowCmd represents the "tools show NAME" subcommand.
var ToolsShowCmd = &cobra.Command{
	Use:   "show NAME",
	Short: "Show detailed tool information.",
	Long: `Shows whatever detailed information is available for the tool.

This is defined by the tool itself, but generally should include the name,
description and input schema.

For important caveats, see the parent command's help text.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := runner.ShowTool(args[0], Stdout); err != nil {
			BailErr(ExitCodeToolsError, err)
		}
	},
}

// ToolsRunCmd represents the "tools run NAME" subcommand.
var ToolsRunCmd = &cobra.Command{
	Use:   "run NAME",
	Short: "Run the named command.",
	Long: `Runs the named command and returns the output as JSON.

Input should be provided on STDIN and *must* be valid JSON conforming to the
tool's input schema.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		data, err := io.ReadAll(cmd.InOrStdin())
		if err != nil {
			fmt.Fprintf(cmd.ErrOrStderr(), "Error reading input: %v\n", err)
			return
		}
		output, err := runner.RunTool(args[0], string(data))
		if err != nil {
			BailErr(ExitCodeToolsError, err)
		}
		// We don't know what the output type is but it should be safe to:
		res := &agent.ToolResult{
			Id:     "N/A",
			Output: output,
		}
		fmt.Println(utils.MustJsonStringPretty(res))

	},
}

func init() {
	// Flags:
	ToolsListCmd.Flags().BoolVar(&toolsListNames, "names", false,
		"Show only tool names.")

	// Registration:
	ToolsCmd.AddCommand(ToolsListCmd)
	ToolsCmd.AddCommand(ToolsShowCmd)
	ToolsCmd.AddCommand(ToolsRunCmd)
	RootCmd.AddCommand(ToolsCmd)
}
