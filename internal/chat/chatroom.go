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
	Messages  []Message // Store all messages
}

func NewChatRoom() *ChatRoom {
	return &ChatRoom{
		Clients:   make(map[string]*Client),
		MessageCh: make(chan Message, 100),
		joinCh:    make(chan *Client, 10),
		leaveCh:   make(chan string, 10),
		Messages:  make([]Message, 0),
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
			r.Mu.Lock()
			r.Messages = append(r.Messages, msg) // Store message
			r.Mu.Unlock()
			r.Mu.RLock()
			for _, c := range r.Clients {
				select {
				case c.MessageCh <- msg:
				default:
				}
			}
			r.Mu.RUnlock()
		}
	}
}

func (r *ChatRoom) HandleJoin(client *Client) {
	r.joinCh <- client
	// Broadcast a system message that this client joined
	r.MessageCh <- Message{
		SenderID: "system",
		Text:     client.ID + " joined",
	}
}

func (r *ChatRoom) HandleLeave(id string) {
	r.leaveCh <- id
	// Broadcast a system message that this client left
	r.MessageCh <- Message{
		SenderID: "system",
		Text:     id + " left",
	}
}
