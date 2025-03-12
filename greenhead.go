// Package greenhead creates and runs AI agents.
//
// This top-level package makes it easy to create your own agents using the
// default code for everything but your specific custom functions.

package greenhead

import (
	"context"

	"github.com/biztos/greenhead/registry"
	"github.com/biztos/greenhead/runner"
	"github.com/biztos/greenhead/tools"
)

// Run exposes the runner submodule's Run function.
var Run = runner.Run

// Register exposes the register submodule's Register function.
var Register = registry.Register

// NewTool exposes the tools submodule's NewTool function.
//
// (We can't use a var for this because of the "generic" type.)
func NewTool[T any](name, desc string, f func(context.Context, T) (any, error)) *tools.Tool[T] {
	return tools.NewTool[T](name, desc, f)
}

// CustomApp sets values in runner for very basic customization of app name,
// title and version.
//
// This should be used when running the standard CLI from the runner submodule
// in cases where the custom tools warrant a different application name.
//
// Note that *only* the strings are affected; no functionality changes.
//
// (Also note that elaborate strings may break the runtime. Keep it simple.)
func CustomApp(name, title, version string) {
	runner.Name = name
	runner.Title = title
	runner.Version = version
}
