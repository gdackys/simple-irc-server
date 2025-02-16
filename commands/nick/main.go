package nick

import (
	s "simple-irc-server/server"
)

var Command = &s.Command{
	Name:         "NICK",
	ParamsParser: paramsParser,
	Handler:      handler,
}
