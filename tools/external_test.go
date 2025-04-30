package tools_test

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/stretchr/testify/require"

	"github.com/biztos/greenhead/registry"
	"github.com/biztos/greenhead/tools"
	"github.com/biztos/greenhead/utils"
)

var ToyCommandValid = false // set in init

func ToyCommand() string {
	cwd, _ := os.Getwd()
	top := filepath.Dir(cwd)
	path, err := filepath.Abs(filepath.Join(top, "testdata", "external_command.pl"))
	if err != nil {
		panic(err)
	}
	return path
}

// Return a full config to exercise the toy command; trim as needed.
func ToyConfig() *tools.ExternalToolConfig {
	return &tools.ExternalToolConfig{
		Name:        "echo_format",
		Description: "Echo args back with formatting.",
		Command:     ToyCommand(),
		Args: []*tools.ExternalToolArg{
			{
				Flag:        "--seed",
				Type:        "number",
				Description: "Seed ID with this real number",
			},
			{
				Flag:        "--header",
				Type:        "string",
				Description: "Header lines to print before echoing.",
				Repeat:      true,
			},
			{
				Flag:        "--indent",
				Type:        "integer",
				Description: "Number of spaces to input the lines.",
			},
			{
				Flag:        "--prefix",
				Type:        "string",
				Description: "Prefix to print after indent on each line.",
			},
			{
				Flag:        "--reverse",
				Description: "Reverse the text of each line, excluding headers.",
			},
			{
				Key:         "line",
				Description: "Line of text to echo back.",
				Repeat:      true,
			},
		},
		PreArgs:       []string{},
		SendInput:     false,
		CombineOutput: true,
	}

}

func TestExternalToolExecOK(t *testing.T) {

	if !ToyCommandValid {
		t.Skip("toy command missing or invalid")
	}

	require := require.New(t)

	cfg := ToyConfig()

	tool, err := tools.NewExternalTool(cfg)
	require.NoError(err, "new")

	input := `{
	"seed": 1.2,
	"indent": 4,
	"prefix": "--",
	"header": ["h1","h2"],
	"reverse": false,
	"line": ["one","two"]
}`

	exp := `691f4bcc60fad8d9f0f8eb5b0189d538
h1
h2
    --one
    --two
`

	res, err := tool.Exec(context.Background(), input)
	require.NoError(err, "exec")
	require.Equal(exp, res)

}

func TestExternalToolExecFailBadInput(t *testing.T) {

	if !ToyCommandValid {
		t.Skip("toy command missing or invalid")
	}

	require := require.New(t)

	cfg := ToyConfig()

	tool, err := tools.NewExternalTool(cfg)
	require.NoError(err, "new")

	input := `{ "foo": 1.2 }`

	_, err = tool.Exec(context.Background(), input)
	require.ErrorIs(err, tools.ErrInvalidInput)

}

func TestExternalToolExecFailTimeout(t *testing.T) {

	if !ToyCommandValid {
		t.Skip("toy command missing or invalid")
	}

	require := require.New(t)

	cfg := ToyConfig()

	cfg.PreArgs = []string{"--stderr", "--sleep", "10.0"} // we get first arg
	tool, err := tools.NewExternalTool(cfg)
	require.NoError(err, "new")

	input := `{
	"seed": 1.2,
	"indent": 0,
	"prefix": "--",
	"header": [],
	"reverse": false,
	"line": ["one","two"]
}`

	// NOTE: tune this if we have to run with very slow Perl for some reason.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second/8)
	defer cancel()

	_, err = tool.Exec(ctx, input)
	require.ErrorIs(err, tools.ErrCommandTimedOut)
	cerr := err.(tools.CommandError)

	// We autoflush so we get some output here.  IRL you might not.
	require.Equal("691f4bcc60fad8d9f0f8eb5b0189d538\n", cerr.Stdout)
	require.Equal("--one\n", cerr.Stderr)

}

func TestNewExternalToolOK(t *testing.T) {

	require := require.New(t)

	tool, err := tools.NewExternalTool(ToyConfig())
	require.NoError(err, "NewExternalTool")

	require.Equal("echo_format", tool.Name())
	require.Equal("Echo args back with formatting.", tool.Description())
	// This is a little ridiculous but it does prove we have a valid Tooler.
	// (At compile time, which is the ridiculous bit.)
	require.NoError(registry.Register(tool))

}

