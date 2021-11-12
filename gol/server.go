package gol

import (
	"errors"
	"flag"
	"fmt"
	"math/rand"
	"net"
	"time"
	//"secretstrings/stubs"
	"net/rpc"
)

var ProcessTurns = "GameOfLife.Process"

type Response struct {
	Message string
}

type Request struct {
	Message string
}

// exported
type GameOfLife struct {}

func handleError(err error) {
	if err != nil {
		panic(err)
	}
}
func (s* GameOfLife) Process (world [][]byte, p Params, c controllerChannels, req Request, res Response) (err error) {
	// process turns of GOL
	for turn < p.Turns {
		initialWorld = calculateNextStep(initialWorld, p)
		turn++
	}
}s

/* Below needs to be the GOL implementation */
func ReverseString(s string, i int) string {
	time.Sleep(time.Duration(rand.Intn(i)) * time.Second)
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
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
