package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/biztos/greenhead/registry"
	"github.com/biztos/greenhead/tools"
)

const ExitCodeWTF = 999

var TheTool *tools.ExternalTool

// WtfCmd is for work in progress.
var WtfCmd = &cobra.Command{
	Hidden: true,
	Use:    "wtf",
	Short:  "WTF? WTF!",
	Long: `What.
The.
(Actual.)
F**K?!

Work-in-progress random experiments.  Use at your own risk!

TODO: equivalent of this, but in an external package.  Hard? Easy?

Maybe just hook into pre/postrun funcs?`,
	RunE: func(cmd *cobra.Command, args []string) error {

		// We set up the command in init below, so can use it on run et al.
		//
		// For the specific wtf-ery we want to actually run it!
		// Remember to test the background timeout for a while.
		fakein := `{
			"seed": 1.2,
			"indent": 4,
			"prefix": "--",
			"header": ["h1","h2"],
			"reverse": false,
			"line": ["one","two"]
		}`

		cmd_args, err := TheTool.CommandArgs(fakein)
		if err != nil {
			return err

		}
		fmt.Println(cmd_args)
		// res, err := tool.Exec(context.Background(), "foo")
		// if err != nil {
		// 	return err
		// }
		// fmt.Println("SUCCESS")
		// fmt.Println(res)
		return nil

	},
}

func init() {
	RootCmd.AddCommand(WtfCmd)

	// Set up an external tool with the toy command.
	cfg := &tools.ExternalToolConfig{
		Name:        "echo_format",
		Description: "Echo args back with formatting.",
		Command:     "testdata/external_command.pl",
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
		PreArgs:   []string{"--stderr"},
		SendInput: false,
	}
	tool, err := tools.NewExternalTool(cfg)
	if err != nil {
		panic(err)
	}
	TheTool = tool
	registry.Register(tool)

}
