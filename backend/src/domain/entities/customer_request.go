package entities

import "time"

type CustomerRequest struct {
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	Description string  `json:"description"`
}

type CustomerRequestResponse struct {
	CustomerRequestID string    `json:"customer_request_id"`
	Latitude          float64   `json:"latitude"`
	Longitude         float64   `json:"longitude"`
	Description       string    `json:"description"`
	Status            string    `json:"status"`
	CancelReason      string    `json:"cancel_reason,omitempty"`
	Distance          float64   `json:"distance"` // Distance in kilometers
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

type CancelCustomerRequest struct {
	Reason string `json:"reason" validate:"required,max=100"`
}

type PaginatedCustomerRequestResponse struct {
	Data       []CustomerRequestResponse `json:"data"`
	Total      int                       `json:"total"`
	Page       int                       `json:"page"`
	Limit      int                       `json:"limit"`
	TotalPages int                       `json:"total_pages"`
}

type WalkInCustomerRequest struct {
	ShopID       string `json:"shop_id"`
	CustomerName string `json:"customer_name"`
	PhoneNumber  string `json:"phone_number"`
}
