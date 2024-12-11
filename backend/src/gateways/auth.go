package gateways

import (
	"fmt"
	"os"
	"recycle-waste-management-backend/src/domain/entities"

	"github.com/gofiber/fiber/v2"
)

func (h *HTTPGateway) AuthGithubCallback(ctx *fiber.Ctx) error {
	code := ctx.Query("code")

	if code == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(entities.ResponseMessage{Message: "invalid query params"})
	}

	data, err := h.AuthService.AuthGithub(code)
	if err != nil {
		return ctx.Status(fiber.StatusForbidden).JSON(entities.ResponseMessage{Message: data.ErrorDesc})
	}
	return ctx.Redirect(fmt.Sprintf("%v/auth?token=%v", os.Getenv("FRONTEND_URL"), data.JWT), fiber.StatusTemporaryRedirect)
}
