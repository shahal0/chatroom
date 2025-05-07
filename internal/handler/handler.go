package handler

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"chat-app/internal/chat"

	"github.com/gin-gonic/gin"
)

// Handler struct to hold the chat room reference
type Handler struct {
	Room *chat.ChatRoom
}

// NewHandler creates a new Handler instance
func NewHandler(room *chat.ChatRoom) *Handler {
	return &Handler{Room: room}
}

// Join handles the /join endpoint for Gin.
func (h *Handler) Join(c *gin.Context) {
	id := c.Query("client_id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Missing client ID"})
		return
	}

	h.Room.Mu.RLock()
	_, exists := h.Room.Clients[id]
	h.Room.Mu.RUnlock()
	if exists {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Client already joined"})
		return
	}

	client := chat.NewClient(id)
	h.Room.HandleJoin(client)

	c.JSON(http.StatusOK, gin.H{"success": true, "message": fmt.Sprintf("Client %s joined the chat room", id)})
}

// SendMessage handles the /send endpoint for Gin.
func (h *Handler) SendMessage(c *gin.Context) {
	id := c.Query("client_id")
	message := c.Query("message")
	if id == "" || message == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Missing client ID or message"})
		return
	}

	h.Room.Mu.RLock()
	client, ok := h.Room.Clients[id]
	h.Room.Mu.RUnlock()
	if !ok || client == nil {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "error": "Client not joined or already left"})
		return
	}

	decodedMessage, err := url.QueryUnescape(message)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": fmt.Sprintf("Error decoding message: %s", err.Error())})
		return
	}

	msg := chat.Message{
		SenderID: id,
		Text:     decodedMessage,
	}
	h.Room.MessageCh <- msg

	c.JSON(http.StatusOK, gin.H{"success": true, "message": fmt.Sprintf("Message sent from client %s", id)})
}

// Leave handles the /leave endpoint for Gin.
func (h *Handler) Leave(c *gin.Context) {
	id := c.Query("client_id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Missing client ID"})
		return
	}

	h.Room.Mu.RLock()
	_, ok := h.Room.Clients[id]
	h.Room.Mu.RUnlock()
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Client not joined or already left"})
		return
	}

	h.Room.HandleLeave(id)
	c.JSON(http.StatusOK, gin.H{"success": true, "message": fmt.Sprintf("Client %s left the chat room", id)})
}

// GetMessages handles the /messages endpoint for Gin.
func (h *Handler) GetMessages(c *gin.Context) {
	id := c.Query("client_id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Missing client ID"})
		return
	}

	h.Room.Mu.RLock()
	client, ok := h.Room.Clients[id]
	h.Room.Mu.RUnlock()
	if !ok || client == nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": "Client not found or already left"})
		return
	}

	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("X-Accel-Buffering", "no")

	flusher, ok := c.Writer.(http.Flusher)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Streaming not supported"})
		return
	}

	timeout := time.After(300 * time.Second)

	for {
		select {
		case message := <-client.MessageCh:
			log.Printf("Sending message to client %s: %s", id, message.Text)
			formatted := fmt.Sprintf("%s:%s", message.SenderID, message.Text)
			fmt.Fprintf(c.Writer, "data: %s\n\n", strings.ReplaceAll(formatted, "\n", "\\n"))
			flusher.Flush()
		case <-c.Request.Context().Done():
			log.Printf("Client %s disconnected", id)
			return
		case <-timeout:
			log.Printf("Timeout for client %s", id)
			fmt.Fprintf(c.Writer, "data: %s\n\n", "timeout")
			flusher.Flush()
			return
		}
	}
}
