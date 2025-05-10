// examples/external/main.go -- example of an external tool.
//
// Also includes a custom "echo" command to run the tool directly, for quick
// and easy demonstration.
//
// NOTE: this only works when run from the greenhead root directory; otherwise
// the path to the executable will be wrong.
package main

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/biztos/greenhead"
	"github.com/biztos/greenhead/tools"
	"github.com/biztos/greenhead/utils"
)

var TheTool *tools.ExternalTool
var TheToolConfig *tools.ExternalToolConfig

type TheToolInput struct {
	Seed    float64  `json:"seed"`
	Indent  int      `json:"indent"`
	Prefix  string   `json:"prefix"`
	Header  []string `json:"header"`
	Reverse bool     `json:"reverse"`
	Line    []string `json:"line"`
}

func main() {

	input := &TheToolInput{}
	echo_cmd := &cobra.Command{
		Use:   "echo ARG1 ARG2...",
		Short: "Echo directly with the external tool.",
		Long: `Echoes args with some options. Full example:

external echo One Two --reverse --header H1 --header H2 --seed 1.1 --indent 3

See option descriptions below.`,
		RunE: func(cmd *cobra.Command, args []string) error {

			input.Line = args
			input_str := utils.MustJsonStringPretty(input)
			res, err := TheTool.Exec(context.Background(), input_str)
			if err != nil {
				return err
			}
			fmt.Println(res)
			return nil

		},
	}

	// For kicks (and demos) we let you customize the input.
	echo_cmd.Flags().Float64Var(&input.Seed, "seed", 0,
		"Seed ID for predictable value.")
	echo_cmd.Flags().IntVar(&input.Indent, "indent", 0,
		"Indent echoed lines this far.")
	echo_cmd.Flags().StringArrayVar(&input.Header, "header", []string{},
		"Header(s) to prepend.")
	echo_cmd.Flags().BoolVar(&input.Reverse, "reverse", false,
		"Reverse echoed lines.")

	// Boilerplate setup and run:
	greenhead.CustomApp("external", "1.0.0", "SuperCorp External Tool",
		"In real life, External means Internal -- to SuperCorp!")
	greenhead.AddCommand(echo_cmd)
	greenhead.Run()
}

func init() {

	// Set up an external tool with the toy command.
	// (tweak above, yay globals)
	TheToolConfig = &tools.ExternalToolConfig{
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
		PreArgs:       []string{},
		SendInput:     false,
		CombineOutput: true,
	}
	tool, err := tools.NewExternalTool(TheToolConfig)
	if err != nil {
		panic(err)
	}
	TheTool = tool
	greenhead.Register(tool)

}
