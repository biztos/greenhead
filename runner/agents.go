// runner/agents.go

package runner

import (
	"context"
	"fmt"
	"io"
	"strconv"
)

// ListAgents instantiates configured agents and prints summary information.
func ListAgents(cfg *Config, w io.Writer) error {
	if len(cfg.Agents) == 0 {
		fmt.Fprintln(w, "<no agents>")
	}
	// We only instantiate in order to catch any errors; for listing we use
	// the config data.
	_, err := cfg.CreateAgents()
	if err != nil {
		return err
	}

	// Get our type and model widths.
	tw := 0
	mw := 0
	for _, a := range cfg.Agents {
		if len(a.Type) > tw {
			tw = len(a.Type)
		}
		if len(a.Model) > mw {
			mw = len(a.Model)
		}
	}
	tw += 2
	mw += 2
	linef := fmt.Sprintf("%%-4s %%-%ds %%-%ds %%s\n", tw, mw)

	fmt.Fprintf(w, linef, "No.", "Type", "Model", "Name")
	for i, a := range cfg.Agents {
		fmt.Fprintf(w, linef, strconv.Itoa(i), a.Type, a.Model, a.Name)
	}
	return nil
}

// CheckAgents instantiated configured agents and runs their ApiClient Check
// commands, thus (presumably) confirming that API keys are valid.
func CheckAgents(cfg *Config, ctx context.Context) error {
	if len(cfg.Agents) == 0 {
		return fmt.Errorf("no agents")
	}
	agents, err := cfg.CreateAgents()
	if err != nil {
		return err
	}
	for i, a := range agents {
		if err := a.Check(ctx); err != nil {
			return fmt.Errorf("error checking agent %d: %w", i, err)
		}
	}
	return nil
}
