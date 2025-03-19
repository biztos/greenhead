# TODO (ordered!)

## FIX DUPE FUCKING TOOL REGISTRATION WTF?

```
go run ./cmd/example tools list
```

## Make log opts work: --debug, --logfile and --silent

Where --silent just devnulls.

Also nice to have a specific log-level option and info is default, but error
level easy to set.

## Clean up runner/agent stuff w/r/t configs, need to take config files.

Ideally want to have an extra set of configs you can set per-agent.

Want a per-app config, a per-agent config, and a per-client extras that can
be anything, inside of per-agent.  Knowing that per-client will need some
special checking TBD... but probably can just r/t the json?

## Set up token limits at Agent level and also in OpenAiClient

## Allow regexen for agent tool lists, want e.g. demo*

Whatever, maybe just say ends-with-* means "this is prefix?"

That's a lot easier so start there.  Want e.g. "demo_*"

## Clean up the chat UI, make it at least somewhat fun with defaults.

Stuff in wtf is a good start.

## Create two-chat setup with tic tac toe as example.

Coordinate how?

Here's one: `ghd team run --config=one.json --config=two.json`

Where you have M:M?  Pair off each possible set and let them run?

InitialPrompt in config?

## Tool for running arbitrary external commands, that could be in config.

Say you have something like

`/opt/bigco/test_server SERVER_NAME`

And you want to make it callable as a tool with just a config.  Do like this:

```[
{
    "name": "test_server",
    "description": "Run tests for the named server.",
    "args": {"hostname":"foobar","iterations":123}, // or specify a schema?
    "argspec": "--i $iterations $hostname",
}
]```

Then it gets called and the result returned.  Easy peasy!  Can do with our
own thing as an example: do Echo with a bunch of options.

Registering checks that the thing exists and is executable.

If no argspec defined then args are delivered in JSON on STDIN and your tool
can be wrapped to handle that, easy enough, give examples.

This lets you set up agents with no recompile!

__This should be a general type of tool__ i.e. "runtime" or something.

There might be others where you could configure them, so could define
configurable tools as something like

```
type ConfigurableTools struct {
    Name string
    Args any
}
```

...but be careful that the configurable tool "master" can not be called by the
agent itself! Unless you build an "anything runner" in which case it's on you.

Tricky part is how to define the configurable tools... well register configurable
I guess, then include them in the config, and _also_ have an easy way to do it
built-in if we don't want to accept config.

(disabling configs in config also has to be a thing)

## Tools that register new tools.

Dangerous AF.  So include a hard stop on that, have to specifically enable it
in compile.

Something like overriding the cmd for instance.  Go that far and it's on you!

Have to consider __When/how does the agent specify tools in the request?__

It would be cool to allow agents' "skills" to grow, but dangerous -- doing so
would require resetting the tools. Maybe require that this be done explicitly?

## Agent config to listen on HTTP for chat.

Can do very simple version first.

How does this work for say two agents not in the same place?

Option for auth token.

But how to do the read/write part?

Would be nice to see a convo of two agents streaming to local output, but one
agent is elsewhere.

## Disable Configs setup somehow.

Make it easy to say "this config can not be set"? Or make them ignored and we
complain about it?

## Examples with more complex Cobra setup

* Remove a subcommand, e.g. disallow pair.
* Add a subcommand.
* Change name, etc.


## UNORDERED BACKLOG

TBD

