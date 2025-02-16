package server

import (
	"fmt"
	"log"
	"net"
	"os"
	c "simple-irc-server/client"
	m "simple-irc-server/message"
)

var logger = log.New(os.Stdout, "SERVER ", log.LstdFlags|log.Lmicroseconds|log.Lmsgprefix)

type Server struct {
	commands Commands
}

type Commands = map[string]Command

type Command interface {
	GetName() string
	Execute(*m.Message) error
}

func NewServer() *Server {
	return &Server{
		commands: make(Commands),
	}
}

func (s *Server) RegisterCommand(cmd Command) {
	s.commands[cmd.GetName()] = cmd
}

func (s *Server) Start() error {
	listener, err := net.Listen("tcp", ":6667")

	defer s.stop(listener)

	if err != nil {
		return fmt.Errorf("error starting: %v", err)
	}

	logger.Printf("Listening on port 6667\n")

	for {
		conn, err := listener.Accept()

		if err != nil {
			logger.Printf("Error accepting connection: %v\n", err)
			continue
		}

		logger.Printf("Accepted new connection: %v\n", conn.RemoteAddr())

		client := c.NewClient(conn)

		go s.handleClient(client)
	}
}

func (s *Server) stop(listener net.Listener) {
	err := listener.Close()

	if err != nil {
		logger.Printf("Error stopping: %v\n", err)
	}

	logger.Printf("Stopped\n")
}

func (s *Server) handleClient(client *c.Client) {
	defer client.Disconnect()

	for {
		message, err := client.ReceiveMessage()

		if err != nil {
			logger.Printf("Error receiving message: %v\n", err)
			break
		}

		if err := s.handleMessage(message); err != nil {
			logger.Printf("Error handling message: %v\n", err)
		}
	}
}

func (s *Server) handleMessage(message *m.Message) error {
	if message == nil {
		return fmt.Errorf("invalid message: %v", message)
	}

	logger.Printf("Received message: %v\n", message)

	return s.handleCommand(message)
}

func (s *Server) handleCommand(message *m.Message) error {
	name := message.Command()
	command, exists := s.commands[name]

	if !exists {
		return fmt.Errorf("unsupported command: %v", name)
	}

	if err := command.Execute(message); err != nil {
		return fmt.Errorf("error executing %v: %v", name, err)
	}

	return nil
}
