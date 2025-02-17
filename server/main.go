package server

import (
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
	Execute(ctx *Context)
}

type Context struct {
	Message *m.Message
	Client  *c.Client
	Server  *Server
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
		return err
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

func (s *Server) handleClient(client *c.Client) {
	defer client.Disconnect()

	for {
		message, err := client.ReceiveMessage()

		if err != nil {
			logger.Printf("Error receiving message: %v\n", err)
			break
		}

		s.handleMessage(message, client)
	}
}

func (s *Server) handleMessage(msg *m.Message, client *c.Client) {
	if msg == nil {
		logger.Printf("Empty message: %v\n", msg)
		return
	}

	logger.Printf("Received message: %v\n", msg)

	s.handleCommand(msg, client)
}

func (s *Server) handleCommand(msg *m.Message, client *c.Client) {
	name := msg.Command()
	command, exists := s.commands[name]

	if !exists {
		logger.Printf("Unsupported command: %v\n", name)
		return
	}

	context := &Context{msg, client, s}

	command.Execute(context)
}

func (s *Server) stop(listener net.Listener) {
	if err := listener.Close(); err != nil {
		logger.Printf("Error stopping: %v\n", err)
	}

	logger.Printf("Stopped\n")
}
