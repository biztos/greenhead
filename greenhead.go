// Package greenhead creates and runs AI agents.
//
// This top-level package makes it easy to create your own agents using the
// default code for everything but your specific custom functions.

package greenhead

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/biztos/greenhead/cmd"
	"github.com/biztos/greenhead/registry"
	"github.com/biztos/greenhead/tools"
)

// Run exposes the cmd submodule's Execute function.
var Run = cmd.Execute

// Register exposes the register submodule's Register function.
var Register = registry.Register

// Clear exposes the register submodule's Clear function.
var Clear = registry.Clear

// NewTool exposes the tools submodule's NewTool function.
//
// (We can't use a var for this because of the "generic" type.)
func NewTool[T any, R any](name, desc string, f func(context.Context, T) (R, error)) *tools.Tool[T, R] {
	return tools.NewTool[T, R](name, desc, f)
}

// CustomApp sets custom values for help messages in cmd.
func CustomApp(name, version, title, description string) {
	if name != "" {
		cmd.Name = name
	}
	if version != "" {
		cmd.Version = version
	}
	if title != "" {
		cmd.Title = title
	}
	if description != "" {
		cmd.Description = description
	}

}

var ErrCommandNotFound = fmt.Errorf("command not found in the root")

// AddCommand adds command c to the Cobra root command.
//
// This is useful for building custom binaries with default subcommands and
// also custom commands.
func AddCommand(c *cobra.Command) {
	cmd.RootCmd.AddCommand(c)
}

// RemoveCommand removes a command from the Cobra root command.
//
// This is useful for building custom binaries with *less* functionality than
// the default.
//
// Panics if the command is not found.
func RemoveCommand(name string) {
	// Create a new slice without the command to be removed
	new_commands := []*cobra.Command{}

	old_commands := cmd.RootCmd.Commands()
	for _, c := range old_commands {
		if c.Name() != name {
			new_commands = append(new_commands, c)
		}
	}
	if len(old_commands) == len(new_commands) {
		panic("no such command: " + name)
	}

	// Clear all commands
	cmd.RootCmd.ResetCommands()

	// Add back all commands except the one to be removed
	for _, c := range new_commands {
		cmd.RootCmd.AddCommand(c)
	}
}

// ResetCommands resets the Cobra root command subcommands.
func ResetCommands() {
	cmd.RootCmd.ResetCommands()
}

// ResetFlags resets the Cobra root command flags, disabling all persistent
// flags.
//
// This is useful for building binaries that rely on a fixed configuration.
func ResetFlags() {
	cmd.RootCmd.ResetFlags()
}

// RemoveFlag removes a flag from the Cobra root command.
//
// This is useful for building binaries that have fewer options; however, take
// care that the options are not simply defined in a config file.
//
// Panics if the flag does not exist.
func RemoveFlag(name string) {

	have := false
	new_flags := pflag.NewFlagSet(cmd.RootCmd.Name(), pflag.ContinueOnError)
	cmd.RootCmd.PersistentFlags().VisitAll(func(flag *pflag.Flag) {
		if flag.Name == name {
			have = true
		} else {
			new_flags.AddFlag(flag)
		}
	})

	// Replace the command's persistent flags
	cmd.RootCmd.ResetFlags()
	cmd.RootCmd.PersistentFlags().AddFlagSet(new_flags)

	if !have {
		panic("no such flag: " + name)
	}
}
