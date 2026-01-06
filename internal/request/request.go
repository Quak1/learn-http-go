package request

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/Quak1/learn-http-go/internal/headers"
)

const bufferSize = 8

type parserState int

const (
	stateInitialized parserState = iota
	stateDone
	stateParsingHeaders
	stateParsingBody
)

type Request struct {
	RequestLine RequestLine
	Headers     headers.Headers
	Body        []byte
	state       parserState
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	buff := make([]byte, bufferSize)
	readToIndex := 0
	request := &Request{
		state: stateInitialized,
	}

	for request.state != stateDone {
		if readToIndex >= len(buff) {
			newBuff := make([]byte, len(buff)*2)
			copy(newBuff, buff)
			buff = newBuff
		}

		n, err := reader.Read(buff[readToIndex:])
		if err != nil {
			if errors.Is(err, io.EOF) {
				if request.state != stateDone {
					return nil, fmt.Errorf("Error: incomplete request")
				}
				break
			}
			return nil, err
		}

		readToIndex += n

		n, err = request.parse(buff[:readToIndex])
		if err != nil {
			return nil, err
		}

		copy(buff, buff[n:])
		readToIndex -= n
	}

	return request, nil
}

func (r *Request) parse(data []byte) (int, error) {
	totalBytesParsed := 0
	for r.state != stateDone {
		n, err := r.parseSingle(data[totalBytesParsed:])
		if err != nil {
			return 0, err
		}
		if n == 0 {
			break
		}
		totalBytesParsed += n
	}

	return totalBytesParsed, nil
}

func (r *Request) parseSingle(data []byte) (int, error) {
	switch r.state {
	case stateInitialized:
		reqLine, n, err := parseRequestLine(data)
		if err != nil {
			return 0, err
		}
		if n == 0 {
			return 0, nil
		}

		r.state = stateParsingHeaders
		r.RequestLine = *reqLine
		r.Headers = headers.NewHeaders()
		return n, nil
	case stateParsingHeaders:
		n, done, err := r.Headers.Parse(data)
		if err != nil {
			return 0, err
		}
		if done {
			r.state = stateParsingBody
		}
		return n, nil
	case stateParsingBody:
		contentLengthStr := r.Headers.Get("Content-Length")
		if contentLengthStr == "" {
			r.state = stateDone
			return len(data), nil
		}

		contentLength, err := strconv.Atoi(contentLengthStr)
		if err != nil {
			return 0, fmt.Errorf("Error: invalid Content-Length value, not a number")
		}
		r.Body = append(r.Body, data...)

		bodyLen := len(r.Body)
		if bodyLen > contentLength {
			return 0, fmt.Errorf("Error: available body data larger than Content-Length value")
		}
		if bodyLen == contentLength {
			r.state = stateDone
		}

		return len(data), nil
	case stateDone:
		return 0, fmt.Errorf("Error: trying to read data in a done state")
	default:
		return 0, fmt.Errorf("Error: unknown state")
	}
}

func parseRequestLine(data []byte) (*RequestLine, int, error) {
	idx := bytes.Index(data, []byte("\r\n"))
	if idx == -1 {
		return nil, 0, nil
	}

	parts := strings.Split(string(data[:idx]), " ")
	if len(parts) != 3 {
		return nil, 0, fmt.Errorf("Request line: wrong number of parts")
	}

	method := parts[0]
	for _, c := range method {
		if c < 'A' || c > 'Z' {
			return nil, 0, fmt.Errorf("Request line: invalid method")
		}
	}

	httpParts := strings.Split(parts[2], "/")
	if len(httpParts) != 2 {
		return nil, 0, fmt.Errorf("Request line: invalid HTTP-version")
	}
	if httpParts[0] != "HTTP" {
		return nil, 0, fmt.Errorf("Request line: invalid HTTP-version name")
	}
	if httpParts[1] != "1.1" {
		return nil, 0, fmt.Errorf("Request line: invalid HTTP-version")
	}

	requestLine := &RequestLine{
		HttpVersion:   httpParts[1],
		RequestTarget: parts[1],
		Method:        parts[0],
	}

	return requestLine, idx + 2, nil
}
