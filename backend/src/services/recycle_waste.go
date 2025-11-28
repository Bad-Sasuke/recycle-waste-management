package services

import (
	"fmt"
	"math/rand"
	"recycle-waste-management-backend/src/domain/entities"
	"recycle-waste-management-backend/src/infrastructure/providers"
	"recycle-waste-management-backend/src/repositories"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

type IRecycleWasteService interface {
	GetRecyclableItems() (*[]entities.RecyclableItemsModel, error)
	GetRecyclableItemsPaginated(page, limit int) (*[]entities.RecyclableItemsModel, int64, error)
	GetRecyclableItemsByShopIDPaginated(shopID string, page, limit int) (*[]entities.RecyclableItemsModel, int64, error)
	GetRecyclableItemsByShopIDAndCategoryPaginated(shopID, category string, page, limit int) (*[]entities.RecyclableItemsModel, int64, error)
	GetRecyclableItemByWasteID(wasteID string) (*entities.RecyclableItemsModel, error)
	AddRecycleWaste(data entities.RecyclableItemsModel, image []byte) error
	DeleteWasteItem(wasteID string) error
	GetCategoryWaste() (*[]entities.CategoryWasteModel, error)
	EditWasteItem(wasteID string, data entities.RecyclableItemsModel, image []byte) error
}

type RecycleWasteService struct {
	RecyclableItemsRepo repositories.IRecyclableItemsRepository
	CategoryWAsteRepo   repositories.ICategoryWasteRepository
	StockService        IStockService // Inject StockService
	AwsS3               providers.IAwsS3Upload
	ImageURLDefault     string
}

func NewRecycleWasteService(recyclableItemsRepo repositories.IRecyclableItemsRepository, CategoryWAsteRepo repositories.ICategoryWasteRepository, stockService IStockService) IRecycleWasteService {
	return &RecycleWasteService{RecyclableItemsRepo: recyclableItemsRepo, CategoryWAsteRepo: CategoryWAsteRepo, StockService: stockService, AwsS3: providers.NewAwsS3(), ImageURLDefault: "https://bucketnaja2.s3.ap-southeast-1.amazonaws.com/images/wastes/DEFAULT.jpg"}
}

func (s *RecycleWasteService) GetRecyclableItems() (*[]entities.RecyclableItemsModel, error) {
	data, err := s.RecyclableItemsRepo.FindAll()
	if err != nil && err != mongo.ErrNoDocuments {
		return nil, err
	}
	if data == nil {
		return nil, fmt.Errorf("no data found")
	}
	return data, nil
}

func (s *RecycleWasteService) GetRecyclableItemsPaginated(page, limit int) (*[]entities.RecyclableItemsModel, int64, error) {
	data, totalCount, err := s.RecyclableItemsRepo.FindAllPaginated(page, limit)
	if err != nil && err != mongo.ErrNoDocuments {
		return nil, 0, err
	}
	if data == nil {
		return &[]entities.RecyclableItemsModel{}, 0, nil
	}
	return data, totalCount, nil
}

func (s *RecycleWasteService) GetRecyclableItemsByShopIDPaginated(shopID string, page, limit int) (*[]entities.RecyclableItemsModel, int64, error) {
	data, totalCount, err := s.RecyclableItemsRepo.FindByShopIDPaginated(shopID, page, limit)
	if err != nil && err != mongo.ErrNoDocuments {
		return nil, 0, err
	}
	if data == nil {
		return &[]entities.RecyclableItemsModel{}, 0, nil
	}
	return data, totalCount, nil
}

func (s *RecycleWasteService) GetRecyclableItemsByShopIDAndCategoryPaginated(shopID, category string, page, limit int) (*[]entities.RecyclableItemsModel, int64, error) {
	data, totalCount, err := s.RecyclableItemsRepo.FindByShopIDAndCategoryPaginated(shopID, category, page, limit)
	if err != nil && err != mongo.ErrNoDocuments {
		return nil, 0, err
	}
	if data == nil {
		return &[]entities.RecyclableItemsModel{}, 0, nil
	}
	return data, totalCount, nil
}

func (s *RecycleWasteService) GetRecyclableItemByWasteID(wasteID string) (*entities.RecyclableItemsModel, error) {
	data, err := s.RecyclableItemsRepo.FindByWasteID(wasteID)
	if err != nil {
		return nil, fmt.Errorf("error finding recyclable item: %v", err)
	}
	return data, nil
}

func (s *RecycleWasteService) AddRecycleWaste(data entities.RecyclableItemsModel, image []byte) error {
	if data.Name == "" || data.Category == "" || data.Price == 0 {
		return fmt.Errorf("invalid data")
	}
	// Validate shop_id if provided
	if data.ShopID != "" {
		// Additional validation can be added here
		// For example: check if the shop exists
	}
	if len(image) > 0 {
		keyname, contentType := s.AwsS3.CreateKeyNameImage(data.Name, "webp")
		linkURL, err := s.AwsS3.UploadS3FromString(image, keyname, contentType)
		if err != nil {
			return err
		}
		data.URL = linkURL
	} else {
		data.URL = s.ImageURLDefault
	}
	data.LastUpdate = time.Now().UTC().Add(7 * time.Hour)
	data.WasteID = generateRandomWasteID()
	if err := s.RecyclableItemsRepo.Create(&data); err != nil {
		return err
	}

	return nil
}

func (s *RecycleWasteService) DeleteWasteItem(wasteID string) error {
	data, err := s.RecyclableItemsRepo.FindByWasteID(wasteID)
	if err != nil && err != mongo.ErrNoDocuments {
		return err
	}
	if data == nil {
		return fmt.Errorf("no data found")
	}
	if err := s.RecyclableItemsRepo.Delete(wasteID); err != nil {
		return err
	}

	// Delete associated stock
	if err := s.StockService.DeleteStockByWasteID(wasteID); err != nil {
		// Log error but don't fail the operation? Or fail?
		// For now, let's return error to ensure consistency
		return fmt.Errorf("error deleting associated stock: %v", err)
	}

	if data.URL != s.ImageURLDefault {
		keyName := fmt.Sprintf("images/wastes/%v.webp", strings.ToUpper(data.Name))
		if err := s.AwsS3.DeleteS3(keyName); err != nil {
			return err
		}
	}

	return nil
}

func (s *RecycleWasteService) GetCategoryWaste() (*[]entities.CategoryWasteModel, error) {
	data, err := s.CategoryWAsteRepo.FindAll()
	if err != nil && err != mongo.ErrNoDocuments {
		return nil, err
	}
	if data == nil {
		return nil, fmt.Errorf("no data found")
	}
	return data, nil
}

func (s *RecycleWasteService) EditWasteItem(wasteID string, data entities.RecyclableItemsModel, image []byte) error {
	if data.Name == "" || data.Category == "" || data.Price == 0 {
		return fmt.Errorf("invalid data")
	}
	// Validate shop_id if provided
	if data.ShopID != "" {
		// Additional validation can be added here
		// For example: check if the shop exists
	}
	if len(image) > 0 {
		keyname, contentType := providers.NewAwsS3().CreateKeyNameImage(data.Name, "webp")
		linkURL, err := providers.NewAwsS3().UploadS3FromString(image, keyname, contentType)
		if err != nil {
			return err
		}
		data.URL = linkURL
	}
	data.LastUpdate = time.Now().UTC().Add(7 * time.Hour)
	if err := s.RecyclableItemsRepo.Update(wasteID, &data); err != nil {
		return err
	}
	return nil
}

func generateRandomWasteID() string {
	characters := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	rand.Seed(time.Now().UnixNano())
	randomStr := make([]byte, 4)
	for i := range randomStr {
		randomStr[i] = characters[rand.Intn(len(characters))]
	}
	return fmt.Sprintf("waste-%s", string(randomStr))
}
