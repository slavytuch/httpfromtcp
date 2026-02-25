package response

import (
	"httpfromtcp/internal/headers"
	"strconv"
)

func GetDefaultHeaders(contentLen int) headers.Headers {
	return headers.Headers{
		"content-length": strconv.Itoa(contentLen),
		"connection":     "close",
		"content-type":   "text/plain",
	}
}
