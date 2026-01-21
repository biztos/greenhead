// examples/minimal/main.go -- smallest meaningfully custom CLI.
package main

import (
	"github.com/biztos/greenhead"
	_ "github.com/biztos/greenhead/ghd/tools/tictactoe"
)

func main() {
	greenhead.CustomApp("minimal", "1.0.0", "SuperCorp Tic Tac Toe", "")
	greenhead.Run()
}
