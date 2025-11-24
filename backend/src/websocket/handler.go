package websocket

import (
	"encoding/json"
	"log"
	"recycle-waste-management-backend/src/domain/entities"
	"time"

	"github.com/gofiber/contrib/websocket"
)

const (
	// Time allowed to write a message to the peer
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer
	pongWait = 60 * time.Second

	// Send pings to peer with this period (must be less than pongWait)
	pingPeriod = (pongWait * 9) / 10
)

func (c *Client) ReadPump() {
	defer func() {
		ChatHub.Unregister <- c
		c.Conn.Close()
	}()

	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error { c.Conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("[Chat] WebSocket error: %v", err)
			}
			break
		}

		log.Printf("[Chat] ReadPump: Received raw message from UserID=%s: %s", c.UserID, string(message))

		// Parse incoming message
		var msgReq entities.ChatMessageRequest
		var messageContent string

		if err := json.Unmarshal(message, &msgReq); err != nil {
			// If JSON parsing fails, treat as plain text
			messageContent = string(message)
			log.Printf("[Chat] ReadPump: Failed to parse JSON, using raw text: %s", messageContent)
		} else {
			messageContent = msgReq.Message
			log.Printf("[Chat] ReadPump: Parsed JSON message: %s", messageContent)
		}

		// Create chat message
		chatMsg := entities.ChatMessage{
			Type:              "message",
			CustomerRequestID: c.CustomerRequestID,
			SenderID:          c.UserID,
			SenderType:        c.UserType,
			Message:           messageContent,
			Timestamp:         time.Now().Format(time.RFC3339),
		}

		// Marshal and broadcast
		if msgBytes, err := json.Marshal(chatMsg); err == nil {
			log.Printf("[Chat] Broadcasting message to room %s from UserID=%s", c.CustomerRequestID, c.UserID)
			ChatHub.Broadcast <- &BroadcastMessage{
				CustomerRequestID: c.CustomerRequestID,
				Message:           msgBytes,
			}
			log.Printf("[Chat] Message from %s (%s) in room %s: %s",
				c.UserID, c.UserType, c.CustomerRequestID, messageContent)
		}
	}
}

func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// Hub closed the channel
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			log.Printf("[Chat] WritePump: Sending to UserID=%s Message=%s", c.UserID, string(message))
			if err := c.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
				log.Printf("[Chat] WritePump Error: %v", err)
				return
			}

		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
