package gol

import (
	"strconv"
	"uk.ac.bris.cs/gameoflife/util"
)

type distributorChannels struct {
	events     chan<- Event
	ioCommand  chan<- ioCommand
	ioIdle     <-chan bool
	ioFilename chan<- string
	ioOutput   chan<- uint8
	ioInput    <-chan uint8
}

//making initial world
func makeWorld(height int, width int, input <-chan uint8, e chan<- Event) [][]byte {

	//making empty world
	world := make([][]byte, height)
	for r := range world {
		world[r] = make([]byte, width)
	}

	for r := 0; r < height; r++ {
		for c := 0; c < width; c++ {

			cellvalue := <-input
			world[c][r] = cellvalue

		}

	}

	return world
}

func calcAliveNeighbourValues(world [][]byte, p Params, r int, c int) int {

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

func getallalivecells(world [][]byte, p Params) []util.Cell {

	var alivecells []util.Cell

	for r := 0; r < p.ImageWidth; r++ {
		for c := 0; c < p.ImageHeight; c++ {
			if world[r][c] == 255 {
				alivecells = append(alivecells, util.Cell{X: r, Y: c})
			}

		}
	}
	return alivecells
}

func calculatenextstep(world [][]byte, p Params) [][]byte {

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

// distributor divides the work between workers and interacts with other goroutines.
func distributor(p Params, c distributorChannels) {

	// TODO: Create a 2D slice to store the world.
	//sending filename to IO
	filename := strconv.Itoa(p.ImageWidth) + "x" + strconv.Itoa(p.ImageHeight)
	c.ioCommand <- ioInput
	c.ioFilename <- filename

	//make initial world
	initialworld := makeWorld(p.ImageHeight, p.ImageWidth, c.ioInput, c.events)
	//calculates next step

	// TODO: Execute all turns of the Game of Life.
	turn := 0

	for turn < p.Turns {

		initialworld = calculatenextstep(initialworld, p)
		turn++

	}

	// TODO: Report the final state using FinalTurnCompleteEvent.

	listofallivecells := getallalivecells(initialworld, p)

	c.events <- FinalTurnComplete{ // Send a final turn complete event to the events channel
		CompletedTurns: turn,
		Alive:          listofallivecells,
	}

	// Make sure that the Io has finished any output before exiting.
	c.ioCommand <- ioCheckIdle
	<-c.ioIdle

	c.events <- StateChange{turn, Quitting}

	// Close the channel to stop the SDL goroutine gracefully. Removing may cause deadlock.
	close(c.events)
}
