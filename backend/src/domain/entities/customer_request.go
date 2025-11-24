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
	Distance          float64   `json:"distance"` // Distance in kilometers
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

type PaginatedCustomerRequestResponse struct {
	Data       []CustomerRequestResponse `json:"data"`
	Total      int                       `json:"total"`
	Page       int                       `json:"page"`
	Limit      int                       `json:"limit"`
	TotalPages int                       `json:"total_pages"`
}
