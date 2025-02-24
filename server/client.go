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

func (c *Client) disconnect() {
	if err := c.conn.Close(); err != nil {
		log.Printf("! Error disconnecting: %v\n", err)
		return
	}

	c.unsetNickname()
	c.unsetUser()
	c.exitChatrooms()

	log.Printf("~ Disconnected from %s\n", c.conn.RemoteAddr())
}

func (c *Client) receiveMessage() (*Message, error) {
	rawMessage, err := c.reader.ReadString('\n')

	if err != nil {
		return nil, err
	}

	log.Printf("< :%s %s\n", c, strings.TrimSpace(rawMessage))

	return NewMessage(rawMessage)
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
	case "PRIVMSG":
		c.handlePrivmsg(params)
	case "PART":
		c.handlePart(params)
	case "QUIT":
		c.handleQuit(params)
	}

	if c.shouldRegister() {
		c.register()
	}
}

/* REGISTRATION */

func (c *Client) shouldRegister() bool {
	return c.registered == false && c.hasUsername() && c.hasNickname()
}

func (c *Client) register() {
	c.sendWelcome()
	c.registered = true
}

func (c *Client) sendWelcome() {
	c.send(fmt.Sprintf(":irc.local 001 %s :Welcome to the Internet Relay Network %s", c.nickname, c))
	c.send(fmt.Sprintf(":irc.local 002 %s :Your host is irc.local, running version 1.00", c.nickname))
	c.send(fmt.Sprintf(":irc.local 003 %s :This server was created Feb 17 2025", c.nickname))
	c.send(fmt.Sprintf(":irc.local 004 %s irc.local 1.0", c.nickname))
}

/* NICK */

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
		c.send(fmt.Sprintf(":%s NICK %s", c, nick))
		c.nickname = nick
	}
}

