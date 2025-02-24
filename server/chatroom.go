package server

import (
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

func (room *Chatroom) addClient(client *Client) {
	room.mtx.Lock()
	defer room.mtx.Unlock()

	room.clients[client] = true
}

func (room *Chatroom) removeClient(client *Client) {
	room.mtx.Lock()
	defer room.mtx.Unlock()

	delete(room.clients, client)
}

func (room *Chatroom) sendToAll(message string) {
	for client := range room.clients {
		client.send(message)
	}
}

func (room *Chatroom) broadcast(source *Client, message string) {
	for client := range room.clients {
		if client != source {
			client.send(message)
		}
	}
}

func (room *Chatroom) nicknames() string {
	room.mtx.RLock()
	defer room.mtx.RUnlock()

	result := make([]string, 0, len(room.clients))

	for client := range room.clients {
		result = append(result, client.nickname)
	}

	return strings.Join(result, " ")
}
