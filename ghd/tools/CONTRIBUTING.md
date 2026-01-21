# Contributing Tools

Tool contributions are welcome, however they must meet the requirements listed
below in order to be accepted into the main project.

If you have an *idea* for a tool, but are not an experienced Go programmer,
feel free to open an issue in GitHub and we can consider it.

## One subpackage per group of tools, in a like-named directory.

If you have a set of tools to manage flubber, they should be in the `flubber`
directory under `tools`.  If there are tools for managing underwater flubber,
and these are distinct from outer-space flubber tools, then they should be in
`tools/flubber/underwater` and `tools/flubber/outerspace` respectively.

## Namespace the tools to match the subpackage directory.

When registering the tools each tool must be prefixed with its subpackage
names.  In the above examples, the `Launch` tool from `outerspace` would be
registered under the name `flubber_outerspace_launch`.

## Descriptions should make sense to laypersons.

Tool descriptions may be quite detailed and technical, but they should also
make it clear and obvious to a layperson what the tool does.

Risks must be identified in the description.

## Include a *working* sample agent config file named `agent.toml`.

It must be possible for anyone with an appropriate API key to run the agent
using that config and see the tools in action, using prompts shown in the
README file.

## Include working prompts in `README.md`.

Where possible these should be prompts that can execute in a single run with
the provided agent config, but where that is not possible, there must be clear
instructions on how to prompt the tool executions from a chat session.

## Test coverage must be 100%.

*Thou Shalt Build For Testing.*

## The tool must not be egregiously dangerous nor illegal.

Anything likely to cause harm will not be accepted, nor will anything that has
or appears to have illegal or harmful intent.

Since one of the original inspirations of Greenhead is to do "crypto stuff,"
this bears some further explanation.  If you want to have the AI do crypto
trades for you, we don't object; but we're also not going to build in any
wallet-emptying tools.

In practice, for anything rejected here on safety grounds, you are probably
better off defining it as an external tool.  That way you have more control
over potential side effects than with built-in tools.

