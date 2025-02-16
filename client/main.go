package client

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	message3 "simple-irc-server/message"
	"strings"
)

var logger = log.New(os.Stdout, "CLIENT ", log.LstdFlags|log.Lmicroseconds|log.Lmsgprefix)

type Client struct {
	conn   net.Conn
	reader *bufio.Reader
}

func NewClient(conn net.Conn) *Client {
	return &Client{
		conn:   conn,
		reader: bufio.NewReader(conn),
	}
}

func (c *Client) Disconnect() {
	if err := c.conn.Close(); err != nil {
		logger.Printf("Error disconnecting: %v\n", err)
		return
	}

	logger.Printf("Disconnected from %s\n", c.conn.RemoteAddr())
}

func (c *Client) ReceiveMessage() (*message3.Message, error) {
	rawMessage, err := c.reader.ReadString('\n')

	if err != nil {
		return nil, fmt.Errorf("error receiving message: %v", err)
	}

	logger.Printf("Raw message: %s\n", strings.TrimSpace(rawMessage))

	return c.parseRawMessage(rawMessage), nil
}

func (c *Client) parseRawMessage(msg string) *message3.Message {
	message, err := message3.Parse(msg)

	if err != nil {
		logger.Printf("Error parsing message: %v\n", err)
		return nil
	}

	return message
}
