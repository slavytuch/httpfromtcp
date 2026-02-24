package request

import (
	"errors"
	"fmt"
	"httpfromtcp/internal/headers"
	"io"
	"regexp"
	"strconv"
	"strings"
)

type Request struct {
	RequestLine RequestLine
	Headers     headers.Headers
	Body        []byte
	status      requestStatus
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

type requestStatus string

const (
	requestStatusInitialized    = requestStatus("initialized")
	requestStatusParsingHeaders = requestStatus("parsing-headers")
	requestStatusParsingBody    = requestStatus("parsing-body")
	requestStatusDone           = requestStatus("done")
	bufferSize                  = 8
)

func (r *Request) parse(data []byte) (int, error) {
	switch r.status {
	case "":
		r.status = requestStatusInitialized
		fallthrough
	case requestStatusInitialized:
		prln, err := parseRequestLine(data, r)
		if err != nil {
			return 0, err
		}
		if prln > 0 {
			r.status = requestStatusParsingHeaders
		}
		return prln, nil
	case requestStatusParsingHeaders:
		prhn, done, err := r.Headers.Parse(data)
		if err != nil {
			return 0, err
		}
		if done {
			r.status = requestStatusParsingBody
		}
		return prhn, nil
	case requestStatusParsingBody:
		bodyLength, err := strconv.Atoi(r.Headers.Get("content-length"))

		if err != nil {
			return 0, err
		}

		if bodyLength == 0 {
			r.status = requestStatusDone
			return 0, nil
		}

		r.Body = append(r.Body, data...)

		if len(r.Body) > bodyLength {
			return 0, errors.New("error content body and header length mismatch")
		}

		if len(r.Body) == bodyLength {
			r.status = requestStatusDone
		}

		return len(data), nil
	case requestStatusDone:
		return 0, errors.New("request is done, cannot parse")
	default:
		return 0, errors.New("unknown request state")
	}
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	buffer := make([]byte, bufferSize)
	pr := Request{
		Headers: headers.NewHeaders(),
	}
	readToIndex := 0
	for pr.status != requestStatusDone {
		n, err := reader.Read(buffer[readToIndex:])

		if errors.Is(io.EOF, err) {
			pr.status = requestStatusDone
			break
		} else if err != nil {
			return nil, err
		}
		readToIndex += n

		if readToIndex >= len(buffer) {
			newBuffer := make([]byte, len(buffer)*2)
			copy(newBuffer, buffer)
			buffer = newBuffer
		}

		pn, err := pr.parse(buffer[:readToIndex])

		if err != nil {
			return nil, err
		}

		if pn > 0 {
			buffer = buffer[pn:]

			if len(buffer) < bufferSize {
				newBuffer := make([]byte, bufferSize)
				copy(newBuffer, buffer)
				buffer = newBuffer
			}

			readToIndex -= pn
		}

		fmt.Printf("Buffer: %s\n", string(buffer))
	}

	bodyLength, _ := strconv.Atoi(pr.Headers.Get("content-length"))

	if bodyLength != 0 && bodyLength > len(pr.Body) {
		return nil, errors.New("body length is shorter that Context-Length header value")
	}

	return &pr, nil
}

func parseRequestLine(data []byte, r *Request) (int, error) {
	if strings.Contains(string(data), "\r\n") != true {
		return 0, nil
	}

	line := strings.Split(string(data), "\r\n")[0]
	lineParts := strings.Split(line, " ")

	if len(lineParts) != 3 {
		return 0, errors.New("invalid request line - length is not 3")
	}

	if ok, _ := regexp.Match("^[A-Z]*$", []byte(lineParts[0])); !ok {
		return 0, errors.New("invalid request line - invalid method")
	}

	versionParams := strings.Split(lineParts[2], "/")

	if len(versionParams) != 2 || versionParams[1] != "1.1" {
		return 0, errors.New("invalid request line - unsupported version")
	}

	r.RequestLine = RequestLine{
		Method:        lineParts[0],
		RequestTarget: lineParts[1],
		HttpVersion:   versionParams[1],
	}

	return len(line) + 2, nil
}
