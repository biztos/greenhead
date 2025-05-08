# Greenhead - Agentic LLM Runner

Greenhead is a framework for building and running AI Agents. It is a Go
package, and a CLI, and a way of making custom CLIs.

There are four main ways to use Greenhead:

1. Run the command-line program to use any built-in or external tools.
2. Easily build your own program with any tools you care to define.
3. Use the Go packages to build anything you want.
4. Contribute tools for others to use.

Greenhead wants these use-cases to be as easy as possible while still being
powerful enough to build Skynet. _(Just kidding?)__

## Motivation

Many people have tasks they would like AI assistance with, but connecting the
tools for Agents is too difficult.

Some of these tools can be defined as functions; others already exist as
binaries on the user's system.

Greenhead presents a uniform way of setting up both kinds of tools and
presenting them to the AI. This, we hope, will make adoption of Agentic AI
easier in situations where programming from scratch is not realistic, but
full control of the environment is required.

## Running the Command-Line Program

The standard CLI as shipped is named `ghd` and includes a standard `--help`
option (also as a "help" subcommand at any level).  In addition, detailed
documentation is available via the `doc` subcommand:

    ghd --help
    ghd config --help
    ghd doc config

## Building Your Own

You can build your own version of `ghd` with as few as N lines of Go:

### Easy Mode

### Power Mode

## Using the Packages

## Contributing Tools

## Acknowledgements

Thanks to Sam York for brainstorming this into liminal existance in Asia.

