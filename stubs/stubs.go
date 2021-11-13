package stubs

var ProcessTurns = "GameOfLife.Process"

type StubP struct {
	Turns       int
	Threads     int
	ImageWidth  int
	ImageHeight int
}

type Response struct {
	World [][]byte
	Turn  int
}

type Request struct {
	World [][]byte
	P     StubP
}
