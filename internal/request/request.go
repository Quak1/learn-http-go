package request

import (
	"fmt"
	"io"
	"strings"
	"unicode"
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
	bytes, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	parts := strings.Split(string(bytes), "\r\n")
	requestLine, err := parseRequestLine(parts[0])
	if err != nil {
		return nil, err
	}

	request := &Request{
		RequestLine: *requestLine,
	}

	return request, nil
}

func parseRequestLine(str string) (*RequestLine, error) {
	parts := strings.Split(str, " ")
	if len(parts) != 3 {
		return nil, fmt.Errorf("Request line: wrong number of parts")
	}

	method := parts[0]
	for _, c := range method {
		if c < 'A' || c > 'z' {
			return nil, fmt.Errorf("Request line: invalid method")
		}
	}

	httpParts := strings.Split(parts[2], "/")
	if len(httpParts) != 2 {
		return nil, fmt.Errorf("Request line: invalid HTTP-version")
	}
	if httpParts[0] != "HTTP" {
		return nil, fmt.Errorf("Request line: invalid HTTP-version name")
	}
	if httpParts[1] != "1.1" {
		return nil, fmt.Errorf("Request line: invalid HTTP-version")
	}

	requestLine := &RequestLine{
		HttpVersion:   httpParts[1],
		RequestTarget: parts[1],
		Method:        parts[0],
	}

	return requestLine, nil
}
