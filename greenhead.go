// Package greenhead runs AI agents.

package greenhead

import (
	_ "github.com/biztos/greenhead/extensions/demo"
	_ "github.com/biztos/greenhead/registry"

	"github.com/biztos/greenhead/runner"
)

// Run calls the runner submodule's Run function.
func Run() {
	runner.Run()
}
