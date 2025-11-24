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
	err, statusCode := h.CustomerRequestService.AddCustomerRequest(token.UserID, *body)
	if statusCode != fiber.StatusOK {
		return ctx.Status(statusCode).JSON(entities.ResponseMessage{Message: err.Error()})
	}
	return ctx.Status(fiber.StatusOK).JSON(entities.ResponseMessage{Message: "customer request added successfully"})
}

func (h *HTTPGateway) GetCustomerRequests(ctx *fiber.Ctx) error {
	token, err := middlewares.DecodeJWTToken(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(entities.ResponseMessage{Message: "unauthorized"})
	}

	// Get pagination parameters from query string
	page := ctx.QueryInt("page", 1)
	limit := ctx.QueryInt("limit", 10)

	// Get max distance parameter (default 20.0 km)
	maxDistance := ctx.QueryFloat("maxDistance", 20.0)

	customerRequests, err := h.CustomerRequestService.GetCustomerRequests(token.UserID, page, limit, maxDistance)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(entities.ResponseMessage{Message: "cannot get customer requests"})
	}
	return ctx.Status(fiber.StatusOK).JSON(customerRequests)
}

func (h *HTTPGateway) GetCustomerRequestByRequestID(ctx *fiber.Ctx) error {
	token, err := middlewares.DecodeJWTToken(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(entities.ResponseMessage{Message: "unauthorized"})
	}
	customerRequests, err := h.CustomerRequestService.GetCustomerRequestByRequestID(token.UserID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(entities.ResponseMessage{Message: "cannot get customer requests"})
	}
	return ctx.Status(fiber.StatusOK).JSON(entities.ResponseModel{Data: customerRequests, Message: "success"})
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

func (h *HTTPGateway) AcceptCustomerRequest(ctx *fiber.Ctx) error {
	// Verify authentication
	_, err := middlewares.DecodeJWTToken(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(entities.ResponseMessage{Message: "unauthorized"})
	}

	// Get customer_request_id from URL params
	customerRequestID := ctx.Params("id")
	if customerRequestID == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(entities.ResponseMessage{Message: "customer_request_id is required"})
	}

	// Accept the request
	err = h.CustomerRequestService.AcceptCustomerRequest(customerRequestID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(entities.ResponseMessage{Message: err.Error()})
	}

	return ctx.Status(fiber.StatusOK).JSON(entities.ResponseMessage{Message: "customer request accepted successfully"})
}
