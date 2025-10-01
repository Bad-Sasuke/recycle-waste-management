package entities

import "time"

type ShopModel struct {
	ShopID      string    `json:"shop_id,omitempty" bson:"shop_id,omitempty"`
	UserID      string    `json:"user_id,omitempty" bson:"user_id,omitempty"`
	Name        string    `json:"name,omitempty" bson:"name,omitempty"`
	Description string    `json:"description,omitempty" bson:"description,omitempty"`
	Address     string    `json:"address,omitempty" bson:"address,omitempty"`
	Phone       string    `json:"phone,omitempty" bson:"phone,omitempty"`
	Email       string    `json:"email,omitempty" bson:"email,omitempty"`
	ImageURL    string    `json:"image_url,omitempty" bson:"image_url,omitempty"`
	OpeningTime string    `json:"opening_time,omitempty" bson:"opening_time,omitempty"`
	ClosingTime string    `json:"closing_time,omitempty" bson:"closing_time,omitempty"`
	Latitude    float64   `json:"latitude,omitempty" bson:"latitude,omitempty"`
	Longitude   float64   `json:"longitude,omitempty" bson:"longitude,omitempty"`
	CreatedAt   time.Time `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt   time.Time `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

type CreateShopRequest struct {
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description"`
	Address     string  `json:"address" validate:"required"`
	Phone       string  `json:"phone"`
	Email       string  `json:"email" validate:"omitempty,email"`
	OpeningTime string  `json:"opening_time"`
	ClosingTime string  `json:"closing_time"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
}

type UpdateShopRequest struct {
	Name        *string  `json:"name,omitempty"`
	Description *string  `json:"description,omitempty"`
	Address     *string  `json:"address,omitempty"`
	Phone       *string  `json:"phone,omitempty"`
	Email       *string  `json:"email,omitempty"`
	OpeningTime *string  `json:"opening_time,omitempty"`
	ClosingTime *string  `json:"closing_time,omitempty"`
	Latitude    *float64 `json:"latitude,omitempty"`
	Longitude   *float64 `json:"longitude,omitempty"`
}
