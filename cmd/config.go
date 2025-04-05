package cmd

import (
	"github.com/spf13/cobra"

	"github.com/biztos/greenhead/runner"
)

// ConfigCmd represents the "config" command set.
var ConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage configs.",
	Long: `The config commands help manage configuration files.

The preferred format for config files is TOML, but JSON and JSON5 are also
supported.`,
}

var configDumpJson bool

// ConfigDumpCmd represents the "config dump" subcommand.
var ConfigDumpCmd = &cobra.Command{
	Use:   "dump [--json]",
	Short: "Dump the currently loaded configuration.",
	Long: `Reads and validates the specified runner and agent config files.

If no errors are found, dumps the resulting unified config which
can be used as a fully equivalent single runner config.

Outputs TOML by default, but with the --json argument JSON will be output.

Comments are *not* preserved at this time.
`,
	Args: cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		r, err := runner.NewRunner(Config)
		if err != nil {
			return err
		}
		if configDumpJson {
			r.Config.DumpJson(Stdout)
		} else {
			r.Config.DumpToml(Stdout)
		}
		return nil
	},
}

// ConfigCheckCmd represents the "config check" subcommand.
var ConfigCheckCmd = &cobra.Command{
	Use:   "check",
	Short: "Check the currently loaded config and exit.",
	Long: `Reads and validates the specified runner and agent config files.

Succeeds silently if there are no errors.`,
	Args: cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		_, err := runner.NewRunner(Config)
		return err
	},
}

func init() {
	ConfigDumpCmd.Flags().BoolVar(&configDumpJson, "json", false,
		"Dump JSON instead of TOML")

	ConfigCmd.AddCommand(ConfigDumpCmd)
	ConfigCmd.AddCommand(ConfigCheckCmd)
	RootCmd.AddCommand(ConfigCmd)
}
