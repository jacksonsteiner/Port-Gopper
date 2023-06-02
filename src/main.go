package main

import (
	"fmt"
	"os"
	"net"
	"math/rand"
	"strconv"
      _ "github.com/google/gopacket"
      _ "github.com/google/gopacket/layers"
)

type fmtNeighbor interface {
	set_server_ip()
	set_port(string, int)
	set_seed()
	set_message()
}

type Neighbor struct {
	fmtNeighbor
	Message   string
	IPAddr    string
	StartPort int
	EndPort   int
	Seed      int64
}

func (neighbor *Neighbor) set_server_ip() {

	fmt.Println("Enter Server IP: ")
	var input string
	fmt.Scanln(&input)

	if net.ParseIP(input) == nil {
		fmt.Println("FATAL ERROR: INVALID IP ADDRESS")
		os.Exit(1)
	}

	neighbor.IPAddr = input

}

func (neighbor *Neighbor) set_port(instruct string, dfalt int) {

	fmt.Println(instruct)
	var input string
	fmt.Scanln(&input)

	var neighborVal *int

	if dfalt == 49152 {
		neighborVal = &neighbor.StartPort
	} else {
		neighborVal = &neighbor.EndPort
	}

	if input == "" {
		*neighborVal = dfalt
	} else {
		if port, err := strconv.Atoi(input); err == nil && port >= 49152 && port <= 65535 {
			*neighborVal = port
		} else {
			fmt.Println(`FATAL ERROR: INVALID STARTING EPHEMERAL PORT.
Choose an integer between 49152 and 65535.`)
			os.Exit(1)
		}
	}

}

func (neighbor *Neighbor) set_seed() {

	fmt.Println("Enter an integer for the port hopping sequence (default, 0, a decimal, or a negative number entered will result in a random number being used as the seed.): ")
	var input string
	fmt.Scanln(&input)

	if input == "" {
		neighbor.Seed = rand.Int63()
	} else {
		if seed, err := strconv.ParseInt(input, 10, 64); err == nil && seed > 0 {
			neighbor.Seed = seed
		} else {
			neighbor.Seed = rand.Int63()
		}
	}

}

func (neighbor *Neighbor) set_message() {

	fmt.Println("Enter message to send: ")
	var input string
	fmt.Scanln(&input)
	neighbor.Message = input

}

func generate_sequence(start int, end int, seed int64) {
	r := rand.New(rand.NewSource(seed))
	var max,min int

	if start <= end {
		max = end
		min = start
	} else {
		max = start
		min = end
	}

	for i := 0; i < 10; i++ {
		fmt.Println(r.Intn(max - min + 1) + min)
	}
}

func main() {

	neighbor := &Neighbor{}

	neighbor.set_server_ip()
	neighbor.set_port("Enter starting port (default is 49152): ", 49152)
	neighbor.set_port("Enter ending port (default is 65535): ", 65535)
	neighbor.set_seed()
	neighbor.set_message()

	generate_sequence(neighbor.StartPort, neighbor.EndPort, neighbor.Seed)
}