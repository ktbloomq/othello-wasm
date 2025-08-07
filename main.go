package main

import (
	"errors"
	"fmt"
	"strconv"

	"honnef.co/go/js/dom/v2"
)

type Board struct {
	Grid       [8][8]rune
	turn       rune
	validMoves [][2]int
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
	board.turn = 'b'

	var cell *dom.HTMLInputElement
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			cell = dom.GetWindow().Document().GetElementByID(fmt.Sprintf("%d%d", i, j)).(*dom.HTMLInputElement)
			cell.AddEventListener("click", false, handleMove)
		}
	}
}

func (board *Board) drawBoard() {
	var cell *dom.HTMLInputElement
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			cell = dom.GetWindow().Document().GetElementByID(fmt.Sprintf("%d%d", i, j)).(*dom.HTMLInputElement)
			if board.Grid[i][j] == 'b' {
				cell.SetValue("⬤")
				cell.Class().Remove("white")
				cell.Class().Add("black")
			} else if board.Grid[i][j] == 'w' {
				cell.SetValue("⬤")
				cell.Class().Remove("black")
				cell.Class().Add("white")
			}
		}
	}
}

func (board *Board) checkValid(row int, col int, opponent int) bool {
	directions := [...][2]int{{-1, -1}, {-1, 0}, {-1, 1}, {0, -1}, {0, 1}, {1, -1}, {1, 0}, {1, 1}}
	for i, dir := range directions {
		fmt.Printf("%d%d", i, dir)
		// TODO
	}
	return true
}

func (board *Board) updateValidMoves() {
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			if board.checkValid(i, j, 1) {
				board.validMoves = append(board.validMoves, [2]int{i, j})
			}
		}
	}
}

func (board *Board) playMove(row int64, col int64) error {
	var err error
	if row >= 0 && row <= 8 && col >= 0 && col <= 8 && board.Grid[row][col] == ' ' {
		board.Grid[row][col] = board.turn
		if board.turn == 'b' {
			board.turn = 'w'
		} else {
			board.turn = 'b'
		}
	} else {
		err = errors.New("Invalid Move")
		return err
	}
	return err
}

func handleMove(event dom.Event) {
	idString := event.Target().ID()
	id, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		fmt.Printf("Move Error: Couldn't parse cell id. %s", err.Error())
	}
	// fmt.Println(idString)

	row := id / 10
	col := id % 10

	err = board.playMove(row, col)
	if err != nil {
		fmt.Println(err.Error())
	}
	board.drawBoard()
}

var board Board

func main() {
	board.constructBoard()
	board.drawBoard()

	select {}
}
