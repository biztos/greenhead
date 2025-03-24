// cmd/example/main.go -- a simple example of custom tools.
//
// This uses the default setup, with a custom function parse_url defined here,
// and the demo tools loaded from their submodule.
//
// This should serve as a demonstration of the "easy" path to deploying custom
// tools in agents.
//
// To check that it registered the tools, use:
//
//	go run ./cmd/example list-tools
package main

import (
	"context"
	"net/url"

	"github.com/biztos/greenhead"

	// Make tools available:
	_ "github.com/biztos/greenhead/tools/demo"
)

func main() {
	// Use a custom name but keep the default description.
	greenhead.CustomApp("example", "1.0.0", "SuperCorp URL Parser", "")
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
