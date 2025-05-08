# IDEAS

## Ship default implementations of e.g. Claud's BASH support.

This is scary?  Yeah, kinda scary!

Any way to run it safely?

## Broker

A broker could be really cool.

Imagine you set up a chat in which you can answer different people, not only
one chat partner.  Will this work?

And you can direct the conversation.

So you can say: `@bob what's today's color?`  And then Bob might ask Jim for
the data, and Jim calls a function.

We could have a broker sitting in the middle.  Broker knows which Agent is for
which name.

But what happens if multiple Agents are talking to one Agent at the same time?

You'd have to make it turn-based somehow, maybe?

Well anyway -- two is plenty to start.

## Database (et al) interfaces.

No reason we can't just let the AI run SQL against some database, but how to
handle defining the tool?

Could do something like `GHD_SQLITE_DB=/path/to/file.db` and if that is not
defined in the tool's init, the tool isn't registered at all.

(This opens a can of interesting worms: might be nice to have ENV switches
for some tools, but how do we make this sufficiently obvious when people are
setting up their environments?)

## Web3 stuff.

Need somebody who knows (and wants support for) Web3 to talk this one up, but
part of the original Greenhead idea was "let crypto bros do crypto things."

