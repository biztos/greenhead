# LAYOUT

```

greenhead.go -> top package which can be imported.
                imports the runner logic which is called with Run()
                imports the extension registry logic and inits that


Then to make a new extension or many of them, you do this:

package main

import (
    "github.com/biztos/greenhead"
)

func main() {
    // could make your own binary, or could just...
    greenhead.Run()
}

func init() {
    // register the extensions
    ex1 := NewMyExtension()
    greenhead.RegisterExtension("My Extension", ex1)
}


```