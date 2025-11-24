package gateways

import (
	"recycle-waste-management-backend/src/domain/entities"
	"recycle-waste-management-backend/src/middlewares"

	"github.com/gofiber/fiber/v2"
)

// GetUserSettings retrieves the current user's settings
func (h *HTTPGateway) GetUserSettings(ctx *fiber.Ctx) error {
	tokenDetails, err := middlewares.DecodeJWTToken(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(entities.ResponseMessage{Message: "unauthorized"})
	}

	settings, err := h.SettingsService.GetSettings(tokenDetails.UserID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(entities.ResponseMessage{Message: "cannot get settings"})
	}

	return ctx.Status(fiber.StatusOK).JSON(entities.ResponseModel{Message: "success", Data: settings})
}

// UpdateUserSettings updates the current user's settings
func (h *HTTPGateway) UpdateUserSettings(ctx *fiber.Ctx) error {
	tokenDetails, err := middlewares.DecodeJWTToken(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(entities.ResponseMessage{Message: "unauthorized"})
	}

	bodyData := new(entities.UpdateSettingsBody)
	if err := ctx.BodyParser(&bodyData); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(entities.ResponseMessage{Message: "invalid json body"})
	}

	settings, err := h.SettingsService.UpdateSettings(tokenDetails.UserID, bodyData)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(entities.ResponseMessage{Message: err.Error()})
	}

	return ctx.Status(fiber.StatusOK).JSON(entities.ResponseModel{Message: "settings updated successfully", Data: settings})
}

// DeleteUserSettings deletes the current user's settings (resets to default)
func (h *HTTPGateway) DeleteUserSettings(ctx *fiber.Ctx) error {
	tokenDetails, err := middlewares.DecodeJWTToken(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(entities.ResponseMessage{Message: "unauthorized"})
	}

	err = h.SettingsService.DeleteSettings(tokenDetails.UserID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(entities.ResponseMessage{Message: "cannot delete settings"})
	}

	return ctx.Status(fiber.StatusOK).JSON(entities.ResponseMessage{Message: "settings deleted successfully"})
}
