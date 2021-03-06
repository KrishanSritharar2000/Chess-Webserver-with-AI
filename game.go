package main

import (
	"fmt"
	"bufio"
	"os"
	"strings"
)

type Game struct {
	Board Board
	IsWhiteTurn bool
	Moves *MoveStack
}

type Move struct {
	From, To Piece
}

func (g *Game) nextTurn() {
	g.IsWhiteTurn = !g.IsWhiteTurn

}

//Pre: Takes in a character
//Post: True if letter in A-G
func isLetter(s string) bool {
	return s == "a" || s == "b" || s == "c" || s == "d" || s == "e" || s == "f" || s == "g" || s == "h"
}

//Pre: Takes in a character
//Post: True if number between 0-7
func isNumber(s string) bool {
	return s == "1" || s == "2" || s == "3" || s == "4" || s == "5" || s == "6" || s == "7" || s == "8"
}

func (g *Game) getTurn(r bufio.Reader) (int, int, int, int) {
	var isCheck bool
	if g.IsWhiteTurn {
		fmt.Println("White's Turn")
		isCheck = g.Board.kingW.isCheck(&g.Board)
	} else {
		fmt.Println("Black's Turn")
		isCheck = g.Board.kingB.isCheck(&g.Board)
	}
	if isCheck {
		fmt.Println("CHECK")
	}
	var str, dst string
	for {
		fmt.Println("\nEnter the piece you want to move and where to: ")
		str, _ = r.ReadString('\n')
		str = strings.ToLower(str)
		str = strings.TrimSpace(str)

		if string(str[2]) == " " {
			dst = str[3:]
			str = str[:2]
			dst = strings.TrimSpace(dst)

			if len(dst) == 2 {
				if isLetter(string(str[0])) && isLetter(string(dst[0])) && isNumber(string(str[1])) && isNumber(string(dst[1])) {
					break
				} 
			}
		}
		fmt.Println("\nThat is an invalid move")
		fmt.Println("Please use the format: [LetterNumber LetterNumber]")
		fmt.Println("To provide the coordinates [A-G 1-8] of the chosen piece, and where to move it")
	}	
	
	startY := int(str[0]) - int('a')
	startX := int(str[1]) - int('1')
	endY := int(dst[0]) - int('a')
	endX := int(dst[1]) - int('1')
	return startX, startY, endX, endY
}

func (g *Game) undoTurn() (bool, Move) {
	move, ok := g.Moves.Pop()
	if ok {
		startPiece := move.From
		endPiece := move.To
		g.Board.Board[startPiece.x][startPiece.y] = startPiece
		g.Board.Board[endPiece.x][endPiece.y] = endPiece
		g.nextTurn()

		//replace king coords 
		if startPiece.Symbol == "K" {
			if startPiece.IsBlack {
				g.Board.kingB = startPiece
			} else {
				g.Board.kingW = startPiece
			}
		}

		//check castle
		if startPiece.Symbol == "K" && abs(startPiece.y - endPiece.y) == 2 {
			var rookX int
			if startPiece.IsBlack {
				rookX = 7
				g.Board.castleCheck[1] = false
			} else {
				rookX = 0
				g.Board.castleCheck[0] = false
			}
			if (endPiece.y == 6) {
				g.Board.Board[rookX][7] = g.Board.Board[rookX][5]
				g.Board.Board[rookX][7].x = rookX
				g.Board.Board[rookX][7].y = 7
				g.Board.Board[rookX][5] = Piece{rookX, 5, " ", false}
				if rookX == 0 {
					g.Board.castleCheck[3] = false
				} else {
					g.Board.castleCheck[2] = false
				}
			} else {
				g.Board.Board[rookX][0] = g.Board.Board[rookX][3]
				g.Board.Board[rookX][0].x = rookX
				g.Board.Board[rookX][0].y = 0
				g.Board.Board[rookX][3] = Piece{rookX, 3, " ", false}
				if rookX == 0 {
					g.Board.castleCheck[5] = false
				} else {
					g.Board.castleCheck[4] = false
				}
			}
		}
		return true, move
	}
	return false, move
}

func (g *Game) makeMove(startX, startY, endX, endY int) bool {

	startPiece := g.Board.Board[startX][startY]	
	if startPiece.Symbol == " " {
		fmt.Println(g.Board)
		fmt.Println(startPiece)
		fmt.Println("There is no piece there!")
		return false
	} else if startPiece.IsBlack == g.IsWhiteTurn {
		fmt.Println("That is not your piece, you cannot move it!")
		return false
	}
	endPiece := g.Board.Board[endX][endY]
	result := g.Board.Board[startX][startY].move(&g.Board, endX, endY)
	//moving white piece need to check piece is white
	if result {
		g.Moves.Push(Move{startPiece, endPiece})
		g.nextTurn()
	} else {
		fmt.Println("That move is not allowed")
	}
	return result
}

