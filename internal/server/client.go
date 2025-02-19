package server

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"regexp"
	"strings"
)

type Client struct {
	conn       net.Conn
	reader     *bufio.Reader
	address    string
	server     *Server
	registered bool
	nickname   string
	username   string
	mode       string
	realname   string
	chatrooms  map[string]*Chatroom
}

func NewClient(conn net.Conn, server *Server) *Client {
	return &Client{
		conn:       conn,
		reader:     bufio.NewReader(conn),
		address:    conn.RemoteAddr().String(),
		server:     server,
		registered: false,
		nickname:   "",
		username:   "",
		mode:       "",
		realname:   "",
		chatrooms:  make(map[string]*Chatroom),
	}
}

func (c *Client) disconnect() {
	if err := c.conn.Close(); err != nil {
		log.Printf("! Error disconnecting: %v\n", err)
		return
	}

	c.unsetNickname()
	c.unsetUser()

	log.Printf("~ Disconnected from %s\n", c.conn.RemoteAddr())
}

func (c *Client) receiveMessage() (*Message, error) {
	rawMessage, err := c.reader.ReadString('\n')

	if err != nil {
		return nil, err
	}

	log.Printf("< %s\n", strings.TrimSpace(rawMessage))

	return NewMessage(rawMessage)
}

func (c *Client) handleConnection() {
	defer c.disconnect()

	for {
		message, err := c.receiveMessage()

		if err != nil {
			log.Printf("! Error receiving message: %v\n", err)
			break
		}

		c.handleMessage(message)
	}
}

func (c *Client) handleMessage(message *Message) {
	if message == nil {
		log.Printf("! Empty message: %v\n", message)
		return
	}

	c.handleCommand(message)
}

func (c *Client) handleCommand(message *Message) {
	command := message.command
	params := message.params

	switch command {
	case "NICK":
		c.handleNickname(params)
	case "USER":
		c.handleUser(params)
	case "JOIN":
		c.handleJoin(params)
	default:
		c.send(fmt.Sprintf(":irc.local 421 %s :Unknown command", command))
	}

	if c.shouldRegister() {
		c.register()
	}
}

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

	if err := c.server.updateNickname(c.nickname, nick); err != nil {
		c.send(fmt.Sprintf(":irc.local 433 %s :Nickname is already in use", nick))
	} else {
		c.send(fmt.Sprintf(":%s NICK %s", c.fullIdentifier(), nick))
		c.nickname = nick
	}
}

func (c *Client) setNickname(nick string) {
	if err := c.server.addNickname(nick); err != nil {
		c.send(":irc.local 433 * :Nickname is already in use")
	} else {
		c.nickname = nick
	}
}

func (c *Client) unsetNickname() {
	if err := c.server.removeNickname(c.nickname); err == nil {
		c.nickname = ""
	}
}

func (c *Client) handleUser(params string) {
	pattern := regexp.MustCompile(`^([^\x00\r\n@ ]+)\s+(\d+)\s+(\S+)\s+:(.+)$`)
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
	if err := c.server.addUsername(username); err != nil {
		c.send(":irc.local 462 :Unauthorized command (already registered)")
	} else {
		c.username = username
		c.mode = mode
		c.realname = realname
	}
}

func (c *Client) unsetUser() {
	if err := c.server.removeUsername(c.username); err == nil {
		c.username = ""
		c.mode = ""
		c.realname = ""
	}
}

func (c *Client) handleJoin(params string) {
	if !c.registered {
		c.send(fmt.Sprintf(":irc.local 451 %s :You have not registered", c.nickname))
		return
	}

	pattern := regexp.MustCompile(`^(#[^\x00\x07\x0a\x0d ,]{1,49}(?:,#[^\x00\x07\x0a\x0d ,]{1,49})*)$`)
	matches := pattern.FindStringSubmatch(params)

	if matches == nil {
		c.send(fmt.Sprintf(":irc.local 403 %s %s :No such channel", c.nickname, params))
		return
	}

	chatrooms := strings.Split(matches[1], ",")

	for _, chatroom := range chatrooms {
		c.joinChatroom(chatroom)
	}
}

func (c *Client) joinChatroom(name string) {
	chatroom, err := c.server.getChatroom(name)

	if err != nil {
		c.send(fmt.Sprintf(":irc.local 403 %s %s :No such channel", c.nickname, name))
		return
	}

	c.addChatroom(chatroom)
	chatroom.addClient(c)

	c.send(fmt.Sprintf(":irc.local 331 %s %s :No topic is set", c.nickname, name))
	c.send(fmt.Sprintf(":irc.local 353 %s = %s :%s", c.nickname, name, strings.Join(chatroom.nicknames(), " ")))
	c.send(fmt.Sprintf(":irc.local 366 %s %s :End of NAMES list", c.nickname, name))

	chatroom.broadcast(fmt.Sprintf(":%s JOIN %s", c.fullIdentifier(), name))
}

func (c *Client) addChatroom(room *Chatroom) {
	c.chatrooms[room.name] = room
}

func (c *Client) shouldRegister() bool {
	return c.registered == false && c.hasUsername() && c.hasNickname()
}

func (c *Client) register() {
	c.sendWelcome()
	c.registered = true
}

func (c *Client) sendWelcome() {
	c.send(fmt.Sprintf(":irc.local 001 %s :Welcome to the Internet Relay Network %s", c.nickname, c.fullIdentifier()))
	c.send(fmt.Sprintf(":irc.local 002 %s :Your host is irc.local, running version 1.00", c.nickname))
	c.send(fmt.Sprintf(":irc.local 003 %s :This server was created Feb 17 2025", c.nickname))
	c.send(fmt.Sprintf(":irc.local 004 %s irc.local 1.0 iwso itkol", c.nickname))
}

func (c *Client) fullIdentifier() string {
	return fmt.Sprintf("%s!%s@%s", c.nickname, c.username, c.address)
}

func (c *Client) send(message string) {
	_, err := fmt.Fprintf(c.conn, "%s\r\n", message)

	if err != nil {
		log.Printf("! Error sending message: %v\n", err)
	}

	log.Printf("> %v", message)
}
