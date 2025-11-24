package gateways

import (
	"recycle-waste-management-backend/src/domain/entities"
	"recycle-waste-management-backend/src/middlewares"

	"github.com/gofiber/fiber/v2"
)

func (h *HTTPGateway) AddCustomerRequest(ctx *fiber.Ctx) error {
	token, err := middlewares.DecodeJWTToken(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(entities.ResponseMessage{Message: "unauthorized"})
	}
	body := new(entities.CustomerRequest)
	if err := ctx.BodyParser(&body); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(entities.ResponseMessage{Message: "invalid json body"})
	}
	err = h.CustomerRequestService.AddCustomerRequest(token.UserID, *body)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(entities.ResponseMessage{Message: "cannot add customer request"})
	}
	return ctx.Status(fiber.StatusOK).JSON(entities.ResponseMessage{Message: "customer request added successfully"})
}

func (h *HTTPGateway) GetCustomerRequest(ctx *fiber.Ctx) error {
	return nil
}

func (h *HTTPGateway) GetCustomerRequests(ctx *fiber.Ctx) error {
	return nil
}

func (h *HTTPGateway) GetCustomerRequestByRequestID(ctx *fiber.Ctx) error {
	return nil
}

func (h *HTTPGateway) DeleteCustomerRequest(ctx *fiber.Ctx) error {
	return nil
}

func (h *HTTPGateway) UpdateCustomerRequest(ctx *fiber.Ctx) error {
	return nil
}

func (h *HTTPGateway) GetCustomerRequestByUserID(ctx *fiber.Ctx) error {
	return nil
}
