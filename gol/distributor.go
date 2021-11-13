package gol

import (
	"flag"
	"net/rpc"
	"strconv"
	"uk.ac.bris.cs/gameoflife/stubs"
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
func makeWorld(height int, width int, input <-chan uint8) [][]byte {

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

func makeCall(client rpc.Client, initialworld [][]byte, p Params) *stubs.Response {

	params := stubs.StubP{Turns: p.Turns, Threads: p.Threads, ImageWidth: p.ImageWidth, ImageHeight: p.ImageHeight}
	request := stubs.Request{World: initialworld, P: params}
	response := new(stubs.Response)
	client.Call(stubs.ProcessTurns, request, response)
	return response

}

// distributor divides the work between workers and interacts with other goroutines.
func distributor(p Params, c distributorChannels) {

	// TODO: Create a 2D slice to store the world.
	//sending filename to IO
	filename := strconv.Itoa(p.ImageWidth) + "x" + strconv.Itoa(p.ImageHeight)
	c.ioCommand <- ioInput
	c.ioFilename <- filename

	//make initial world
	initialworld := makeWorld(p.ImageHeight, p.ImageWidth, c.ioInput)
	//calculates next step

	// TODO: Execute all turns of the Game of Life.

	//client side boiler code copied from secretstrings
	server := flag.String("server", "3.91.157.15:8030", "IP:port string to connect to as server")
	flag.Parse()
	client, _ := rpc.Dial("tcp", *server)
	defer client.Close()

	//calls client and assigns the recieved world to initialworld
	turn := 0
	responseval := makeCall(*client, initialworld, p)
	turn = responseval.Turn
	initialworld = responseval.World

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
