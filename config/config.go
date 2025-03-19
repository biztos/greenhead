// Package config consolidates configuration types.
package config

// RunConfig describes the configuration used by runners.  Not all runners
// make use of all configs; see the Run* functions for details.
type RunConfig struct {
	// Informational:
	Debug   bool   // Log at DEBUG level instead of INFO.
	LogFile string // Write logs to this file instead of os.StdErr.
	Silent  bool   // Suppress LLM output (overrides stream).
	Stream  bool   // Stream LLM output if supported (overrides silent).
	DumpDir string // Dump all completions to JSON files in this dir.

	// Safety and limits:  (Zero generally means "no limit.")
	MaxTokens      int  // Maximum number of total tokens for all operations.
	MaxToolChain   int  // Max number of tool calls allowed in a row.
	AbortOnRefusal bool // Abort if a completion is refused by an LLM.

}

// ExtraConfig is an arbitrary struct used by specific implementations.
type ExtraConfig map[string]any

// AgentConfig describes the
type AgentConfig struct {
	// TODO
}
