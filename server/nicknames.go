package server

import (
	"fmt"
	"sync"
)

type Nicknames struct {
	list map[string]bool
	mtx  sync.RWMutex
}

func NewNicknames() *Nicknames {
	return &Nicknames{
		list: make(map[string]bool),
	}
}

func (nicks *Nicknames) Insert(nickname string) error {
	nicks.mtx.Lock()
	defer nicks.mtx.Unlock()

	if _, exists := nicks.list[nickname]; exists {
		return fmt.Errorf("nickname exists: %v", nickname)
	}

	nicks.list[nickname] = true

	return nil
}

func (nicks *Nicknames) Rename(nickname, newNick string) error {
	nicks.mtx.Lock()
	defer nicks.mtx.Unlock()

	if _, exists := nicks.list[newNick]; exists {
		return fmt.Errorf("nickname exists: %v", newNick)
	}

	delete(nicks.list, nickname)
	nicks.list[newNick] = true

	return nil
}

func (nicks *Nicknames) Remove(nickname string) error {
	nicks.mtx.Lock()
	defer nicks.mtx.Unlock()

	if _, exists := nicks.list[nickname]; !exists {
		return fmt.Errorf("nickname does not exist: %v", nickname)
	}

	delete(nicks.list, nickname)

	return nil
}
