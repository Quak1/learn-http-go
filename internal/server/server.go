package server

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"sync/atomic"

	"github.com/Quak1/learn-http-go/internal/request"
	"github.com/Quak1/learn-http-go/internal/response"
)

type Handler func(w io.Writer, req *request.Request) *HandlerError

type HandlerError struct {
	StatusCode response.StatusCode
	Message    string
}

func (h HandlerError) Write(w io.Writer) {
	response.WriteStatusLine(w, h.StatusCode)
	headers := response.GetDefaultHeaders(len(h.Message))
	response.WriteHeaders(w, headers)
	w.Write([]byte(h.Message))
}

type Server struct {
	handler      Handler
	listener     net.Listener
	serverClosed atomic.Bool
}

func Serve(port int, handler Handler) (*Server, error) {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}

	s := Server{
		listener: ln,
		handler:  handler,
	}

	go s.listen()

	return &s, nil
}

func (s *Server) Close() error {
	s.serverClosed.Store(true)
	return s.listener.Close()
}

func (s *Server) listen() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			if s.serverClosed.Load() {
				return
			}

			log.Println("Error: couldn't accept connection")
			continue
		}

		go s.handle(conn)
	}
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()

	req, err := request.RequestFromReader(conn)
	if err != nil {
		requestError := &HandlerError{
			StatusCode: response.StatusBadRequest,
			Message:    err.Error(),
		}
		requestError.Write(conn)
		return
	}

	var buff bytes.Buffer
	handlerError := s.handler(&buff, req)
	if handlerError != nil {
		handlerError.Write(conn)
		return
	}

	err = response.WriteStatusLine(conn, response.StatusOK)
	if err != nil {
		log.Println("Error: couldn't write status line.", err)
		return
	}

	b := buff.Bytes()
	headers := response.GetDefaultHeaders(len(b))
	err = response.WriteHeaders(conn, headers)
	if err != nil {
		log.Println("Error: couldn't write headers.", err)
		return
	}

	conn.Write(b)
}
