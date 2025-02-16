package message

import "fmt"

type Message struct {
	parts *Parts
}

func (m *Message) String() string {
	return fmt.Sprintf("Message{%v}", m.parts)
}

func (m *Message) Command() string {
	return m.parts.command
}

func (m *Message) Params() string {
	return m.parts.params
}
