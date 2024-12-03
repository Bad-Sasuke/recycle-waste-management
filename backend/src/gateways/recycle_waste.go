package gateways

import (
	"recycle-waste-management-backend/src/domain/entities"

	"github.com/gofiber/fiber/v2"
)

func (h *HTTPGateway) GetRecycleWaste(ctx *fiber.Ctx) error {
	data, err := h.RecycleService.GetRecyclableItems()
	if err != nil {
		return ctx.Status(fiber.StatusForbidden).JSON(entities.ResponseModel{Message: err.Error()})
	}
	return ctx.Status(fiber.StatusOK).JSON(entities.ResponseModel{Message: "success", Data: data})
}
