package server

import (
	"fmt"
	"sync"
)

type Usernames struct {
	list map[string]bool
	mtx  sync.RWMutex
}

func NewUsernames() *Usernames {
	return &Usernames{
		list: make(map[string]bool),
	}
}

func (names *Usernames) Insert(username string) error {
	names.mtx.Lock()
	defer names.mtx.Unlock()

	if _, exists := names.list[username]; exists {
		return fmt.Errorf("username exists: %v", username)
	}

	names.list[username] = true

	return nil
}

func (names *Usernames) Remove(username string) error {
	names.mtx.Lock()
	defer names.mtx.Unlock()

	if _, exists := names.list[username]; !exists {
		return fmt.Errorf("username does not exist: %v", username)
	}

	delete(names.list, username)

	return nil
}
