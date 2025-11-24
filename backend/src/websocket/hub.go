package websocket

import (
	"encoding/json"
	"log"
	"recycle-waste-management-backend/src/domain/entities"
	"sync"
	"time"

	"github.com/gofiber/contrib/websocket"
)

type Client struct {
	Conn              *websocket.Conn
	UserID            string
	UserType          string // "customer" or "shop"
	CustomerRequestID string // Room ID
	Send              chan []byte
}

type Hub struct {
	// Registered clients per room
	Rooms map[string]map[*Client]bool

	// Register requests from clients
	Register chan *Client

	// Unregister requests from clients
	Unregister chan *Client

	// Broadcast messages to clients in a room
	Broadcast chan *BroadcastMessage

	// Mutex for thread-safe operations
	mutex sync.RWMutex
}

type BroadcastMessage struct {
	CustomerRequestID string
	Message           []byte
}

var ChatHub *Hub

func NewHub() *Hub {
	return &Hub{
		Rooms:      make(map[string]map[*Client]bool),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan *BroadcastMessage),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.mutex.Lock()
			if h.Rooms[client.CustomerRequestID] == nil {
				h.Rooms[client.CustomerRequestID] = make(map[*Client]bool)
			}
			h.Rooms[client.CustomerRequestID][client] = true
			roomSize := len(h.Rooms[client.CustomerRequestID])
			h.mutex.Unlock()

			log.Printf("[Chat] Client registered: UserID=%s, Type=%s, Room=%s (Total in room: %d)",
				client.UserID, client.UserType, client.CustomerRequestID, roomSize)

			// Send join notification in a goroutine to avoid blocking
			go func(c *Client) {
				joinMsg := entities.ChatMessage{
					Type:              "join",
					CustomerRequestID: c.CustomerRequestID,
					SenderID:          c.UserID,
					SenderType:        c.UserType,
					Message:           "joined the chat",
					Timestamp:         time.Now().Format(time.RFC3339),
				}
				if msgBytes, err := json.Marshal(joinMsg); err == nil {
					h.Broadcast <- &BroadcastMessage{
						CustomerRequestID: c.CustomerRequestID,
						Message:           msgBytes,
					}
				}
			}(client)

			// Send existing users presence to the new client
			// We need to do this carefully to avoid race conditions, but since we have the list of clients
			// from the map (which we accessed under lock), we can iterate and send.
			// However, we need to be careful not to send to closed channels if they disconnect in between.
			// But here we are just sending to 'client.Send' which is our own channel.

			// We need to capture the existing clients info, not the pointers themselves if we want to be super safe,
			// but for this simple case, we can just iterate the map inside the lock or copy it.
			// We already have the map access inside the lock above.
			// Let's modify the logic above to capture existing clients.

			// Since I cannot easily modify the code above the block I'm replacing without replacing the whole block,
			// I will assume I can't access 'existingUsers' unless I captured them.
			// I'll have to re-acquire lock or just do it in a separate goroutine with a new lock acquisition?
			// No, re-acquiring lock is fine.

			go func(target *Client) {
				h.mutex.RLock()
				var others []*Client
				if room, ok := h.Rooms[target.CustomerRequestID]; ok {
					for c := range room {
						if c != target {
							others = append(others, c)
						}
					}
				}
				h.mutex.RUnlock()

				for _, other := range others {
					presenceMsg := entities.ChatMessage{
						Type:              "join",
						CustomerRequestID: other.CustomerRequestID,
						SenderID:          other.UserID,
						SenderType:        other.UserType,
						Message:           "is online",
						Timestamp:         time.Now().Format(time.RFC3339),
					}
					if msgBytes, err := json.Marshal(presenceMsg); err == nil {
						target.Send <- msgBytes
					}
				}
			}(client)

		case client := <-h.Unregister:
			h.mutex.Lock()
			if clients, ok := h.Rooms[client.CustomerRequestID]; ok {
				if _, exists := clients[client]; exists {
					delete(clients, client)
					close(client.Send)

					// Clean up empty rooms
					if len(clients) == 0 {
						delete(h.Rooms, client.CustomerRequestID)
					}
				}
			}
			h.mutex.Unlock()

			log.Printf("[Chat] Client unregistered: UserID=%s, Type=%s, Room=%s",
				client.UserID, client.UserType, client.CustomerRequestID)

			// Send leave notification in a goroutine to avoid blocking
			go func(c *Client) {
				leaveMsg := entities.ChatMessage{
					Type:              "leave",
					CustomerRequestID: c.CustomerRequestID,
					SenderID:          c.UserID,
					SenderType:        c.UserType,
					Message:           "left the chat",
					Timestamp:         time.Now().Format(time.RFC3339),
				}
				if msgBytes, err := json.Marshal(leaveMsg); err == nil {
					h.Broadcast <- &BroadcastMessage{
						CustomerRequestID: c.CustomerRequestID,
						Message:           msgBytes,
					}
				}
			}(client)

		case broadcast := <-h.Broadcast:
			h.mutex.RLock()
			roomClients := h.Rooms[broadcast.CustomerRequestID]
			clientCount := len(roomClients)

			// Copy clients to a slice to avoid holding the lock during broadcast
			// and to avoid concurrent map modification issues
			var clients []*Client
			for client := range roomClients {
				clients = append(clients, client)
			}
			h.mutex.RUnlock()

			log.Printf("[Chat] Broadcasting message to room %s (Clients: %d)", broadcast.CustomerRequestID, clientCount)

			for _, client := range clients {
				select {
				case client.Send <- broadcast.Message:
					log.Printf("[Chat] Sent to client UserID=%s", client.UserID)
				default:
					log.Printf("[Chat] Failed to send to client UserID=%s (Channel full)", client.UserID)
					// Client's send channel is full, unregister
					h.mutex.Lock()
					if _, ok := h.Rooms[broadcast.CustomerRequestID][client]; ok {
						delete(h.Rooms[broadcast.CustomerRequestID], client)
						close(client.Send)
					}
					h.mutex.Unlock()
				}
			}
		}
	}
}

func InitChatHub() {
	ChatHub = NewHub()
	go ChatHub.Run()
	log.Println("[Chat] WebSocket Hub initialized and running")
}
