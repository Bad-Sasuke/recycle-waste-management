package services

import (
	"context"
	"errors"
	"math"
	"recycle-waste-management-backend/src/domain/entities"
	"recycle-waste-management-backend/src/domain/models"
	"recycle-waste-management-backend/src/repositories"
)

type IReviewService interface {
	CreateReview(ctx context.Context, userID string, req *entities.CreateReviewRequest) error
	SkipReview(ctx context.Context, userID string, req *entities.SkipReviewRequest) error
	GetReviewsByShopID(ctx context.Context, shopID string, page, pageSize int) (*entities.GetReviewsResponse, error)
	CheckReviewExists(ctx context.Context, customerRequestID string) (bool, error)
}

type reviewService struct {
	reviewRepo          repositories.IReviewRepository
	customerRequestRepo repositories.ICustomerRequestRepository
}

func NewReviewService(reviewRepo repositories.IReviewRepository, customerRequestRepo repositories.ICustomerRequestRepository) IReviewService {
	return &reviewService{
		reviewRepo:          reviewRepo,
		customerRequestRepo: customerRequestRepo,
	}
}

func (s *reviewService) CreateReview(ctx context.Context, userID string, req *entities.CreateReviewRequest) error {
	// Validate rating
	if req.Rating < 1 || req.Rating > 5 {
		return errors.New("rating must be between 1 and 5")
	}

	// Validate comment length
	if len(req.Comment) > 200 {
		return errors.New("comment must not exceed 200 characters")
	}

	// Check if customer request exists and belongs to user
	customerRequest, err := s.customerRequestRepo.GetCustomerRequestByID(req.CustomerRequestID)
	if err != nil {
		return err
	}
	if customerRequest == nil {
		return errors.New("customer request not found")
	}
	if customerRequest.UserID != userID {
		return errors.New("unauthorized: customer request does not belong to user")
	}

	// Check if customer request is done
	if customerRequest.Status != models.CR_DONE {
		return errors.New("can only review completed requests")
	}

	// Check if review already exists
	exists, err := s.reviewRepo.CheckReviewExists(ctx, req.CustomerRequestID)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("review already exists for this request")
	}

	// Get shop_id from accepted shop (assuming it's stored somewhere)
	// For now, we'll need to get this from the customer request
	// You may need to add shop_id field to CustomerRequestModel
	shopID := customerRequest.UserID // This might need adjustment based on your schema

	review := &models.ReviewModel{
		CustomerRequestID: req.CustomerRequestID,
		ShopID:            shopID,
		UserID:            userID,
		Rating:            req.Rating,
		Comment:           req.Comment,
		IsSkipped:         false,
	}

	return s.reviewRepo.CreateReview(ctx, review)
}

func (s *reviewService) SkipReview(ctx context.Context, userID string, req *entities.SkipReviewRequest) error {
	// Check if customer request exists and belongs to user
	customerRequest, err := s.customerRequestRepo.GetCustomerRequestByID(req.CustomerRequestID)
	if err != nil {
		return err
	}
	if customerRequest == nil {
		return errors.New("customer request not found")
	}
	if customerRequest.UserID != userID {
		return errors.New("unauthorized: customer request does not belong to user")
	}

	// Check if review already exists
	exists, err := s.reviewRepo.CheckReviewExists(ctx, req.CustomerRequestID)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("review action already taken for this request")
	}

	// Get shop_id from customer request
	shopID := customerRequest.UserID // This might need adjustment based on your schema

	// Create a skipped review (no rating, no comment)
	review := &models.ReviewModel{
		CustomerRequestID: req.CustomerRequestID,
		ShopID:            shopID,
		UserID:            userID,
		Rating:            0,
		Comment:           "",
		IsSkipped:         true,
	}

	return s.reviewRepo.CreateReview(ctx, review)
}

func (s *reviewService) GetReviewsByShopID(ctx context.Context, shopID string, page, pageSize int) (*entities.GetReviewsResponse, error) {
	// Validate pagination parameters
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	reviews, totalCount, err := s.reviewRepo.GetReviewsByShopID(ctx, shopID, page, pageSize)
	if err != nil {
		return nil, err
	}

	// Convert to response format
	reviewResponses := make([]entities.ReviewResponse, len(reviews))
	for i, review := range reviews {
		reviewResponses[i] = entities.ReviewResponse{
			ReviewID:          review.ReviewID,
			CustomerRequestID: review.CustomerRequestID,
			ShopID:            review.ShopID,
			UserID:            review.UserID,
			Rating:            review.Rating,
			Comment:           review.Comment,
			IsSkipped:         review.IsSkipped,
			CreatedAt:         review.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:         review.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
	}

	totalPages := int(math.Ceil(float64(totalCount) / float64(pageSize)))

	return &entities.GetReviewsResponse{
		Reviews:    reviewResponses,
		TotalCount: totalCount,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

func (s *reviewService) CheckReviewExists(ctx context.Context, customerRequestID string) (bool, error) {
	return s.reviewRepo.CheckReviewExists(ctx, customerRequestID)
}
