// cmd/example/minimal/main.go -- smallest custom CLI.
package main

import (
	"github.com/biztos/greenhead"
	_ "github.com/biztos/greenhead/tools/tictactoe"
)

func main() {
	// greenhead.CustomApp("minimal", "1.0.0", "SuperCorp Tic Tac Toe Player", "")
	greenhead.Run()
}
