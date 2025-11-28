package gateways

import (
	"strconv"

	"recycle-waste-management-backend/src/services"

	"github.com/gofiber/fiber/v2"
)

type StockGateway struct {
	StockService services.IStockService
}

func NewStockGateway(stockService services.IStockService) *StockGateway {
	return &StockGateway{
		StockService: stockService,
	}
}

func (h *StockGateway) GetStocksByShopID(ctx *fiber.Ctx) error {
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

	allStocks, err := h.StockService.GetStocksByShopID(shopID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to get stocks",
			"error":   err.Error(),
		})
	}

	// Calculate pagination
	total := len(allStocks)
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

	var paginatedStocks = allStocks[start:end]

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"success":     true,
		"message":     "Stocks retrieved successfully",
		"data":        paginatedStocks,
		"page":        page,
		"page_size":   pageSize,
		"total":       total,
		"total_pages": totalPages,
	})
}
