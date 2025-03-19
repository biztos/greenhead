// Package runner defines the logic to run ghd, the Greenhead CLI tool.
package runner

import (
	"context"
	"fmt"

	"github.com/biztos/greenhead/registry"
)

// ListTools prints all registered tools to standard output.
func ListTools() error {
	names := registry.Names()
	if len(names) == 0 {
		return fmt.Errorf("no tools registered")
	}
	for _, name := range names {
		fmt.Println(name)
	}
	return nil
}

// ShowTool shows tool help.
func ShowTool(name string) error {
	t := registry.Get(name)
	if t == nil {
		// TODO (arguably) exit nonzero here
		return fmt.Errorf("tool not registered: %s", name)
	}
	fmt.Println(t.Help())
	return nil
}

// RunTool runs a tool with args as a string to be converted to the input
// type of the tool.
func RunTool(name, args string) (any, error) {
	t := registry.Get(name)
	if t == nil {
		return nil, fmt.Errorf("tool not registered: %s", name)
	}

	// We do not JSON-ify the result here, the caller can deal with that.
	output, err := t.Exec(context.Background(), args)
	if err != nil {
		return nil, err
	}
	return output, nil

}
