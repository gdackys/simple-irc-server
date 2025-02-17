package main

import (
	"log"
	"simple-irc-server/command"
	s "simple-irc-server/server"
)

func main() {
	server := s.NewServer()

	server.RegisterCommand(command.Nick())

	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}
