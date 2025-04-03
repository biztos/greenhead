# Useful Commands for Development and Testing

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




