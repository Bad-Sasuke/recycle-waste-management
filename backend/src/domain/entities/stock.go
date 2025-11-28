package entities

import "time"

type Stock struct {
	ID         string    `json:"id" bson:"_id,omitempty"`
	ShopID     string    `json:"shop_id" bson:"shop_id"`
	WasteID    string    `json:"waste_id" bson:"waste_id"`
	Quantity   float64   `json:"quantity" bson:"quantity"` // Current balance
	UpdatedAt  time.Time `json:"updated_at" bson:"updated_at"`
	Category   string    `json:"category" bson:"category"`         // Snapshot from waste
	Name       string    `json:"name" bson:"name"`                 // Snapshot from waste
	PricePerKg float64   `json:"price_per_kg" bson:"price_per_kg"` // Snapshot from waste
}

// StockWithDetails includes waste information for display
type StockWithDetails struct {
	ID            string    `json:"id"`
	ShopID        string    `json:"shop_id"`
	WasteID       string    `json:"waste_id"`
	Quantity      float64   `json:"quantity"`
	UpdatedAt     time.Time `json:"updated_at"`
	Category      string    `json:"category"`
	Name          string    `json:"name"`
	PurchasePrice float64   `json:"purchase_price"` // ราคาที่รับซื้อ (จาก Stock Snapshot)
	CurrentPrice  float64   `json:"current_price"`  // ราคาปัจจุบัน (จาก Waste Live Data)
	Profit        float64   `json:"profit"`         // ส่วนต่าง (กำไร/ขาดทุน)
}
