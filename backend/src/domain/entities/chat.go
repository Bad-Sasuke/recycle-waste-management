package entities

type ChatMessage struct {
	Type              string `json:"type"`                // "message", "join", "leave"
	CustomerRequestID string `json:"customer_request_id"` // Room ID
	SenderID          string `json:"sender_id"`
	SenderType        string `json:"sender_type"` // "customer" or "shop"
	Message           string `json:"message"`
	Timestamp         string `json:"timestamp"`
}

type ChatMessageRequest struct {
	CustomerRequestID string `json:"customer_request_id"`
	Message           string `json:"message"`
}
