package socket

import "github.com/olahol/melody"

type Channel struct {
	Clients map[*melody.Session]bool
}

func NewChannel() *Channel {
	return &Channel{
		Clients: make(map[*melody.Session]bool),
	}
}

func (c *Channel) AddClient(s *melody.Session) {
	c.Clients[s] = true
}

func (c *Channel) RemoveClient(s *melody.Session) {
	delete(c.Clients, s)
}

func (c *Channel) Broadcast(message []byte) {
	for client := range c.Clients {
		err := client.Write(message)
		if err != nil {
			c.RemoveClient(client)
		}
	}
}
