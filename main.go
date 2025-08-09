package main

import (
	"fmt"
	"syscall/js"
)

func handleMove(this js.Value, args []js.Value) interface{} {
	var result interface{}
	row := args[0].Int()
	col := args[1].Int()

	isValid := false
	for _, valid := range board.ValidMoves {
		if valid[0] == row && valid[1] == col {
			isValid = true
			break
		}
	}
	if isValid {
		err := board.playMove(row, col)
		if err != nil {
			fmt.Println(err.Error())
		}
		board.drawBoard()
	}
	return result
}

var board Board

func main() {
	js.Global().Set("handleMove", js.FuncOf(handleMove))
	board.constructBoard()
	board.drawBoard()

	<-make(chan struct{})
}
