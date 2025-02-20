package server

import (
	"fmt"
	"sync"
)

type Usernames struct {
	list map[string]*Client
	mtx  sync.RWMutex
}

func NewUsernames() *Usernames {
	return &Usernames{
		list: make(map[string]*Client),
	}
}

func (names *Usernames) AddUsername(name string, client *Client) error {
	return names.add(name, client)
}

func (names *Usernames) RemoveUsername(name string) error {
	return names.remove(name)
}

func (names *Usernames) add(username string, client *Client) error {
	names.mtx.Lock()
	defer names.mtx.Unlock()

	if _, exists := names.list[username]; exists {
		return fmt.Errorf("username exists: %v", username)
	}

	names.list[username] = client

	return nil
}

func (names *Usernames) remove(username string) error {
	names.mtx.Lock()
	defer names.mtx.Unlock()

	if _, exists := names.list[username]; !exists {
		return fmt.Errorf("username does not exist: %v", username)
	}

	delete(names.list, username)

	return nil
}
