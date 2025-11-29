package gateways

import (
	"strconv"

	"recycle-waste-management-backend/src/middlewares"
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

	tokenDetails, err := middlewares.DecodeJWTToken(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})

	}
	userID := tokenDetails.UserID
	var req services.CreateReceiptRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid request body",
			"error":   err.Error(),
		})
	}

	receipt, err := h.ReceiptService.CreateReceipt(userID, req)
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

func (h *ReceiptGateway) GetReceiptByID(ctx *fiber.Ctx) error {
	receiptID := ctx.Params("receipt_id")
	if receiptID == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Receipt ID is required",
		})
	}

	receipt, err := h.ReceiptService.GetReceiptByID(receiptID)
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

func (h *ReceiptGateway) GetReceiptsByShopID(ctx *fiber.Ctx) error {
	shopID := ctx.Params("shop_id")
	if shopID == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "shop_id is required",
		})
	}

	// Get pagination params
	page, _ := strconv.Atoi(ctx.Query("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.Query("page_size", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	allReceipts, err := h.ReceiptService.GetReceiptsByShopID(shopID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to get receipts",
			"error":   err.Error(),
		})
	}

	// Calculate pagination
	total := len(allReceipts)
	totalPages := (total + pageSize - 1) / pageSize
	start := (page - 1) * pageSize
	end := start + pageSize

	if start >= total {
		start = 0
		end = 0
	}
	if end > total {
		end = total
	}

	var paginatedReceipts = allReceipts[start:end]

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"success":     true,
		"message":     "Receipts retrieved successfully",
		"data":        paginatedReceipts,
		"page":        page,
		"page_size":   pageSize,
		"total":       total,
		"total_pages": totalPages,
	})
}
