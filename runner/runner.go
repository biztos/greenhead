// Package runner defines the logic to run ghd, the Greenhead CLI tool.

package runner

import (
	"log"
	"os"

	"github.com/biztos/greenhead/registry"

	"github.com/docopt/docopt-go"
)

var Version = "ghd v0.1.0"

var Args = os.Args[1:]

var Exit = os.Exit

var Usage = `Greenhead

Usage:
  ghd [options]
  ghd -h | --help
  ghd --version

Options:
  -h --help           Show this screen.
  --debug             Debug mode (more verbose logging).

This program runs AI agents.

The rest is TBD.

(c) 2025 The Greenhead Authors; distributed under MIT License.

`

type Config struct {
	Debug bool
}

func Run() {

	cfg := &Config{}
	opts, err := docopt.ParseArgs(Usage, Args, Version)
	if err != nil {
		log.Fatal(err)
	}
	if err := opts.Bind(cfg); err != nil {
		log.Fatal(err)
	}
	log.Println("debug", cfg.Debug)

	// OK, let's see what's registered!
	log.Println(len(registry.Names()), "extensions")
	for _, name := range registry.Names() {
		log.Println(name)
	}

	if ext := registry.Get("Demo"); ext != nil {
		log.Println(ext.Description())

	}

}
