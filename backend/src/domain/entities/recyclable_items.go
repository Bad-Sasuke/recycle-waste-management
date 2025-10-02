package entities

import "time"

type RecyclableItemsModel struct {
	WasteID    string    `json:"waste_id,omitempty" bson:"waste_id,omitempty"`
	ShopID     string    `json:"shop_id,omitempty" bson:"shop_id,omitempty"`
	Name       string    `json:"name,omitempty" bson:"name,omitempty"`
	Category   string    `json:"category,omitempty" bson:"category,omitempty"`
	Price      float64   `json:"price,omitempty" bson:"price,omitempty"`
	LastUpdate time.Time `json:"last_update,omitempty" bson:"last_update,omitempty"`
	Hours      string    `json:"hours,omitempty" bson:"hours,omitempty"`
	URL        string    `json:"url,omitempty" bson:"url,omitempty"`
}
