package main

import (
	"fmt"
	"strconv"

	"honnef.co/go/js/dom/v2"
)

func handleMove(event dom.Event) {
	idString := event.Target().ID()
	id, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		fmt.Printf("Move Error: Couldn't parse cell id. %s", err.Error())
	}
	// fmt.Println(idString)

	row := int(id / 10)
	col := int(id % 10)

	isValid := false
	for _, valid := range board.validMoves {
		if valid[0] == row && valid[1] == col {
			isValid = true
			break
		}
	}
	if isValid {
		err = board.playMove(row, col)
		if err != nil {
			fmt.Println(err.Error())
		}
		board.drawBoard()
	}
}

var board Board

func main() {
	board.constructBoard()
	board.drawBoard()

	select {}
}
