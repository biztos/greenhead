// runner/agents.go

package runner

import (
	"context"
	"fmt"
	"io"
	"strconv"

	"github.com/biztos/greenhead/utils"
)

// RunListAgents prints summary information for the configured agents.
//
// Note that the runner's Config is used here, not the instantiated agents.
// In normal operation these will be the same.
func (r *Runner) RunListAgents(w io.Writer) {
	if len(r.Config.Agents) == 0 {
		fmt.Fprintln(w, "<no agents>")
		return
	}

	// Get our type and model widths from the config.
	tw := 0
	mw := 0
	for _, a := range r.Config.Agents {
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
	for i, a := range r.Config.Agents {
		fmt.Fprintf(w, linef, strconv.Itoa(i), a.Type, a.Model, a.Name)
	}
}

// RunCheckAgents runs the ApiClient Check command on configured agents, in
// order.
//
// Agents must be instantiated, and at least one agent must be present.
func (r *Runner) RunCheckAgents(w io.Writer) error {
	if len(r.Agents) == 0 {
		return fmt.Errorf("no agents")
	}
	for i, a := range r.Agents {
		if err := a.Check(context.Background()); err != nil {
			return fmt.Errorf("error checking agent %d: %w", i, err)
		}
	}
	fmt.Fprintln(w, "OK")
	return nil
}

// RunRunAgents runs the single-prompt completion on all agents, in order.
func (r *Runner) RunRunAgents(w io.Writer, prompt string, json bool) error {
	if len(r.Agents) == 0 {
		return fmt.Errorf("no agents")
	}
	for _, a := range r.Agents {
		if !json {
			fmt.Fprintln(w, a.Ident())
		}
		res, err := a.RunCompletion(context.Background(), prompt)
		if err != nil {
			return fmt.Errorf("error running completion: %w", err)
		}
		if json {
			v := map[string]any{
				"agent":    a.Ident(),
				"prompt":   prompt,
				"response": res.Content,
			}
			fmt.Fprintln(w, utils.MustJsonString(v))

		}

	}
	return nil
}
