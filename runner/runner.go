// Package runner defines the logic to run ghd, the Greenhead CLI tool.
package runner

// Config describes the configuration used to run the application.
//
// Further configuration is at the Agent level, and can be included here.
type Config struct {
	// Informational:
	Debug   bool   `toml:"debug"`    // Log at DEBUG level instead of INFO.
	LogFile string `toml:"log_file"` // Write logs to this file instead of os.StdErr.
	Silent  bool   `toml:"silent"`   // Suppress LLM output (overrides stream).
	Stream  bool   `toml:"stream"`   // Stream LLM output if supported (overrides silent).
	DumpDir string `toml:"dump_dir"` // Dump all completions to JSON files in this dir.

}

// By convention, runner functions should follow the Cobra subcommand
// hierarchy, e.g. for "ghd tools list" see ListTools in tools.go.
