package services

import (
	"recycle-waste-management-backend/src/domain/entities"
	"recycle-waste-management-backend/src/repositories"
)

type IStockService interface {
	AddStock(shopID, wasteID string, quantity float64) error
	DeleteStockByWasteID(wasteID string) error
	GetStocksByShopID(shopID string) ([]entities.StockWithDetails, error)
}

type StockService struct {
	StockRepo           repositories.IStockRepository
	RecyclableItemsRepo repositories.IRecyclableItemsRepository
}

func NewStockService(stockRepo repositories.IStockRepository, recyclableItemsRepo repositories.IRecyclableItemsRepository) IStockService {
	return &StockService{
		StockRepo:           stockRepo,
		RecyclableItemsRepo: recyclableItemsRepo,
	}
}

func (s *StockService) AddStock(shopID, wasteID string, quantity float64) error {
	// Get waste details to store as snapshot
	waste, err := s.RecyclableItemsRepo.FindByWasteID(wasteID)
	if err != nil || waste == nil {
		// If waste not found, use default/unknown values
		return s.StockRepo.UpdateStock(shopID, wasteID, quantity, "Unknown", "Unknown", 0)
	}

	return s.StockRepo.UpdateStock(shopID, wasteID, quantity, waste.Category, waste.Name, waste.Price)
}

func (s *StockService) DeleteStockByWasteID(wasteID string) error {
	return s.StockRepo.DeleteByWasteID(wasteID)
}

func (s *StockService) GetStocksByShopID(shopID string) ([]entities.StockWithDetails, error) {
	stocks, err := s.StockRepo.GetStocksByShopID(shopID)
	if err != nil {
		return nil, err
	}

	// Build result with waste details
	var result []entities.StockWithDetails
	for _, stock := range stocks {
		var category, name string
		var currentPrice float64

		// 1. Get Live Data (Current Price & Details)
		waste, err := s.RecyclableItemsRepo.FindByWasteID(stock.WasteID)
		if err == nil && waste != nil {
			category = waste.Category
			name = waste.Name
			currentPrice = waste.Price
		} else {
			// If waste not found, fallback details to snapshot, but current price is 0 or unknown
			category = stock.Category
			name = stock.Name
			currentPrice = 0 // Or handle as N/A
		}

		// Fallback for details if still empty
		if category == "" {
			category = "Unknown"
		}
		if name == "" {
			name = "Unknown"
		}

		// Calculate profit/loss
		// Profit = (Current Price - Purchase Price) * Quantity
		profit := (currentPrice - stock.PricePerKg) * stock.Quantity

		result = append(result, entities.StockWithDetails{
			ID:            stock.ID,
			ShopID:        stock.ShopID,
			WasteID:       stock.WasteID,
			Quantity:      stock.Quantity,
			UpdatedAt:     stock.UpdatedAt,
			Category:      category,
			Name:          name,
			PurchasePrice: stock.PricePerKg,
			CurrentPrice:  currentPrice,
			Profit:        profit,
		})
	}

	return result, nil
}
