package server

import (
	"fmt"
	"strings"
	"sync"
)

type Chatroom struct {
	name    string
	clients map[*Client]bool
	mtx     sync.RWMutex
}

func NewChatroom(name string) *Chatroom {
	return &Chatroom{
		name:    name,
		clients: make(map[*Client]bool),
	}
}

func (cr *Chatroom) addClient(client *Client) error {
	cr.mtx.Lock()
	defer cr.mtx.Unlock()

	_, exists := cr.clients[client]

	if exists {
		return fmt.Errorf("client already exist: %s", client.id())
	}

	cr.clients[client] = true

	return nil
}

func (cr *Chatroom) removeClient(client *Client) error {
	cr.mtx.Lock()
	defer cr.mtx.Unlock()

	_, exists := cr.clients[client]

	if !exists {
		return fmt.Errorf("client does not exist: %s", client.id())
	}

	delete(cr.clients, client)

	return nil
}

func (cr *Chatroom) sendToAll(message string) {
	for client := range cr.clients {
		client.send(message)
	}
}

func (cr *Chatroom) broadcast(source *Client, message string) {
	for client := range cr.clients {
		if client != source {
			client.send(message)
		}
	}
}

func (cr *Chatroom) nicknames() string {
	cr.mtx.RLock()
	defer cr.mtx.RUnlock()

	result := make([]string, 0, len(cr.clients))

	for client := range cr.clients {
		result = append(result, client.nickname)
	}

	return strings.Join(result, " ")
}
