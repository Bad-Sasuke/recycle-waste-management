package services

import (
	"fmt"
	"recycle-waste-management-backend/src/domain/entities"
	"recycle-waste-management-backend/src/domain/repositories"

	"go.mongodb.org/mongo-driver/mongo"
)

type IRecycleWasteService interface {
	GetRecyclableItems() (*[]entities.RecyclableItemsModel, error)
}

type RecycleWasteService struct {
	RecyclableItemsRepo repositories.IRecyclableItemsRepository
}

func NewRecycleWasteService(recyclableItemsRepo repositories.IRecyclableItemsRepository) IRecycleWasteService {
	return &RecycleWasteService{RecyclableItemsRepo: recyclableItemsRepo}
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
