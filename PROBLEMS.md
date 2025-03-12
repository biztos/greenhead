# PROBLEMS

## How to structure tool output?

GPT at least just wants a fucking string.

## How to instantiate new objects on the fly?

Want to be able to do something like:

```
new_foo() --> id
```

...and then have foo's func available where they maybe were not before.
Sane? Messed up?  Doable but...

## How to keep track of e.g. turns?

The tic tac toe will need this.  Call "clear" and it says it's not your turn.

## Need an abstraction layer so not always using openai API

E.g. Anthropic discourages it, and Claude is quite useful!

(Groq and Ollama are mostly compatible.)

## Not every LLM needs the same tools!  So register how?

Global registration doesn't really make sense then?  Or does it?

Case A: simplest use.  Want all tools available everywhere.

So just register at init the way it does now.

Case B: different tools each bot, arguably same-named.

Ooh, much harder.

So e.g. you might be running two bots in parallel and each one has a
"fire(args)" tool and you want them to be completely different functions.

What about then aliasing the tools?  Could do that.

Sure, give option for tool alias in bot config.

## What is configurable and how much?

