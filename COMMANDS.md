# Useful Commands for Development and Testing

## Run a chatty two-agent chat for 6 completions total.

Evens will end with first == second, odds with first > second.  Look for
the original prompt in the context.

```
just run pair run "Hello!" -s -d --agent=testdata/agent-chatty.toml --agent=testdata/agent-chatty.toml --log-file=tmp.log --dump-dir=build --show-calls --max-completions=6
```

## Run a game of tic tac toe.

```
go run ./cmd/ghd pair run "Please start a game of Tic Tac Toe and make the first move."  -s -d --agent=tools/tictactoe/agent.toml --agent=tools/tictactoe/agent.toml --log-file=tmp.log --dump-dir=build --show-calls  --max-completions=10
```


## Prettify the JSON in build

```
find build -type f -name "*.json" -exec sh -c 'jq "." {} > {}.formatted.json' \;
```
