# Greenhead - Agentic LLM Runner


## *WARNING: this is pre-release software and may change at any time.*

### *NB: examples may only work with OpenAI at this time. More APIs soon!*

Greenhead is a framework for building and running AI Agents. It is a Go
package, a command-line utility, a Web API, and a way of making custom
command-line utilities and Web APIs.

There are four main ways to use Greenhead:

1. Run the `ghd` program to use any built-in or external tools.
2. Easily build your own program with any tools you care to define.
3. Use the Go packages to build anything you want.
4. Contribute tools for others to use.

Greenhead wants these use-cases to be as easy as possible while still being
powerful enough to build Skynet. _(Just kidding?)_

## Motivation

Many people have tasks they would like AI assistance with, but connecting the
tools for Agents is often too difficult or insufficiently secure.

Some of these tools can be defined as functions; others already exist as
binaries or scripts on the user's system.

Greenhead presents a uniform way of setting up both kinds of tools and
presenting them to the AI. This, we hope, will make adoption of Agentic AI
easier in situations where programming from scratch is not realistic, but
full control of the environment is required.

## Running the Command-Line Program

The standard CLI as shipped is named `ghd` and includes a standard `--help`
option which is also available for every subcommand.  In addition, detailed
documentation is available via the `doc` subcommand:

```sh
ghd --help
ghd config --help
ghd doc
ghd doc config
ghd doc api
```

A good starting point is the prompt runner:

```sh
ghd agents run Howdy --log-file=tmp.log --agent=chatty --agent=pirate
```

## Running the Web API

The Web API exposes persistent chat conversations over an HTTP interface. The
server maintains chat context.  Full documentation of the API is available
online at the project [website][ghd] as well as from the `ghd doc api`
command.

To run the HTTP API with two built-in agents available:

```sh
ghd api serve --agent=chatty --agent=pirate
```

Then direct your browser to [localhost:3030](http://localhost:3030/), enter
the temporary API Key printed by the above command, and click `Submit Query`
to play with the primitive Web UI.

Don't forget to kill the server when you're done, John Connor!

## Building Your Own

You can build your own version of `ghd` with very little code.  Here is the
"minimal" example:

```go
package main

import (
    "github.com/biztos/greenhead"
    _ "github.com/biztos/greenhead/tools/tictactoe"
)

func main() {
    greenhead.CustomApp("minimal", "1.0.0", "SuperCorp Tic Tac Toe", "")
    greenhead.Run()
}
```

The `pair run` command could then be used to play a game of Tic Tac Toe
between two OpenAI agents:

```sh
minimal run "Start game." --agent=tictactoe --agent=tictactoe -l tmp.log 
```

The top-level `greenhead` package exposes functions for easy customization of
your app.  In addition to the name, you will usually want to customize the
tools, and sometimes also the command-line options.

For more advanced customization, see the `examples` subdirectory of the source
code.

## Configuring External Tools

In order to expose your own programs as tools available to the LLM, you must
configure them in the main (runner) config.  Consider the UNIX `host` command,
which should be readily available.  Let's assume we want to use it to look up
a host, optionally specifying the type and verbosity.

The usage string for this limited use of the command would be:

```
host [-v] [-t type] {name}
```

The config section would look like this:

```toml
[[external_tools]]
  name = "host"
  description = "Look up an internet host.."
  command = "/usr/bin/host" # or wherever it lives on your system!

  [[external_tools.args]]
    flag = "-t"
    key = "type"
    type = "string"
    description = "Query type: CNAME, NS, SOA, TXT, DNSKEY, AXFR, etc."
    optional = true

  [[external_tools.args]]
    flag = "-v"
    key = "verbose"
    type = "boolean"
    description = "Output verbose query results.."
    optional = true

  [[external_tools.args]]
    key = "hostname" # <-- Note: no flag specified here!
    type = "string"
    description = "Name of the host to look up, e.g. 'google.com'."
    optional = false
```

For more details, see the config [documentation][ghd] or use:

```sh
ghd doc config
```

## Using the Packages

## Contributing Tools

## Further Reading

Additional documentation is available at the project [website][ghd].

## Acknowledgements

Thanks to Sam York for helping brainstorm this into liminal existance in Asia.

Greenhead is built with
[Cobra](https://cobra.dev/),
[Fiber](https://gofiber.io/),
[go-openai](https://github.com/sashabaranov/go-openai),
and other great open-source software packages.


[ghd]: https://ghd.biztos.com/
