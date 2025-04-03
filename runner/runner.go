// Package runner defines the logic to run ghd, the Greenhead CLI tool.
package runner

import (
	"fmt"

	"github.com/biztos/greenhead/agent"
	"github.com/biztos/greenhead/registry"
)

// Runner is the runner of commands.
type Runner struct {
	Config *Config
	Agents []*agent.Agent
}

// NewRunner returns a new runner with the configuration processed.
//
// It is a thin wrapper around SetupTools and CreateAgents.
func NewRunner(cfg *Config) (*Runner, error) {

	if err := SetupTools(cfg); err != nil {
		return nil, err
	}
	agents, err := CreateAgents(cfg)
	if err != nil {
		return nil, err
	}
	return &Runner{
		Config: cfg,
		Agents: agents,
	}, nil

}

// SetupTools creates, registers, and/or deregisters tools based on cfg.
//
// Note that there is no concept of "allow nothing" -- set NoTools to achieve
// that result.
func SetupTools(cfg *Config) error {

	// NoTools is the easiest thing!
	if cfg.NoTools {
		registry.Clear()
		return nil
	}

	// TODO: handle CustomTools before applying allow/remove logic.

	// Save mutexes if nothing to see here.
	if len(cfg.AllowTools) == 0 && len(cfg.RemoveTools) == 0 {
		return nil
	}

	// Get allow and remove lists.
	allow, err := registry.MatchingNames(cfg.AllowTools)
	if err != nil {
		return fmt.Errorf("error in allowed tools: %w", err)
	}
	remove, err := registry.MatchingNames(cfg.RemoveTools)
	if err != nil {
		return fmt.Errorf("error in remove tools: %w", err)
	}

	// No overlap allowed.
	allowed := map[string]bool{}
	for _, n := range allow {
		allowed[n] = true
	}
	for _, n := range remove {
		if allowed[n] {
			return fmt.Errorf("can not both allow and remove tool: %q", n)
		}
	}

	// Allow is actually removal-based.
	if len(allow) > 0 {
		remove = []string{}
		for _, n := range registry.Names() {
			if !allowed[n] {
				remove = append(remove, n)
			}
		}
	}

	for _, n := range remove {
		if err := registry.Remove(n); err != nil {
			return fmt.Errorf("error removing tools: %w", err)
		}
	}

	// TODO: config to leave it unlocked for managing custom runtime tools.
	// (need to define that better first)
	registry.Lock()

	return nil
}

// CreateAgents creates agents from cfg.
func CreateAgents(cfg *Config) ([]*agent.Agent, error) {
	agents := make([]*agent.Agent, 0, len(cfg.Agents))
	for i, c := range cfg.Agents {
		a, err := agent.NewAgent(c)
		if err != nil {
			return nil, fmt.Errorf("error creating agent %d: %w", i+1, err)
		}
		agents = append(agents, a)
	}
	return agents, nil
}
