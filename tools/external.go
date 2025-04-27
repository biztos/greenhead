package tools

import (
	"fmt"
	"slices"
	"strings"

	"github.com/biztos/greenhead/utils"
)

// ExternalToolArg represents an argument (or option) for a command.
//
// Note that the Connector affects how options are passed to the command.
// If there is a Connector, the Flag and incoming value will be sent as a
// single arg.  This is usually not advisable.
type ExternalToolArg struct {
	Flag        string `toml:"spec"`        // Flag to use in option, e.g. "-foo"; empty for args.
	Key         string `toml:"key"`         // Key for schema; defaults to Flag, ergo required for args.
	Type        string `toml:"type"`        // Type for schema; defaults to boolean.
	Connector   string `toml:"connector"`   // Connector of flag to value.
	Description string `toml:"description"` // Description, used in the schema, no default.
	Optional    bool   `toml:"optional"`    // Is the arg optional or required?
	Repeat      bool   `toml:"repeat"`      // Can repeat this arg (use array in schema).
}

var ExternalToolArgTypes = []string{"string", "number", "integer", "boolean"}

var ErrExternalToolArgInvalid = fmt.Errorf("external tool arg invalid")

// Validate sanity-checks the values of arg a, setting defaults as needed.
//
// It is called from the ExternalToolConfig.Validate function.
func (a *ExternalToolArg) Validate() error {

	// Either flag or key must be set, otherwise we can't make a schema.
	if a.Key == "" {
		a.Key = a.Flag
	}
	if a.Key == "" {
		return fmt.Errorf("%w: neither key nor flag specified",
			ErrExternalToolArgInvalid)
	}

	// Type must be one of SupportedTypes.
	if a.Type == "" {
		a.Type = "boolean"
	}
	if !slices.Contains(ExternalToolArgTypes, a.Type) {
		return fmt.Errorf("%w: unsupported type for %q: %q",
			ErrExternalToolArgInvalid, a.Key, a.Type)
	}

	// Connector is for options only, i.e. requires flag.
	if a.Connector != "" && a.Flag == "" {
		return fmt.Errorf("%w: connector without flag for %q: %q",
			ErrExternalToolArgInvalid, a.Key, a.Connector)
	}
	// NB: we do *not* default the connector to anything, because the call
	// will be much cleaner without a connector.

	return nil

}

// TODO: send incoming payload as JSON text, validation *optional*
// TODO: optionals have type of ["real-type","null"]
//   https://platform.openai.com/docs/guides/function-calling?api-mode=chat

// hmm, we will need to round-trip this, from the input schema back...
// also we need a type for the input, but we can't dynamically create a type
// ...not matter if we gen the schema without the type, right?

// // config example
// var x = ExternalToolConfig{
// 	Name:        "list_files",
// 	Description: "The UNIX `ls` command, used to list files and directories.",
// 	Command:     "/bin/ls",
// 	Args: []string{
// 		{"-l", "long", "List files in the long format."},
// 		{"-A", "all", "List all files, including dotfiles."},
// 	},
// 	Args: {"[FILE...]"},
// }

// ExternalToolConfig represents the configuration of an ExternalTool.
//
// This is used within a Runner config.
type ExternalToolConfig struct {
	Name        string
	Description string
	Command     string // Path to the executable command.
	Args        []*ExternalToolArg

	PreArgs []string // Args to include before any specific tool args.
}

var ErrExternalToolConfigInvalid = fmt.Errorf("invalid external tool config")

// Validate checks that c has correct values:
//
// - Name and Description must not be empty.
// - Command must point to an executable file.
// - Args must all have allowed values, and not have redundant keys.
func (c *ExternalToolConfig) Validate() error {

	if strings.TrimSpace(c.Name) == "" {
		return fmt.Errorf("%w: empty name",
			ErrExternalToolConfigInvalid)
	}
	if strings.TrimSpace(c.Description) == "" {
		return fmt.Errorf("%w: empty description for %q",
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
