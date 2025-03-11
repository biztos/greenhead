// Package greenhead runs AI agents.

package greenhead

import (
	_ "github.com/biztos/greenhead/registry"
	_ "github.com/biztos/greenhead/tools/demo"

	"github.com/biztos/greenhead/runner"
)

// Run calls the runner submodule's Run function.
func Run() {
	runner.Run()
}
