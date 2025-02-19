package main

import (
	"log"
	s "simple-irc-server/internal/server"
)

func main() {
	server := s.NewServer(6667)

	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}
