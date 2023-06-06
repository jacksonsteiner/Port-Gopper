package main

import (
	"fmt"
	"os"
	"net"
	"strings"
	"strconv"
	"bytes"
	"math/rand"
    neighbor "github.com/Port-Gopper/src/pkg"
    _ "github.com/google/gopacket"
    _ "github.com/google/gopacket/layers"
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

	buf = make([]byte, 1024)
	_, addr, err = udpServerMaster.ReadFrom(buf)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	neighbor.IPAddr = addr.String()

	portRange := strings.Split(string(bytes.Trim(buf, "\x00")), ":")
	neighbor.StartPort, err = strconv.Atoi(portRange[0])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	neighbor.EndPort, err = strconv.Atoi(portRange[1])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	neighbor.Seed = rand.Int63()

	udpServerMaster.WriteTo([]byte(strconv.FormatInt(neighbor.Seed, 10)), addr)

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
