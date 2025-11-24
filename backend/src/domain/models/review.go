package models

import "time"

type ReviewModel struct {
	ReviewID          string    `json:"review_id" bson:"review_id"`
	CustomerRequestID string    `json:"customer_request_id" bson:"customer_request_id"`
	ShopID            string    `json:"shop_id" bson:"shop_id"`
	UserID            string    `json:"user_id" bson:"user_id"`
	Rating            int       `json:"rating" bson:"rating"`         // 1-5 stars
	Comment           string    `json:"comment" bson:"comment"`       // max 200 characters
	IsSkipped         bool      `json:"is_skipped" bson:"is_skipped"` // true if user clicked "skip"
	CreatedAt         time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt         time.Time `json:"updated_at" bson:"updated_at"`
}
