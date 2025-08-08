package main

import (
	"errors"
	"fmt"

	"honnef.co/go/js/dom/v2"
)

type Board struct {
	Grid        [8][8]rune
	displayGrid [8][8]*dom.HTMLButtonElement
	self        rune
	opponent    rune
	validMoves  [][2]int
	whites      int
	blacks      int
}

func (board *Board) constructBoard() {
	board.Grid = [8][8]rune{
		{' ', ' ', ' ', ' ', ' ', ' ', ' ', ' '},
		{' ', ' ', ' ', ' ', ' ', ' ', ' ', ' '},
		{' ', ' ', ' ', ' ', ' ', ' ', ' ', ' '},
		{' ', ' ', ' ', 'w', 'b', ' ', ' ', ' '},
		{' ', ' ', ' ', 'b', 'w', ' ', ' ', ' '},
		{' ', ' ', ' ', ' ', ' ', ' ', ' ', ' '},
		{' ', ' ', ' ', ' ', ' ', ' ', ' ', ' '},
		{' ', ' ', ' ', ' ', ' ', ' ', ' ', ' '}}
	board.self = 'b'
	board.opponent = 'w'
	board.whites = 2
	board.blacks = 2

	var cell *dom.HTMLButtonElement
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			cell = dom.GetWindow().Document().GetElementByID(fmt.Sprintf("%d%d", i, j)).(*dom.HTMLButtonElement)
			board.displayGrid[i][j] = cell
			cell.AddEventListener("click", false, handleMove)
		}
	}
}

func (board *Board) drawBoard() {
	var cell *dom.HTMLButtonElement
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			cell = board.displayGrid[i][j]
			cell.Class().Remove("hilight")
			switch board.Grid[i][j] {
			case 'b':
				cell.Class().Remove("white")
				cell.Class().Add("black")
				board.blacks++
			case 'w':
				cell.Class().Remove("black")
				cell.Class().Add("white")
				board.whites++
			}
		}
	}
	board.updateValidMoves()
}

func (board *Board) checkFlankDir(row int, col int, dir [2]int) bool {
	if row+dir[0] >= 0 && row+dir[0] <= 7 && col+dir[1] >= 0 && col+dir[1] <= 7 {
		// fmt.Printf("%d%d\n", row+dir[0], col+dir[1])
		switch board.Grid[row+dir[0]][col+dir[1]] {
		case board.opponent:
			return board.checkFlankDir(row+dir[0], col+dir[1], dir)
		case board.self:
			return true
		}
	}
	return false
}

func (board *Board) checkValid(row int, col int) bool {
	directions := [...][2]int{{-1, -1}, {-1, 0}, {-1, 1}, {0, -1}, {0, 1}, {1, -1}, {1, 0}, {1, 1}}
	result := false
	// fmt.Printf("checking %d-%d\n", row, col)
	for _, dir := range directions {
		if row+dir[0] >= 0 && row+dir[0] <= 7 && col+dir[1] >= 0 && col+dir[1] <= 7 {
			if board.Grid[row+dir[0]][col+dir[1]] == board.opponent {
				// fmt.Printf("dir%d\n", dir)
				result = board.checkFlankDir(row+dir[0], col+dir[1], dir)
				if result {
					return true
				}
			}
		}
	}
	return false
}

func (board *Board) updateValidMoves() {
	board.validMoves = [][2]int{}
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			if board.Grid[i][j] == ' ' {
				if board.checkValid(i, j) {
					board.validMoves = append(board.validMoves, [2]int{i, j})
					board.displayGrid[i][j].Class().Add("hilight")
				}
			}
		}
	}
}

func (board *Board) flank(row int, col int, dir [2]int) {
	board.Grid[row][col] = board.self
	cell := board.Grid[row+dir[0]][col+dir[1]]
	switch cell {
	case board.opponent:
		board.flank(row+dir[0], col+dir[1], dir)
	}
}

func (board *Board) playMove(row int, col int) error {
	var err error
	if row >= 0 && row <= 7 && col >= 0 && col <= 7 && board.Grid[row][col] == ' ' {
		directions := [...][2]int{{-1, -1}, {-1, 0}, {-1, 1}, {0, -1}, {0, 1}, {1, -1}, {1, 0}, {1, 1}}
		for _, dir := range directions {
			if board.checkFlankDir(row, col, dir) {
				board.flank(row, col, dir)
			}
		}

		tmp := board.opponent
		board.opponent = board.self
		board.self = tmp
	} else {
		err = errors.New("Invalid Move")
		return err
	}
	return err
}
