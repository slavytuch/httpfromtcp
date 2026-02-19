package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	addr, err := net.ResolveUDPAddr("udp", "localhost:42069")

	if err != nil {
		log.Fatal("Error resolving udp addr " + err.Error())
	}

	conn, err := net.DialUDP("udp", nil, addr)

	if err != nil {
		log.Fatal("Error dialing conn " + err.Error())
	}

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print(" > ")
		input, err := reader.ReadString('\n')

		if err != nil {
			log.Fatal("Error reading from console " + err.Error())
		}

		_, err = conn.Write([]byte(input))

		if err != nil {
			log.Fatal("Error writing to connection " + err.Error())
		}
	}
}
