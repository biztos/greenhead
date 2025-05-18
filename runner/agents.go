// runner/agents.go

package runner

import (
	"context"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/biztos/greenhead/agent"
	"github.com/biztos/greenhead/utils"
)

// ListAgents prints summary information for the configured agents; or if
// no agents are configured, the available named agents.
func (r *Runner) ListAgents(w io.Writer) error {

	var agents []*agent.Agent
	if len(r.Agents) == 0 {
		if len(NamedAgentConfigs) == 0 {
			fmt.Fprintln(w, "<no agents>")
			return nil
		}
		fmt.Fprintln(w, "No agents loaded; showing available named agents.")
		names := []string{}
		for name := range NamedAgentConfigs {
			names = append(names, name)
		}
		sort.Strings(names)
		for _, name := range names {
			agent, err := agent.NewAgent(NamedAgentConfigs[name])
			if err != nil {
				// TODO: trigger bad config to prove this works.
				return fmt.Errorf("bad named config for %s: %w", name, err)
			}
			agents = append(agents, agent)
		}
	} else {
		agents = r.Agents
	}

	// Get our type, model, name widths from the config.
	tw := 4
	mw := 5
	nw := 4
	for _, a := range agents {
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
	for i, a := range agents {
		// Take only first line of desc.
		// TODO (maybe): JSON option, which could include everything including
		// potentially tools?
		// TODO: trim to the terminal width!  (which we now have in utils)
		desc := strings.TrimSpace(strings.SplitN(a.Description, "\n", 2)[0])
		if desc != a.Description {
			desc += "..."
		}
		fmt.Fprintf(w, linef, strconv.Itoa(i), a.Type, a.Model, a.Name, desc)
	}

	return nil
}

// CheckAgents runs the ApiClient Check command on configured agents, in
// order.
//
// Agents must be instantiated, and at least one agent must be present.
func (r *Runner) CheckAgents(w io.Writer) error {
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

// RunAgents runs the single-prompt completion on all agents, in order.
//
// If prompt starts with @ then it is read from a file, e.g. `@file.txt`.
func (r *Runner) RunAgents(w io.Writer, prompt string, json bool) error {
	if len(r.Agents) == 0 {
		return fmt.Errorf("no agents")
	}
	if strings.HasPrefix(prompt, "@") {
		file := strings.TrimPrefix(prompt, "@")
		b, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("error reading prompt file: %w", err)
		}
		prompt = string(b)
	}

	for _, a := range r.Agents {
		if !json {
			a.Print(a.Ident() + "\n")
		}
		res_content, err := a.RunCompletionPrompt(prompt)
		if err != nil {
			return fmt.Errorf("error running completion: %w", err)
		}
		if json {
			v := map[string]any{
				"agent":    a.Ident(),
				"prompt":   prompt,
				"response": res_content,
			}
			fmt.Fprintln(w, utils.MustJsonString(v))

		}

	}
	return nil
}

// PrintColors prints the colors and their pair colors for all Agents, and
// also for any extra arg colors passed in.
func (r *Runner) PrintColors(w io.Writer, args ...string) error {

	for i, c := range r.Config.Agents {

		prefix := fmt.Sprintf("Agent %d: %s - ", i, c.Name)
		agent.PrintColorPairSample(w, c.Color, c.BgColor, prefix)

	}
	for i, n := range args {
		fg := n
		bg := ""
		parts := strings.SplitN(n, "/", 2)
		if len(parts) > 1 {
			fg = parts[0]
			bg = parts[1]
		}
		prefix := fmt.Sprintf("Arg %d: ", i+1)
		if err := agent.PrintColorPairSample(w, fg, bg, prefix); err != nil {
			return err
		}
	}

	return nil
}
