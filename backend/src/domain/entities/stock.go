package entities

import "time"

type Stock struct {
	ID        string    `json:"id" bson:"_id,omitempty"`
	ShopID    string    `json:"shop_id" bson:"shop_id"`
	WasteID   string    `json:"waste_id" bson:"waste_id"`
	Quantity  float64   `json:"quantity" bson:"quantity"` // Current balance
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}
