# Greenhead HTTP API

_Work in progress!_

## Rough Plan

Want to run chats with context memory living on the server.

Want some kind of external API keys for users.

Want to BYO API Key for users, so you could have a tool-running API and
users have to provide their own `OPENAI_API_KEY` or equivalent.

Want to *maybe* allow custom agents from user. Any good reason why to do this?
Makes it easy to experiment.  Don't have to restart the server.

Let's see this in WTF, how would it work?

So the API is run from a runner, and it knows the runner's agents, and if you
ask for a new session with an agent it gives you a clone of that agent. OK.

Can we add agents? Seems like we should be able to do that.

I like the idea that an agent is *not* tied to the runner per se, even if it
is cloned from the runner's agents in many cases.  Thus we can allow a user to
create an agent and if it's valid it runs for them but not for others, others
have just the agents from the config (or have a publish-agents permission).

This really should be backed by a database.  But also should work in-memory if
you don't care!

```
api, err := NewApi(r)
if err != nil {
	return err
}
if err := api.Serve(); err != nil {
	return err
}
```

hm

api genkeys --> generate a set of sample keys for regular user, admin, et al.

## PROBLEMS

- Want to require auth keys, but also want to get them from somewhere!
	- Read from a file maybe? If file configged then yes, to add you just append.

## Endpoints

Unless the server is configured with the `NoKeys` option, most endpoints
require a bearer auth header.

### GET /v1/agents/list

List the named agents available for use.

	Authorization: Bearer <key>
	Returns <agent-list> struct.

### POST /v1/agents/new

Create an agent (clone from the runner's agents by name).

	Authorization: Bearer <key>
	Payload:
	{
		"agent": "<name>"
	}
	Returns:
	{
		"id": "<agent_id>",
		"name": "<agent_name",
		"description": "<agent_description>"
	}

### POST /v1/agents/<id>/chat

Send a chat completion prompt to an agent.
Context is maintained on the server.

	Authorization: Bearer <key>
	Payload:
	{
		"prompt": "<user_prompt>"
	}
	Returns:
	{
		"content": "<completion_text>",
		"tool_calls": [<tool_calls>]
	}

TODO: make `tool_calls` subject to permission or config.

### POST /v1/agents/<id>/completion

Send a chat completion prompt to an agent, returning the full response.

The response can get quite long; unless the client plans to dig in for
debugging, it is better to use the chat endpoint.

	Authorization: Bearer <key>
	Payload:
	{
		"prompt": "<user_prompt>"
	}
	Returns:
	{
		"finish_reason": "<string>",
		"content": "<completion_text">,
		"tool_calls": [<tool_calls>],
		"usage": [<usage>],
		"raw_completions": [<raw_completions>]
	}

### POST /v1/agents/<ulid>/end

End a conversation with an agent, making the agent unavailable for chat.
If no data store is in use, this frees the agent's memory.

Inactive agents may be reaped, subject to runner configs.

	Authorization: Bearer <key>
	Returns success.

### POST /v1/agents/create *LOW-PRIORITY, SPECULATIVE*

Create an agent from a config. Permission-based.

	Authorization: Bearer <key>
	Payload:
	{
		...valid agent config...
	}
	Returns:
	{
		"name": "<new_agent_name>"
	}

### POST /v1/agents/publish *LOW-PRIORITY, SPECULATIVE*

Publish an agent from a config so others may use it. Permission-based.
(Argument against this: have to manage permissions, who gets to use this
agent?  Complicated.)

	Authorization: Bearer <key>
	Payload:
	{
		...valid agent config...
	}
	Returns:
	{
		"name": "<new_agent_name>"
	}
