package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"chat-app/internal/chat"
	"chat-app/internal/handler"

	"github.com/gin-gonic/gin"
)

func main() {
	// Create a new ChatRoom.
	room := chat.NewChatRoom()
	go room.Run()

	// Create a new Handler instance.
	h := handler.NewHandler(room)

	// Create a new Gin router.
	router := gin.Default()

	// Use handler methods directly.
	router.GET("/join", h.Join)
	router.GET("/send", h.SendMessage)
	router.GET("/leave", h.Leave)
	router.GET("/messages", h.GetMessages)

	// Start the HTTP server.
	port := "3000"
	server := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	go func() {
		fmt.Printf("Server listening on port %s\n", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Handle graceful shutdown.
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh
	fmt.Println("Shutting down server...")

	if err := server.Close(); err != nil {
		log.Printf("Error closing server: %v", err)
	}
	fmt.Println("Server stopped.")
}
