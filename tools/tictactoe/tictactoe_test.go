package tictactoe_test

// TODO: FIX STUFF HERE AND FINISH COVERAGE, messy messy! KFKF!

import (
	"testing"

	"github.com/stretchr/testify/require"

	// "github.com/biztos/greenhead/runner"
	ttt "github.com/biztos/greenhead/tools/tictactoe"
)

func moveOk(r *require.Assertions, g *ttt.Game, row, col int, player string) {
	err := g.Move(row, col, player)
	r.NoError(err)
	r.Equal("OK", g.State())
}

func TestGameWinRows(t *testing.T) {
	require := require.New(t)
	for r := 0; r < 3; r++ {
		// O takes another row
		rr := 0
		if r == 0 {
			rr = 1
		}
		g := ttt.NewGame()
		moveOk(require, g, r, 0, "X")
		moveOk(require, g, rr, 0, "O")
		moveOk(require, g, r, 1, "X")
		moveOk(require, g, rr, 1, "O")
		err := g.Move(r, 2, "X") // the winning move
		require.NoError(err)
		require.Equal("X won", g.State())
	}
}

func TestGameWinCols(t *testing.T) {
	require := require.New(t)
	for c := 0; c < 3; c++ {
		// O takes another col
		cc := 0
		if c == 0 {
			cc = 1
		}
		g := ttt.NewGame()
		moveOk(require, g, 0, c, "X")
		moveOk(require, g, 0, cc, "O")
		moveOk(require, g, 1, c, "X")
		moveOk(require, g, 1, cc, "O")
		err := g.Move(2, c, "X") // the winning move
		require.NoError(err)
		require.Equal("X won", g.State())
	}
}

func TestGameWinDiagLR(t *testing.T) {

	require := require.New(t)
	g := ttt.NewGame()
	moveOk(require, g, 0, 0, "X")
	moveOk(require, g, 0, 1, "O")
	moveOk(require, g, 1, 1, "X")
	moveOk(require, g, 0, 2, "O")
	err := g.Move(2, 2, "X") // the winning move
	require.NoError(err)
	require.Equal("X won", g.State())

}

func TestGameWinDiagRL(t *testing.T) {

	require := require.New(t)
	g := ttt.NewGame()
	moveOk(require, g, 2, 0, "X")
	moveOk(require, g, 0, 1, "O")
	moveOk(require, g, 1, 1, "X")
	moveOk(require, g, 1, 2, "O")
	err := g.Move(0, 2, "X") // the winning move
	require.NoError(err)
	require.Equal("X won", g.State())

}

func xxTestGameWinDiagRL(t *testing.T) {

	require := require.New(t)
	g := ttt.NewGame()
	err := g.Move(2, 0, "X")
	require.NoError(err)
	require.Equal("OK", g.State())
	err = g.Move(1, 1, "X")
	require.NoError(err)
	require.Equal("OK", g.State())
	err = g.Move(0, 2, "X")
	require.NoError(err)
	require.Equal("X won", g.State())

}
