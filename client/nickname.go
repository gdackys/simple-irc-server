package client

import (
	"fmt"
	"regexp"
)

func (c *Client) handleNickname(params string) {
	pattern := regexp.MustCompile(`^[a-zA-Z\[\]\\` + "`" + `_^{|}][a-zA-Z0-9\[\]\\` + "`" + `_^{|}-]{0,8}$`)
	matches := pattern.FindStringSubmatch(params)

	if matches == nil {
		c.send(fmt.Sprint(":irc.local 432 * :Erroneous nickname"))
		return
	}

	nickname := matches[0]

	if nickname == "" {
		c.send(":irc.local 431 * :No nickname given")
		return
	}

	if c.hasNickname() {
		c.changeNickname(nickname)
	} else {
		c.setNickname(nickname)
	}
}

func (c *Client) hasNickname() bool {
	return len(c.nickname) > 0
}

func (c *Client) changeNickname(nick string) {
	if c.nickname == nick {
		return
	}

	if err := c.server.UpdateNickname(c.nickname, nick); err != nil {
		c.send(fmt.Sprintf(":irc.local 433 %s :Nickname is already in use", nick))
	} else {
		c.send(fmt.Sprintf(":%s NICK %s", c.fullIdentifier(), nick))
		c.nickname = nick
	}
}

func (c *Client) setNickname(nick string) {
	if err := c.server.InsertNickname(nick); err != nil {
		c.send(":irc.local 433 * :Nickname is already in use")
	} else {
		c.nickname = nick
	}
}

func (c *Client) unsetNickname() {
	if err := c.server.RemoveNickname(c.nickname); err == nil {
		c.nickname = ""
	}
}
