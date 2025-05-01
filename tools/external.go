package tools

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"slices"
	"strings"

	"github.com/sashabaranov/go-openai"
	"github.com/sashabaranov/go-openai/jsonschema"

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
	Name          string             `toml:"name"`           // Name, required.
	Description   string             `toml:"description"`    // Description, required.
	Command       string             `toml:"command"`        // Path to the executable command.
	Args          []*ExternalToolArg `toml:"args"`           // Options/args as defined above.
	PreArgs       []string           `toml:"pre_args"`       // Args to include verbatim in every call.
	SendInput     bool               `toml:"send_input"`     // Send raw input JSON on STDIN instead of args.
	CombineOutput bool               `toml:"combine_output"` // Include STDERR after STDOUT in result.
}

var ErrExternalToolConfigInvalid = fmt.Errorf("invalid external tool config")

// Validate checks that c has correct values:
//
// - Name and Description must not be empty.
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

// ExternalTool represents a Tooler that executes an external binary.
//
// Nonzero return codes are considered errors.
type ExternalTool struct {
	cfg     *ExternalToolConfig
	argMap  map[string]*ExternalToolArg
	argList []string
}

// NewExternalTool creates an ExternalTool from cfg.
func NewExternalTool(cfg *ExternalToolConfig) (*ExternalTool, error) {

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	argMap := make(map[string]*ExternalToolArg, len(cfg.Args))
	argList := make([]string, 0, len(cfg.Args))
	for _, arg := range cfg.Args {
		argMap[arg.Key] = arg // already normalized in validation.
		argList = append(argList, arg.Key)
	}

	return &ExternalTool{
		cfg:     cfg,
		argMap:  argMap,
		argList: argList,
	}, nil
}

// Name implements Tooler.
func (t *ExternalTool) Name() string {
	return t.cfg.Name
}

// Description implements Tooler.
func (t *ExternalTool) Description() string {
	return t.cfg.Description
}

// Help implements Tooler.
func (t *ExternalTool) Help() string {

	s := fmt.Sprintf("%s\n\n%s\n\n", t.cfg.Name, t.cfg.Description)

	// Hmm, how to handle this?  Can we just use NewTool?
	// no because no way to coerce the type is there?

	return s
}

var ErrInvalidInput = fmt.Errorf("command input invalid")

// ValidateInput validates input against the InputSchema and returns the
// cleaned object if input is valid.
//
// Cleanup currently addresses only the conversion of float64 types to int
// where the arg specifies an "integer" type.
func (t *ExternalTool) ValidateInput(input string) (map[string]any, error) {

	// TODO: cache schemas! Here and for InputSchema.
	schema := t.InputSchema().(jsonschema.Definition)

	var v map[string]any
	b := []byte(input)
	if err := jsonschema.VerifySchemaAndUnmarshal(schema, b, &v); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInvalidInput, err)
	}

	return v, nil

}

// InputSchema implements Tooler, returning the "parameters" object for the
// tool definition.
func (t *ExternalTool) InputSchema() any {

	// TODO: handle optionals.
	// Spec is pretty weird about it, AI not helping right now.
	props := map[string]jsonschema.Definition{}
	for _, n := range t.argList {
		a := t.argMap[n]
		p := jsonschema.Definition{}
		if a.Repeat {
			p.Type = jsonschema.Array
			p.Items = &jsonschema.Definition{Type: jsonschema.DataType(a.Type)}
		} else {
			p.Type = jsonschema.DataType(a.Type)
		}
		props[a.Key] = p
	}

	return jsonschema.Definition{
		Type:                 jsonschema.Object,
		AdditionalProperties: false,
		Properties:           props,
		Required:             t.argList,
	}
}

var ErrExternalToolInputBadType = fmt.Errorf("wrong type for arg in tool input")

