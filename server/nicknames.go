package server

import (
	"fmt"
	"sync"
)

type Nicknames struct {
	list map[string]*Client
	mtx  sync.RWMutex
}

func NewNicknames() *Nicknames {
	return &Nicknames{
		list: make(map[string]*Client),
	}
}

func (nicks *Nicknames) GetClients() []*Client {
	return nicks.getClients()
}

func (nicks *Nicknames) GetClientByNickname(nick string) (*Client, error) {
	return nicks.get(nick)
}

func (nicks *Nicknames) AddNickname(nick string, client *Client) error {
	return nicks.add(nick, client)
}

func (nicks *Nicknames) UpdateNickname(nick, newNick string) error {
	return nicks.rename(nick, newNick)
}

func (nicks *Nicknames) RemoveNickname(nick string) error {
	return nicks.remove(nick)
}

func (nicks *Nicknames) getClients() []*Client {
	nicks.mtx.RLock()
	defer nicks.mtx.RUnlock()

	clients := make([]*Client, len(nicks.list))

	for _, client := range nicks.list {
		clients = append(clients, client)
	}

	return clients
}

func (nicks *Nicknames) add(nickname string, client *Client) error {
	nicks.mtx.Lock()
	defer nicks.mtx.Unlock()

	if _, exists := nicks.list[nickname]; exists {
		return fmt.Errorf("nickname already exists: %v", nickname)
	}

	nicks.list[nickname] = client

	return nil
}

func (nicks *Nicknames) get(nickname string) (*Client, error) {
	nicks.mtx.RLock()
	defer nicks.mtx.RUnlock()

	client, exists := nicks.list[nickname]

	if !exists {
		return nil, fmt.Errorf("nickname does not exist: %v", nickname)
	}

	return client, nil
}

func (nicks *Nicknames) rename(nickname, newNick string) error {
	nicks.mtx.Lock()
	defer nicks.mtx.Unlock()

	if _, exists := nicks.list[newNick]; exists {
		return fmt.Errorf("nickname already exists: %v", newNick)
	}

	client := nicks.list[nickname]

	delete(nicks.list, nickname)

	nicks.list[newNick] = client

	return nil
}

func (nicks *Nicknames) remove(nickname string) error {
	nicks.mtx.Lock()
	defer nicks.mtx.Unlock()

	if _, exists := nicks.list[nickname]; !exists {
		return fmt.Errorf("nickname does not exist: %v", nickname)
	}

	delete(nicks.list, nickname)

	return nil
}
