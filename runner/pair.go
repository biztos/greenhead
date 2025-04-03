package runner

import (
	"fmt"

	"github.com/biztos/greenhead/agent"
)

// RunPair sets up and runs a Pair of Agents from the Runner.
func (r *Runner) RunPair() error {

	if len(r.Agents) != 2 {
		return fmt.Errorf("exactly two agents required; got %d", len(r.Agents))
	}

	pair := agent.NewPair(r.Agents[0], r.Agents[1])
	if pair == nil {
		return nil
	}
	return nil

}
