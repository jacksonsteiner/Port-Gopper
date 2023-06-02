package main

import (
	"fmt"
	"os"
	"net"
	"math/rand"
	"strconv"
	"bytes"
        "time"
	"bufio"
      _ "github.com/google/gopacket"
      _ "github.com/google/gopacket/layers"
)

type fmtNeighbor interface {
	set_server_ip()
	set_port(string, int)
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

func (neighbor *Neighbor) set_port(instruct string, dflt int) {

	fmt.Println(instruct)
	var input string
	fmt.Scanln(&input)

	var neighborVal *int

	if dflt == 49152 {
		neighborVal = &neighbor.StartPort
	} else {
		neighborVal = &neighbor.EndPort
	}

	if input == "" {
		*neighborVal = dflt
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

func (neighbor *Neighbor) set_message() {

	reader := bufio.NewReader(os.Stdin)
		fmt.Println("Enter message to send: ")
	neighbor.Message, _ = reader.ReadString('\n')

}

func run_udp(neighbor *Neighbor) {

	udpServerInit, err := net.ResolveUDPAddr("udp", neighbor.IPAddr + ":6095")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	conn, err := net.DialUDP("udp", nil, udpServerInit)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	_, err = conn.Write([]byte(strconv.Itoa(neighbor.StartPort) + ":" + strconv.Itoa(neighbor.EndPort)))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	recvFrom := make([]byte, 1024)
	_, err = conn.Read(recvFrom)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	neighbor.Seed, err = strconv.ParseInt(string(bytes.Trim(recvFrom, "\x00")), 10, 64)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	conn.Close()

	r := rand.New(rand.NewSource(neighbor.Seed))

	var max,min int

	if neighbor.StartPort <= neighbor.EndPort {
		max = neighbor.EndPort
		min = neighbor.StartPort
	} else {
		max = neighbor.StartPort
		min = neighbor.EndPort
	}

	buf := bytes.Buffer{}
	buf.Write([]byte(neighbor.Message))

	for buf.Len() > 0 {

		time.Sleep(2 * time.Second)

		newPort := strconv.Itoa(r.Intn(max - min + 1) + min)
		udpServer, err := net.ResolveUDPAddr("udp", neighbor.IPAddr + ":" + newPort)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		newConn, err := net.DialUDP("udp", nil, udpServer)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		newBuf := buf.Next(10)
		n, err := newConn.Write(newBuf)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println("Sent " + strconv.Itoa(n) + " bytes on port " + newPort)

		newConn.Close()

	}

}

func main() {

	neighbor := &Neighbor{}

	neighbor.set_server_ip()
	neighbor.set_port("Enter starting port (default is 49152): ", 49152)
	neighbor.set_port("Enter ending port (default is 65535): ", 65535)
	neighbor.set_message()

	run_udp(neighbor)

}
