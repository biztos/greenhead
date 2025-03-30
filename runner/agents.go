// runner/agents.go

package runner

import (
	"context"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/biztos/greenhead/utils"
)

// RunListAgents prints summary information for the configured agents.
func (r *Runner) RunListAgents(w io.Writer) {
	if len(r.Config.Agents) == 0 {
		fmt.Fprintln(w, "<no agents>")
		return
	}

	// Get our type, model, name widths from the config.
	tw := 4
	mw := 5
	nw := 4
	for _, a := range r.Agents {
		if len(a.Type) > tw {
			tw = len(a.Type)
		}
		if len(a.Model) > mw {
			mw = len(a.Model)
		}
		if len(a.Name) > nw {
			nw = len(a.Name)
		}
	}
	tw += 2
	mw += 2
	nw += 2
	linef := fmt.Sprintf("%%-4s %%-%ds %%-%ds %%-%ds %%s\n", tw, mw, nw)

	fmt.Fprintf(w, linef, "No.", "Type", "Model", "Name", "Description")
	for i, a := range r.Agents {
		// Take only first line of desc.
		// TODO (maybe): JSON option, which could include everything including
		// potentially tools?
		// TODO: trim to the terminal width!
		// 	if fd := int(os.Stdout.Fd()); term.IsTerminal(fd) {
		// if w, _, err := term.GetSize(fd); err == nil && w > 0 {
		// 	width = w
		// }
		desc := strings.TrimSpace(strings.SplitN(a.Description, "\n", 2)[0])
		if desc != a.Description {
			desc += "..."
		}
		fmt.Fprintf(w, linef, strconv.Itoa(i), a.Type, a.Model, a.Name, desc)
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
