package entities

import "time"

type NewUserBody struct {
	Username string `json:"username" bson:"username"`
	Email    string `json:"email" bson:"email"`
}

type UserDataFormat struct {
	UserID    string    `json:"user_id,omitempty" bson:"user_id,omitempty"`
	Username  string    `json:"username,omitempty" bson:"username,omitempty"`
	Email     string    `json:"email,omitempty" bson:"email,omitempty"`
	ImageURL  string    `json:"image_url,omitempty" bson:"image_url,omitempty"`
	JWT       string    `json:"jwt,omitempty" bson:"jwt,omitempty"`
	Role      string    `json:"role,omitempty" bson:"role,omitempty"` // User role: admin, user, moderator, etc.
	CreatedAt time.Time `json:"created_at,omitempty" bson:"created_at,omitempty"`
	LastLogin time.Time `json:"last_login,omitempty" bson:"last_login,omitempty"`
}

// UserRole represents the possible user roles in the system
type UserRole string

const (
	UserRoleUser    UserRole = "user"
	UserRoleAdmin   UserRole = "admin"
	UserRoleModerator UserRole = "moderator"
)
