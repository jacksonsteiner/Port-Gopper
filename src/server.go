package main

import (
	"fmt"
	"os"
	"net"
	"strconv"
	"bytes"
	"math/rand"
	"encoding/json"
	neighbor "github.com/Port-Gopper/src/pkg"
)

func run_udp(neighbor *neighbor.Neighbor) {

	udpServerMaster, err := net.ListenPacket("udp", ":6095")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer udpServerMaster.Close()

	var buf []byte
	var addr net.Addr
	portRange := make(map[string]int)

	buf = make([]byte, 1024)
	_, addr, err = udpServerMaster.ReadFrom(buf)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	neighbor.IPAddr = addr.String()
	
	buf = bytes.Trim(buf, "\x00")

	err = json.Unmarshal(buf, &portRange)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	neighbor.StartPort = portRange["start"]
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	neighbor.EndPort = portRange["end"]
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	neighbor.Seed = rand.Int63()

	a, _ := json.Marshal(neighbor.Seed)
	udpServerMaster.WriteTo([]byte(a), addr)

	r := rand.New(rand.NewSource(neighbor.Seed))

	var max,min int

	if neighbor.StartPort <= neighbor.EndPort {
		max = neighbor.EndPort
		min = neighbor.StartPort
	} else {
		max = neighbor.StartPort
		min = neighbor.EndPort
	}

	for {
		newPort := strconv.Itoa(r.Intn(max - min + 1) + min)
		udpServer, err := net.ListenPacket("udp", ":" + newPort)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		newBuf := make([]byte, 1024)
		n, addr, err := udpServer.ReadFrom(newBuf)

		udpServer.Close()

		if err != nil {
			fmt.Println(err)
			continue
		}

		fmt.Println("Received " + strconv.Itoa(n) + " bytes on port " + newPort + " from " + addr.String())
		neighbor.Message += string(newBuf)

		if n < 10 {
			fmt.Println(neighbor.Message)
			break
		}

	}

}

func main() {

	neighbor := &neighbor.Neighbor{}
	run_udp(neighbor)

}
