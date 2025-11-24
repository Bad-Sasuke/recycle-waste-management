package gateways

import (
	"log"
	ws "recycle-waste-management-backend/src/websocket"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

func (h *HTTPGateway) HandleWebSocketChat(c *websocket.Conn) {
	// Get query parameters
	userID := c.Query("user_id")
	customerRequestID := c.Query("request_id")

	if userID == "" || customerRequestID == "" {
		log.Printf("[Chat] Missing required parameters")
		c.Close()
		return
	}

	log.Printf("[Chat] New connection: UserID=%s, Room=%s", userID, customerRequestID)

	// Create client
	client := &ws.Client{
		Conn:              c,
		UserID:            userID,
		UserType:          "user", // Default type since we removed it from query
		CustomerRequestID: customerRequestID,
		Send:              make(chan []byte, 256),
	}

	// Register client
	ws.ChatHub.Register <- client

	// Use a channel to wait for completion
	done := make(chan struct{})

	// Start read and write pumps in goroutines
	go func() {
		defer close(done)
		client.ReadPump()
	}()

	go client.WritePump()

	// Wait for ReadPump to finish (it will finish when connection closes)
	<-done

	log.Printf("[Chat] Connection closed: UserID=%s, Room=%s", userID, customerRequestID)
}

func (h *HTTPGateway) WebSocketChatUpgrade(c *fiber.Ctx) error {
	// Check if connection is WebSocket upgrade
	if websocket.IsWebSocketUpgrade(c) {
		return c.Next()
	}

	return c.Status(fiber.StatusUpgradeRequired).JSON(fiber.Map{
		"message": "WebSocket upgrade required",
	})
}
