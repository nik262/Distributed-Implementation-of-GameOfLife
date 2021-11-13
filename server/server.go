package main

import (
	"flag"
	"math/rand"
	"net"
	"net/rpc"
	"time"
	"uk.ac.bris.cs/gameoflife/stubs"
)

/* Below needs to be the GOL implementation */

func calcAliveNeighbourValues(world [][]byte, p stubs.StubP, r int, c int) int {

	alivemeter := 0

	//inner for loop calculates the state of the neighbours
	for i := r - 1; i <= r+1; i++ {
		for j := c - 1; j <= c+1; j++ {

			if i == r && j == c {
				continue
			}
			if world[((i + p.ImageWidth) % p.ImageWidth)][(j+p.ImageHeight)%p.ImageHeight] == 255 {
				alivemeter++
			}

		}
	}

	return alivemeter
}

func calculatenextstep(world [][]byte, p stubs.StubP) [][]byte {

	//replicating world so we can work on testerworld without disturbing world
	testerworld := make([][]byte, len(world))
	for i := range world {
		testerworld[i] = make([]byte, len(world[i]))
		copy(testerworld[i], world[i])
	}

	for r := 0; r < p.ImageHeight; r++ {
		for c := 0; c < p.ImageWidth; c++ {

			numberofaliveneighbours := calcAliveNeighbourValues(testerworld, p, r, c)

			//changing initial world with GOL conditions
			if numberofaliveneighbours < 2 || numberofaliveneighbours > 3 {
				world[r][c] = 0
			}
			if numberofaliveneighbours == 3 {
				world[r][c] = 255
			}

		}
	}

	return world
}

// exported
type GameOfLife struct{}

func handleError(err error) {
	if err != nil {
		panic(err)
	}
}

func (s *GameOfLife) Process(req stubs.Request, res *stubs.Response) (err error) {

	// process turns of GOL
	turn := 0
	for turn < req.P.Turns {
		req.World = calculatenextstep(req.World, req.P)
		turn++
	}
	res.World = req.World
	res.Turn = turn

	return
}

func main() {
	pAddr := flag.String("port", "8030", "Port to listen on")
	flag.Parse()
	rand.Seed(time.Now().UnixNano())
	rpc.Register(&GameOfLife{})
	listener, _ := net.Listen("tcp", ":"+*pAddr)
	defer listener.Close()
	rpc.Accept(listener)
}
