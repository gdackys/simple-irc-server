package message

import (
	"fmt"
	"regexp"
)

var messagePattern = regexp.MustCompile(`^(?:[:]((?:[a-zA-Z0-9\[\]\\` + "`" + `_^{|}][a-zA-Z0-9\[\]\\` + "`" + `_^{|}-]*))(?:(?:!([^@]+))?@([^ ]+))? )?([a-zA-Z]+|[0-9]{3}) (.+?)\r\n$`)

func Parse(msg string) (*Message, error) {
	matches := messagePattern.FindStringSubmatch(msg)

	if matches == nil {
		return nil, fmt.Errorf("invalid message: %v", msg)
	}

	parts := &Parts{
		prefix:  matches[1],
		user:    matches[2],
		host:    matches[3],
		command: matches[4],
		params:  matches[5],
	}

	message := &Message{parts}

	return message, nil
}
