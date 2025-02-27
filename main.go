package main

import (
	"log"
	s "simple-irc-server/server"
)

func main() {
	server := s.NewServer(6667)

	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}
