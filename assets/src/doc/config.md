# Configuration

All configuration files are in TOML format.

The top-level configuration is referred to as the _runner config_ and is the
only config required as it can contain agent configs, as well as configuration
of the HTTP API.  Agent configs may also be specified individually.

At startup, the program reads command-line options which take precedence over
config options; then any runner config provided with the `--config` option,
which may itself contain agent configurations; and finally any agent
configurations specified with the `--agent` option.  Agent configurations can
be built-in agent names or agent config files.  For a list of built-in agents
use the `agents builtin` command.

Runner configs take precedence over agent configs of the same name.  Consider
this simple example of a `config.toml`:

```toml
# SuperCorp Greenhead Config
log_file = "/var/log/greenhead.log"
stream = false
max_completions = 100
max_toolchain = 10
[[agents]]
  name = "The Riddler"
  description = "An agent that likes to talk in riddles it gets from its tools."
  type = "openai"
  model = "gpt-4o"
  tools = ["/^get_riddle/"] # include only the various riddle tools.
  stream = true
  max_toolchain = 9999
  max_completions = 3
```

Here the agent will be run with `stream = false`, `max_completions = 100` and
`max_toolchain = 10` because the top-level config has those values.

Note that zero-values do *not* override.

## Runner Configs

### Output Control

### Safety

### Tool Selection

### API Config

## Agent Configs

