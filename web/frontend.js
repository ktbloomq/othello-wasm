function updateBoard(boardString) {
    console.log("updating board")
    board = JSON.parse(boardString)
    console.log(board)
    grid = board.Grid
    for (let i = 0; i < 8; i++) {
        for (let j = 0; j < 8; j++) {
            cell = displayBoard[i][j]
            cell.classList.remove("hilight")
            switch (grid[i][j]) {
                case 1:
                    if(cell.classList.contains("white")) {
                        cell.classList.add("flip")
                    }
                    cell.classList.remove("white")
                    cell.classList.add("black")
                    break;
                case 2:
                    if(cell.classList.contains("black")) {
                        cell.classList.add("flip")
                    }
                    cell.classList.remove("black")
                    cell.classList.add("white")
                    break;
            }
        }
    }

    for (let i=0; i<board.ValidMoves.length; i++) {
        valid = board.ValidMoves[i]
        displayBoard[valid[0]][valid[1]].classList.add("hilight")
    }

    if (board.Won !== "") {
        window.alert(board.Won)
    }
}

function moveEvent(event) {
    idString = event.target.id;
    id = parseInt(idString)
    handleMove(id/10, id%10)
}

let displayBoard = []
window.addEventListener('load', function () {
  for (let i = 0; i < 8; i++) {
      displayBoard[i]=[]
      for (let j = 0; j < 8; j++) {
          displayBoard[i][j] = document.getElementById(i.toString()+j.toString())
          displayBoard[i][j].addEventListener("click",moveEvent)
      }
  }
  // console.log(userBoard)
})