func (c *Client) setNickname(nick string) {
	if err := c.server.AddNickname(nick, c); err != nil {
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

/* USER */

func (c *Client) handleUser(params string) {
	pattern := regexp.MustCompile(`^([^\x00\r\n@ ]+)\s+(\d+)\s+(\S+)\s+:(.+)$`)
	matches := pattern.FindStringSubmatch(params)

	if matches == nil {
		c.send(fmt.Sprintf(":irc.local 461 %s USER :Not enough parameters", c.nickname))
		return
	}

	username := matches[1]
	mode := matches[2]
	realname := matches[4]

	if c.hasUsername() {
		c.send(fmt.Sprintf(":irc.local 462 %s :Unauthorized command (already registered)", c.nickname))
	} else {
		c.setUser(username, mode, realname)
	}
}

func (c *Client) hasUsername() bool {
	return c.username != ""
}

func (c *Client) setUser(username, mode, realname string) {
	if err := c.server.AddUsername(username, c); err != nil {
		c.send(fmt.Sprintf(":irc.local 462 %s :Unauthorized command (already registered)", c.nickname))
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

/* JOIN */

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

	roomNames := strings.Split(matches[1], ",")

	for _, name := range roomNames {
		c.joinChatroom(name)
	}
}

func (c *Client) joinChatroom(name string) {
	chatroom, err := c.server.GetChatroom(name)

	if err != nil {
		c.send(fmt.Sprintf(":irc.local 403 %s %s :No such channel", c.nickname, name))
		return
	}

	chatroom.addClient(c)
	c.addChatroom(chatroom)

	chatroom.sendToAll(fmt.Sprintf(":%s JOIN %s", c, name))

	c.send(fmt.Sprintf(":irc.local 331 %s %s :No topic is set", c.nickname, name))
	c.send(fmt.Sprintf(":irc.local 353 %s = %s :%s", c.nickname, name, chatroom.nicknames()))
	c.send(fmt.Sprintf(":irc.local 366 %s %s :End of NAMES list", c.nickname, name))
}

/* PRIVMSG */

func (c *Client) handlePrivmsg(params string) {
	if !c.registered {
		c.send(fmt.Sprintf(":irc.local 451 * :You have not registered"))
		return
	}

	pattern := regexp.MustCompile(`^([#\w][^\s,]+)\s+:(.+)$`)
	matches := pattern.FindStringSubmatch(params)

	if matches == nil {
		c.send(fmt.Sprintf(":irc.local 411 %s :No recipient given (PRIVMSG)", c.nickname))
		return
	}

	target, message := matches[1], matches[2]

	if target == "" {
		c.send(fmt.Sprintf(":irc.local 411 %s :No recipient given (PRIVMSG)", c.nickname))
		return
	}

	if message == "" {
		c.send(fmt.Sprintf(":irc.local 412 %s :No text to send", c.nickname))
		return
	}

	if strings.HasPrefix(target, "#") {
		c.sendToChatroom(target, message)
	} else {
		c.sendToClient(target, message)
	}
}

func (c *Client) sendToChatroom(name, message string) {
	chatroom, exists := c.chatrooms[name]

	if !exists {
		c.send(fmt.Sprintf(":irc.local 404 %s %s :Cannot send to channel", c.nickname, name))
		return
	}

	chatroom.broadcast(c, fmt.Sprintf(":%s PRIVMSG %s :%s", c, name, message))
}

func (c *Client) sendToClient(nickname, message string) {
	client, err := c.server.GetClientByNickname(nickname)

	if err != nil {
		c.send(fmt.Sprintf(":irc.local 401 %s %s :No such nick/channel", c.nickname, nickname))
		return
	}

	client.send(fmt.Sprintf(":%s PRIVMSG %s :%s", c, client.nickname, message))
}

/* PART */

func (c *Client) handlePart(params string) {
	if !c.registered {
		c.send(fmt.Sprintf(":irc.local 451 %s :You have not registered", c.nickname))
		return
	}

	pattern := regexp.MustCompile(`^(#[^\x00\x07\x0a\x0d ,]{1,49}(?:,#[^\x00\x07\x0a\x0d ,]{1,49})*)(?:\s+:(.+))?$`)
	matches := pattern.FindStringSubmatch(params)

	if matches == nil {
		c.send(fmt.Sprintf(":irc.local 461 %s PART :Not enough parameters", c.nickname))
		return
	}

	roomNames := strings.Split(matches[1], ",")
	partNotice := ""

	if len(matches) > 2 && matches[2] != "" {
		partNotice = matches[2]
	}

	for _, name := range roomNames {
		c.partChatroom(name, partNotice)
	}
}

func (c *Client) partChatroom(name, partNotice string) {
	chatroom, exists := c.chatrooms[name]

	if !exists {
		c.send(fmt.Sprintf(":irc.local 442 %s %s :You're not on that channel", c.nickname, name))
		return
	}

	quitMessage := fmt.Sprintf(":%s PART %s", c, chatroom.name)

	if partNotice != "" {
		quitMessage = fmt.Sprintf("%s :%s", quitMessage, partNotice)
	}

	chatroom.sendToAll(quitMessage)

	c.exitChatroom(chatroom)
}

/* QUIT */

func (c *Client) handleQuit(params string) {
	quitMessage := "Quit"

	pattern := regexp.MustCompile(`^:(.+)$`)
	matches := pattern.FindStringSubmatch(params)

	if matches != nil && len(matches) > 1 {
		quitMessage = matches[1]
	}

	quitNotice := fmt.Sprintf(":%s QUIT :%s", c, quitMessage)

	for _, chatroom := range c.chatrooms {
		chatroom.sendToAll(quitNotice)
	}

	c.send(fmt.Sprintf("ERROR :Closing Link: %s (%s)", c.address, quitMessage))

	c.disconnect()
}

/* MISC */

func (c *Client) id() string {
	return fmt.Sprintf("%s!%s@%s", c.nickname, c.username, c.address)
}

func (c *Client) String() string {
	return c.id()
}

func (c *Client) send(message string) {
	_, err := fmt.Fprintf(c.conn, "%s\r\n", message)

	if err != nil {
		log.Printf("! Error sending message: %v\n", err)
	}

	log.Printf("> %v", message)
}

func (c *Client) exitChatrooms() {
	for _, chatroom := range c.chatrooms {
		c.exitChatroom(chatroom)
	}
}

func (c *Client) exitChatroom(chatroom *Chatroom) {
	chatroom.removeClient(c)
	c.removeChatroom(chatroom)
}

func (c *Client) addChatroom(room *Chatroom) {
	c.chatrooms[room.name] = room
}

func (c *Client) removeChatroom(chatroom *Chatroom) {
	delete(c.chatrooms, chatroom.name)
}
