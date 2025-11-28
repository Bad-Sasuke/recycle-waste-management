package services

import (
	"recycle-waste-management-backend/src/repositories"
)

type IStockService interface {
	AddStock(shopID, wasteID string, quantity float64) error
	DeleteStockByWasteID(wasteID string) error
	// Future: DeductStock, CheckAvailability, etc.
}

type StockService struct {
	StockRepo repositories.IStockRepository
}

func NewStockService(stockRepo repositories.IStockRepository) IStockService {
	return &StockService{
		StockRepo: stockRepo,
	}
}

func (s *StockService) AddStock(shopID, wasteID string, quantity float64) error {
	// Business logic can be added here (e.g., validation, logging)
	return s.StockRepo.UpdateStock(shopID, wasteID, quantity)
}

func (s *StockService) DeleteStockByWasteID(wasteID string) error {
	return s.StockRepo.DeleteByWasteID(wasteID)
}
