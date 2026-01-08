package response

import (
	"fmt"
	"io"
	"strconv"

	"github.com/Quak1/learn-http-go/internal/headers"
)

type writerState int

const (
	WriterStateStatusLine writerState = iota
	WriterStateHeaders
	WriterStateBody
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

	reason := getStatusReason(statusCode)
	HTTPVersion := "HTTP/1.1"

	_, err := fmt.Fprintf(w.writer, "%s %d %s\r\n", HTTPVersion, statusCode, reason)
	w.state = WriterStateHeaders

	return err
}

func (w *Writer) WriteHeaders(headers headers.Headers) error {
	if w.state != WriterStateHeaders {
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
	if w.state != WriterStateBody {
		return 0, fmt.Errorf("Error: cannot write body on state %d", w.state)
	}

	return w.writer.Write(p)
}

func (w *Writer) WriteChunkedBody(p []byte) (int, error) {
	if w.state != WriterStateBody {
		return 0, fmt.Errorf("Error: cannot write body on state %d", w.state)
	}

	sizeHex := strconv.FormatInt(int64(len(p)), 16)
	n, err := fmt.Fprintf(w.writer, "%s\r\n", sizeHex)
	if err != nil {
		return 0, err
	}

	n2, err := w.writer.Write(p)
	if err != nil {
		return 0, nil
	}

	n3, err := w.writer.Write([]byte("\r\n"))
	if err != nil {
		return 0, nil
	}

	return n + n2 + n3, nil
}

func (w *Writer) WriteChunkedBodyDone() (int, error) {
	if w.state != WriterStateBody {
		return 0, fmt.Errorf("Error: cannot write body on state %d", w.state)
	}

	return w.writer.Write([]byte("0\r\n\r\n"))
}
