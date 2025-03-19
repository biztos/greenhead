// Package cmd holds the Cobra-based command definitions.
//
// It is designed to be easily customizable.
package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"

	"github.com/biztos/greenhead/config"
	"github.com/biztos/greenhead/registry"
)

var Version = "0.1.0"
var Name = "ghd"
var Title = "Greenhead Agent Runner"
var Description = `Runs AI (LLM) agents.

More soon!
`

var Config = &config.RunConfig{}

var Exit = os.Exit // for testability
var Stdout io.Writer = os.Stdout
var Stderr io.Writer = os.Stderr

// RootCmd is the Cobra Root and may be changed to suit after initialization.
var RootCmd = &cobra.Command{
	Use:     Name,        // Only the first word of "Use" applies to root.
	Version: Version,     // Nb: Version is *just* the number e.g. "1.2.3".
	Long:    Description, // Short is ignored in root.
}

func Execute() {
	RootCmd.SetOut(Stdout)
	RootCmd.SetErr(Stderr)
	// TODO: probably put this in the runners?
	// We want to enable stuff out of the config files but not after that.
	// Anyway this lets us stop a rogue tool from registering new tools, and
	// that's enough for now until we get the "factory" worked out.
	registry.Lock()
	err := RootCmd.Execute()
	if err != nil {
		Exit(1) // error message is printed above already.
	}
}

// BailErr bails out with the given error code and error message to Stderr.
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

	// The persistent flags cover the main things in Config.
	RootCmd.PersistentFlags().BoolVarP(&Config.Debug, "debug", "d", false,
		"Log at DEBUG level (default is INFO).")
	RootCmd.PersistentFlags().BoolVarP(&Config.Stream, "stream", "s", false,
		"Stream LLM output where applicable.")
	RootCmd.PersistentFlags().BoolVar(&Config.Silent, "silent", false,
		"Stream LLM output where applicable.")
	RootCmd.PersistentFlags().StringVar(&Config.DumpDir, "dump-dir", "",
		"Dump all LLM interactions into subdirectories of this dir.")

}
