package gateways

import (
	"recycle-waste-management-backend/src/services"

	"github.com/gofiber/fiber/v2"
)

type ReceiptGateway struct {
	ReceiptService services.IReceiptService
}

func NewReceiptGateway(receiptService services.IReceiptService) *ReceiptGateway {
	return &ReceiptGateway{
		ReceiptService: receiptService,
	}
}

func (h *ReceiptGateway) CreateReceipt(ctx *fiber.Ctx) error {
	// Check user role authorization
	role := ctx.Locals("role")
	if role == nil || (role != "shop" && role != "admin") {
		return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"success": false,
			"message": "Access denied. Only shop owners and admins can create receipts.",
		})
	}

	var req services.CreateReceiptRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid request body",
			"error":   err.Error(),
		})
	}

	receipt, err := h.ReceiptService.CreateReceipt(req)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to create receipt",
			"error":   err.Error(),
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "Receipt created successfully",
		"data":    receipt,
	})
}

func (h *ReceiptGateway) GetReceiptByRequestID(ctx *fiber.Ctx) error {
	requestID := ctx.Params("request_id")
	if requestID == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Request ID is required",
		})
	}

	receipt, err := h.ReceiptService.GetReceiptByCustomerRequestID(requestID)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Receipt not found",
			"error":   err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    receipt,
	})
}