func TestNewExternalToolFailsBadConfig(t *testing.T) {

	require := require.New(t)

	_, err := tools.NewExternalTool(&tools.ExternalToolConfig{})
	require.ErrorIs(err, tools.ErrExternalToolConfigInvalid)

}

func TestExternalToolInputSchemaOK(t *testing.T) {

	require := require.New(t)

	tool, err := tools.NewExternalTool(ToyConfig())
	require.NoError(err, "NewExternalTool")

	exp := `{
  "type": "object",
  "properties": {
    "seed": {
      "type": "number"
    },
    "header": {
      "type": "array",
      "items": {
        "type": "string"
      }
    },
    "indent": {
      "type": "integer"
    },
    "prefix": {
      "type": "string"
    },
    "reverse": {
      "type": "boolean"
    },
    "line": {
      "type": "array",
      "items": {
        "type": "string"
      }
    }
  },
  "additionalProperties": false,
  "required": [
    "seed",
    "header",
    "indent",
    "prefix",
    "reverse",
    "line"
  ]
}
`
	got := utils.MustJsonStringPretty(tool.InputSchema())
	require.JSONEq(exp, got) // random hash order could bit us otherwise.
}

func TestExternalToolOpenAiToolOK(t *testing.T) {

	require := require.New(t)

	tool, err := tools.NewExternalTool(ToyConfig())
	require.NoError(err, "NewExternalTool")

	exp := `{
  "type": "function",
  "function": {
    "name": "echo_format",
    "description": "Echo args back with formatting.",
    "strict": true,
    "parameters": {
      "type": "object",
      "properties": {
        "seed": {
          "type": "number"
        },
        "header": {
          "type": "array",
          "items": {
            "type": "string"
          }
        },
        "indent": {
          "type": "integer"
        },
        "prefix": {
          "type": "string"
        },
        "reverse": {
          "type": "boolean"
        },
        "line": {
          "type": "array",
          "items": {
            "type": "string"
          }
        }
      },
      "additionalProperties": false,
      "required": [
        "seed",
        "header",
        "indent",
        "prefix",
        "reverse",
        "line"
      ]
    }
  }
}
`
	got := utils.MustJsonStringPretty(tool.OpenAiTool())
	// require.Equal("", got)
	require.JSONEq(exp, got) // random hash order could bit us otherwise.
}

func TestExternalToolConfigValidateOK(t *testing.T) {

	require := require.New(t)

	// our toy config must always be valid, duh.
	require.NoError(ToyConfig().Validate())

}

func TestExternalToolConfigValidateFailsNoName(t *testing.T) {

	require := require.New(t)
	cfg := &tools.ExternalToolConfig{}
	err := cfg.Validate()
	require.ErrorIs(err, tools.ErrExternalToolConfigInvalid)
	require.ErrorContains(err, "empty name")
}

func TestExternalToolConfigValidateFailsNoDescription(t *testing.T) {

	require := require.New(t)
	cfg := &tools.ExternalToolConfig{Name: "foo"}
	err := cfg.Validate()
	require.ErrorIs(err, tools.ErrExternalToolConfigInvalid)
	require.ErrorContains(err, `empty description for "foo"`)
}

func TestExternalToolConfigValidateFailsNoCommand(t *testing.T) {

	require := require.New(t)
	cfg := &tools.ExternalToolConfig{Name: "foo", Description: "bar"}
	err := cfg.Validate()
	require.ErrorIs(err, tools.ErrExternalToolConfigInvalid)
	require.ErrorContains(err, `empty command for "foo"`)
}

func TestExternalToolConfigValidateFailsCommandNotFound(t *testing.T) {

	require := require.New(t)

	cmd := filepath.Join("no", "such", "thing", ulid.Make().String())
	cfg := &tools.ExternalToolConfig{
		Name:        "foo",
		Description: "bar",
		Command:     cmd,
	}
	err := cfg.Validate()
	require.ErrorIs(err, tools.ErrExternalToolConfigInvalid)
	require.ErrorContains(err, `command error for "foo"`)
	require.ErrorContains(err, "no such file or directory")

}

