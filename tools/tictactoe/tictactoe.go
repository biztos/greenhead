// Package tictactoe models (very simply) the game of Tic Tac Toe.
package tictactoe

import (
	"context"
	"fmt"
	"strings"

	"github.com/oklog/ulid/v2"

	"github.com/biztos/greenhead/registry"
	"github.com/biztos/greenhead/tools"
)

// NewGameInput represents the input to a new Game.
type NewGameInput struct{}

// MoveInput represents a move in a Game.
type MoveInput struct {
	GameId string `json:"game_id"`
	Row    int    `json:"row"`
	Col    int    `json:"col"`
	Player string `json:"player"`
}

// MoveResult represents a non-error move result in a Game.
type MoveResult struct {
	Result string `json:"result"`
	Board  string `json:"board"`
}

// Game represents a game of Tic Tac Toe on a 3x3 grid of cells,
// each of which can have a value of "X", "O", or "-" (unoccupied).
type Game struct {
	Id          string `json:"game_id"`
	last_player string
	cells       [][]string
}

// NewGame returns a new game with the grid all set to "-".
func NewGame() *Game {
	cells := [][]string{
		{"-", "-", "-"},
		{"-", "-", "-"},
		{"-", "-", "-"},
	}
	return &Game{
		Id:    ulid.Make().String(),
		cells: cells,
	}
}

// Board draws the game board in text.
func (g *Game) Board() string {
	s := ""
	for _, r := range g.cells {
		s += strings.Join(r, " ") + "\n"
	}
	return s
}

// State returns the game state, as one of:
//
// - "X won" if X has won.
// - "Y won" if Y has won.
// - "Stalemate" if the board is full.
// - "OK" if the game can continue.
//
// TODO: make Stalemate work for either player or both, and not just do the
// full board since that's (obviously) not the only stalemate!
func (g *Game) State() string {
	// Brute-force the wining state, since there are only 8 possibilities.
	for i := 0; i < 3; i++ {
		if g.cells[i][0] == g.cells[i][1] && g.cells[i][1] == g.cells[i][2] {
			if g.cells[i][0] != "-" {
				return fmt.Sprintf("%s won", g.cells[i][0]) // row match
			}
		}
		if g.cells[0][i] == g.cells[1][i] && g.cells[1][i] == g.cells[2][i] {
			if g.cells[0][i] != "-" {
				return fmt.Sprintf("%s won", g.cells[0][i]) // col match
			}
		}
	}
	if g.cells[0][0] == g.cells[1][1] && g.cells[1][1] == g.cells[2][2] {
		if g.cells[0][0] != "-" {
			return fmt.Sprintf("%s won", g.cells[0][0]) // diagonal L-R
		}
	}
	if g.cells[2][0] == g.cells[1][1] && g.cells[1][1] == g.cells[0][2] {
		if g.cells[2][0] != "-" {
			return fmt.Sprintf("%s won", g.cells[2][0]) // diagonalR-L
		}
	}

	// How to determine stalemate?  Meaning no winnable move.  Depends on who
	// has the next move.  For starters we just check if it's full.
	for r := 0; r < 3; r++ {
		for c := 0; c < 3; c++ {
			if g.cells[r][c] == "-" {
				return "OK" // except where actually not, cf TODO above.
			}
		}
	}
	return "Stalemate"

}

// Move attempts to make a move for player "X" or "Y".
//
// If the move is allowed, State is called and the result returned.
func (g *Game) Move(row, col int, player string) (string, error) {

	if player != "X" && player != "O" {
		return "", fmt.Errorf("player must be either X or O")
	}
	if row < 0 || col < 0 || row > 2 || col > 2 {
		return "", fmt.Errorf("out of bounds: Tic Tac Toe is played on a 3x3 grid")
	}
	if g.last_player == player {
		return "", fmt.Errorf("it is not %s's turn", player)
	}
	taken := g.cells[row][col]
	if taken != "-" {
		return "", fmt.Errorf("that cell is already taken by %s", taken)
	}
	g.cells[row][col] = player
	g.last_player = player

	return g.State(), nil

}

func init() {
	// Wrap the functions.
	games := map[string]*Game{}
	new := tools.NewTool[NewGameInput, *Game](
		"tictactoe_new_game",
		"Starts a new game of Tic Tac Toe.",
		func(ctx context.Context, in NewGameInput) (*Game, error) {
			g := NewGame()
			games[g.Id] = g
			return g, nil
		},
	)
	move := tools.NewTool[MoveInput, MoveResult](
		"tictactoe_move",
		"Make a move in a Tic Tac Toe game, occupying a cell for player X or O.",
		func(ctx context.Context, in MoveInput) (MoveResult, error) {
			g, have := games[in.GameId]
			if !have {
				return MoveResult{}, fmt.Errorf("game not found for game_id %q", in.GameId)
			}
			s, err := g.Move(in.Row, in.Col, in.Player)
			if err != nil {
				return MoveResult{}, err
			}
			return MoveResult{
				Result: s,
				Board:  g.Board(),
			}, nil
		},
	)

	if err := registry.Register(new); err != nil {
		panic(err)
	}
	if err := registry.Register(move); err != nil {
		panic(err)
	}
}
