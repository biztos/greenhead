// Package runner defines the logic to run ghd, the Greenhead CLI tool.

package runner

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/biztos/greenhead/registry"

	"github.com/docopt/docopt-go"
)

type Config struct {
	Debug     bool
	ListTools bool `docopt:"list-tools"`
	ShowTool  bool `docopt:"show-tool"`
	RunTool   bool `docopt:"run-tool"`
	Tool      string
	Input     string
}

var Version = "ghd v0.1.0"

var Args = os.Args[1:]

var Exit = os.Exit

var Name = "ghd"

var Title = "Greenhead"

// TODO:
// gonna have to rewrite this using Cobra sooner or later anyway, so might as
// well get started soon.
//
// HOWEVER, first do some other stuff.
var UsageSource = `<title> (<name>)

Usage:
  <name> [options]
  <name> list-tools
  <name> show-tool <TOOL>
  <name> run-tool <TOOL> <INPUT>
  <name> -h | --help
  <name> --version

Options:
  -h --help           Show this screen.
  --debug             Debug mode (more verbose logging).

This program runs AI agents.

The rest is TBD.

Registered tools:

<tools>

(c) 2025 The Greenhead Authors; distributed under MIT License.

`

func Run() {

	cfg := &Config{}
	opts, err := docopt.ParseArgs(Usage(), Args, Version)
	if err != nil {
		log.Fatal(err)
	}
	if err := opts.Bind(cfg); err != nil {
		log.Fatal(err)
	}

	// Basic dispatcher as long as we're smol:
	switch {
	case cfg.ListTools:
		ListTools()
	case cfg.ShowTool:
		ShowTool(cfg.Tool)
	case cfg.RunTool:
		RunTool(cfg.Tool, cfg.Input)
	default:
		Default()
	}

}

func Default() {
	// TBD what the default's gonna be.
	fmt.Println("Hello Greenhead")
}

// ListTools prints all registered tools to standard output, or "<no tools>".
func ListTools() {
	names := registry.Names()
	if len(names) == 0 {
		fmt.Println("<no tools>")
	}
	for _, name := range names {
		fmt.Println(name)
	}
}

// ShowTool shows tool help.
func ShowTool(name string) {
	t := registry.Get(name)
	if t == nil {
		// TODO (arguably) exit nonzero here
		fmt.Println("<not found>")
		return
	}
	fmt.Println(t.Help())

}

// RunTool runs a tool with JSON input and prints its JSON output.
func RunTool(name string, input string) {
	t := registry.Get(name)
	if t == nil {
		// TODO (arguably) exit nonzero here
		fmt.Println("<not found>")
		return
	}

	// TODO: we will need this anyway sooner rather than later: run with JSON
	//
	// So we need to marshal JSON into T or fail, but it has to work if valid.
	fmt.Println("TODO")
}

// Usage returns the Usage string based on UsageSource with:
//
// * <name> replaced by Name
// * <title> replace by Title.
// * <tools> replaced by registered tools
func Usage() string {
	s := strings.ReplaceAll(UsageSource, "<name>", Name)
	s = strings.ReplaceAll(s, "<title>", Title)
	s = strings.ReplaceAll(s, "<tools>", registry.Display())
	return s
}
