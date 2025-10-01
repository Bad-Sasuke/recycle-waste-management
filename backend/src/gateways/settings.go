package gateways

import (
	"recycle-waste-management-backend/src/domain/entities"
	"recycle-waste-management-backend/src/middlewares"

	"github.com/gofiber/fiber/v2"
)

// GetUserSettings retrieves the current user's settings
// @Summary Get user settings
// @Description Get the current user's settings
// @Tags settings
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} entities.ResponseModel{data=entities.UserSettings}
// @Failure 401 {object} entities.ResponseMessage
// @Failure 500 {object} entities.ResponseMessage
// @Router /api/settings [get]
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
// @Summary Update user settings
// @Description Update the current user's settings
// @Tags settings
// @Accept json
// @Produce json
// @Security Bearer
// @Param settings body entities.UpdateSettingsBody true "Settings to update"
// @Success 200 {object} entities.ResponseModel{data=entities.UserSettings}
// @Failure 400 {object} entities.ResponseMessage
// @Failure 401 {object} entities.ResponseMessage
// @Failure 500 {object} entities.ResponseMessage
// @Router /api/settings [put]
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
// @Summary Delete user settings
// @Description Delete the current user's settings (will reset to default on next access)
// @Tags settings
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} entities.ResponseMessage
// @Failure 401 {object} entities.ResponseMessage
// @Failure 500 {object} entities.ResponseMessage
// @Router /api/settings [delete]
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
