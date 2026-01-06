package server

import (
	"fmt"
	"log"
	"net"
	"sync/atomic"

	"github.com/Quak1/learn-http-go/internal/response"
)

type Server struct {
	listener     net.Listener
	serverClosed atomic.Bool
}

func Serve(port int) (*Server, error) {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}

	s := Server{
		listener: ln,
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

	err := response.WriteStatusLine(conn, response.StatusOK)
	if err != nil {
		log.Println("Error: couldn't write status line.", err)
		return
	}

	headers := response.GetDefaultHeaders(0)
	err = response.WriteHeaders(conn, headers)
	if err != nil {
		log.Println("Error: couldn't write headers.", err)
		return
	}
}
