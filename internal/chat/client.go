package chat

type Client struct {
	ID        string
	MessageCh chan Message
}

func NewClient(id string) *Client {
	return &Client{
		ID:        id,
		MessageCh: make(chan Message, 20),
	}
}
