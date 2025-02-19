package server

import (
	"fmt"
	"regexp"
)

var messagePattern = regexp.MustCompile(`^(?:[:]((?:[a-zA-Z0-9\[\]\\` + "`" + `_^{|}][a-zA-Z0-9\[\]\\` + "`" + `_^{|}-]*))(?:(?:!([^@]+))?@([^ ]+))? )?([a-zA-Z]+|[0-9]{3}) (.+?)\r\n$`)

type Parts struct {
	prefix  string
	user    string
	host    string
	command string
	params  string
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
			prefix:  matches[1],
			user:    matches[2],
			host:    matches[3],
			command: matches[4],
			params:  matches[5],
		},
	}

	return message, nil
}
