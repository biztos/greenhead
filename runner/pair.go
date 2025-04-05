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

// RunPair runs a Pair of Agents from the Runner.
//
// On successful completion, "<DONE>" will be printed to w.
func (r *Runner) RunPair(prompt string, w io.Writer) error {

	if len(r.Agents) != 2 {
		return fmt.Errorf("exactly two agents required; got %d", len(r.Agents))
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

	pair := agent.NewPair(r.Agents[0], r.Agents[1], r.Config.MaxCompletions)
	err := pair.Run(context.Background(), prompt)
	if !errors.Is(err, agent.ErrMaxCompletions) {
		return err
	}

	fmt.Fprintln(w, "<DONE>")
	return nil

}
