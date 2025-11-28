package entities

import "time"

type Receipt struct {
	ID                string    `json:"id" bson:"_id,omitempty"`
	ShopID            string    `json:"shop_id" bson:"shop_id"`
	PaymentMethod     string    `json:"payment_method" bson:"payment_method"` // e.g., "cash"
	TotalAmount       float64   `json:"total_amount" bson:"total_amount"`
	VatRate           float64   `json:"vat_rate" bson:"vat_rate"` // e.g. 0.07
	VAT               float64   `json:"vat" bson:"vat"`
	NetTotal          float64   `json:"net_total" bson:"net_total"` // ยอดสุทธิ (Total + VAT)
	Status            string    `json:"status" bson:"status"`       // เช่น "completed", "cancelled"
	CustomerRequestID string    `json:"customer_request_id,omitempty" bson:"customer_request_id,omitempty"`
	CreatedAt         time.Time `json:"created_at" bson:"created_at"`
}

type ReceiptItem struct {
	ID        string  `json:"id" bson:"_id,omitempty"`
	ReceiptID string  `json:"receipt_id" bson:"receipt_id"`
	ShopID    string  `json:"shop_id" bson:"shop_id"`
	WasteID   string  `json:"waste_id" bson:"waste_id"`
	Name      string  `json:"name" bson:"name"`             // Snapshot of name
	Category  string  `json:"category" bson:"category"`     // Snapshot of category
	Weight    float64 `json:"weight" bson:"weight"`         // kg
	UnitPrice float64 `json:"unit_price" bson:"unit_price"` // Price per unit at that time
	Price     float64 `json:"price" bson:"price"`           // Total price for this item
}
