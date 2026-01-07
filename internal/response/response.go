package response

import (
	"fmt"
	"io"
	"strconv"

	"github.com/Quak1/learn-http-go/internal/headers"
)

type StatusCode int

const (
	StatusOK                  StatusCode = 200
	StatusBadRequest          StatusCode = 400
	StatusInternalServerError StatusCode = 500
)

const HTTPVersion = "HTTP/1.1"

type writerState int

var stateError = fmt.Errorf("Error: writer is in the wrong state for this operation make sure to use the methods in order. Writer.WriteStatusLine -> Writer.WriteHeaders -> Writer.WriteBody")

const (
	WriterStateStatusLine writerState = iota
	WriterStateHeaders
	WriterStateBody
	WriterStateDone
)

type Writer struct {
	writer io.Writer
	state  writerState
}

func NewResponseWriter(w io.Writer) *Writer {
	return &Writer{
		writer: w,
		state:  WriterStateStatusLine,
	}
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	if w.state != WriterStateStatusLine {
		return fmt.Errorf("Error: cannot write status line on state %d", w.state)
	}

	var reason string
	switch statusCode {
	case StatusOK:
		reason = "OK"
	case StatusBadRequest:
		reason = "Bad Request"
	case StatusInternalServerError:
		reason = "Internal Server Error"
	}

	_, err := fmt.Fprintf(w.writer, "%s %d %s\r\n", HTTPVersion, statusCode, reason)
	w.state = WriterStateHeaders

	return err
}

func (w *Writer) WriteHeaders(headers headers.Headers) error {
	if w.state != WriterStateStatusLine {
		return fmt.Errorf("Error: cannot write headers on state %d", w.state)
	}

	for v, k := range headers {
		_, err := fmt.Fprintf(w.writer, "%s: %s\r\n", v, k)
		if err != nil {
			return err
		}
	}

	_, err := w.writer.Write([]byte("\r\n"))
	w.state = WriterStateBody
	return err
}

func (w *Writer) WriteBody(p []byte) (int, error) {
	if w.state != WriterStateStatusLine {
		return 0, fmt.Errorf("Error: cannot write body on state %d", w.state)
	}

	return w.writer.Write(p)
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	h := headers.NewHeaders()
	h.Set("Content-Length", strconv.Itoa(contentLen))
	h.Set("Connection", "close")
	h.Set("Content-Type", "text/plain")
	return h
}
