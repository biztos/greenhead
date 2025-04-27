package tools_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/oklog/ulid/v2"
	"github.com/stretchr/testify/require"

	"github.com/biztos/greenhead/registry"
	"github.com/biztos/greenhead/tools"
	"github.com/biztos/greenhead/utils"
)

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
		PreArgs:   []string{"--stdin"},
		SendInput: true,
	}

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
