package services

import (
	"fmt"
	"math"
	"recycle-waste-management-backend/src/domain/entities"
	"recycle-waste-management-backend/src/domain/models"
	"recycle-waste-management-backend/src/repositories"
	"sort"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type ICustomerRequestService interface {
	AddCustomerRequest(userID string, body entities.CustomerRequest) (error, int)
	GetCustomerRequestByRequestID(userID string) ([]entities.CustomerRequestResponse, error)
	GetCustomerRequests(userID string, page, limit int, maxDistanceKm float64) (*entities.PaginatedCustomerRequestResponse, error)
	AcceptCustomerRequest(customerRequestID string) error
	CancelCustomerRequest(customerRequestID string, cancelReason string) error
	CompleteCustomerRequest(customerRequestID string, shopID string) error
}
type customerRequestService struct {
	customerRequestRepository repositories.ICustomerRequestRepository
	shopRepository            repositories.IShopRepository
}

func NewCustomerRequestService(customerRequestRepository repositories.ICustomerRequestRepository, shopRepository repositories.IShopRepository) ICustomerRequestService {
	return &customerRequestService{
		customerRequestRepository: customerRequestRepository,
		shopRepository:            shopRepository,
	}
}

func (s *customerRequestService) AddCustomerRequest(userID string, body entities.CustomerRequest) (error, int) {
	if userID == "" {
		return fmt.Errorf("user id is empty"), fiber.StatusUnauthorized
	}
	if err := s.customerRequestRepository.CheckCustomerRequestAlreadyExist(userID); err != nil {
		return fmt.Errorf("customer request already exist"), fiber.StatusBadRequest
	}
	modelData := models.CustomerRequestModel{
		CustomerRequestID: uuid.New().String(),
		UserID:            userID,
		Latitude:          body.Latitude,
		Longitude:         body.Longitude,
		Description:       body.Description,
		Status:            models.CR_PENDING,
	}
	return s.customerRequestRepository.AddCustomerRequest(modelData), fiber.StatusOK
}

func (s *customerRequestService) GetCustomerRequestByRequestID(userID string) ([]entities.CustomerRequestResponse, error) {
	if userID == "" {
		return nil, fmt.Errorf("user id is empty")
	}
	// if err := s.customerRequestRepository.DeleteAllCustomerRequest(userID); err != nil {
	// 	return nil, err
	// }
	data, err := s.customerRequestRepository.GetCustomerRequests(userID)
	if err != nil {
		return nil, err
	}
	var response []entities.CustomerRequestResponse
	for _, v := range *data {
		response = append(response, entities.CustomerRequestResponse{
			CustomerRequestID: v.CustomerRequestID,
			Latitude:          v.Latitude,
			Longitude:         v.Longitude,
			Description:       v.Description,
			Status:            string(v.Status),
			CancelReason:      v.CancelReason,
			CreatedAt:         v.CreatedAt,
			UpdatedAt:         v.UpdatedAt,
		})
	}
	return response, nil
}

// haversineDistance calculates the distance between two points on earth (in kilometers)
func haversineDistance(lat1, lon1, lat2, lon2 float64) float64 {
	const earthRadiusKm = 6371.0

	// Convert degrees to radians
	lat1Rad := lat1 * math.Pi / 180
	lat2Rad := lat2 * math.Pi / 180
	deltaLat := (lat2 - lat1) * math.Pi / 180
	deltaLon := (lon2 - lon1) * math.Pi / 180

	// Haversine formula
	a := math.Sin(deltaLat/2)*math.Sin(deltaLat/2) +
		math.Cos(lat1Rad)*math.Cos(lat2Rad)*
			math.Sin(deltaLon/2)*math.Sin(deltaLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return earthRadiusKm * c
}

func (s *customerRequestService) GetCustomerRequests(userID string, page, limit int, maxDistanceKm float64) (*entities.PaginatedCustomerRequestResponse, error) {
	// Get shop location
	shopData, err := s.shopRepository.GetByUserID(userID)
	if err != nil {
		fmt.Printf("Error getting shop data for user %s: %v\n", userID, err)
		return nil, err
	}

	fmt.Printf("Shop location: lat=%.6f, lon=%.6f\n", shopData.Latitude, shopData.Longitude)

	// Get all customer requests
	data, err := s.customerRequestRepository.GetCustomerRequestsPublic()
	if err != nil {
		return nil, err
	}

	// Set default max distance if not provided or invalid
	if maxDistanceKm <= 0 {
		maxDistanceKm = 20.0
	}

	fmt.Printf("Max distance filter: %.2f km\n", maxDistanceKm)
	fmt.Printf("Total customer requests: %d\n", len(*data))

	// Calculate distances and filter within specified distance
	var requestsWithDistance []entities.CustomerRequestResponse

	for _, v := range *data {
		distance := haversineDistance(
			shopData.Latitude,
			shopData.Longitude,
			v.Latitude,
			v.Longitude,
		)

		// Round distance to 2 decimal places
		distance = math.Round(distance*100) / 100

		fmt.Printf("Customer request %s: lat=%.6f, lon=%.6f, distance=%.2f km\n",
			v.CustomerRequestID, v.Latitude, v.Longitude, distance)

		// Only include requests within specified distance
		if distance <= maxDistanceKm {
			requestsWithDistance = append(requestsWithDistance, entities.CustomerRequestResponse{
				CustomerRequestID: v.CustomerRequestID,
				Latitude:          v.Latitude,
				Longitude:         v.Longitude,
				Description:       v.Description,
				Status:            string(v.Status),
				CancelReason:      v.CancelReason,
				Distance:          distance,
				CreatedAt:         v.CreatedAt,
				UpdatedAt:         v.UpdatedAt,
			})
			fmt.Printf("  ✓ INCLUDED (distance %.2f km <= %.2f km)\n", distance, maxDistanceKm)
		} else {
			fmt.Printf("  ✗ EXCLUDED (distance %.2f km > %.2f km)\n", distance, maxDistanceKm)
		}
	}

	fmt.Printf("Filtered customer requests: %d\n", len(requestsWithDistance))

	// Sort by distance (nearest first)
	sort.Slice(requestsWithDistance, func(i, j int) bool {
		return requestsWithDistance[i].Distance < requestsWithDistance[j].Distance
	})

	// Apply pagination
	total := len(requestsWithDistance)

	// Set default values
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	// Calculate pagination
	totalPages := int(math.Ceil(float64(total) / float64(limit)))
	startIndex := (page - 1) * limit
	endIndex := startIndex + limit

	// Get paginated data
	var paginatedData []entities.CustomerRequestResponse
	if startIndex < total {
		if endIndex > total {
			endIndex = total
		}
		paginatedData = requestsWithDistance[startIndex:endIndex]
	} else {
		paginatedData = []entities.CustomerRequestResponse{}
	}

	return &entities.PaginatedCustomerRequestResponse{
		Data:       paginatedData,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

func (s *customerRequestService) AcceptCustomerRequest(customerRequestID string) error {
	return s.customerRequestRepository.UpdateCustomerRequestStatus(customerRequestID, models.CR_ACCEPTED)
}

func (s *customerRequestService) CancelCustomerRequest(customerRequestID string, cancelReason string) error {
	return s.customerRequestRepository.CancelCustomerRequest(customerRequestID, cancelReason)
}

func (s *customerRequestService) CompleteCustomerRequest(customerRequestID string, shopID string) error {
	return s.customerRequestRepository.CompleteCustomerRequest(customerRequestID, shopID)
}