// CommandArgs returns the full set of arguments to send to the command based
// on the input JSON.  ValidateInput is called if SendInput is false.
//
// Note that the args do not include the t.Command.
func (t *ExternalTool) CommandArgs(input string) ([]string, error) {

	// Base args are always the same set.
	all_args := make([]string, len(t.cfg.PreArgs))
	copy(all_args, t.cfg.PreArgs)

	// If you want your raw input, you get your raw input.
	if t.cfg.SendInput {
		return all_args, nil
	}

	// The validation guarantees our type safety for casts below.
	// (Fun task: try to find a case where it doesn't.)
	input_map, err := t.ValidateInput(input)
	if err != nil {
		return nil, err
	}

	// Now we build our list, erroring out if we need to on the way.
	// TODO: revisit the valGetter[T] idea b/c probably faster. But bench it.
	// TODO: handle weird optional stuff openai-style.

	for _, k := range t.argList {
		arg := t.argMap[k]
		prop := input_map[k]
		vals := []any{}
		if arg.Repeat {
			// We trust the validation we did above, and cast to array.
			// If for some reason this doesn't work, figure out why and handle
			// it (no sense having something we can't test).
			array := prop.([]any)
			for _, val := range array {
				vals = append(vals, val)
			}
		} else {
			vals = append(vals, prop)
		}

		// Boolean options get special handling; all else is same-same.
		if arg.Type == "boolean" && arg.Flag != "" {
			for _, v := range vals {
				if v.(bool) == true {
					all_args = append(all_args, arg.Flag)
				}
			}
			continue

		}

		// All other values are stringified with the default Go style; if this
		// this turns out to be a problem, we revisit.
		// (And presumably we could make it a lot faster if we care, since
		// we know the types and could just use a big ugly switch.)
		svals := make([]string, 0, len(vals))
		for _, val := range vals {
			svals = append(svals, fmt.Sprint(val))
		}
		if arg.Flag == "" {
			// Plan arg, use as-is.
			all_args = append(all_args, svals...)
		} else {
			// Flag arg, use pair.
			for _, s := range svals {
				all_args = append(all_args, arg.Flag, s)
			}
		}

	}

	return all_args, nil
}

// CommandError represents an error returned from an external command.
//
// To obtain the exit code you must examine the Unwrap return value, which
// may be a *exec.ExitError.
type CommandError struct {
	err    error
	Stdout string
	Stderr string
}

// NewCommandError returns a CommandError with its underlying error set to
// err, or err wrapped with Stderr if Stderr is not empty.
//
// There is no limit to the amount of Stderr included, but it is
// space-trimmed.
func NewCommandError(err error, stdout string, stderr string) CommandError {

	if stderr != "" {
		err = fmt.Errorf("%w: %s", err, strings.TrimSpace(stderr))
	}
	return CommandError{err, stdout, stderr}
}

// Error implements the error interface.
func (e CommandError) Error() string {
	return e.err.Error()
}

// Unwrap returns the underlying error in e.
func (e CommandError) Unwrap() error {
	return e.err
}

// Detail returns a detailed human-readable version of e.
func (e CommandError) Detail() string {
	s := fmt.Sprintf("CommmandError: %s\n", e.Error())
	s += "-------------------------------------------------------------\n"
	if e.Stdout == "" {
		s += "<No Stdout>\n"
	} else {
		s += "Stdout:\n"
		s += "-------------------------------------------------------------\n"
		s += e.Stdout + "\n"
	}
	s += "-------------------------------------------------------------\n"
	if e.Stderr == "" {
		s += "<No Stderr>\n"
	} else {
		s += "Stderr:\n"
		s += "-------------------------------------------------------------\n"
		s += e.Stderr + "\n"
	}
	s += "-------------------------------------------------------------\n"
	return s
}

var ErrCommandTimedOut = fmt.Errorf("command timed out")
var ErrCommandCanceled = fmt.Errorf("command canceled")
var ErrCommandFailed = fmt.Errorf("command failed")

// Exec implements Tooler.
func (t *ExternalTool) Exec(ctx context.Context, input string) (any, error) {

	// Prepare the command
	args, err := t.CommandArgs(input)
	if err != nil {
		return nil, err
	}
	cmd := exec.CommandContext(ctx, t.cfg.Command, args...)
	if t.cfg.SendInput {
		cmd.Stdin = strings.NewReader(input)
	}

	// Capture stdout and stderr
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// Run the command
	err = cmd.Run()

	// Check for errors.
	ctx_err := ctx.Err()
	var cmd_err error
	if ctx_err == context.Canceled {
		cmd_err = ErrCommandCanceled
	} else if ctx_err == context.DeadlineExceeded {
		cmd_err = ErrCommandTimedOut
	} else if err != nil {
		cmd_err = ErrCommandFailed
	}
	if cmd_err != nil {
		return nil, NewCommandError(
			fmt.Errorf("%w: %w", cmd_err, err),
			stdout.String(),
			stderr.String(),
		)
	}
	if t.cfg.CombineOutput {
		return stdout.String() + stderr.String(), nil
	} else {
		return stdout.String(), nil
	}

}

// OpenAiTool implements Tooler.
// TODO: move this elsewhere!
func (t *ExternalTool) OpenAiTool() openai.Tool {
	return openai.Tool{
		Type: openai.ToolTypeFunction,
		Function: &openai.FunctionDefinition{
			Name:        t.cfg.Name,
			Description: t.cfg.Description,
			Strict:      true, // TODO: what does this mean?
			// TODO: prove this works, it *should* be good to go.
			Parameters: t.InputSchema(),
		},
	}
}
