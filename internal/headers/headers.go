package headers

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

type Headers map[string]string

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	dataStr := string(data)
	rnIdx := strings.Index(dataStr, "\r\n")
	if rnIdx == -1 {
		return 0, false, nil
	}

	if rnIdx == 0 {
		return 2, true, nil
	}

	if endLineIdx := strings.Index(dataStr, "\r\n\r\n"); endLineIdx != -1 {
		dataStr = dataStr[:endLineIdx+2]
	}

	bytes := 0
	for _, line := range strings.Split(dataStr, "\r\n") {

		parts := strings.SplitN(line, ":", 2)

		if len(parts) < 2 {
			fmt.Println(dataStr)
			fmt.Println(line)
			return 0, false, errors.New("invalid header name - missing value")
		}

		if parts[0][len(parts[0])-1] == ' ' {
			return 0, false, errors.New("invalid header name - space before colon")
		}

		if ok, _ := regexp.MatchString("^[A-Za-z0-9!#$%&'*+\\-.^_`|~]*$", parts[0]); !ok {
			return 0, false, errors.New("invalid header name - invalid character")
		}

		key := strings.ToLower(strings.Trim(parts[0], " "))
		value := strings.Trim(parts[1], " ")

		if value == "" {
			return 0, false, errors.New("invalid header name - missing value")
		}

		if _, ok := h[key]; ok {
			h[key] += ", " + value
		} else {
			h[key] = value
		}

		bytes += len(line) + 2
	}

	return bytes, false, nil
}

func NewHeaders() Headers {
	return make(Headers)
}

func (h Headers) Get(name string) string {
	return h[strings.ToLower(name)]
}
