package tools

import (
	"fmt"
	"slices"
	"strings"

	"github.com/biztos/greenhead/utils"
)

// ExternalToolArg represents an argument (or option) for a command.
//
// An option has a flag; an argument has no flag.
type ExternalToolArg struct {
	Flag        string `toml:"spec"`        // Flag to use in option, e.g. "-foo"; empty for args.
	Key         string `toml:"key"`         // Key for schema; defaults to trimmed Flag; required for args.
	Type        string `toml:"type"`        // Type for schema; defaults to boolean.
	Description string `toml:"description"` // Description, used in the schema, no default.
	Optional    bool   `toml:"optional"`    // Is the arg optional or required?
	Repeat      bool   `toml:"repeat"`      // Can repeat this arg (use array in schema).
}

var ExternalToolArgTypes = []string{"string", "number", "integer", "boolean"}

var ErrExternalToolArgInvalid = fmt.Errorf("external tool arg invalid")

// Validate sanity-checks the values of arg a, setting defaults as needed.
//
// It is called from the config's Validate and need not be called separately.
func (a *ExternalToolArg) Validate() error {

	// Either flag or key must be set, otherwise we can't make a schema.
	if a.Key == "" {
		a.Key = strings.Trim(a.Flag, "-")
	}
	if a.Key == "" {
		return fmt.Errorf("%w: neither key nor flag specified",
			ErrExternalToolArgInvalid)
	}

	// Type must be one of SupportedTypes.
	if a.Type == "" {
		// For an arg with no flag, it defaults to string (just an arg).
		if a.Flag == "" {
			a.Type = "string"
		} else {
			// Otherwise it's treated as a classic on/off flag.
			a.Type = "boolean"
		}
	}
	if !slices.Contains(ExternalToolArgTypes, a.Type) {
		return fmt.Errorf("%w: unsupported type for %q: %q",
			ErrExternalToolArgInvalid, a.Key, a.Type)
	}

	return nil

}

// ExternalToolConfig represents the configuration of an ExternalTool.
//
// This is used within a Runner config.
type ExternalToolConfig struct {
	Name        string             `toml:"name"`        // Name, required.
	Description string             `toml:"description"` // Description, required.
	Command     string             `toml:"command"`     // Path to the executable command.
	Args        []*ExternalToolArg `toml:"args"`        // Options/args as defined above.
	PreArgs     []string           `toml:"pre_args"`    // Args to include verbatim in every call.
	SendInput   bool               `toml:"send_input"`  // Send command input JSON on STDIN.
	NoArgs      bool               `toml:"no_args"`     // Do *NOT* send Args if SendInput is true.
	NoValidate  bool               `toml:"no_validate"` // Do *NOT* validate if SendInput is true.
}

var ErrExternalToolConfigInvalid = fmt.Errorf("invalid external tool config")

// Validate checks that c has correct values:
//
// - Name and Description must not be empty.
// - NoArgs and NoValidate require SendInput.
// - Command must point to an executable file.
// - Args must all validate, and not have redundant keys.
func (c *ExternalToolConfig) Validate() error {

	if strings.TrimSpace(c.Name) == "" {
		return fmt.Errorf("%w: empty name",
			ErrExternalToolConfigInvalid)
	}
	if strings.TrimSpace(c.Description) == "" {
		return fmt.Errorf("%w: empty description for %q",
			ErrExternalToolConfigInvalid, c.Name)
	}
	if strings.TrimSpace(c.Command) == "" {
		return fmt.Errorf("%w: empty command for %q",
			ErrExternalToolConfigInvalid, c.Name)
	}
	can_exec, err := utils.IsExecutable(c.Command)
	if err != nil {
		return fmt.Errorf("%w: command error for %q: %w",
			ErrExternalToolConfigInvalid, c.Name, err)
	}
	if !can_exec {
		return fmt.Errorf("%w: command not executable for %q: %w",
			ErrExternalToolConfigInvalid, c.Name, err)
	}
	if !c.SendInput {
		if c.NoArgs {
			return fmt.Errorf("%w: args must be sent if input not sent for %q",
				ErrExternalToolConfigInvalid, c.Name)
		}
		if c.NoValidate {
			return fmt.Errorf("%w: arg validation can not be disabled for %q",
				ErrExternalToolConfigInvalid, c.Name)
		}
	}

	have_key := map[string]bool{}
	for i, arg := range c.Args {
		if err := arg.Validate(); err != nil {
			return fmt.Errorf("%w: %q arg %d: %w",
				ErrExternalToolConfigInvalid, c.Name, i, err)
		}
		if have_key[arg.Key] {
			return fmt.Errorf("%w: %q arg %d: duplicate key %q",
				ErrExternalToolConfigInvalid, c.Name, i, arg.Key)
		}
		have_key[arg.Key] = true

	}
	return nil

}

// Schema returns a JSON schema from the c.
func (c *ExternalToolConfig) Schema() (map[string]any, error) {

	// "parameters": {
	//   "type": "object",
	//   "additionalProperties": false,
	//   "properties": {
	//     "game_id": {
	//       "type": "string"
	//     }
	//   },
	//   "required": [
	//     "game_id"
	//   ],
	// }
	obj := map[string]any{
		"type":                 "object",
		"additionalProperties": false,
	}

	return obj, nil

}

// ExternalTool represents a Tooler that executes an external binary.
//
// Nonzero return codes are considered errors.
type ExternalTool struct {
	cfg    *ExternalToolConfig
	argMap map[string]*ExternalToolArg
}

// NewExternalTool creates an ExternalTool from cfg.
func NewExternalTool(cfg *ExternalToolConfig) (*ExternalTool, error) {

	// Check that command exists and can be run.

	// Convert Options and Args to a schema.

	// Keep a copy of the config, we need it when running.

	return nil, nil
}
