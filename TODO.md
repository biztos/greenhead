# TODO (ordered!)

## Run a tool from JSON string input

Marshal or fail; write standard-format output { result, error }

Gonna need some output type guessing I think.  Something we could control in
the runner config at some point.

- error: "ERROR: err"
- string: string
- array of strings: joined "\n"
- array of anything else? sprintf as value or what?  TBD

Anything object-y gets JSONned.  Or just everything?  Well not straight prim.

## Run a r/t to prove the tool format is right for OpenAI

Need to make sure the tools will be called correctly before going further down
this path.

## Make a CLI chat i/o thing with streaming

Gonna need it sooner or later anyway, should be Claude-able.

Main point is I really want to see stuff happening real time.  And have some
colors as an option.

Bonus: up-arrow for history, special chars or :help or something.

## Make that run the function calls (at first don't do that b/c complicated).

At this point it's ready to demo to Sam I think.

## Convert from docopt to cobra because it's gonna get long w/subcommands.

    ghd tools list
    ghd tools help foo
    ghd tools run foo INPUT
    ghd rools run foo --file=INPUT_FILE

Then agent, whatever else.

## UNORDERED BACKLOG

TBD

