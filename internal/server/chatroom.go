package server

type Chatroom struct {
	name    string
	clients map[string]*Client
}

func NewChatroom(name string) *Chatroom {
	return &Chatroom{
		name:    name,
		clients: make(map[string]*Client),
	}
}

func (room *Chatroom) addClient(client *Client) {
	room.clients[client.username] = client
}

func (room *Chatroom) nicknames() []string {
	result := make([]string, 0, len(room.clients))

	for _, client := range room.clients {
		result = append(result, client.nickname)
	}

	return result
}

func (room *Chatroom) broadcast(message string) {
	for _, client := range room.clients {
		client.send(message)
	}
}
