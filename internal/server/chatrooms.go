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

func (rms *Chatrooms) include(name string) bool {
	rms.mtx.RLock()
	defer rms.mtx.RUnlock()

	_, exists := rms.list[name]

	return exists
}

func (rms *Chatrooms) get(name string) (*Chatroom, error) {
	rms.mtx.RLock()
	defer rms.mtx.RUnlock()

	room, exists := rms.list[name]

	if !exists {
		return nil, fmt.Errorf("channel does not exist: %s", name)
	}

	return room, nil
}

func (rms *Chatrooms) create(name string) (*Chatroom, error) {
	rms.mtx.Lock()
	defer rms.mtx.Unlock()

	_, exists := rms.list[name]

	if exists {
		return nil, fmt.Errorf("channel already exists: %s", name)
	}

	room := NewChatroom(name)

	rms.list[name] = room

	return room, nil
}
