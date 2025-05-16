# TODO (ordered! sorta!)

## HTTP API

1. Logger setup in Runner not just in Agents.
2. Get Logger working for API (fiber) and Agents, same one.
3. Make sure agents not streaming!
4. Make sure agents not printing, just want res!
5. Add middleware for API keys: also to check agents not just endpoints

fuck fuck fuck, I really want to see the tool calls!

## Serialize agents, duh.

Lots of good reasons to pause/restart an agent, and for that we need to be
able to serialize.  WTF was I thinking?

## Make a "multi-api" type of API Client, e.g. can run OpenAI or Llama.

This would just branch at creation based on what env vars it finds, and use
the other constructors.  *Ideally* this would be the default type for all the
built-in agents.

Maybe just call the type "flex"?

## Better API for clearing/adding named agents from top level.

OK probably not urgently needed.

## Use Glamour for rendering incoming stuff, also for streaming!

Glamour has nice Markdown to ASNI rendering.  Ideally want to back up and
re-render whenever we have both:

* reached the end of a block we think is useful
* finished getting a chunk to render

One problem is keeping track of the stuff we already printed, e.g. by lines.
(Fun problem but what priority?)

Other problem is what represents a point at which to re-render -- and should
we maybe add stuff as temporary markdown?

For instance we get a code block, it's nice to start rendering it as code.
So add on a closing marker.  Same maybe true of other things?  Unsure.

Printing the streaming response is neat but this would be neater.

__Starting point is maybe using this for non-streaming output__ and then...
what about the print funcs? Could just start with render to ASCII and then
print in color, could work.

Those change to styles maybe?  Nice to keep the colors, I like that. But then
I have a whole fucking color scheme for keeping the output aligned with the
original color... yikes.

__THIS IS A RABBIT HOLE AND YOU HAVE MORE IMPORTANT SHIT TO WORRY ABOUT__

## Make "tools run" take an arg instead of stdin

Stdin is stupid, let it take an arg but file-style like @foo, same logic is
in agents run.

## Make sure all the runner commands take io.Writer first

Otherwise shitshow.


## DECIDE: open up or lock down the Agent?

Currently half-half.  Probably lock down, make it opaque.

But it needs to be possible to make one with a custom client.


## Run a script that can do.... whatever.... with the responses.

That is to say set up another program (or JS/etc code at some point?) that
will interact with an agent as the other end of a chat.

Most of this logic will be same as using pair chat but with one half of the
pair being outside (or inside an interpreter?).

__How is this different than running a server?__

Basically it's the same logic as running an agent over HTTP?

__USE CASE IS WHAT?__ maybe just reformat stuff?

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

## Limit tools at runner config level.

This avoids having to do it at the agent level.  Just override everything if
tools are set.

Probably rather two:

* allow_tools -- these tools can be used, no others can be used.
    * impl: make good/bad list, bail if:
        * anything from good list missing
        * anything on both good and bad list
        * then go remove from bad list
* remove_tools -- just remove the named tools, bail if not exist.
* agent_tools -- override all agents, everybody has these tools.
* can do any combo; but run AFTER any configged tools are read in.

## UNORDERED BACKLOG

### Support comments in config dumps.

Nice to have!  But by no means urgent.  And if doing it, make it optional.

github.com/pelletier/go-toml/v2

## MAYBE Default config file w/complex logic

No hurry here, *really* no hurry!  Normal use is probably just to define
agents on the command line.

__PROBLEM:__ if you say default is /foo.toml and want to ignore it if not
found, you have to handle --config=/foo.toml which *should* throw if not
found.

Also problem: no idea where to look for that file, i.e. we are not saying
(nor want to say) where you should put the binary.

One solution is to have a settable default but not use it out of the box.
Custom binaries only.  But then can't easily have it in the flag settings b/c
that's in init and the var is used then.

## MAYBE tools.ExternalToolArg.Connector

We had an `ExternalToolArg.Connector` field to cause options to be coupled
with their values, as in --foo=RECEIVED_FOO_VAL.  The need for this is
speculative though, for POSIX we should not have to worry about it.

__DO NOT DO THIS WITHOUT A REAL-WORLD USE-CASE__

