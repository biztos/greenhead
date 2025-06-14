# Greenhead HTTP API

_Work in progress!_

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

## Possible Future Endpoints *LOW-PRIORITY, SPECULATIVE*

### POST /v1/agents/create

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

	Authorization: Bearer <key>
	Payload:
	{
		...valid agent config...
	}
	Returns:
	{
		"name": "<new_agent_name>"
	}
