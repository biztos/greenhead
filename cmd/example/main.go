// cmd/example/main.go -- a simple example of custom tools.
//
// This uses the default setup, with a custom function parse_url.
//
// This should serve as a demonstration of the "easy" path to deploying custom
// tools in agents.
//
// To check that it registered the tool, use:
//
//	go run ./cmd/example list-tool
package main

import (
	"context"
	"net/url"

	"github.com/biztos/greenhead"
)

func main() {
	greenhead.CustomApp("example", "Greenhead URL Parser", "example 1.0.0")
	greenhead.Run()
}

type ParseUrlInput struct {
	Url string
}

func ParseUrl(ctx context.Context, in ParseUrlInput) (any, error) {
	return url.Parse(in.Url) // context ignored
}

func init() {

	parse_url := greenhead.NewTool[ParseUrlInput](
		"parse_url",
		"Parses an URL and returns its parts in a struct.",
		ParseUrl,
	)
	if err := greenhead.Register(parse_url); err != nil {
		panic(err)
	}
}
