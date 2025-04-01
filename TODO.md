# TODO (ordered!)

## PAIR CHAT

Need this to feel excited, too much redesign happening right now.

Pre-task: __make limits work__ b/c could get runaway and cost $$.

Basic MO is like...

`ghd pair run "FIRST PROMPT"`

## Make sure logging to file works OK for multiple agents

Might not?  Multiple fh's open is not the right way to go, what else can we
do?

Probably: prove a fail first.  Then try to define logfiles in the runner.
Though it would be better to NOT have that specific dep.

__Prove this with a test__

## Make the rest of runner stuff work like chat.

## Support tools defined at load in config (or theoretically at runtime)

Basic version is a tool that calls a binary with args using os.Exec with a
timeout. (Or with no timeout?  Timeout might be nice.)

Let that be set in runner config which gets tools, agent can then be configged
to call it.

Bigger/more-complex idea: run code in Javascript or maybe other interpreter.

Goja for JS would be great.  What else?  If it's easy can just add...

Any way to make that pluggable?  Should work actually.  Want a type that is
"code interpreter" and then the concrete type is say JS interp built-in, can
add others in the standard way.

These two should both be supersets of Tooler!

Another one, kinda scary but why TF not? An actual code interpreter the agent
can call as a tool.  Prolog maybe: https://github.com/ichiban/prolog

(Or just reuse the JS thing above.)

__FOR NOW, PUNT ON SELF-DEFINING TOOLS BUT DO SUPPORT CONFIGURED TOOLS__

Make it something you have to turn on, or something you can turn off?

## Run a script that can do.... whatever.... with the responses.

That is to say set up another program (or JS/etc code at some point?) that
will interact with an agent as the other end of a chat.

Most of this logic will be same as using pair chat but with one half of the
pair being outside (or inside an interpreter?).

__How is this different than running a server?__

Basically it's the same logic as running an agent over HTTP, except the

## Set up token limits at Agent level and also in OpenAiClient


## Clean up the chat UI, make it at least somewhat fun with defaults.

Some ideas:

```
/spool FILE --> start spooling to FILE, write prompt/content pairs only
    detect json, otherwise spool text
/hist --> list prompt history (ergo keep history)
/tools --> list tools (maybe ls for short?)
/call TOOL_NAME [TOOL_ARGS] --> call a tool; sanely expand args for UX
/dump --> dump last interaction to temp file
/!cmd --> run shell command then return
/r!cmd --> run shell command and take stdout as prompt
/ed --> edit current prompt in $EDITOR
/q --> quit
/c --> run agent check command
/logs --> run logs thru $PAGER
```

Basic command structure should be "^/" == command and if you really need to
start a line with "/" you can start with "\/" ...

What about catching ESC instead?  Can do that?

## Namespacing tools?

Nah.  Just whatever you put in the repo has to be "toolname_foo_bar" or I do
not accept it.

Want other people to have total freedom to define "boo" or whatever.

So the TODO here is __Document This__.


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

