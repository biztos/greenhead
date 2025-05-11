package runner

import (
	"errors"
	"fmt"
	"io"
	"path/filepath"

	"github.com/BurntSushi/toml"

	"github.com/biztos/greenhead/agent"
	"github.com/biztos/greenhead/assets"
	"github.com/biztos/greenhead/rgxp"
	"github.com/biztos/greenhead/utils"
)

// Config describes the configuration used to run the application.
//
// Further configuration is at the Agent level, and can be included here.
type Config struct {

	// Output control:
	Debug       bool   `toml:"debug"`         // Log at DEBUG level instead of INFO.
	LogFile     string `toml:"log_file"`      // Write logs to this file instead of os.StdErr.
	NoLog       bool   `toml:"no_log"`        // Do not log at all.
	Silent      bool   `toml:"silent"`        // Suppress LLM output if not already streamed.
	Stream      bool   `toml:"stream"`        // Stream LLM output if supported.
	ShowCalls   bool   `toml:"show_calls"`    // Stream (or print) tool calls (potentially leaking data).
	DumpDir     string `toml:"dump_dir"`      // Dump all completions to JSON files in this dir.
	LogToolArgs bool   `toml:"log_tool_args"` // Log the tool args (potentially leaking data).

	// Usage limits:
	MaxCompletions int `toml:"max_completions"` // Max number of completions to run.
	MaxToolChain   int `toml:"max_toolchain"`   // Max number of tool calls in a row.

	// Custom tool definitions:
	CustomTools []any `toml:"custom_tools"` // TBD, will be custom type

	// Tool access control:
	// (Can use /regexp/ syntax.)
	NoTools     bool                 `toml:"no_tools"`     // Unregister all tools and remove from agents.
	AllowTools  []*rgxp.OptionalRgxp `toml:"allow_tools"`  // Only these tools will remain registered.
	RemoveTools []*rgxp.OptionalRgxp `toml:"remove_tools"` // These tools will be unregistered.
	AgentTools  []*rgxp.OptionalRgxp `toml:"agent_tools"`  // Override all agent Tools with this if set.

	// Safety:
	StopMatches []*rgxp.Rgxp `toml:"stop_matches"` // Stop if any output matches any of these.

	// Agent configs:
	Agents []*agent.Config `toml:"agents"` // Multiple Agents.

	// API config:
	ApiListenAddress string `toml:"api_listen_address"` // Listen address, usually as ":3000".
}

var ErrNamedAgentNotAvailable = errors.New("named agent not available")
var ErrBadAgentConfig = errors.New("agent config error")

// LoadConfigs loads a runner config, followed by any agent configs, with the
// nonzero/non-nil values of c taking precedence.
//
// Agents are appended.  ConformAgents and Validate are called before
// returning.
//
// If agents have no extension then they are treated as built-in agents.
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
		// Bools, ints and strings from the original take precedence.
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
		if c.MaxCompletions == 0 {
			c.MaxCompletions = r.MaxCompletions
		}
		if c.MaxToolChain == 0 {
			c.MaxToolChain = r.MaxToolChain
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

		// StopMatches are handled the same as AllowTools.
		if c.StopMatches == nil {
			c.StopMatches = r.StopMatches
		}

		// We keep all agents!
		c.Agents = append(c.Agents, r.Agents...)

	}

	for _, file := range agentFiles {

		a := &agent.Config{}
		if filepath.Ext(file) == "" {
			a = NamedAgentConfigs[file]
			if a == nil {
				return fmt.Errorf("%w: %q", ErrNamedAgentNotAvailable, file)
			}
		} else if err := utils.UnmarshalFile(file, a); err != nil {
			return fmt.Errorf("%w: %q: %w", ErrBadAgentConfig, file, err)
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
// Special cases:
//
// - MaxCompletions and MaxToolChain only override if nonzero.
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
		if c.MaxCompletions != 0 {
			a.MaxCompletions = c.MaxCompletions
		}
		if c.MaxToolChain != 0 {
			a.MaxToolChain = c.MaxToolChain
		}
		if c.NoTools {
			a.Tools = nil
		} else if len(c.AgentTools) > 0 {
			a.Tools = c.AgentTools
		}
		if len(c.StopMatches) > 0 {
			a.StopMatches = c.StopMatches
		}
	}
}

// DumpJson dumps the config as indented JSON to w, or panics trying.
func (c *Config) DumpJson(w io.Writer) {

	// In order to respect the TOML struct tags, we round-trip into data.
	b := utils.MustToml(c)
	var v map[string]any
	if err := toml.Unmarshal(b, &v); err != nil {
		panic(err)
	}
	fmt.Fprintln(w, utils.MustJsonStringPretty(v))
}

// DumpToml dumps the config as uncommented TOML to w, or panics trying.
func (c *Config) DumpToml(w io.Writer) {

	fmt.Fprintln(w, utils.MustTomlString(c))

}

// NamedAgentConfigs holds the configs for agents that can be specified by
// name.  The built-ins are loaded at init.
var NamedAgentConfigs = map[string]*agent.Config{}

func init() {

	// We assume well-named agents in our assets.
	for _, name := range assets.PrefixNames("agents", false) {
		a := &agent.Config{}
		utils.MustUnToml(assets.MustAsset(name), a)
		NamedAgentConfigs[a.Name] = a
	}

}
