package main

import (
	"fmt"
	"httpfromtcp/internal/request"
	"log"
	"net"
)

func main() {
	listener, err := net.Listen("tcp", "127.0.0.1:42069")

	if err != nil {
		log.Fatal("Error opening listener " + err.Error())
	}

	defer listener.Close()

	for {
		conn, err := listener.Accept()

		if err != nil {
			fmt.Println("error accepting conn " + err.Error())
			break
		}

		fmt.Println("Conn is accepted")

		r, err := request.RequestFromReader(conn)

		if err != nil {
			fmt.Println("error parsing r " + err.Error())
		} else {
			fmt.Printf(
				"Request line:\n- Method: %s\n- Target: %s\n- Version: %v\nHeaders:\n",
				r.RequestLine.Method,
				r.RequestLine.RequestTarget,
				r.RequestLine.HttpVersion,
			)

			for k, v := range r.Headers {
				fmt.Printf("- %s: %s\n", k, v)
			}

			fmt.Printf("Body:\n%s", string(r.Body))
		}

		fmt.Println("Conn is closed")
	}

}
