// Package runner defines the logic to run ghd, the Greenhead CLI tool.
package runner

import (
	"fmt"
	"io"
	"log/slog"
	"os"

	"github.com/fabien-marty/slog-helpers/pkg/human"

	"github.com/biztos/greenhead/agent"
	"github.com/biztos/greenhead/registry"
	"github.com/biztos/greenhead/tools"
)

// Runner is the runner of commands.
type Runner struct {
	Config *Config
	Agents []*agent.Agent
	Logger *slog.Logger
}

// NewRunner returns a new runner with the configuration processed.
//
// It is a thin wrapper around SetupTools, CreateLogger and CreateAgents.
//
// Logger is set as the slog default.
func NewRunner(cfg *Config) (*Runner, error) {

	if err := SetupTools(cfg); err != nil {
		return nil, err
	}
	logger, err := CreateLogger(cfg)
	if err != nil {
		return nil, err
	}
	slog.SetDefault(logger)
	agents, err := CreateAgents(cfg)
	if err != nil {
		return nil, err
	}
	return &Runner{
		Config: cfg,
		Agents: agents,
		Logger: logger,
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

	// Register any external tools before dealing with other limits.
	if err := RegisterExternalTools(cfg.ExternalTools); err != nil {
		return err
	}

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

// RegisterExternalTools registers all the external tools defined in configs,
// and should be called before setting available tools for a runner or agent.
//
// Note that overriding same-named built-in tools is explicitly allowed
// unless disabled with registry.LockForReplace() -- but duplicate names
// within the same call to this function are not allowed.
func RegisterExternalTools(configs []*tools.ExternalToolConfig) error {

	if len(configs) == 0 {
		return nil
	}

	// Check names before trying to register anything.
	ext_tools := make([]*tools.ExternalTool, 0, len(configs))
	have := map[string]bool{}
	for _, cfg := range configs {
		tool, err := tools.NewExternalTool(cfg)
		if err != nil {
			return err
		}
		if have[cfg.Name] {
			return fmt.Errorf("%w: %q", ErrExternalToolDupeName, cfg.Name)
		}
		have[cfg.Name] = true
		ext_tools = append(ext_tools, tool)
	}

	// Now register them.
	for _, tool := range ext_tools {
		if err := registry.Register(tool); err != nil {
			return fmt.Errorf("failed to register %q: %s", tool.Name(), err)
		}
	}
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

// CreateLogger creates a master logger according to cfg.
//
// This should be expanded downstream by using With.
func CreateLogger(cfg *Config) (*slog.Logger, error) {

	if cfg.NoLog {
		return slog.New(slog.NewTextHandler(io.Discard, nil)), nil
	}

	var handler slog.Handler
	level := slog.LevelInfo
	if cfg.Debug {
		level = slog.LevelDebug
	}
	out := os.Stderr
	if cfg.LogFile != "" {
		// Log to a file.
		f, err := os.OpenFile(cfg.LogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return nil, fmt.Errorf("failed to open log file: %w", err)
		}
		out = f
	}
	if cfg.LogHuman {
		// Human format is very nice for watching logs in the console, but
		// it should not do color to a file, that makes the file very hard to
		// read (though arguably fun to cat).
		// NOTE: ReplaceAttr doesn't work here, presumably it's overriding
		// with internal stuff.  But fine, don't care for now.
		// NOTE: kinda surprising there aren't more interesting sloggers!
		handler = human.New(out, &human.Options{
			HandlerOptions: slog.HandlerOptions{
				Level: level,
			},
			UseColors: cfg.LogFile == "",
		})
	} else if cfg.LogText {
		handler = slog.NewTextHandler(out, &slog.HandlerOptions{
			Level: level,
		})
	} else {
		handler = slog.NewJSONHandler(out, &slog.HandlerOptions{
			Level: level,
		})
	}

	// TODO: maybe give the runner an ULID of its own so we could always track
	// runs of the app.  Nice for observability.  Overkill?
	return slog.New(handler), nil

}
