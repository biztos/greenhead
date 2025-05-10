// examples/custom/main.go -- example of a custom tool.
//
// This uses the default setup, with a custom function parse_url defined here,
// and some modifications to the commands.  It includes all the built-in tools
// as well.
//
// This should serve as a demonstration of the "easy" path to deploying custom
// runners with tools programmed in Go.
//
// To check that it registered the tools, use:
//
//	go run ./examples/custom tools show parse_url
package main

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/spf13/cobra"

	"github.com/biztos/greenhead"
	_ "github.com/biztos/greenhead/tools/all"
)

func main() {
	// Use a custom name but keep the default description.
	greenhead.CustomApp("custom", "1.0.0", "SuperCorp URL Parser", "")

	// Disallow chat, we don't trust you!
	greenhead.RemoveCommand("chat")

	// Never let the program run in silent mode.
	//
	// Note: to lock down and only use a set config file, use ResetFlags.
	greenhead.RemoveFlag("silent")

	// Add a custom command.
	greenhead.AddCommand(&cobra.Command{
		Use:   "hello [NAME...]",
		Short: "Say hello.",
		Long:  `Greets the provided name.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			greeting := "HELLO"
			if len(args) > 0 {
				greeting += " " + strings.Join(args, " ")
			}
			fmt.Println(strings.ToUpper(greeting))
			return nil
		},
	})

	greenhead.Run()
}

type ParseUrlInput struct {
	Url string
}

func ParseUrl(ctx context.Context, in ParseUrlInput) (*url.URL, error) {
	return url.Parse(in.Url) // context ignored
}

func init() {

	parse_url := greenhead.NewTool[ParseUrlInput, *url.URL](
		"parse_url",
		"Parses an URL and returns its parts in a struct.",
		ParseUrl,
	)
	if err := greenhead.Register(parse_url); err != nil {
		panic(err)
	}
}
