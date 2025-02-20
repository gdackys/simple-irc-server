package server

import (
	"fmt"
	"log"
	"net"
)

type Server struct {
	port int
	*Nicknames
	*Usernames
	*Chatrooms
}

func NewServer(port int) *Server {
	return &Server{
		port,
		NewNicknames(),
		NewUsernames(),
		NewChatrooms(),
	}
}

func (s *Server) stop(listener net.Listener) {
	if err := listener.Close(); err != nil {
		log.Printf("! Error stopping: %v\n", err)
	}

	log.Printf("~ Bye!\n")
}

func (s *Server) Start() error {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", s.port))

	if err != nil {
		return err
	}

	defer s.stop(listener)

	log.Printf("~ Listening on port %d\n", s.port)

	for {
		conn, err := listener.Accept()

		if err != nil {
			log.Printf("! Error accepting connection: %v\n", err)
			continue
		}

		log.Printf("~ Accepted new connection: %v\n", conn.RemoteAddr())

		client := NewClient(conn, s)

		go client.handleConnection()
	}
}
