package main

import (
	"fmt"
	"os"
	"net"
	"math/rand"
	"strconv"
	"bytes"
	"time"
	"encoding/json"
	neighbor "github.com/Port-Gopper/src/pkg"
	_ "github.com/google/gopacket"
	_ "github.com/google/gopacket/layers"
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

	marsh, _ := json.Marshal(map[string]int{"start": neighbor.StartPort,"end": neighbor.EndPort})
	_, err = conn.Write([]byte(marsh))
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

	recvFrom = bytes.Trim(recvFrom, "\x00")
	err = json.Unmarshal(recvFrom, &neighbor.Seed)

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

	sendBuf := bytes.Buffer{}
	sendBuf.Write([]byte(neighbor.Message))
	nextTenBuf := sendBuf.Next(10)
	
	timeoutCount := 0

	for sendBuf.Len() > 0 {

		if timeoutCount >= 5 {
			break
		}

		newPortInt := r.Intn(max - min + 1) + min
		newPort := strconv.Itoa(newPortInt)
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

		defer newConn.Close()

		n, err := newConn.Write(nextTenBuf)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		err = newConn.SetReadDeadline(time.Now().Add(10 * time.Second))
    	if err != nil {
        	fmt.Println("Case for read deadline failing not handled! Exiting")
        	os.Exit(1)
    	}

    	recvBuf := make([]byte, 1024)
    	_, err = newConn.Read(recvBuf)
    	if err != nil {
    		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
    			timeoutCount = timeoutCount + 1
				fmt.Println("Timeout error count: " + strconv.Itoa(timeoutCount))
    			continue
    		} else {
    			timeoutCount = timeoutCount + 1
    			fmt.Println("Timeout error count: " + strconv.Itoa(timeoutCount))
    			fmt.Println("Server confirmation error. Retrying.")
    			continue
    		}
    	}

		var servPortInt int
		recvBuf = bytes.Trim(recvBuf, "\x00")
		servPortInt, err = strconv.Atoi(string(recvBuf))

		if servPortInt == newPortInt {
			timeoutCount = 0
			nextTenBuf = sendBuf.Next(10)
			fmt.Println("Sent " + strconv.Itoa(n) + " bytes on port " + newPort)
		} else {
			fmt.Println("Server out of sync. Exiting.")
			os.Exit(1)
		}
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
