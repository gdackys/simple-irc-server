package client

import (
	"regexp"
)

func (c *Client) handleUser(params string) {
	pattern := regexp.MustCompile(`^(\S+)\s+(\d+)\s+(\S+)\s+:(.+)$`)
	matches := pattern.FindStringSubmatch(params)

	if matches == nil {
		c.send(":irc.local 461 USER :Not enough parameters")
		return
	}

	username := matches[1]
	mode := matches[2]
	realname := matches[4]

	if c.hasUsername() {
		c.send(":irc.local 462 :Unauthorized command (already registered)")
	} else {
		c.setUser(username, mode, realname)
	}
}

func (c *Client) hasUsername() bool {
	return c.username != ""
}

func (c *Client) setUser(username, mode, realname string) {
	if err := c.server.InsertUsername(username); err != nil {
		c.send(":irc.local 462 :Unauthorized command (already registered)")
	} else {
		c.username = username
		c.mode = mode
		c.realname = realname
	}
}

func (c *Client) unsetUser() {
	if err := c.server.RemoveUsername(c.username); err == nil {
		c.username = ""
		c.mode = ""
		c.realname = ""
	}
}