func StartGame() {
	g := Game{Board{}, true, &MoveStack{}}
	SetupBoard(&g.Board)
	fmt.Println(g.Board)
	reader := bufio.NewReader(os.Stdin)
	for {
		g.makeMove(g.getTurn(*reader))
		fmt.Println(g.Board)
	}
	// g.testBoard()
    // fmt.Print("Enter text: ")
    // text, _ := reader.ReadString('\n')
	// fmt.Println(text, len(text))
	// text = strings.TrimSpace(text)
	// fmt.Println(len(text))
	// for i := 0; i < len(text); i++ {
	// 	fmt.Println(string(text[i]))
	// }
}

func (g Game) testBoard() {
	Board := Board{}

	fmt.Println("Hello World3!")
	SetupBoard(&Board)
	// fmt.Println(Board)
	// Board.Board[1][0].move(&Board, 2, 0)
	fmt.Println(Board)
	Board.Board[0][1].move(&Board, 2, 2)
	fmt.Println(Board)
	// Board.Board[2][2].move(&Board, 4, 3)
	// Board.Board[7][1].move(&Board, 5, 2)
	// Board.Board[5][2].move(&Board, 5, 4)
	// Board.Board[4][3].move(&Board, 5, 4)
	Board.Board[1][3].move(&Board, 3, 3)
	Board.Board[3][3].move(&Board, 4, 3)
	Board.Board[0][3].move(&Board, 2, 3)
	Board.Board[2][3].move(&Board, 4, 5)
	Board.Board[4][5].move(&Board, 6, 4)
	Board.Board[0][4].move(&Board, 1, 3)
	Board.Board[1][3].move(&Board, 3, 3)
	Board.Board[6][4].move(&Board, 5, 4)
	Board.Board[7][3].move(&Board, 5, 5)
	Board.Board[4][4].move(&Board, 3, 4)
	
	fmt.Println("END Is black king currently in check: ", Board.kingB.isCheck(&Board))
	fmt.Println("END Is white king currently in check: ", Board.kingW.isCheck(&Board))
	Board.Board[3][4].move(&Board, 2, 4)
	Board.Board[1][3].move(&Board, 2, 4)

	Board.Board[5][5].move(&Board, 5,6)
	// Board.Board[4][4].move(&Board, 3,5)
	// Board.Board[3][5].move(&Board, 2,5)
	// Board.Board[3][4].move(&Board, 3,3)
	// Board.Board[3][3].move(&Board, ,3)
	// Board.Board[1][3].move(&Board, 1,4)
	// Board.Board[1][4].move(&Board, 1,5)
	// Board.Board[1][5].move(&Board, 3,5)
	Board.Board[7][6].move(&Board, 5,5)
	Board.Board[5][5].move(&Board, 3,6)
	Board.Board[3][6].move(&Board, 5,7)
	Board.Board[5][7].move(&Board, 4,5)
	Board.Board[4][5].move(&Board, 3,3)
	Board.Board[3][3].move(&Board, 1,2)
	Board.Board[1][2].move(&Board, 0,4)

	Board.Board[0][4].move(&Board, 1,6)
	Board.Board[1][6].move(&Board, 3,7)

	Board.Board[1][1].move(&Board, 3,1)
	Board.Board[3][1].move(&Board, 4,1)
	fmt.Println("---------------------")
	Board.Board[6][2].move(&Board, 4,2)
	fmt.Println("---------------------")
	Board.Board[4][1].move(&Board, 7,2)
	fmt.Println("---------------------")
	Board.Board[1][5].move(&Board, 3,5)
	Board.Board[6][7].move(&Board, 4,7)
	fmt.Println("---------------------")
	Board.Board[4][1].move(&Board, 7,2)
	Board.Board[2][4].move(&Board, 2,5)
	fmt.Println("---------------------")

	Board.Board[7][5].move(&Board, 5,3)
	fmt.Println("---------------------")

	Board.Board[7][4].move(&Board, 7,6)
	

	fmt.Println("END Is black king currently in check: ", Board.kingB.isCheck(&Board))
	fmt.Println("END Is white king currently in check: ", Board.kingW.isCheck(&Board))



	// fmt.Println(Board)
	// Board.Board[6][3].move(&Board, 2, 3)
	// Board.Board[1][2].move(&Board, 3, 2)
	fmt.Println(Board)

	fmt.Println("---------------------")
	// var fir string 
	// fmt.Scanln(&fir)
	// fmt.Println("This is the string:" , fir, string(fir[4]), len(fir))
	// colorReset := "\033[0m"

    // colorRed := "\033[31m"
    // colorGreen := "\033[32m"
    // colorYellow := "\033[33m"
    // colorBlue := "\033[34m"
    // colorPurple := "\033[35m"
    // colorCyan := "\033[36m"
	// colorWhite := "\033[37m"
	// boardDark := "\033[48;2;181;136;99m"
	// boardLight := "\033[48;2;240;217;181m"
    
    // fmt.Println(boardDark,string(colorRed), "test", string(colorReset))
    // fmt.Println(boardLight, string(colorGreen), "test", string(colorReset))
    // fmt.Println(boardDark, string(colorYellow), "test")
    // fmt.Println(string(colorBlue), "test")
    // fmt.Println(string(colorPurple), "test")
    // fmt.Println(string(colorWhite), "test")
    // fmt.Println(string(colorCyan), "test", string(colorReset))
    // fmt.Println("next")
}