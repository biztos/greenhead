// Package cmd holds the Cobra-based command definitions.
//
// It is designed to be easily customizable.
package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"

	"github.com/biztos/greenhead/rgxp"
	"github.com/biztos/greenhead/runner"
)

var Version = "0.1.0"
var Name = "ghd"
var Title = "Greenhead Agent Runner"
var Description = `Runs AI (LLM) agents.

More soon!

Special flag types:

*Array:       flag allowing multiple uses, e.g. --foo=A --foo=B

rgxpArray:    regular expressions in the format /pattern/[ism]

optRgxpArray: as rgxpArray, but plain strings (not in the above format) can
              be used for exact matches.
`

var Config = &runner.Config{}

var Exit = os.Exit // for testability
var Stdout io.Writer = os.Stdout
var Stderr io.Writer = os.Stderr

var runnerConfigFile string
var agentConfigFiles []string

// RootCmd is the Cobra Root and may be changed to suit after initialization.
var RootCmd = &cobra.Command{
	Use:          Name,        // Only the first word of "Use" applies to root.
	Version:      Version,     // Nb: Version is *just* the number e.g. "1.2.3".
	Long:         Description, // Short is ignored in root.
	SilenceUsage: true,        // Do not spam usage all the time!
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {

		// We already have a config with flags; load any config files, letting
		// the flags override.
		if err := Config.LoadConfigs(runnerConfigFile, agentConfigFiles...); err != nil {
			return fmt.Errorf("error loading config: %w", err)
		}

		return nil
	},
}

func Execute() {
	RootCmd.SetOut(Stdout)
	RootCmd.SetErr(Stderr)
	err := RootCmd.Execute()
	if err != nil {
		Exit(1) // error message is printed above already.
	}
}

// BailErr bails out with the given error code and error message to Stderr.
//
// Use this when you need more control over the error behavior than RunE
// provides.
func BailErr(code int, err error) {
	fmt.Fprintln(Stderr, err.Error())
	Exit(code)
}

// UpdateInfo updates the help and usage info from the package variables.
func UpdateInfo() {
	RootCmd.Use = Name
	RootCmd.Version = Version
	RootCmd.Long = fmt.Sprintf("%s (%s)\n\n%s", Title, Name, Description)

}

func init() {

	UpdateInfo()

	// Config files:
	RootCmd.PersistentFlags().StringVar(&runnerConfigFile, "config", "",
		"Config file from which to read the master configuration.")
	RootCmd.PersistentFlags().StringArrayVar(&agentConfigFiles, "agent", []string{},
		"Config file from which to read the agent configuration.")

	// Output control:
	RootCmd.PersistentFlags().BoolVarP(&Config.Debug, "debug", "d", false,
		"Log at DEBUG level (default is INFO).")
	RootCmd.PersistentFlags().BoolVarP(&Config.Stream, "stream", "s", false,
		"Stream LLM output to the console.")
	RootCmd.PersistentFlags().BoolVar(&Config.Silent, "silent", false,
		"Suppress LLM output.")
	RootCmd.PersistentFlags().StringVar(&Config.LogFile, "log-file", "",
		"Log to this file instead of STDERR.")
	RootCmd.PersistentFlags().BoolVar(&Config.NoLog, "no-log", false,
		"Do not log at all.")
	RootCmd.PersistentFlags().StringVar(&Config.DumpDir, "dump-dir", "",
		"Dump all LLM interactions into this dir.")
	RootCmd.PersistentFlags().BoolVar(&Config.ShowCalls, "show-calls", false,
		"Show tool calls with output (experimental; can leak data!).")

	// Limits:
	RootCmd.PersistentFlags().IntVarP(&Config.MaxCompletions, "max-completions", "m", 100,
		"Maximum number of completions to run (tool calls not included).")
	RootCmd.PersistentFlags().IntVar(&Config.MaxToolChain, "max-toolchain", 10,
		"Maximum number of tool calls allowed in a completion.")

	// Safety:
	RootCmd.PersistentFlags().Var(&rgxp.RgxpArrayValue{Rgxps: &Config.StopMatches},
		"stop-match",
		"Stop running if any completion content matches.")

	// Tools:
	// Note: tool selection is a bit complicated and should be covered in the
	// main help text.  It's important that these default to nil, not an empty
	// array, as a config-file empty array might override.
	RootCmd.PersistentFlags().Var(&rgxp.OptionalRgxpArrayValue{OptionalRgxps: &Config.AllowTools},
		"allow-tool",
		"Allow tool use. See main help for details.")
	RootCmd.PersistentFlags().Var(&rgxp.OptionalRgxpArrayValue{OptionalRgxps: &Config.RemoveTools},
		"remove-tool",
		"Disallow tool use. See main help for details.")
	RootCmd.PersistentFlags().Var(&rgxp.OptionalRgxpArrayValue{OptionalRgxps: &Config.AgentTools},
		"agent-tool",
		"Override all agent tool lists.")
	RootCmd.PersistentFlags().BoolVar(&Config.NoTools, "no-tools", false,
		"Remove all tools before running (agents will have no tools).")

}
