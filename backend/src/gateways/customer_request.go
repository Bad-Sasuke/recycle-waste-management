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

func (h *HTTPGateway) CancelCustomerRequest(ctx *fiber.Ctx) error {
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

	// Parse request body
	body := new(entities.CancelCustomerRequest)
	if err := ctx.BodyParser(&body); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(entities.ResponseMessage{Message: "invalid json body"})
	}

	// Validate reason
	if body.Reason == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(entities.ResponseMessage{Message: "reason is required"})
	}

	// Check reason length (max 100 characters)
	if len(body.Reason) > 100 {
		return ctx.Status(fiber.StatusBadRequest).JSON(entities.ResponseMessage{Message: "reason must not exceed 100 characters"})
	}

	// Verify that the request belongs to the user
	// (You might want to add this verification in the service layer)

	// Cancel the request
	err = h.CustomerRequestService.CancelCustomerRequest(customerRequestID, body.Reason)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(entities.ResponseMessage{Message: err.Error()})
	}

	return ctx.Status(fiber.StatusOK).JSON(entities.ResponseMessage{Message: "customer request cancelled successfully"})
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

func (h *HTTPGateway) CompleteCustomerRequest(ctx *fiber.Ctx) error {
	// Verify authentication
	token, err := middlewares.DecodeJWTToken(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(entities.ResponseMessage{Message: "unauthorized"})
	}

	// Get customer_request_id from URL params
	customerRequestID := ctx.Params("id")
	if customerRequestID == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(entities.ResponseMessage{Message: "customer_request_id is required"})
	}

	// Get shop by user ID
	shop, err := h.ShopService.GetShopByUserID(token.UserID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(entities.ResponseMessage{Message: "cannot get shop info"})
	}

	// Complete the request
	err = h.CustomerRequestService.CompleteCustomerRequest(customerRequestID, shop.ShopID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(entities.ResponseMessage{Message: err.Error()})
	}

	return ctx.Status(fiber.StatusOK).JSON(entities.ResponseMessage{Message: "customer request completed successfully"})
}

func (h *HTTPGateway) CreateWalkInRequest(ctx *fiber.Ctx) error {
	// Verify authentication (Shop Owner)
	token, err := middlewares.DecodeJWTToken(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(entities.ResponseMessage{Message: "unauthorized"})
	}

	// Parse body
	body := new(entities.WalkInCustomerRequest)
	if err := ctx.BodyParser(&body); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(entities.ResponseMessage{Message: "invalid json body"})
	}

	// Get Shop ID from User ID (to ensure security)
	shop, err := h.ShopService.GetShopByUserID(token.UserID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(entities.ResponseMessage{Message: "cannot get shop info"})
	}
	body.ShopID = shop.ShopID

	// Call Service
	requestID, err := h.CustomerRequestService.CreateWalkInRequest(*body)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(entities.ResponseMessage{Message: err.Error()})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"success":             true,
		"message":             "Walk-in request created successfully",
		"customer_request_id": requestID,
	})
}
