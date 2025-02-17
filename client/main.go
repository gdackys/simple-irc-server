package client

import (
	"bufio"
	"log"
	"net"
	"os"
	m "simple-irc-server/message"
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

func (c *Client) ReceiveMessage() (*m.Message, error) {
	rawMessage, err := c.reader.ReadString('\n')

	if err != nil {
		return nil, err
	}

	logger.Printf("Incoming message: %s\n", strings.TrimSpace(rawMessage))

	return m.Parse(rawMessage)
}
