package response

import (
	"fmt"
	"httpfromtcp/internal/headers"
	"io"
)

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	var err error
	switch statusCode {
	case StatusOK:
		_, err = w.Write([]byte("HTTP/1.1 200 OK\r\n"))
	case StatusBadRequest:
		_, err = w.Write([]byte("HTTP/1.1 400 Bad Request\r\n"))
	case StatusError:
		_, err = w.Write([]byte("HTTP/1.1 500 Internal Server Error\r\n"))
	default:
		_, err = w.Write([]byte(fmt.Sprintf("HTTP/1.1 %d\r\n", statusCode)))
	}

	return err
}

func WriteHeaders(w io.Writer, headers headers.Headers) error {
	for k, v := range headers {
		_, err := w.Write([]byte(fmt.Sprintf("%s: %s\r\n", k, v)))

		if err != nil {
			return err
		}
	}

	_, err := w.Write([]byte("\r\n"))
	return err
}
