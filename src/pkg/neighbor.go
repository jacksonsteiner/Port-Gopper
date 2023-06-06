package neighbor

import (
	"fmt"
	"strconv"
	"net"
	"os"
	"bufio"
)

type fmtNeighbor interface {
	set_server_ip()
	set_start_port()
	set_end_port()
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

func (neighbor *Neighbor) Set_Server_IP() {

	var input string

	for len(neighbor.IPAddr) == 0 {
		fmt.Println("Enter Server IP: ")
		fmt.Scanln(&input)

		if net.ParseIP(input) == nil {
			fmt.Println("ERROR: INVALID IP ADDRESS")
		} else {
			neighbor.IPAddr = input
		}
	}

}

func (neighbor *Neighbor) Set_Start_Port() {

	for neighbor.StartPort == 0 {
		fmt.Println("Enter starting port (default is 49152): ")
		var input string
		fmt.Scanln(&input)

		if input == "" {
			neighbor.StartPort = 49152
		} else {
			if port, err := strconv.Atoi(input); err == nil && port >= 49152 && port <= 65535 {
				neighbor.StartPort = port
			} else {
				fmt.Println(`ERROR: INVALID STARTING EPHEMERAL PORT.
	Choose an integer between 49152 and 65535.`)
			}
		}
	}
}

func (neighbor *Neighbor) Set_End_Port() {

	for neighbor.EndPort == 0 {
		fmt.Println("Enter ending port (default is 65535): ")
		var input string
		fmt.Scanln(&input)

		if input == "" {
			neighbor.EndPort = 65535
		} else {
			if port, err := strconv.Atoi(input); err == nil && port >= 49152 && port <= 65535 && port != neighbor.StartPort {
				neighbor.EndPort = port
			} else {
				fmt.Println(`ERROR: INVALID ENDING EPHEMERAL PORT.
	Choose an integer between 49152 and 65535 that is not the same as the starting port.`)
			}
		}
	}
}

func (neighbor *Neighbor) Set_Message() {

	reader := bufio.NewReader(os.Stdin)
		fmt.Println("Enter message to send: ")
	neighbor.Message, _ = reader.ReadString('\n')

}