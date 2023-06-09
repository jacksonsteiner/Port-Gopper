package main

import (
	"fmt"
	"os"
	"net"
	"math/rand"
	"strconv"
	"bytes"
    "time"
    neighbor "github.com/Port-Gopper/src/pkg"
)

func run_udp(neighbor *neighbor.Neighbor) {

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

	neighbor := &neighbor.Neighbor{}

	neighbor.Set_Server_IP()
	neighbor.Set_Start_Port()
	neighbor.Set_End_Port()
	neighbor.Set_Message()

	run_udp(neighbor)

}
