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

## v0.2

* HTTP API
	* Main use-case is other bots e.g. a Slackbot
	* "chat" over HTTP with context on server
	* BYO API Key
	* Secondary API Keys tied to agents
		* ...so key ABC gets agent Foo, BCA gets Bar, and so on.
		* reasonable way to refresh that so e.g. you can assign keys
	* Graceful shutdown
* Useful built-in tools
	* pro-bing
	* others TBD
* Support for other LLMs besides OpenAI.
	* LLAMA et al, will test via Groq
	* Claude ideally.
	* Others based on feedback.
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
