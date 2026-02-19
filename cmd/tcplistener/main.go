package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"strings"
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

		for l := range getLinesChannel(conn) {
			fmt.Println(l)
		}

		fmt.Println("Conn is closed")
	}

}

func getLinesChannel(f io.ReadCloser) <-chan string {
	ch := make(chan string)

	go func() {
		var currentLine []string
		data := make([]byte, 8)

		for {
			n, err := f.Read(data)

			if err != nil {
				break
			}

			dataStr := string(data[:n])

			if !strings.Contains(dataStr, "\n") {
				currentLine = append(currentLine, dataStr)
				continue
			}

			parts := strings.Split(dataStr, "\n")

			currentLine = append(currentLine, parts[0])

			ch <- strings.Join(currentLine, "")

			currentLine = nil

			if len(parts) > 2 {
				for _, p := range parts[1 : len(parts)-2] {
					ch <- p
				}
			}

			currentLine = append(currentLine, parts[len(parts)-1])
		}

		if len(currentLine) > 0 {
			ch <- strings.Join(currentLine, "")
		}

		close(ch)

		f.Close()
	}()

	return ch
}
