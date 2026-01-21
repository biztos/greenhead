// examples/lockdown/main.go -- a maximally locked-down custom binary.
//
// ** WORK IN PROGRESS **
//
// This CLI is as locked-down as possible, demonstrating the ability to make
// a custom binary that explicitly disallows user customization.
//
// Note that we load all tools then prune the registered tools.  This is in
// real life harder than just loading the tools you want, but IRL you also
// might not want to bother loading them individually when a regexp is easier
// to read.  Imagine for example wanting all of /^net/ recursively.
//
// TODO: Best practice for API Keys here is TBD!
//
// To check that it registered the tools, use:
//
//	go run ./examples/lockdown -h
package main

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/biztos/greenhead"
	"github.com/biztos/greenhead/ghd/runner"
	_ "github.com/biztos/greenhead/ghd/tools/all"
	"github.com/biztos/greenhead/ghd/utils"
)

var Stdout = os.Stdout // For testing.

var ConfigToml = `# The Lockdown Agent Built-In Config

# Top-level tool control:
allow_tools = ["/^demo/"]

# The Agent:
[[agents]]
  # Basics:
  name = "LockdownAgent"
  description = "An agent that thinks it's safe."
  type = "openai"
  model = "gpt-4o"

  # Output:
  stream = true
  show_calls = true
  color = "lightblue"
  log_file = "/tmp/lockdown-agent.log"

  # Safety:
  max_completions = 20
  max_toolchain = 3
  tools = ["/^demo/"]
  stop_matches = ["/abort/"]
  
  # System Prompt:
  [[agents.context]]
    Role = "system"
    content =  """\
You are a helpful assistant, but you are stern and prefer to stay focused \
on the task at hand.  That task is performing storing and retrieving \
information using the demo commands. \
  """

`

func main() {
	// Use a custom name but keep the default description.
	greenhead.CustomApp("lockdown", "1.0.0", "SuperCorp Secret Agent",
		`Lockdown Agent

Runs an OpenAI agent with tool-calling abilities for the various demo
functions, which can be listed with "list".`)

	// No default commands, no flags.
	greenhead.ResetCommands()
	greenhead.ResetFlags()

	// Custom config built in, with one agent, see above.
	cfg := &runner.Config{}
	utils.MustUnTomlString(ConfigToml, cfg)

	// One runner named "run" will run a single prompt a la "agents run".
	run := &cobra.Command{
		Use:   "run $MY_PROMPT",
		Short: "Run the agent with the user prompt provided.",
		Long: `Runs a single completion with the Lockdown Agent.

If the prompt begins with '@' then it will be read from a file, e.g. @foo.txt.
`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			r, err := runner.NewRunner(cfg)
			if err != nil {
				return err
			}
			return r.RunAgents(Stdout, args[0], false)
		},
	}
	greenhead.AddCommand(run)

	// Custom version of tools list.
	list := &cobra.Command{
		Use:   "list",
		Short: "List tools.",
		Long:  `Lists the tools that are available to the Lockdown Agent.`,
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			r, err := runner.NewRunner(cfg)
			if err != nil {
				return err
			}
			return r.ListTools(Stdout, false)
		},
	}
	greenhead.AddCommand(list)

	greenhead.Run()
}
