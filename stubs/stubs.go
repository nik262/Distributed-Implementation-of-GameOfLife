package gol

var ProcessTurns = "GameOfLife.Process"

type Response struct {
	world [][]byte
	Turn  int
}

type Request struct {
	World [][]byte
	P     Params
}
