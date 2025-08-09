package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"syscall/js"
)

type Board struct {
	Grid       [8][8]int
	Self       int
	Opponent   int
	ValidMoves [][2]int
	Whites     int
	Blacks     int
	Skipped    int
	Won        string
}

func (board *Board) constructBoard() {
	board.Grid = [8][8]int{
		{0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 2, 1, 0, 0, 0},
		{0, 0, 0, 1, 2, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0}}
	board.Self = 1
	board.Opponent = 2
	board.Whites = 2
	board.Blacks = 2
	board.Skipped = 0
	board.Won = ""
	board.updateValidMoves()
}

func (board *Board) json() string {
	result, err := json.Marshal(board)
	if err != nil {
		return "{}"
	}
	return string(result)
}

func (board *Board) drawBoard() {
	board.updateValidMoves()
	js.Global().Call("updateBoard", board.json())
}

func (board *Board) checkFlankDir(row int, col int, dir [2]int, depth int) bool {
	if row+dir[0] >= 0 && row+dir[0] <= 7 && col+dir[1] >= 0 && col+dir[1] <= 7 {
		// fmt.Printf("%d%d\n", row+dir[0], col+dir[1])
		switch board.Grid[row+dir[0]][col+dir[1]] {
		case board.Opponent:
			return board.checkFlankDir(row+dir[0], col+dir[1], dir, depth+1)
		case board.Self:
			if depth > 0 {
				// fmt.Printf("dir: %d\n", dir)
				return true
			} else {
				return false
			}
		}
	}
	return false
}

func (board *Board) checkValid(row int, col int) bool {
	directions := [...][2]int{{-1, -1}, {-1, 0}, {-1, 1}, {0, -1}, {0, 1}, {1, -1}, {1, 0}, {1, 1}}
	result := false
	// fmt.Printf("checking %d-%d\n", row, col)
	for _, dir := range directions {
		// fmt.Printf("dir%d\n", dir)
		result = board.checkFlankDir(row, col, dir, 0)
		if result {
			// fmt.Printf("Valid: %d-%d\n", row, col)
			return true
		}
	}
	return false
}

func (board *Board) updateValidMoves() {
	board.ValidMoves = [][2]int{}
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			if board.Grid[i][j] == 0 {
				if board.checkValid(i, j) {
					// fmt.Printf("Valid: %d-%d\n", i, j)
					board.ValidMoves = append(board.ValidMoves, [2]int{i, j})
				}
			}
		}
	}
	if len(board.ValidMoves) == 0 {
		board.skipMove()
	}
}

func (board *Board) flank(row int, col int, dir [2]int) {
	board.Grid[row][col] = board.Self
	if board.Self == 1 {
		// fmt.Printf("flipped %d-%d from white to black\n", row, col)
		board.Whites--
		board.Blacks++
	} else {
		// fmt.Printf("flipped %d-%d from black to white\n", row, col)
		board.Blacks--
		board.Whites++
	}
	cell := board.Grid[row+dir[0]][col+dir[1]]
	switch cell {
	case board.Opponent:
		board.flank(row+dir[0], col+dir[1], dir)
	}
}

func (board *Board) playMove(row int, col int) error {
	var err error
	if row >= 0 && row <= 7 && col >= 0 && col <= 7 && board.Grid[row][col] == 0 {
		board.Skipped = 0
		board.Grid[row][col] = board.Self
		if board.Self == 1 {
			board.Blacks++
		} else {
			board.Whites++
		}

		directions := [...][2]int{{-1, -1}, {-1, 0}, {-1, 1}, {0, -1}, {0, 1}, {1, -1}, {1, 0}, {1, 1}}
		for _, dir := range directions {
			if board.checkFlankDir(row, col, dir, 0) {
				board.flank(row+dir[0], col+dir[1], dir)
			}
		}

		tmp := board.Opponent
		board.Opponent = board.Self
		board.Self = tmp
	} else {
		err = errors.New("Invalid Move")
		return err
	}
	return err
}

func (board *Board) skipMove() {
	// dom.GetWindow().Alert("You have no moves. Skipping turn")
	fmt.Println("You have no moves. Skipping turn")
	tmp := board.Opponent
	board.Opponent = board.Self
	board.Self = tmp
	board.Skipped++

	if board.checkWin() {
		if board.Blacks > board.Whites {
			board.Won = "Black Wins!"
		} else if board.Blacks < board.Whites {
			board.Won = "White Wins!"
		} else {
			board.Won = "Tie Game!"

		}
	} else {
		board.drawBoard()
	}
}

func (board *Board) checkWin() bool {
	return board.Skipped >= 2 || board.Blacks == 0 || board.Whites == 0 || board.Blacks+board.Whites == 64
}
