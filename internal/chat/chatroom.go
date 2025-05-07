package chat

import (
	"sync"
)

type ChatRoom struct {
	Clients   map[string]*Client
	MessageCh chan Message
	Mu        sync.RWMutex
	joinCh    chan *Client
	leaveCh   chan string
}

func NewChatRoom() *ChatRoom {
	return &ChatRoom{
		Clients:   make(map[string]*Client),
		MessageCh: make(chan Message, 100),
		joinCh:    make(chan *Client, 10),
		leaveCh:   make(chan string, 10),
	}
}

func (r *ChatRoom) Run() {
	for {
		select {
		case client := <-r.joinCh:
			r.Mu.Lock()
			r.Clients[client.ID] = client
			r.Mu.Unlock()
		case id := <-r.leaveCh:
			r.Mu.Lock()
			if c, ok := r.Clients[id]; ok {
				close(c.MessageCh)
				delete(r.Clients, id)
			}
			r.Mu.Unlock()
		case msg := <-r.MessageCh:
			r.Mu.RLock()
			for _, c := range r.Clients {
				// Don't send to sender
				if c.ID != msg.SenderID {
					select {
					case c.MessageCh <- msg:
					default:
					}
				}
			}
			r.Mu.RUnlock()
		}
	}
}

func (r *ChatRoom) HandleJoin(client *Client) {
	r.joinCh <- client
}

func (r *ChatRoom) HandleLeave(id string) {
	r.leaveCh <- id
}
