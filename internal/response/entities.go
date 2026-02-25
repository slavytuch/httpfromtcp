package response

import (
	"fmt"
	"httpfromtcp/internal/headers"
	"io"
)

type StatusCode int
type WriterState string

const (
	StatusOK         = StatusCode(200)
	StatusBadRequest = StatusCode(400)
	StatusError      = StatusCode(500)
	WriterStateH     = WriterState("headers")
	WriterStateB     = WriterState("body")
)

type Writer struct {
	Target io.Writer
	State  WriterState
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	if w.State != "" {
		return fmt.Errorf("invalid writer state - got: %s need: empty", w.State)
	}

	err := WriteStatusLine(w.Target, statusCode)

	if err != nil {
		return err
	}

	w.State = WriterStateH

	return nil
}

func (w *Writer) WriteHeaders(headers headers.Headers) error {
	if w.State != WriterStateH {
		return fmt.Errorf("invalid writer state - got: %s need: %s", w.State, WriterStateH)
	}

	err := WriteHeaders(w.Target, headers)

	if err != nil {
		return err
	}

	w.State = WriterStateB

	return nil
}

func (w *Writer) WriteBody(p []byte) (int, error) {
	if w.State != WriterStateB {
		return 0, fmt.Errorf("invalid writer state - got: %s need: %s", w.State, WriterStateB)
	}

	return w.Target.Write(p)
}

func (w *Writer) WriteChunkedBody(p []byte) (int, error) {
	if w.State != WriterStateB {
		return 0, fmt.Errorf("invalid writer state - got: %s need: %s", w.State, WriterStateB)
	}

	return w.Target.Write([]byte(fmt.Sprintf("%x\r\n%s\r\n", len(p), p)))
}

func (w *Writer) WriteChunkedBodyDone() (int, error) {
	if w.State != WriterStateB {
		return 0, fmt.Errorf("invalid writer state - got: %s need: %s", w.State, WriterStateB)
	}

	return w.Target.Write([]byte("0\r\n"))
}

func (w *Writer) WriteTrailers(h headers.Headers) error {
	if w.State != WriterStateB {
		return fmt.Errorf("invalid writer state - got: %s need: %s", w.State, WriterStateB)
	}

	return WriteHeaders(w.Target, h)
}
