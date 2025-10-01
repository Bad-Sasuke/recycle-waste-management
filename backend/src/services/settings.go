package services

import (
	"fmt"
	"recycle-waste-management-backend/src/domain/entities"
	"recycle-waste-management-backend/src/domain/repositories"
)

type settingsService struct {
	SettingsRepository repositories.ISettingsRepository
}

type ISettingsService interface {
	GetSettings(userID string) (*entities.UserSettings, error)
	UpdateSettings(userID string, updates *entities.UpdateSettingsBody) (*entities.UserSettings, error)
	DeleteSettings(userID string) error
}

func NewSettingsService(repo repositories.ISettingsRepository) ISettingsService {
	return &settingsService{
		SettingsRepository: repo,
	}
}

func (sv *settingsService) GetSettings(userID string) (*entities.UserSettings, error) {
	settings, err := sv.SettingsRepository.GetSettings(userID)
	if err != nil {
		return nil, err
	}
	return settings, nil
}

func (sv *settingsService) UpdateSettings(userID string, updates *entities.UpdateSettingsBody) (*entities.UserSettings, error) {
	// Get current settings
	currentSettings, err := sv.SettingsRepository.GetSettings(userID)
	if err != nil {
		return nil, err
	}

	// Apply updates (only update fields that are provided)
	if updates.EmailNotifications != nil {
		currentSettings.EmailNotifications = *updates.EmailNotifications
	}
	if updates.PushNotifications != nil {
		currentSettings.PushNotifications = *updates.PushNotifications
	}
	if updates.MarketplaceUpdates != nil {
		currentSettings.MarketplaceUpdates = *updates.MarketplaceUpdates
	}
	if updates.ShopUpdates != nil {
		currentSettings.ShopUpdates = *updates.ShopUpdates
	}
	if updates.ProfileVisibility != nil {
		// Validate profile visibility
		validVisibility := map[string]bool{"public": true, "private": true, "friends": true}
		if !validVisibility[*updates.ProfileVisibility] {
			return nil, fmt.Errorf("invalid profile visibility: %s", *updates.ProfileVisibility)
		}
		currentSettings.ProfileVisibility = *updates.ProfileVisibility
	}
	if updates.ShowEmail != nil {
		currentSettings.ShowEmail = *updates.ShowEmail
	}
	if updates.ShowPhone != nil {
		currentSettings.ShowPhone = *updates.ShowPhone
	}
	if updates.Theme != nil {
		// Validate theme
		validThemes := map[string]bool{"light": true, "dark": true, "auto": true}
		if !validThemes[*updates.Theme] {
			return nil, fmt.Errorf("invalid theme: %s", *updates.Theme)
		}
		currentSettings.Theme = *updates.Theme
	}
	if updates.Language != nil {
		// Validate language
		validLanguages := map[string]bool{"th": true, "en": true}
		if !validLanguages[*updates.Language] {
			return nil, fmt.Errorf("invalid language: %s", *updates.Language)
		}
		currentSettings.Language = *updates.Language
	}

	// Save updated settings
	err = sv.SettingsRepository.UpsertSettings(currentSettings)
	if err != nil {
		return nil, err
	}

	return currentSettings, nil
}

func (sv *settingsService) DeleteSettings(userID string) error {
	err := sv.SettingsRepository.DeleteSettings(userID)
	if err != nil {
		return err
	}
	return nil
}
