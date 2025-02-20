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

func (rooms *Chatrooms) GetChatroom(name string) (*Chatroom, error) {
	if rooms.include(name) {
		return rooms.get(name)
	} else {
		return rooms.create(name)
	}
}

func (rooms *Chatrooms) include(name string) bool {
	rooms.mtx.RLock()
	defer rooms.mtx.RUnlock()

	_, exists := rooms.list[name]

	return exists
}

func (rooms *Chatrooms) get(name string) (*Chatroom, error) {
	rooms.mtx.RLock()
	defer rooms.mtx.RUnlock()

	room, exists := rooms.list[name]

	if !exists {
		return nil, fmt.Errorf("channel does not exist: %s", name)
	}

	return room, nil
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
