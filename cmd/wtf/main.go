// wtf -- my work-in-progress binary to test shit out
//
// obviously: once this is published, keep this kind of thing elsewhere in its
// own repo.
//
// the runner func is set up in wtf.go, this here is just boilerplate
package main

import (
	"context"
	"net/url"

	"github.com/biztos/greenhead"
)

func main() {
	greenhead.CustomApp("wtf", "WTAF", "zerp!")
	greenhead.Run()
}

type ParseUrlInput struct {
	Url string `json:"url"`
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
