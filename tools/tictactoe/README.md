# Tic Tac Toe Tool

Get agents to play Tic Tac Toe.

(Also, rough cut at what a README for a tool should look like.)

## Synopsis

In the agent's `config.toml`:

```
tools = ["/tictactoe.*/"]
```

In custom runners:

```
import (
    _ "github.com/biztos/greenhead/tools/tictactoe"
)
```

## Usage

These examples are run from the root directory, with `OPENAI_API_KEY` set.

### Run a simple agent that should start a game:

```
go run ./cmd/ghd agents run -s -d --agent=tools/tictactoe/agent.toml \
--log-file=tmp.log --dump-dir=build --show-calls "start game"
```

### Run a chat in which you can play against the agent to test its responses:

```
go run ./cmd/ghd chat -s -d --agent=tools/tictactoe/agent.toml \
--log-file=x --dump-dir=build --show-calls "start game"
```

The agent's context window must be tuned accordingly; see `agent.toml` in this
directory for an example that *should* work with ChatGPT 4o.

There are two tools the agent should call:

- tictactoe_new_game {}
    - returns a game_id (ULID) which must be used for moves.
- tictactoe_move { game_id: "ULID", row": R, "col": C, "player": "X"|"O" }

## Description

Tic Tac Toe is an attempt at proving two-player agentic competitive behavior
with stateful objects.

It is a work in progress.

## NOTES


## Run the tic tac toe agent to confirm it can start a game.

`go run ./cmd/ghd agents run "start game"  -s -d --agent=testdata/agent-ttt.toml --log-file=x --dump-dir=build --show-calls`

Responses seen so far -- sometimes it moves, sometimes not:

I've started a new Tic Tac Toe game. You can join the game using the ID: 01JQWYA42174NVZQCNMN05D809. I'll make my first move as soon as you join.

---

I've made the first move as X. Here's the current board:

```
X - -
- - -
- - -
```

You need to join the game using the game ID: `01JQWYBXRCJ7Z5814D1BVP3K2P`. It's your turn!




## Run the tic tac toe agent to confirm it can join a game -- which attempt will fail b/c no such game!



 just run chat -s -d --agent=testdata/agent-ttt.toml --log-file=x --dump-dir=build --show-calls
 9822  2025-03-31 00:25:20 go run ./cmd/ghd run "start game"  -s -d --agent=testdata/agent-ttt.toml --log-file=x --dump-dir=build --show-calls



