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

// ShopInfo contains information about a shop
type ShopInfo struct {
	ShopID       string `json:"shop_id"`
	ShopName     string `json:"shop_name"`
	ShopImageURL string `json:"shop_image_url"`
}

// GroupedRecyclableItem represents a product grouped by name with multiple shop_ids
type GroupedRecyclableItem struct {
	Name       string             `json:"name"`
	Category   string             `json:"category"`
	Price      float64            `json:"price"`
	LastUpdate time.Time          `json:"last_update"`
	Hours      string             `json:"hours,omitempty"`
	URL        string             `json:"url,omitempty"`
	Shops      []ShopInfo         `json:"shops"`     // Array of shop info that have this product
	WasteIDs   []string           `json:"waste_ids"` // Array of waste IDs that have this product name
}
