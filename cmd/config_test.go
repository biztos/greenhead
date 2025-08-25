package cmd_test

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/biztos/greenhead/agent"
	"github.com/biztos/greenhead/cmd"
	"github.com/biztos/greenhead/runner"
)

var BadConfig = &runner.Config{
	Agents: []*agent.Config{&agent.Config{Type: "nonesuch"}},
}

func TestConfigCheckErrorViaExecute(t *testing.T) {

	// This may or may not be a better way to do things:
	// run the whole Execute, so we know the args are parsed
	// correctly.
	// However, it really wants a test suite so the buffer and
	// exit stuff can be done on setup/teardown, it'll be the same
	// for every command set we test.
	//
	// TODO: setup/teardown for this stuff.
	require := require.New(t)

	var outBuf bytes.Buffer
	var outWriter io.Writer = &outBuf
	var errBuf bytes.Buffer
	var errWriter io.Writer = &errBuf
	var didCallExit bool
	var exitCode int
	var exitFunc = func(code int) {
		didCallExit = true
		exitCode = code
	}
	cmd.Stdout = outWriter
	cmd.Stderr = errWriter
	cmd.Exit = exitFunc
	defer func() {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Exit = os.Exit
	}()

	cmd.Config = BadConfig
	os.Args = []string{"foo", "config", "check"}
	cmd.Execute()
	require.True(didCallExit, "exited")
	require.Equal(exitCode, 1, "exit code 1")

	require.Zero(outBuf.String(), "nothing to Stdout")
	require.Equal(`Error: error creating agent 1: no client for type "nonesuch"`,
		strings.TrimSpace(errBuf.String()),
		"error to Stderr")

}

func TestConfigCheckRunnerError(t *testing.T) {

	require := require.New(t)

	cmd.Config = BadConfig

	err := cmd.ConfigCheckCmd.RunE(nil, []string{})
	require.EqualError(err, `error creating agent 1: no client for type "nonesuch"`)

}

func TestConfigCheckRunnerOK(t *testing.T) {

	require := require.New(t)

	cmd.Config = &runner.Config{}

	err := cmd.ConfigCheckCmd.RunE(nil, []string{})
	require.NoError(err, "check OK")

}

func TestConfigDumpCmdRunnerError(t *testing.T) {

	require := require.New(t)

	cmd.Config = BadConfig

	err := cmd.ConfigDumpCmd.RunE(nil, []string{})
	require.EqualError(err, `error creating agent 1: no client for type "nonesuch"`)

}

func TestConfigDumpCmdRunnerOKToml(t *testing.T) {

	// TOML is default.
	//
	require := require.New(t)

	var buf bytes.Buffer
	var writer io.Writer = &buf

	cmd.Config = &runner.Config{} // TODO: add some stuff!
	cmd.Stdout = writer
	defer func() {
		cmd.Stdout = os.Stdout
	}()

	err := cmd.ConfigDumpCmd.RunE(nil, []string{})
	require.NoError(err, "no error on command")

	exp := `debug = false
log_file = ""
log_text = false
log_human = false
no_log = false
silent = false
stream = false
show_calls = false
dump_dir = ""
log_tool_args = false
max_completions = 0
max_toolchain = 0
no_tools = false

`
	require.Equal(exp, buf.String(), "toml as expected")

}

func TestConfigDumpCmdRunnerOKJson(t *testing.T) {

	// JSON requires setting a flag var.
	require := require.New(t)

	var buf bytes.Buffer
	var writer io.Writer = &buf

	cmd.Config = &runner.Config{} // TODO: add some stuff!
	cmd.Stdout = writer
	cmd.ConfigDumpJson = true
	defer func() {
		cmd.Stdout = os.Stdout
		cmd.ConfigDumpJson = false
	}()

	err := cmd.ConfigDumpCmd.RunE(nil, []string{})
	require.NoError(err, "no error on command")

	exp := `{
  "debug": false,
  "dump_dir": "",
  "log_file": "",
  "log_human": false,
  "log_text": false,
  "log_tool_args": false,
  "max_completions": 0,
  "max_toolchain": 0,
  "no_log": false,
  "no_tools": false,
  "show_calls": false,
  "silent": false,
  "stream": false
}
`
	require.Equal(exp, buf.String(), "json as expected")

}
