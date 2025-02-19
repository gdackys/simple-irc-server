package server

import (
	"fmt"
	"log"
	"net"
)

type Server struct {
	port      int
	nicknames *Nicknames
	usernames *Usernames
	chatrooms *Chatrooms
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

		client := NewClient(conn, s)

		go client.handleConnection()
	}
}

func (s *Server) addNickname(nick string) error {
	return s.nicknames.add(nick)
}

func (s *Server) updateNickname(nick, newNick string) error {
	return s.nicknames.rename(nick, newNick)
}

func (s *Server) removeNickname(nick string) error {
	return s.nicknames.remove(nick)
}

func (s *Server) addUsername(name string) error {
	return s.usernames.add(name)
}

func (s *Server) removeUsername(name string) error {
	return s.usernames.remove(name)
}

func (s *Server) getChatroom(name string) (*Chatroom, error) {
	if s.chatrooms.contain(name) {
		return s.chatrooms.get(name)
	} else {
		return s.chatrooms.create(name)
	}
}
