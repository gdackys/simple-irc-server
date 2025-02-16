package nick

import (
	"log"
	s "simple-irc-server/server"
)

func handler(ctx *s.Context) error {
	nickname, err := ctx.Param("nickname")

	if err != nil {
		log.Printf("invalid nickname: %v", err)
	}

	log.Printf("NICK: %s\n", nickname)

	return nil
}
