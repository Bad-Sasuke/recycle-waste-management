package entities

import "time"

type CreateEmployeeRequest struct {
	ShopID    string `json:"shop_id" bson:"shop_id"`
	FirstName string `json:"first_name" bson:"first_name" validate:"required"`
	LastName  string `json:"last_name" bson:"last_name" validate:"required"`
	Username  string `json:"username" bson:"username" validate:"required"`
	Password  string `json:"password" bson:"password" validate:"required,min=6"`
}

type UpdateEmployeeRequest struct {
	FirstName string `json:"first_name,omitempty" bson:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty" bson:"last_name,omitempty"`
	Username  string `json:"username,omitempty" bson:"username,omitempty"`
	Password  string `json:"password,omitempty" bson:"password,omitempty"`
}

type EmployeeResponse struct {
	EmployeeID string    `json:"employee_id"`
	ShopID     string    `json:"shop_id"`
	FirstName  string    `json:"first_name"`
	LastName   string    `json:"last_name"`
	Username   string    `json:"username"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type EmployeeListResponse struct {
	Employees  []EmployeeResponse `json:"employees"`
	TotalPages int                `json:"total_pages"`
	TotalItems int                `json:"total_items"`
	Page       int                `json:"page"`
	PageSize   int                `json:"page_size"`
}
