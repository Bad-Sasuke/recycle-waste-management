package entities

import "time"

type RecyclableItemsModel struct {
	Name       string    `json:"name,omitempty" bson:"name,omitempty"`
	Category   string    `json:"category,omitempty" bson:"category,omitempty"`
	Price      float64   `json:"price,omitempty" bson:"price,omitempty"`
	LastUpdate time.Time `json:"lastUpdate,omitempty" bson:"lastUpdate,omitempty"`
	Hours      string    `json:"hours,omitempty" bson:"hours,omitempty"`
	URL        string    `json:"url,omitempty" bson:"url,omitempty"`
}
