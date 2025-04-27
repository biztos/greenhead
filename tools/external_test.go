package tools_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/biztos/greenhead/tools"
	"github.com/biztos/greenhead/utils"
)

var AllOptsToml = `# TOML to exercise all options on the toy command.

options:
`

func ToyConfig() *tools.ExternalToolConfig {
	cfg := &tools.ExternalToolConfig{}
	utils.MustUnTomlString(AllOptsToml, cfg)
	return cfg
}

func SingleArgConfig(arg *tools.ExternalToolArg) *tools.ExternalToolConfig {

	return &tools.ExternalToolConfig{
		Name:        "name",
		Description: "desc",
		Command:     "cmd",
		Args:        []*tools.ExternalToolArg{arg},
		PreArgs:     []string{},
	}
}

func TestExternalToolArgValidateOK(t *testing.T) {

	require := require.New(t)

	arg := &tools.ExternalToolArg{
		Flag:        "--indent",
		Key:         "indent_level",
		Type:        "integer",
		Connector:   "",
		Description: "Indent with this many spaces.",
		Optional:    true,
		Repeat:      false,
	}
	err := arg.Validate()
	require.NoError(err)

}
