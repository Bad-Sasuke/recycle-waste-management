package models

import "time"

type STATUS_REQUEST string

const (
	CR_PENDING   STATUS_REQUEST = "pending"
	CR_ACCEPTED  STATUS_REQUEST = "accepted"
	CR_REJECTED  STATUS_REQUEST = "rejected"
	CR_DONE      STATUS_REQUEST = "done"
	CR_CANCELLED STATUS_REQUEST = "cancelled"
)

type CustomerRequestModel struct {
	UserID            string         `json:"user_id" bson:"user_id"`
	CustomerRequestID string         `json:"customer_request_id" bson:"customer_request_id"`
	Latitude          float64        `json:"latitude" bson:"latitude"`
	Longitude         float64        `json:"longitude" bson:"longitude"`
	Description       string         `json:"description" bson:"description"`
	ShopID            string         `json:"shop_id,omitempty" bson:"shop_id,omitempty"` // Shop that completed the request
	Status            STATUS_REQUEST `json:"status" bson:"status"`
	CancelReason      string         `json:"cancel_reason,omitempty" bson:"cancel_reason,omitempty"`
	CreatedAt         time.Time      `json:"created_at" bson:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at" bson:"updated_at"`
}
