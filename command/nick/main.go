package nick

import (
	"fmt"
	"log"
	"regexp"
	"simple-irc-server/command/base"
	m "simple-irc-server/message"
	s "simple-irc-server/server"
)

type Nick struct {
	*base.Base
}

type Params struct {
	nickname string
}

const (
	ErrNoNicknameGiven = 431
	ErrNicknameInUse   = 433
)

func New() *Nick {
	return &Nick{
		base.New("NICK"),
	}
}

func (cmd *Nick) Execute(ctx *s.Context) {
	params, err := cmd.parseParams(ctx.Message)

	if err != nil {
		log.Printf("Error parsing command params: %s\n", err)
		return
	}

	cmd.run(params)
}

func (cmd *Nick) parseParams(msg *m.Message) (*Params, error) {
	rawParams := msg.Params()
	pattern := regexp.MustCompile(`^[a-zA-Z\[\]\\` + "`" + `_^{|}][a-zA-Z0-9\[\]\\` + "`" + `_^{|}-]{0,8}$`)
	matches := pattern.FindStringSubmatch(rawParams)

	if matches == nil {
		return nil, fmt.Errorf("invalid params: %v", rawParams)
	}

	return &Params{matches[0]}, nil
}

func (cmd *Nick) run(params *Params) {
	log.Printf("NICK: %s\n", params.nickname)
}
