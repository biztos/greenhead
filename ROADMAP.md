# Greenhead Feature Roadmap

The plan forward, roughly in order.

## v0.1 MVP for Feedback

_You're soaking in it!_

* Support OpenAI
* Built-in demo tools
* Configurable external tools
* Single-prompt runner (any number of agents)
* Pair runner (two-agent conversation)
* Simple chat
* Built-in agents (simple demonstration agents)
* HTTP API
* Simple chat web UI

## v0.2

* Support for other LLMs besides OpenAI.
	* LLAMA et al, will test via Groq
	* Claude ideally.
	* Others based on feedback.
* Flexibile agents and tools.
* Tool config.
* HTTP API Enhancements
	* BYO API Key
	* Graceful shutdown
* Useful built-in tools
	* pro-bing
	* miekg/dns
	* others TBD
* Improvements/changes based on feedback.

## v0.3

* Add MCP support. (Ouch, probably a lot of work.)
	* May be necessary for support of non-OpenAI LLMs.
* Improvements/changes based on feedback.

## v1.0

* TBD -- the above is a "1.0 of the mind" but need to see what users want!

## FUTURE

* HTTP API for running pairs, sets, etc.
* "Broker" support for many agents.
* Pass-through API? Maybe.
* Model-switcher: let the LLM change the model.
* Tool: more general KV store, maybe Badger?
* Tool: SQLite3 database interface.
* cf TODO.md for more speculative things.
