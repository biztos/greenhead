# Configuration

Configuration is at two levels: the runner (top level) and the agent. At
startup, the program reads command-line options which take precedence over
config options; then any runner config provided with the `--config` option,
which may itself contain agent configurations; and finally any agent config
files specified with the `--agent` option.

Runner configs take precedence over agent configs.  Consider this simple
example of a `config.toml`:

```toml
# SuperCorp Greenhead Config
log_file = "/var/log/greenhead.log"
max_completions = 100
max_toolchain = 10
[[agents]]
  name = "The Riddler"
  description = "An agent that likes to talk in riddles it gets from its tools."
  type = "openai"
  model = "gpt-4o"
  tools = ["/^get_riddle/"] # include only the various riddle tools.
  max_toolchain = 100
  max_completions = 3
```

Here the agent will be run with `max_completions=100` and `max_toolchain=10`
because the top-level config has those values.

## Runner Configs

### Output Control

### Safety

### Tool Selection

## Agent Configs

