package entities

import "time"

// UserSettings represents the user settings stored in database
type UserSettings struct {
	UserID       string    `json:"user_id" bson:"user_id"`

	// Notification settings
	EmailNotifications    bool `json:"email_notifications" bson:"email_notifications"`
	PushNotifications     bool `json:"push_notifications" bson:"push_notifications"`
	MarketplaceUpdates    bool `json:"marketplace_updates" bson:"marketplace_updates"`
	ShopUpdates           bool `json:"shop_updates" bson:"shop_updates"`

	// Privacy settings
	ProfileVisibility string `json:"profile_visibility" bson:"profile_visibility"` // public, private, friends
	ShowEmail         bool   `json:"show_email" bson:"show_email"`
	ShowPhone         bool   `json:"show_phone" bson:"show_phone"`

	// Appearance settings
	Theme    string `json:"theme" bson:"theme"`       // light, dark, auto
	Language string `json:"language" bson:"language"` // th, en

	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}

// UpdateSettingsBody represents the request body for updating settings
type UpdateSettingsBody struct {
	// Notification settings
	EmailNotifications *bool `json:"email_notifications,omitempty"`
	PushNotifications  *bool `json:"push_notifications,omitempty"`
	MarketplaceUpdates *bool `json:"marketplace_updates,omitempty"`
	ShopUpdates        *bool `json:"shop_updates,omitempty"`

	// Privacy settings
	ProfileVisibility *string `json:"profile_visibility,omitempty"`
	ShowEmail         *bool   `json:"show_email,omitempty"`
	ShowPhone         *bool   `json:"show_phone,omitempty"`

	// Appearance settings
	Theme    *string `json:"theme,omitempty"`
	Language *string `json:"language,omitempty"`
}

// DefaultUserSettings returns default settings for a new user
func DefaultUserSettings(userID string) *UserSettings {
	return &UserSettings{
		UserID:             userID,
		EmailNotifications: true,
		PushNotifications:  false,
		MarketplaceUpdates: true,
		ShopUpdates:        true,
		ProfileVisibility:  "public",
		ShowEmail:          false,
		ShowPhone:          false,
		Theme:              "light",
		Language:           "th",
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}
}
