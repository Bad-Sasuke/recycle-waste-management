package services

import (
	"context"
	"fmt"
	"math/rand"
	"recycle-waste-management-backend/src/domain/entities"
	"recycle-waste-management-backend/src/infrastructure/providers"
	"recycle-waste-management-backend/src/repositories"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

type IShopService interface {
	CreateShop(userID string, data entities.CreateShopRequest, image []byte) error
	GetShopByShopID(shopID string) (*entities.ShopResponse, error)
	GetShopByUserID(userID string) (*entities.ShopModel, error)
	GetAllShops(page, limit int) (*[]entities.ShopModel, int64, error)
	UpdateShop(shopID string, data entities.UpdateShopRequest, image []byte) error
	DeleteShop(shopID string) error
}

type ShopService struct {
	ShopRepository   repositories.IShopRepository
	ReviewRepository repositories.IReviewRepository
	AwsS3            providers.IAwsS3Upload
	ImageURLDefault  string
}

func NewShopService(shopRepo repositories.IShopRepository, reviewRepo repositories.IReviewRepository) IShopService {
	return &ShopService{
		ShopRepository:   shopRepo,
		ReviewRepository: reviewRepo,
		AwsS3:            providers.NewAwsS3(),
		ImageURLDefault:  "https://bucketnaja2.s3.ap-southeast-1.amazonaws.com/images/shops/DEFAULT.jpg",
	}
}

func (s *ShopService) CreateShop(userID string, data entities.CreateShopRequest, image []byte) error {
	// Validate required fields
	if data.Name == "" || data.Address == "" {
		return fmt.Errorf("name and address are required")
	}

	// Check if user already has a shop
	_, err := s.ShopRepository.GetByUserID(userID)
	if err != nil {
		// If there's an error and it's not "no document found", return the error
		if err != mongo.ErrNoDocuments {
			return err
		}
		// If it's "no document found" error, it means user doesn't have a shop yet,
		// so we continue with creating the shop
	} else {
		// If we found an existing shop (no error returned), return an error
		return fmt.Errorf("user already has a shop")
	}

	shopModel := &entities.ShopModel{
		ShopID:      generateRandomShopID(),
		UserID:      userID,
		Name:        data.Name,
		Description: data.Description,
		Address:     data.Address,
		Phone:       data.Phone,
		Email:       data.Email,
		OpeningTime: data.OpeningTime,
		ClosingTime: data.ClosingTime,
		Latitude:    data.Latitude,
		Longitude:   data.Longitude,
		CreatedAt:   time.Now().UTC().Add(7 * time.Hour),
		UpdatedAt:   time.Now().UTC().Add(7 * time.Hour),
	}

	if len(image) > 0 {
		keyname, contentType := s.AwsS3.CreateKeyNameImage(shopModel.Name, "webp")
		linkURL, err := s.AwsS3.UploadS3FromString(image, keyname, contentType)
		if err != nil {
			return err
		}
		shopModel.ImageURL = linkURL
	} else {
		shopModel.ImageURL = s.ImageURLDefault
	}

	if err := s.ShopRepository.Create(shopModel); err != nil {
		return err
	}

	return nil
}

func (s *ShopService) GetShopByShopID(shopID string) (*entities.ShopResponse, error) {
	shop, err := s.ShopRepository.GetByShopID(shopID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("shop not found")
		}
		return nil, err
	}

	// Calculate rating
	avgRating, totalReviews, err := s.ReviewRepository.GetShopRatingStats(context.Background(), shopID)
	if err != nil {
		// Log error but continue
		fmt.Printf("Error calculating shop rating: %v\n", err)
	}

	return &entities.ShopResponse{
		ShopModel:     *shop,
		AverageRating: avgRating,
		TotalReviews:  totalReviews,
	}, nil
}

func (s *ShopService) GetShopByUserID(userID string) (*entities.ShopModel, error) {
	shop, err := s.ShopRepository.GetByUserID(userID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("no shop found for this user")
		}
		return nil, err
	}
	return shop, nil
}

func (s *ShopService) GetAllShops(page, limit int) (*[]entities.ShopModel, int64, error) {
	shops, totalCount, err := s.ShopRepository.GetAll(page, limit)
	if err != nil && err != mongo.ErrNoDocuments {
		return nil, 0, err
	}
	if shops == nil {
		return &[]entities.ShopModel{}, 0, nil
	}
	return shops, totalCount, nil
}

func (s *ShopService) UpdateShop(shopID string, data entities.UpdateShopRequest, image []byte) error {
	// Get existing shop
	existingShop, err := s.ShopRepository.GetByShopID(shopID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return fmt.Errorf("shop not found")
		}
		return err
	}

	// Update fields if provided
	if data.Name != nil {
		existingShop.Name = *data.Name
	}
	if data.Description != nil {
		existingShop.Description = *data.Description
	}
	if data.Address != nil {
		existingShop.Address = *data.Address
	}
	if data.Phone != nil {
		existingShop.Phone = *data.Phone
	}
	if data.Email != nil {
		existingShop.Email = *data.Email
	}
	if data.OpeningTime != nil {
		existingShop.OpeningTime = *data.OpeningTime
	}
	if data.ClosingTime != nil {
		existingShop.ClosingTime = *data.ClosingTime
	}
	if data.Latitude != nil {
		existingShop.Latitude = *data.Latitude
	}
	if data.Longitude != nil {
		existingShop.Longitude = *data.Longitude
	}

	existingShop.UpdatedAt = time.Now().UTC().Add(7 * time.Hour)

	if len(image) > 0 {
		keyname, contentType := s.AwsS3.CreateKeyNameImage(existingShop.Name, "webp")
		linkURL, err := s.AwsS3.UploadS3FromString(image, keyname, contentType)
		if err != nil {
			return err
		}
		existingShop.ImageURL = linkURL
	}

	if err := s.ShopRepository.Update(shopID, existingShop); err != nil {
		return err
	}

	return nil
}

func (s *ShopService) DeleteShop(shopID string) error {
	shop, err := s.ShopRepository.GetByShopID(shopID)
	if err != nil && err != mongo.ErrNoDocuments {
		return err
	}
	if shop == nil {
		return fmt.Errorf("shop not found")
	}

	if err := s.ShopRepository.Delete(shopID); err != nil {
		return err
	}

	// Delete image from S3 if it's not the default image
	if shop.ImageURL != s.ImageURLDefault {
		keyName := fmt.Sprintf("images/shops/%v.webp", strings.ToUpper(shop.Name))
		if err := s.AwsS3.DeleteS3(keyName); err != nil {
			return err
		}
	}

	return nil
}

func generateRandomShopID() string {
	characters := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	rand.Seed(time.Now().UnixNano())
	randomStr := make([]byte, 6)
	for i := range randomStr {
		randomStr[i] = characters[rand.Intn(len(characters))]
	}
	return fmt.Sprintf("shop-%s", string(randomStr))
}
