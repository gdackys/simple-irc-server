package message

import (
	"fmt"
	"regexp"
)

var messagePattern = regexp.MustCompile(`^(?:[:]((?:[a-zA-Z0-9\[\]\\` + "`" + `_^{|}][a-zA-Z0-9\[\]\\` + "`" + `_^{|}-]*))(?:(?:!([^@]+))?@([^ ]+))? )?([a-zA-Z]+|[0-9]{3}) (.+?)\r\n$`)

type Parts struct {
	Prefix  string
	User    string
	Host    string
	Command string
	Params  string
}

type Message struct {
	*Parts
}

func NewMessage(msg string) (*Message, error) {
	matches := messagePattern.FindStringSubmatch(msg)

	if matches == nil {
		return nil, fmt.Errorf("invalid message: %v", msg)
	}

	message := &Message{
		&Parts{
			Prefix:  matches[1],
			User:    matches[2],
			Host:    matches[3],
			Command: matches[4],
			Params:  matches[5],
		},
	}

	return message, nil
}
