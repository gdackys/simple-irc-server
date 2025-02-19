package server

import (
	"fmt"
	"sync"
)

type Chatrooms struct {
	list map[string]*Chatroom
	mtx  sync.RWMutex
}

func NewChatrooms() *Chatrooms {
	return &Chatrooms{
		list: make(map[string]*Chatroom),
	}
}

func (rooms *Chatrooms) contain(name string) bool {
	rooms.mtx.RLock()
	defer rooms.mtx.RUnlock()

	_, exists := rooms.list[name]

	return exists
}

func (rooms *Chatrooms) get(name string) (*Chatroom, error) {
	rooms.mtx.RLock()
	defer rooms.mtx.RUnlock()

	channel, exists := rooms.list[name]

	if !exists {
		return nil, fmt.Errorf("channel does not exist: %s", name)
	}

	return channel, nil
}

func (rooms *Chatrooms) create(name string) (*Chatroom, error) {
	rooms.mtx.Lock()
	defer rooms.mtx.Unlock()

	_, exists := rooms.list[name]

	if exists {
		return nil, fmt.Errorf("channel already exists: %s", name)
	}

	room := NewChatroom(name)

	rooms.list[name] = room

	return room, nil
}
