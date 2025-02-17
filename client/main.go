package client

import (
	"bufio"
	"fmt"
	"log"
	"net"
	m "simple-irc-server/message"
	"strings"
)

type Client struct {
	conn       net.Conn
	reader     *bufio.Reader
	address    string
	server     Server
	registered bool
	nickname   string
	username   string
	mode       string
	realname   string
}

func NewClient(conn net.Conn, server Server) *Client {
	return &Client{
		conn:       conn,
		reader:     bufio.NewReader(conn),
		address:    conn.RemoteAddr().String(),
		server:     server,
		registered: false,
		nickname:   "",
		username:   "",
		mode:       "",
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

func (c *Client) receiveMessage() (*m.Message, error) {
	rawMessage, err := c.reader.ReadString('\n')

	if err != nil {
		return nil, err
	}

	log.Printf("< %s\n", strings.TrimSpace(rawMessage))

	return m.NewMessage(rawMessage)
}

func (c *Client) HandleConnection() {
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

func (c *Client) handleMessage(message *m.Message) {
	if message == nil {
		log.Printf("! Empty message: %v\n", message)
		return
	}

	c.handleCommand(message)
}

func (c *Client) handleCommand(message *m.Message) {
	command := message.Command
	params := message.Params

	switch command {
	case "NICK":
		c.handleNickname(params)
	case "USER":
		c.handleUser(params)
	default:
		c.send(fmt.Sprintf(":irc.local 421 %s :Unknown command", command))
	}

	if c.shouldRegister() {
		c.register()
	}
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
	payload := fmt.Sprintf("%s\r\n", message)
	_, err := c.conn.Write([]byte(payload))

	if err != nil {
		log.Printf("! Error sending message: %v\n", err)
	}

	log.Printf("> %v", message)
}
