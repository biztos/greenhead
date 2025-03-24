// Package runner defines the logic to run ghd, the Greenhead CLI tool.
package runner

import (
	"fmt"

	"github.com/biztos/greenhead/agent"
	"github.com/biztos/greenhead/utils"
)

// Config describes the configuration used to run the application.
//
// Further configuration is at the Agent level, and can be included here.
type Config struct {

	// Output control:
	Debug       bool   `toml:"debug"`         // Log at DEBUG level instead of INFO.
	LogFile     string `toml:"log_file"`      // Write logs to this file instead of os.StdErr.
	Silent      bool   `toml:"silent"`        // Suppress LLM output if not already streamed.
	Stream      bool   `toml:"stream"`        // Stream LLM output if supported.
	ShowCalls   bool   `toml:"stream_calls"`  // Stream (or print) tool calls (potentially leaking data).
	DumpDir     string `toml:"dump_dir"`      // Dump all completions to JSON files in this dir.
	LogToolArgs bool   `toml:"log_tool_args"` // Log the tool args (potentially leaking data).

	// Tools from which the agents may choose; any others are deregistered.
	Tools       []string `toml:"tools,omitempty"` // TBD, how to define them?
	CustomTools []string `toml:"custom_tools"`    // TBD, typed

	// Agent configs:
	Agents []*agent.Config `toml:"agents,omitempty"` // Multiple Agents.

}

// LoadConfigs loads a runner config, followed by any agent configs, adding
// any new nonzero values to c.  Agents are appended.  ConformAgents and
// Validate are called before returning.
//
// This is normally used when c holds flag values and config files are loaded
// before executing runner functions.
func (c *Config) LoadConfigs(runnerFile string, agentFiles ...string) error {

	if runnerFile != "" {
		r := &Config{}
		if err := utils.UnmarshalFile(runnerFile, r); err != nil {
			return err
		}
		// Bools and strings from the original take precedence.
		c.Debug = c.Debug || r.Debug
		c.Silent = c.Silent || r.Silent
		c.Stream = c.Stream || r.Stream
		c.ShowCalls = c.ShowCalls || r.ShowCalls
		c.LogToolArgs = c.LogToolArgs || r.LogToolArgs
		if c.LogFile == "" {
			c.LogFile = r.LogFile
		}
		if c.DumpDir == "" {
			c.DumpDir = r.DumpDir
		}

		// Tools and agents are added (in the flags case they would be empty).
		c.Tools = append(c.Tools, r.Tools...)
		c.Agents = append(c.Agents, r.Agents...)

	}

	for _, file := range agentFiles {
		a := &agent.Config{}
		if err := utils.UnmarshalFile(file, a); err != nil {
			return err
		}
		c.Agents = append(c.Agents, a)
	}

	c.ConformAgents()
	return c.Validate()

}

// Validate checks the config for internal consistency.
func (c *Config) Validate() error {
	if c.Stream && c.Silent {
		return fmt.Errorf("Stream and Silent can not both be enabled.")
	}
	// TODO: other things perhaps!
	return nil
}

// ConformAgents applies values from c to any configs in c.Agents, so that
// agents conform to the runner values that can be set with command-line
// flags.
//
// This is not strictly necessary, but one would expect havoc to ensue if the
// values differ.  If you find a compelling use-case for that, please open
// an issue.
func (c *Config) ConformAgents() {
	for _, a := range c.Agents {
		a.Stream = c.Stream
		a.Silent = c.Silent
		a.Debug = c.Debug
		a.ShowCalls = c.ShowCalls
		a.LogToolArgs = c.LogToolArgs
		a.LogFile = c.LogFile
		a.DumpDir = c.DumpDir
	}
}

// CreateAgents creates an Agent for each one configured.
func (c *Config) CreateAgents() ([]*agent.Agent, error) {
	agents := make([]*agent.Agent, 0, len(c.Agents))
	for i, cfg := range c.Agents {
		a, err := agent.NewAgent(cfg)
		if err != nil {
			return nil, fmt.Errorf("error creating agent %d: %w", i, err)
		}
		agents = append(agents, a)
	}
	return agents, nil

}

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
