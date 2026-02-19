package request

import (
	"errors"
	"io"
	"regexp"
	"strings"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	rd, err := io.ReadAll(reader)

	if err != nil {
		return nil, err
	}

	rlines := strings.Split(string(rd), "\r\n")

	if len(rlines) == 0 {
		return nil, errors.New("empty request line")
	}

	rl, err := parseRequestLine(rlines[0])

	if err != nil {
		return nil, err
	}

	return &Request{
		RequestLine: rl,
	}, nil
}

func parseRequestLine(rls string) (RequestLine, error) {
	lines := strings.Split(rls, " ")

	if len(lines) != 3 {
		return RequestLine{}, errors.New("invalid request line - length is not 3")
	}

	if ok, _ := regexp.Match("^[A-Z]*$", []byte(lines[0])); !ok {
		return RequestLine{}, errors.New("invalid request line - invalid method")
	}

	versionParams := strings.Split(lines[2], "/")

	if len(versionParams) != 2 || versionParams[1] != "1.1" {
		return RequestLine{}, errors.New("invalid request line - unsupported version")
	}

	return RequestLine{
		Method:        lines[0],
		RequestTarget: lines[1],
		HttpVersion:   versionParams[1],
	}, nil
}
