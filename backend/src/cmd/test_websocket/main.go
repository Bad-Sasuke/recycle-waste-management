package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
)

type ChatMessageRequest struct {
	Message string `json:"message"`
}

type ChatMessage struct {
	Type              string `json:"type"`
	CustomerRequestID string `json:"customer_request_id"`
	SenderID          string `json:"sender_id"`
	SenderType        string `json:"sender_type"`
	Message           string `json:"message"`
	Timestamp         string `json:"timestamp"`
}

func main() {
	serverAddr := flag.String("addr", "127.0.0.1:1818", "http service address")
	flag.Parse()

	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)

	// User 1
	u1 := url.URL{Scheme: "ws", Host: *serverAddr, Path: "/ws/chat", RawQuery: "user_id=5LAYNlL3MiaDsQ6GS9huVP81qvB3&request_id=8c12bcd9-ee03-43fb-98cf-cfd96b4ce713"}
	log.Printf("🔌 Connecting User 1 to %s", u1.String())

	c1, _, err := websocket.DefaultDialer.Dial(u1.String(), nil)
	if err != nil {
		log.Fatal("❌ Failed to dial User 1:", err)
	}
	defer c1.Close()
	log.Println("✅ User 1 connected")

	// User 2 (Simulating a second connection, maybe same user ID but different connection)
	// The user provided the SAME URL for both.
	u2 := url.URL{Scheme: "ws", Host: *serverAddr, Path: "/ws/chat", RawQuery: "user_id=5LAYNlL3MiaDsQ6GS9huVP81qvB3&request_id=8c12bcd9-ee03-43fb-98cf-cfd96b4ce713"}
	log.Printf("🔌 Connecting User 2 to %s", u2.String())

	c2, _, err := websocket.DefaultDialer.Dial(u2.String(), nil)
	if err != nil {
		log.Fatal("❌ Failed to dial User 2:", err)
	}
	defer c2.Close()
	log.Println("✅ User 2 connected")

	done := make(chan struct{})
	user2ReceivedFromUser1 := false
	user1ReceivedFromUser2 := false

	// Read loop for User 1
	go func() {
		for {
			_, message, err := c1.ReadMessage()
			if err != nil {
				log.Println("❌ User 1 read error:", err)
				return
			}
			log.Printf("📨 User 1 received: %s", message)

			var chatMsg ChatMessage
			if err := json.Unmarshal(message, &chatMsg); err == nil {
				if chatMsg.Message == "Hello from User 2" {
					log.Println("🎉 SUCCESS: User 1 received message from User 2")
					user1ReceivedFromUser2 = true
				}
			}
		}
	}()

	// Read loop for User 2
	go func() {
		defer close(done)
		for {
			_, message, err := c2.ReadMessage()
			if err != nil {
				log.Println("❌ User 2 read error:", err)
				return
			}
			log.Printf("📨 User 2 received: %s", message)

			var chatMsg ChatMessage
			if err := json.Unmarshal(message, &chatMsg); err == nil {
				if chatMsg.Message == "Hello from User 1" {
					log.Println("🎉 SUCCESS: User 2 received message from User 1")
					user2ReceivedFromUser1 = true
				}
			}
		}
	}()

	// Wait a bit for join messages
	time.Sleep(1 * time.Second)

	// Send message from User 1
	msg1 := ChatMessageRequest{Message: "Hello from User 1"}
	msgBytes1, _ := json.Marshal(msg1)
	err = c1.WriteMessage(websocket.TextMessage, msgBytes1)
	if err != nil {
		log.Println("❌ User 1 write error:", err)
		return
	}
	log.Println("📤 User 1 sent: Hello from User 1")

	// Wait a bit
	time.Sleep(1 * time.Second)

	// Send message from User 2
	msg2 := ChatMessageRequest{Message: "Hello from User 2"}
	msgBytes2, _ := json.Marshal(msg2)
	err = c2.WriteMessage(websocket.TextMessage, msgBytes2)
	if err != nil {
		log.Println("❌ User 2 write error:", err)
		return
	}
	log.Println("📤 User 2 sent: Hello from User 2")

	// Wait for messages to be received
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	select {
	case <-done:
		log.Println("⏹️  Done channel closed")
	case <-time.After(5 * time.Second):
		log.Println("⏱️  Timeout after 5 seconds")
	case <-interrupt:
		log.Println("⚠️  Interrupted")
	}

	// Summary
	log.Println("\n========== TEST RESULTS ==========")
	log.Printf("User 1 received from User 2: %v", user1ReceivedFromUser2)
	log.Printf("User 2 received from User 1: %v", user2ReceivedFromUser1)

	if user1ReceivedFromUser2 && user2ReceivedFromUser1 {
		log.Println("✅ TEST PASSED: Both users can see each other's messages")
	} else {
		log.Println("❌ TEST FAILED: Users cannot see each other's messages")
		log.Println("This is the bug that needs to be fixed!")
	}
	log.Println("===================================")
}
