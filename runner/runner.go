// Package runner defines the logic to run ghd, the Greenhead CLI tool.
package runner

import (
	"github.com/biztos/greenhead/agent"
)

// Runner is the runner of commands.
type Runner struct {
	Config *Config
	Agents []*agent.Agent
}

// NewRunner returns a new runner with the configuration processed.
func NewRunner(cfg *Config) (*Runner, error) {
	// TODO: this seems like the place to register tools from the config!
	agents, err := cfg.CreateAgents()
	if err != nil {
		return nil, err
	}
	return &Runner{
		Config: cfg,
		Agents: agents,
	}, nil

}
