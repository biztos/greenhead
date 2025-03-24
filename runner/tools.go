package runner

import (
	"context"
	"fmt"
	"io"

	"github.com/biztos/greenhead/registry"
)

// ListTools prints all registered tools to w.
func ListTools(names_only bool, w io.Writer) {
	var lines []string
	if names_only {
		lines = registry.Names()
	} else {
		lines = registry.DisplayLines()
	}
	if len(lines) == 0 {
		fmt.Fprintln(w, "<no tools>")
	}
	for _, line := range lines {
		fmt.Fprintln(w, line)
	}

}

// ShowTool prints tool help to w, or returns an error if no tool is
// registered for name.
func ShowTool(name string, w io.Writer) error {
	t, err := registry.Get(name)
	if err != nil {
		return err
	}
	fmt.Fprintln(w, t.Help())
	return nil
}

// RunTool runs a tool with args as a string to be converted to the input
// type of the tool.
func RunTool(name, args string) (any, error) {
	t, err := registry.Get(name)
	if err != nil {
		return nil, err
	}

	// We do not JSON-ify the result here, the caller can deal with that.
	output, err := t.Exec(context.Background(), args)
	if err != nil {
		return nil, err
	}
	return output, nil

}
