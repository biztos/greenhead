![provisional greenhead mascot](ghd/assets/src/webui/greenhead-150x225.png)

# Greenhead - Agentic AI Framework in Go.

<!-- cut -->

## 📢 Important Note on Project Status

Greenhead still has a lot of promise, because -- AFAICT -- it's the only
framework that allows for config-only integration with your legacy
command-line tools while also allowing for full customizability of the
binary.

That said, Your Humble [Author][biztos] has other 🐟 to 🍳 at
the moment, and further work on Greenhead is likely to be sporadic at best,
unless someone forks it or [hires me][hireme] to customize
the software to their needs.

In particular, considerations include:

* Finding a more flexible LLM API abstraction similar to [Vercel AI][vercel-ai]
* Possibly moving to [MCP][mcp] (or maybe not; function calls are a good fit here)
* Rewriting the crappy web UI to be modern.
* Fix the access API, it's a mess.
* Surely other things besides!

Read on if you'd like to play with this or [copy][LICENSE] any of the ideas
herein.

---

Greenhead is a framework for building and running AI Agents. It is a Go
package, a command-line utility, a Web API, and a way of making custom
command-line utilities and Web APIs.

There are four main ways to use Greenhead:

1. Run the `ghd` program to use any built-in or external tools.
2. Easily build your own program with any tools you care to define.
3. Use the Go packages to build anything you want.
4. Contribute tools for others to use.

Greenhead wants these use-cases to be as easy as possible while still being
powerful enough to build [Skynet][skynet].

## Motivation

Many people have tasks they would like AI assistance with, but connecting the
__tools__ for Agents is often too difficult or insufficiently secure, and/or
comes with framework lock-in.

Some of these tools can be defined as native functions; others already exist
as binaries or scripts on the user's system.

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

## Configuring External Tools

In order to expose your own programs as tools available to the LLM, you must
configure them in the main (runner) config.  Consider the UNIX `host` command,
which should be readily available.  Let's assume we want to use it to look up
a host, optionally specifying the type and verbosity.

The usage string for this limited use of the command would be:

```sh
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

## Building Your Own

You can build your own version of `ghd` with very little code.  Here is the
"minimal" example:

```go
package main

import (
    "github.com/biztos/greenhead"
    _ "github.com/biztos/greenhead/ghd/tools/tictactoe"
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

For more advanced customization, see the `ghd/examples` subdirectory.

## Using the Packages

If your goal is to incorporate the agent-runner logic into your own project,
you will eventually need to use more than the top-level `greenhead` package.

To copy the full run cycle, especially regarding config management, you will
want to examine the `ghd/runner` subpackage.  To manage agents and tools
directly, the `ghd/agent`, `ghd/tools`, and `ghd/registry` packages should be
consulted.

Subpackages in a nutshell, all under `ghd/`:

* agent - agent logic and API clients.
* api - the HTTP API; __needs work__
* assets - assets built from the `src` subdir with [binsanity][binsanity].
* cmd - the [Cobra][cobra] setup for the CLI; uses `runner` for command logic.
* registry - global registry of tools.
* rgxp - optional-regexp format for some config values. _NB: may be factored out!_
* runner - the command runner, including top-level config logic.
* tools - tools types and logic; subdirs contain the built-in tools.
* utils - misc utils; _may be factored out at some point_.
* version - canonical version numbers.

## Contributing Tools

Tool contributions are very welcome!  Pull requests will be evaluated first
for safety and utility: it should be impossible for naïve users to cause
damage, and the tool should be useful to a reasonable number of people.
(Entertainment counts as "useful" -- games are encouraged!)

Please be sure that:

1. The tool name follows the namespacing conventions.
2. Descriptive text and documentation is (also) in English.
3. A working `agent.config` is included.
4. The `README.md` is easy for a layperson to understand.

For a canonical example, see the `ghd/tools/tictactoe` directory.

If you have an idea for a useful tool but are not a Go programmer, feel free
to open a ticket describing your idea.  Someone else may find it compelling
and code it for you.

### Namespacing

Tools may only have simple ASCII names like `mytool_zebra_renoberator`.
The name should follow the subdirectory structure:

```sh
tools
└── mytool
    └── zebra
        └── renoberator
            ├── README.md
            ├── agent.toml
            ├── renoberator.go
            └── renoberator_test.go
```

A single subpackage may add multiple tools.  For instance the `renoberator`
package might define three tools:

* `mytool_zebra_renoberator_obtain`
* `mytool_zebra_renoberator_stripe`
* `mytool_zebra_renoberator_destripe`

These should be clearly described in the README file.

Especially for higher-level packages, avoid names that might conflict with
future implementations at a lower level.  The bias is for longer and more
descriptive tool names.

### External Binaries

Submitted tools may call external binaries, within reasonable safety bounds.

Such tools _must_ check for the presence of the binaries before registering.

For some cases the easiest method is to wrap a standard `ExternalTool` in the
normal tool packaging.  This is legitimate, but the tool name must reflect any
special usage or restrictions, to differentiate itself from other wrappings of
the same binary.

For example a `host` command wrapped to take a `-t` option but nothing else
might be named `unix_host_lookup_by_type`.

## Acknowledgements

Thanks to Sam York for helping brainstorm this into liminal existance in Asia.

Greenhead is built with
[Cobra][cobra],
[Fiber][fiber],
[go-openai][go-openai],
and other great open-source software packages.

[ghd]: https://ghd.biztos.com/
[binsanity]: https://pkg.go.dev/github.com/biztos/binsanity
[cobra]: https://cobra.dev/
[fiber]: https://gofiber.io/
[go-openai]: https://github.com/sashabaranov/go-openai
[vercel-ai]: https://ai-sdk.dev/docs/introduction
[biztos]: https://biztos.com
[hireme]: https://biztos.com/ghcv/
[mcp]: https://modelcontextprotocol.io/docs/getting-started/intro
[license]: ./LICENSE
[skynet]: https://skynetobserver.substack.com/
