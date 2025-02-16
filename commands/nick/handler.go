package nick

import (
	"log"
	c "simple-irc-server/command"
)

func Handler(ctx *c.Context[Params]) error {
	nickname := ctx.Params.nickname

	// WIP
	log.Printf("NICK: %s\n", nickname)

	return nil
}
