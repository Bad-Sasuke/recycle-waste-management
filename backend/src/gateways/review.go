package gateways

import (
	"recycle-waste-management-backend/src/domain/entities"
	"recycle-waste-management-backend/src/middlewares"
	"recycle-waste-management-backend/src/services"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type ReviewGateway struct {
	reviewService services.IReviewService
}

func NewReviewGateway(reviewService services.IReviewService) *ReviewGateway {
	return &ReviewGateway{
		reviewService: reviewService,
	}
}

// CreateReview handles POST /api/reviews
func (g *ReviewGateway) CreateReview(c *fiber.Ctx) error {
	// Get user ID from JWT token
	tokenDetails, err := middlewares.DecodeJWTToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	userID := tokenDetails.UserID

	var req entities.CreateReviewRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := g.reviewService.CreateReview(c.Context(), userID, &req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Review created successfully",
	})
}

// SkipReview handles POST /api/reviews/skip
func (g *ReviewGateway) SkipReview(c *fiber.Ctx) error {
	// Get user ID from JWT token
	tokenDetails, err := middlewares.DecodeJWTToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	userID := tokenDetails.UserID

	var req entities.SkipReviewRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := g.reviewService.SkipReview(c.Context(), userID, &req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Review skipped",
	})
}

// GetReviewsByShopID handles GET /api/reviews/shop/:shop_id
func (g *ReviewGateway) GetReviewsByShopID(c *fiber.Ctx) error {
	shopID := c.Params("shop_id")
	if shopID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "shop_id is required",
		})
	}

	// Get pagination parameters
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("page_size", "10"))

	reviews, err := g.reviewService.GetReviewsByShopID(c.Context(), shopID, page, pageSize)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(reviews)
}

// CheckReviewExists handles GET /api/reviews/check/:customer_request_id
func (g *ReviewGateway) CheckReviewExists(c *fiber.Ctx) error {
	customerRequestID := c.Params("customer_request_id")
	if customerRequestID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "customer_request_id is required",
		})
	}

	exists, err := g.reviewService.CheckReviewExists(c.Context(), customerRequestID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"exists": exists,
	})
}
