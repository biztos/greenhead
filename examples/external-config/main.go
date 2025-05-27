// examples/external-config/main.go -- example of a configured external tool.
//
// See init() for the TOML config source.
//
// NOTE: this only works when run from the greenhead root directory; otherwise
// the path to the executable will be wrong.
//
// For an example of the same thing baked into the Go proper, which provides
// extra flexibility for calling the tool from the CLI itself, see:
//
//	examples/external/main.go
//
// To run from the project root:
//
//	go run ./examples/external-config/ agents run 'hello world' --show-calls
package main

import (
	"github.com/biztos/greenhead"
	"github.com/biztos/greenhead/cmd"
	"github.com/biztos/greenhead/utils"

	_ "github.com/biztos/greenhead/tools/all"
)

func main() {

	// Boilerplate setup and run:
	greenhead.CustomApp("external-config", "1.0.0", "SuperCorp External Tool",
		"In real life, External means Internal -- to SuperCorp!")
	greenhead.Run()
}

func init() {

	var ConfigToml = `# Config including the echo_format tool.
[[external_tools]]
  name = "echo_format"
  description = "Echo args back with formatting."
  command = "testdata/external_command.pl"

  [[external_tools.args]]
    flag = "--seed"
    key = "seed"
    type = "number"
    description = "Seed ID with this real number"
    optional = false
    repeat = false

  [[external_tools.args]]
    flag = "--header"
    key = "header"
    type = "string"
    description = "Header lines to print before echoing."
    optional = false
    repeat = true

  [[external_tools.args]]
    flag = "--indent"
    key = "indent"
    type = "integer"
    description = "Number of spaces to input the lines."
    optional = false
    repeat = false

  [[external_tools.args]]
    flag = "--prefix"
    key = "prefix"
    type = "string"
    description = "Prefix to print after indent on each line."
    optional = false
    repeat = false

  [[external_tools.args]]
    flag = "--reverse"
    key = "reverse"
    type = "boolean"
    description = "Reverse the text of each line, excluding headers."
    optional = false
    repeat = false

  [[external_tools.args]]
    flag = ""
    key = "line"
    type = "string"
    description = "Line of text to echo back."
    optional = false
    repeat = true

[[agents]]
  name = "echobot"
  description = "An agent that likes to echo things with formatting."
  type = "openai"
  model = "gpt-4o"
  endpoint = ""
  tools = ["echo_format","/^demo/"]
  max_completion_tokens = 0
  max_completions = 100
  max_tokens = 0
  max_toolchain = 10
  abort_on_refusal = false
  color = "darkkhaki"
  bg_color = ""
  stream = false
  show_calls = false
  silent = false
  debug = false
  log_file = ""
  no_log = false
  dump_dir = ""
  log_tool_args = false

  [[agents.context]]
    role = "system"
    content =  """\
  You are a helpful assistant.  You have a set of tools which you use as \
  instructed.  By default, if there is no further instruction, you use the \
  echo_format tool to echo back your input with an indent value of 4, \
  a header of "goobers" and reverse set to true. \
  """
`
	utils.MustUnTomlString(ConfigToml, cmd.Config)

}
