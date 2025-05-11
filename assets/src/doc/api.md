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

## Endpoints

### POST /v1/new

Create an agent (clone from the runner's agents by name).

	Auth: <key>
	{
		"agent": "<name>"
	}
	Returns <ulid> string.

### POST /v1/agent/<ulid>/chat

Send a chat completion prompt to an agent.
Context is maintained on the server.

	Auth: <key>
	{
		"prompt": "<user_prompt>"
	}
	Returns <response> struct.

### POST /v1/agent/<ulid>/end

End a conversation with an agent, making the agent unavailable for chat.
If no data store is in use, this frees the agent's memory.

Inactive agents may be reaped, subject to runner configs.

	Auth: <key>
	{
		"reason": "<any_string>"
	}
	Returns success.

### POST /v1/agent/create *LOW-PRIORITY, SPECULATIVE*

Create an agent from a config. Permission-based.

	Auth: <key>
	{
		...valid agent config...
	}
	Returns <ulid> string.

### POST /v1/agent/publish *LOW-PRIORITY, SPECULATIVE*

Publish an agent from a config so others may use it. Permission-based.

	Auth: <key>
	{
		...valid agent config...
	}
	Returns <ulid> string.
