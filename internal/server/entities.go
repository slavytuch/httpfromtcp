package server

import (
	"fmt"
	"httpfromtcp/internal/request"
	"httpfromtcp/internal/response"
	"net"
	"sync/atomic"
)

type Server struct {
	listener net.Listener
	handler  Handler
	closed   atomic.Bool
}

func Serve(port int, handler Handler) (*Server, error) {
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))

	if err != nil {
		return nil, err
	}

	server := Server{
		listener: l,
		handler:  handler,
	}

	go server.listen()

	return &server, nil
}

func (s *Server) Close() error {
	s.closed.Store(true)
	return s.listener.Close()
}

func (s *Server) listen() {
	for {
		conn, err := s.listener.Accept()

		if err != nil {
			if s.closed.Load() {
				return
			}
			fmt.Println("Error accepting conn " + err.Error())
			break
		}

		go s.handle(conn)
	}
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()
	req, err := request.RequestFromReader(conn)
	rw := response.Writer{Target: conn}

	if err != nil {
		HandlerError{
			StatusCode: response.StatusBadRequest,
			Message:    err.Error(),
		}.write(&rw)
		return
	}

	s.handler(&rw, req)
}

type Handler func(w *response.Writer, req *request.Request)

type HandlerError struct {
	StatusCode response.StatusCode
	Message    string
}

func (he HandlerError) write(w *response.Writer) {
	w.WriteStatusLine(he.StatusCode)
	headers := response.GetDefaultHeaders(len(he.Message))
	headers.Set("Content-Type", "text/html")
	w.WriteHeaders(headers)

	if he.Message != "" {
		w.WriteBody([]byte(he.Message))
	}
}
