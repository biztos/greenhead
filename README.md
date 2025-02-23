# greenhead

LLM thing based on convo with Sam in เชียงใหม่

## Pluggable Headless Agentic AI

Basic idea is:

* You can plug in functions that are callable by the AIs.
* You can run "headless" -- put the AI in a simple endpoint.
* You can connect two (or more?) to each other.

(Maybe "headless" isn't the best description for that?)

## Demo MVP Features

Just to prove it out, seems like we'll need:

1. Pluggable, i.e. you can add functions without touching other code.
2. Useful set of "toy" functions to start with -- things easy to demo.
3. Simple CLI so you can play with the thing.
4. Simple HTTP API.  (TBD -- not just re-exposing the OpenAI one.)
5. Simple HTML wrapper for that so you can 
6. Ability to connect two agents to each other.

Also very cool but not really "minimal" -- ability to have more than the two,
i.e. a converstation, potentially brokered?

## Difficulties

### Plug-ins

First hard thing is the plug-in thing considering no dynamic loading in Go.
Probably acceptable to have a registry of these which is compiled and then
set up which funcs are available, with default rigging to do that via a TOML
file.

### API

Next obviously hard thing is the API.  What exactly are we exposing? Consider:

- Setup: have func `color_for_day(d) -> color`
- User: What is today's color?
- LLM:
    1. "Calling the function."
    2. call func(today)
- System: func(today) result = "blue"
- LLM: "Today's color is blue."

In that example we do or do not want to expose the LLM-to-system round trip to
the user?  We need to be able to see it in some debug way, but we also want to
be able to keep the agent's activities under wraps.

OK so probably *do* need an API that includes everything, but build in a way
to deal with the function calls.  Problem being that by the time you get the
text in `1.` above, the func calls are already underway.  Streaming not good
in this case.

### Demo Funcs that actually do something.

Ideally they should do something, however useless, when run between two Agents
with no human.

### Broker

A broker could be really cool.

Imagine you set up a chat in which you can answer different people, not only
one chat partner.  Will this work?

And you can direct the conversation.

So you can say: `@bob what's today's color?`  And then Bob might ask Jim for
the data, and Jim calls a function.

We could have a broker sitting in the middle.  Broker knows which Agent is for
which name.

But what happens if multiple Agents are talking to one Agent at the same time?

You'd have to make it turn-based somehow, maybe?

Well anyway -- two is plenty to start.

