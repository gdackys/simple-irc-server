package nick

import (
	c "simple-irc-server/command"
)

type Params struct {
	nickname string
}

var Command = &c.Command[Params]{
	Name:         "NICK",
	ParamsParser: ParamsParser,
	Handler:      Handler,
}
