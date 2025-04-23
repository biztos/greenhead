package tools

import (
	"fmt"
)

// ExternalToolOption represents an option for a command.
//
// To force a numeric type, use NUM in Spec, as in: "--max=NUM".
type ExternalToolOption struct {
	Spec        string // Spec as in "-l" or "--key=API_KEY" or "-o FILE"
	Key         string // Key for input schema; must be unique within the config.
	Description string // Description of the option, passed to the LLM.
}

// Schema parses o into JSON schema format.
func (o *ExternalToolOption) Schema() (map[string]any, error) {

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

// hmm, we will need to round-trip this, from the input schema back...
// also we need a type for the input, but we can't dynamically create a type
// ...not matter if we gen the schema without the type, right?

// // config example
// var x = ExternalToolConfig{
// 	Name:        "list_files",
// 	Description: "The UNIX `ls` command, used to list files and directories.",
// 	Command:     "/bin/ls",
// 	Options: []string{
// 		{"-l", "long", "List files in the long format."},
// 		{"-A", "all", "List all files, including dotfiles."},
// 	},
// 	Args: {"[FILE...]"},
// }

// ExternalToolConfig represents the configuration of an ExternalTool.
//
// This is normally used within an Agent or Runner config.
type ExternalToolConfig struct {
	Name        string
	Description string
	Command     string // Path to the executable command.
	Options     []ExternalToolOption
	Args        []string

	PreArgs []string // Args to include before any specific tool args.
}

// Schema creates an input schema from c.
func (c *ExternalToolConfig) Schema() (map[string]any, error) {
	return nil, fmt.Errorf("TODO")
}

// Copy performs a deep copy of c.
func (c *ExternalToolConfig) Copy() *ExternalToolConfig {

	n := &ExternalToolConfig{
		Name:        c.Name,
		Description: c.Description,
		Command:     c.Command,
	}
	n.Options = make([]ExternalToolOption, len(c.Options))
	copy(n.Options, c.Options)
	n.Args = make([]string, len(c.Args))
	copy(n.Args, c.Args)
	return n
}

// ExternalTool represents a Tooler that executes an external binary.
//
// Nonzero return codes are considered errors.
type ExternalTool struct {
	cfg *ExternalToolConfig
}

// NewExternalTool creates an ExternalTool from cfg.
func NewExternalTool(cfg *ExternalToolConfig) (*ExternalTool, error) {

	// Check that command exists and can be run.

	// Convert Options and Args to a schema.

	// Keep a copy of the config, we need it when running.

	return nil, nil
}
