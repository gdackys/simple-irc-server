package server

import (
	"fmt"
	"log"
	"net"
	c "simple-irc-server/client"
)

type Server struct {
	port      int
	nicknames *Nicknames
	usernames *Usernames
}

func NewServer(port int) *Server {
	return &Server{
		port:      port,
		nicknames: NewNicknames(),
		usernames: NewUsernames(),
	}
}

func (s *Server) stop(listener net.Listener) {
	if err := listener.Close(); err != nil {
		log.Printf("! Error stopping: %v\n", err)
	}

	log.Printf("~ Stopped\n")
}

func (s *Server) Start() error {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", s.port))

	defer s.stop(listener)

	if err != nil {
		return err
	}

	log.Printf("~ Listening on port %d\n", s.port)

	for {
		conn, err := listener.Accept()

		if err != nil {
			log.Printf("! Error accepting connection: %v\n", err)
			continue
		}

		log.Printf("~ Accepted new connection: %v\n", conn.RemoteAddr())

		client := c.NewClient(conn, s)

		go client.HandleConnection()
	}
}

func (s *Server) InsertNickname(nick string) error {
	return s.nicknames.Insert(nick)
}

func (s *Server) UpdateNickname(nick, newNick string) error {
	return s.nicknames.Rename(nick, newNick)
}

func (s *Server) RemoveNickname(nick string) error {
	return s.nicknames.Remove(nick)
}

func (s *Server) InsertUsername(name string) error {
	return s.usernames.Insert(name)
}

func (s *Server) RemoveUsername(name string) error {
	return s.usernames.Remove(name)
}
