package main

import (
	"crypto/sha256"
	"fmt"
	"httpfromtcp/internal/headers"
	"httpfromtcp/internal/request"
	"httpfromtcp/internal/response"
	"httpfromtcp/internal/server"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
)

const port = 42069

func main() {
	s, err := server.Serve(port, handler)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer s.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}

func handler(w *response.Writer, req *request.Request) {
	if strings.HasPrefix(req.RequestLine.RequestTarget, "/httpbin/") {
		proxyHandler(w, req)
		return
	}

	if req.RequestLine.RequestTarget == "/video" {
		data, err := os.ReadFile("assets/vim.mp4")
		h := response.GetDefaultHeaders(len(data))
		if err != nil {
			w.WriteStatusLine(response.StatusError)
			w.WriteHeaders(h)
			w.WriteBody([]byte("error opening file - " + err.Error()))
			return
		}

		h.Set("Content-Type", "video/mp4")
		w.WriteStatusLine(response.StatusOK)
		w.WriteHeaders(h)
		w.WriteBody(data)

		return
	}

	var status response.StatusCode
	h := response.GetDefaultHeaders(0)
	h.Set("Content-Type", "text/html")
	var message string
	switch req.RequestLine.RequestTarget {
	case "/yourproblem":
		status = response.StatusBadRequest
		message = "<html>\n  <head>\n    <title>400 Bad Request</title>\n  </head>\n  <body>\n    <h1>Bad Request</h1>\n    <p>Your request honestly kinda sucked.</p>\n  </body>\n</html>"
	case "/myproblem":
		status = response.StatusError
		message = "<html>\n  <head>\n    <title>500 Internal Server Error</title>\n  </head>\n  <body>\n    <h1>Internal Server Error</h1>\n    <p>Okay, you know what? This one is on me.</p>\n  </body>\n</html>"
	default:
		status = response.StatusOK
		message = "<html>\n  <head>\n    <title>200 OK</title>\n  </head>\n  <body>\n    <h1>Success!</h1>\n    <p>Your request was an absolute banger.</p>\n  </body>\n</html>"
	}

	w.WriteStatusLine(status)

	h.Set("Content-Length", strconv.Itoa(len(message)))

	w.WriteHeaders(h)

	w.WriteBody([]byte(message))
}

func proxyHandler(w *response.Writer, req *request.Request) {
	resp, err := http.Get("https://httpbin.org/" + strings.TrimPrefix(req.RequestLine.RequestTarget, "/httpbin/"))

	if err != nil {
		w.WriteStatusLine(response.StatusError)
		w.WriteBody([]byte("Internal error - " + err.Error()))
		return
	}

	defer resp.Body.Close()

	w.WriteStatusLine(response.StatusOK)
	h := response.GetDefaultHeaders(0)
	h.Delete("Content-Length")
	h.Set("Transfer-Encoding", "chunked")
	h.Set("Trailer", "X-Content-SHA256, X-Content-Length")
	w.WriteHeaders(h)

	buff := make([]byte, 1024)
	var body []byte

	for {
		n, err := resp.Body.Read(buff)

		if err != nil {
			fmt.Println("Reading error: " + err.Error())
			break
		}

		w.WriteChunkedBody(buff[:n])

		body = append(body, buff[:n]...)
	}

	w.WriteChunkedBodyDone()

	err = w.WriteTrailers(headers.Headers{
		"X-Content-SHA256": fmt.Sprintf("%x", sha256.Sum256(body)),
		"X-Content-Length": fmt.Sprintf("%d", len(body)),
	})

	if err != nil {
		fmt.Println("error writing trailers: " + err.Error())
	}
}
