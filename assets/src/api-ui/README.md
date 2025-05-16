# ghd-web -- a crack at a VERY minimal web UI for greenhead chats

Package it all up into one simple SPA and serve under, say, the API itself!

	/v1/ui

Sure, why not!  It'll be dirt simple anyway.  But if doing this we'll need a
first GET that sets up the Key first, which is then sent to the other bit
that will seed it.  Because without a key you *might* get fuck-all.

Ugh, need to serve from same thing, not other command, because same-site.

```

	/v1/ui -->

	root.html -- super-basic root to get you started with a key.

	This goes in a POST to same.  If no Agent in header, goes to agent select.

	Aha, OK, same UI does agent select: no agent means start from new.


	For starters serve the js file separately but then bundle it.

	ghd.html -- the main SPA (bundle in ghd.css and ghd.js, one request)

```