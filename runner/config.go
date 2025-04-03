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
	Debug       bool   `toml:"debug"`            // Log at DEBUG level instead of INFO.
	LogFile     string `toml:"log_file"`         // Write logs to this file instead of os.StdErr.
	NoLog       bool   `toml:"no_log,omitempty"` // Do not log at all.
	Silent      bool   `toml:"silent"`           // Suppress LLM output if not already streamed.
	Stream      bool   `toml:"stream"`           // Stream LLM output if supported.
	ShowCalls   bool   `toml:"stream_calls"`     // Stream (or print) tool calls (potentially leaking data).
	DumpDir     string `toml:"dump_dir"`         // Dump all completions to JSON files in this dir.
	LogToolArgs bool   `toml:"log_tool_args"`    // Log the tool args (potentially leaking data).

	// Custom tool definitions:
	CustomTools []any `toml:"custom_tools,omitempty"` // TBD, will be custom type

	// Tool access control:
	// (Can use /regexp/ syntax.)
	NoTools     bool     `toml:"no_tools,omitempty"`     // Unregister all tools and remove from agents.
	AllowTools  []string `toml:"allow_tools,omitempty"`  // Only these tools will remain registered.
	RemoveTools []string `toml:"remove_tools,omitempty"` // These tools will be unregistered.
	AgentTools  []string `toml:"agent_tools,omitempty"`  // Override all agent Tools with this if set.

	// Agent configs:
	Agents []*agent.Config `toml:"agents,omitempty"` // Multiple Agents.

}

// LoadConfigs loads a runner config, followed by any agent configs, with the
// nonzero/non-nil values of c taking precedence.
//
// Agents are appended.  ConformAgents and Validate are called before
// returning.
//
// This is normally used when c holds flag values and config files are loaded
// before executing runner functions.
//
// No action is taken here outside the config structs themselves.
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
		c.NoLog = c.NoLog || r.NoLog
		c.NoTools = c.NoTools || r.NoTools
		if c.LogFile == "" {
			c.LogFile = r.LogFile
		}
		if c.DumpDir == "" {
			c.DumpDir = r.DumpDir
		}

		// Tool selection lists are taken from the original if non-nil, else
		// from the file.  In normal operation you will only have values here
		// if they were set with flags, but you *could* override that.
		if c.AllowTools == nil {
			c.AllowTools = r.AllowTools
		}
		if c.RemoveTools == nil {
			c.RemoveTools = r.RemoveTools
		}
		if c.AgentTools == nil {
			c.AgentTools = r.AgentTools
		}

		// We keep all agents!
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
	if (c.LogFile != "" || c.Debug) && c.NoLog {
		return fmt.Errorf("Logging can not be both specified and disabled.")
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
		a.NoLog = c.NoLog
		a.DumpDir = c.DumpDir
		if c.NoTools {
			a.Tools = nil
		} else if len(c.AgentTools) > 0 {
			a.Tools = c.AgentTools
		}
	}
}
