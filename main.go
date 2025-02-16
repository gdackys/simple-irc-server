package main

import (
	"log"
	"simple-irc-server/commands"
	s "simple-irc-server/server"
)

func main() {
	server := s.NewServer()

	server.RegisterCommand(commands.Nick)

	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}
