package server

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"sync/atomic"

	"github.com/Quak1/learn-http-go/internal/request"
	"github.com/Quak1/learn-http-go/internal/response"
)

type Handler func(w *response.Writer, req *request.Request)

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
		return
	}

	var buff bytes.Buffer
	writer := response.NewResponseWriter(&buff)
	s.handler(writer, req)

	b := buff.Bytes()
	conn.Write(b)
}
