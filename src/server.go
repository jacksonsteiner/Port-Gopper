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

func talk_to_neighbor(udpServerMaster net.PacketConn, buf []byte, addr net.Addr) {

	neighbor := &neighbor.Neighbor{}
	portRange := make(map[string]int)

	neighbor.IPAddr = addr.String()

	buf = bytes.Trim(buf, "\x00")

	err := json.Unmarshal(buf, &portRange)
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

	marsh, _ := json.Marshal(neighbor.Seed)
	udpServerMaster.WriteTo([]byte(marsh), addr)

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
		laddr, err := net.ResolveUDPAddr("udp", ":" + newPort)
		udpServer, err := net.ListenUDP("udp", laddr)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		newBuf := make([]byte, 1024)
		n, addr, err := udpServer.ReadFrom(newBuf)

		if err != nil {
			fmt.Println(err)
			udpServer.Close()
			continue
		}

		newPortInt, _ := strconv.Atoi(newPort)
		marsh, _ = json.Marshal(newPortInt)
		udpServer.WriteTo([]byte(marsh), addr)

		udpServer.Close()

		fmt.Println("Received " + strconv.Itoa(n) + " bytes on port " + newPort + " from " + addr.String())
		neighbor.Message += string(newBuf)

		if n < 10 {
			fmt.Println(neighbor.Message)
			break
		}
	}

}

func run_master_server() {

	laddr, err := net.ResolveUDPAddr("udp", ":6095")
	udpServerMaster, err := net.ListenUDP("udp", laddr)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer udpServerMaster.Close()

	for {
		buf := make([]byte, 1024)
		n, addr, err := udpServerMaster.ReadFrom(buf)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if n != 0 && addr != nil {
			go talk_to_neighbor(udpServerMaster, buf, addr)
		}
	}

}

func main() {

	run_master_server()

}
