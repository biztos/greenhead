# TODO (ordered!)

## Make log opts work: --debug, --logfile and --silent

Where --silent just devnulls.

Also nice to have a specific log-level option and info is default, but error
level easy to set.

## Clean up runner/agent stuff w/r/t configs, need to take config files.

Ideally want to have an extra set of configs you can set per-agent.

Want a per-app config, a per-agent config, and a per-client extras that can
be anything, inside of per-agent.  Knowing that per-client will need some
special checking TBD... but probably can just r/t the json?

## Convert from docopt to cobra because it's gonna get long w/subcommands.

    ghd tools list
    ghd tools help foo
    ghd tools run foo INPUT
    ghd rools run foo --file=INPUT_FILE
    ghd chat --file=config_file

Then agent, whatever else.

Check is a nice thing to have.  `ghd agent check --config=file.json`

Maybe best to start there b/c can validate configs and so on.

## Set up token limits at Agent level and also in OpenAiClient

## Clean up the chat UI, make it at least somewhat fun with defaults.

Stuff in wtf is a good start.

## Create two-chat setup with tic tac toe as example.

Coordinate how?

Here's one: `ghd team run --config=one.json --config=two.json`

Where you have M:M?  Pair off each possible set and let them run?

InitialPrompt in config?

## UNORDERED BACKLOG

TBD

