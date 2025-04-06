package runner

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/biztos/greenhead/agent"
)

var ErrNotPair = fmt.Errorf("exactly two agents are required for a pair")

// RunPair runs a Pair of Agents from the Runner.
//
// On successful completion, "<DONE>" will be printed to w.
func (r *Runner) RunPair(prompt string, w io.Writer) error {

	if len(r.Agents) != 2 {
		return fmt.Errorf("%w: got %d", ErrNotPair, len(r.Agents))
	}

	// TODO: centralize this, have twice already!
	if strings.HasPrefix(prompt, "@") {
		file := strings.TrimPrefix(prompt, "@")
		b, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("error reading prompt file: %w", err)
		}
		prompt = string(b)
	}

	// If we're printing output to the console we strongly prefer to show the
	// two agents in different colors, it's just way more useful that way.
	// TODO: have some other differentiator if you can't use/see/want colors!
	// And then have a no-color option maybe to force that to be respected?
	if !r.Config.Silent {
		r.ColorizePair()
	}

	first := r.Agents[0]
	second := r.Agents[1]

	pair := agent.NewPair(first, second, r.Config.MaxCompletions)
	err := pair.Run(context.Background(), prompt)
	if !errors.Is(err, agent.ErrMaxCompletions) {
		return err
	}

	fmt.Fprintln(w, "<DONE>")
	return nil

}

// ColorizePair sets the colors for the two agents if needed.
func (r *Runner) ColorizePair() error {

	if len(r.Agents) != 2 {
		return fmt.Errorf("%w: got %d", ErrNotPair, len(r.Agents))
	}

	// TODO: rethink allowing agent config access, annoying to not have.
	fg1 := r.Config.Agents[0].Color
	fg2 := r.Config.Agents[1].Color
	bg1 := r.Config.Agents[0].BgColor
	bg2 := r.Config.Agents[1].BgColor
	if fg1 == fg2 && bg1 == bg2 {
		// Force a difference.  Start with a decent fg/bg scheme if not
		// set.
		if fg1 == "" && bg1 == "" {
			fg1 = "black"
			bg1 = "cornsilk"
		}
		fg2, _ = agent.FindComplementaryColor(fg1) // known good, TODO: really?
		bg2, _ = agent.FindComplementaryColor(bg1)

		c1, _ := agent.PrintColor(fg1, bg1)
		c2, _ := agent.PrintColor(fg2, bg2)
		r.Agents[0].SetPrintFunc(c1.PrintFunc())
		r.Agents[1].SetPrintFunc(c2.PrintFunc())

	}
	return nil
}
