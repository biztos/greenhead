// Package greenhead creates and runs AI agents.
//
// This top-level package makes it easy to create your own agents using the
// default code for everything but your specific custom functions.

package greenhead

import (
	"context"

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
	cmd.UpdateInfo()

}
