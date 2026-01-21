# Tic Tac Toe Tool

Tic Tac Toe is an attempt at proving two-player agentic competitive behavior
with stateful objects.

It is a work in progress.

(Also, rough cut at what a README for a tool should look like.)

## Tool Functions

* tictactoe_new_game {}
    * returns a game_id (ULID) which must be used for moves.
* tictactoe_state { game_id: "ULID" }
    * returns the game state consisting of:
        * state: "OK", "X|O won", or "Stalemate"
        * next_player: "X" or "O"
        * board: the board layout as text
* tictactoe_move { game_id: "ULID", row": R, "col": C, "player": "X"|"O" }
    * returns the state as above if the move did not result in error.

## Caveats

0. This has been tested only on OpenAI's "4o" model (not "o4!").
1. Sometimes the agent makes illegal moves; it usually self-corrects.
2. Runaway games are a possibility, however remote. Always set max-completions!
3. The game logic is rather basic.
4. Dead games are not cleared in a runner, nor is there an upper limit on games.
5. This has not undergone *any* adversarial testing!

## Configuration

In the agent's `config.toml`:

```
tools = ["/tictactoe.*/"]
```

In custom runners:

```
import (
    _ "github.com/biztos/greenhead/ghd/tools/tictactoe"
)
```

## Usage

These examples are run from the root directory, with `OPENAI_API_KEY` set.

### Run a simple agent that should start a game:

```
go run ./cmd/ghd agents run -s --agent=tools/tictactoe/agent.toml \
--show-calls "start game and make the first move"
```

This is the easiest way to show that the agent is able to call the tools.

### Run a chat in which you can play against the agent to test its responses:

```
go run ./cmd/ghd chat -s --agent=tools/tictactoe/agent.toml \
--log-file=tmp.log --show-calls "start game"
```

### Have two (identical) agents play against each other.

It is assumed that ten completions should be enough.  If not, tune as needed.

```
go run ./cmd/ghd pair run -s --show-calls --log-file=tmp.log \
--agent=tools/tictactoe/agent.toml --agent=tools/tictactoe/agent.toml \
"Please start a game of Tic Tac Toe and make the first move." \
--max-completions=10
```


