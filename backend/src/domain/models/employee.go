package models

import "time"

type Employee struct {
	EmployeeID string    `json:"employee_id,omitempty" bson:"employee_id,omitempty"`
	ShopID     string    `json:"shop_id,omitempty" bson:"shop_id,omitempty"`
	FirstName  string    `json:"first_name,omitempty" bson:"first_name,omitempty"`
	LastName   string    `json:"last_name,omitempty" bson:"last_name,omitempty"`
	Username   string    `json:"username,omitempty" bson:"username,omitempty"`
	Password   string    `json:"password,omitempty" bson:"password,omitempty"` // Hashed password
	CreatedAt  time.Time `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt  time.Time `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}