func TestExternalToolConfigValidateFailsCommandNotExecutable(t *testing.T) {

	require := require.New(t)

	cmd := "external.go"
	cfg := &tools.ExternalToolConfig{
		Name:        "foo",
		Description: "bar",
		Command:     cmd,
	}
	err := cfg.Validate()
	require.ErrorIs(err, tools.ErrExternalToolConfigInvalid)
	require.ErrorContains(err, `command not executable for "foo"`)

}

func TestExternalToolConfigValidateFailsBadToolArg(t *testing.T) {

	require := require.New(t)

	cfg := ToyConfig()
	cfg.Args = []*tools.ExternalToolArg{
		{}, // ergo no name
	}
	err := cfg.Validate()
	require.ErrorIs(err, tools.ErrExternalToolConfigInvalid)
	require.ErrorIs(err, tools.ErrExternalToolArgInvalid)
	require.ErrorContains(err, `"echo_format" arg 0`)
	require.ErrorContains(err, "neither key nor flag specified")

}

func TestExternalToolConfigValidateFailsDupeToolKey(t *testing.T) {

	require := require.New(t)

	cfg := ToyConfig()
	cfg.Args = []*tools.ExternalToolArg{
		{Key: "foo"},
		{Key: "foo"},
	}
	err := cfg.Validate()
	require.ErrorIs(err, tools.ErrExternalToolConfigInvalid)
	require.ErrorContains(err, `"echo_format" arg 1`)
	require.ErrorContains(err, `duplicate key "foo"`)

}

func TestExternalToolArgValidateFailsNoKey(t *testing.T) {

	require := require.New(t)

	arg := &tools.ExternalToolArg{
		Flag:        "--",
		Key:         "",
		Type:        "",
		Description: "",
		Optional:    true,
		Repeat:      true,
	}
	err := arg.Validate()
	require.ErrorIs(err, tools.ErrExternalToolArgInvalid)
	require.ErrorContains(err, "neither key nor flag specified")

}

func TestExternalToolArgValidateFailsBadType(t *testing.T) {

	require := require.New(t)

	arg := &tools.ExternalToolArg{
		Flag:        "",
		Key:         "arg1",
		Type:        "flubber",
		Description: "",
		Optional:    true,
		Repeat:      true,
	}
	err := arg.Validate()
	require.ErrorIs(err, tools.ErrExternalToolArgInvalid)
	require.ErrorContains(err, `unsupported type for "arg1": "flubber"`)

}

func TestExternalToolArgValidateOK(t *testing.T) {

	require := require.New(t)

	arg := &tools.ExternalToolArg{
		Flag:        "--indent",
		Key:         "indent_level",
		Type:        "integer",
		Description: "Indent with this many spaces.",
		Optional:    true,
		Repeat:      true,
	}
	require.NoError(arg.Validate())

}

func TestExternalToolArgValidateUntypedOptDefaultsOK(t *testing.T) {

	require := require.New(t)

	arg := &tools.ExternalToolArg{
		Flag:        "--debug",
		Key:         "",
		Type:        "",
		Description: "Run in debug mode.",
		Optional:    true,
		Repeat:      false,
	}
	exp := &tools.ExternalToolArg{
		Flag:        "--debug",
		Key:         "debug",
		Type:        "boolean",
		Description: "Run in debug mode.",
		Optional:    true,
		Repeat:      false,
	}
	require.NoError(arg.Validate())
	require.EqualValues(exp, arg)

}

func TestExternalToolArgValidateUntypedArgDefaultsOK(t *testing.T) {

	require := require.New(t)

	arg := &tools.ExternalToolArg{
		Flag:        "",
		Key:         "INPUT_FILE",
		Type:        "",
		Description: "File to read.",
		Optional:    true,
		Repeat:      false,
	}
	exp := &tools.ExternalToolArg{
		Flag:        "",
		Key:         "INPUT_FILE",
		Type:        "string",
		Description: "File to read.",
		Optional:    true,
		Repeat:      false,
	}
	require.NoError(arg.Validate())
	require.EqualValues(exp, arg)

}

func init() {
	// Make sure the toy command is available and the config valid.
	// (This exercises some code that is also unit-tested, which is fine.
	// If this fails, the test will fail with (possibly) more info.)
	cfg := ToyConfig()
	_, err := tools.NewExternalTool(cfg)
	if err != nil {
		fmt.Println("TOY COMMAND:", err)
		ToyCommandValid = false
	} else {
		fmt.Println("TOY COMMAND:", cfg.Command)
		ToyCommandValid = true
	}

}
