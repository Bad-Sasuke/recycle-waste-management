package services

import (
	"recycle-waste-management-backend/src/domain/entities"
	"recycle-waste-management-backend/src/repositories"
)

type CustomerRequestService interface {
	AddCustomerRequest(userID string, body entities.CustomerRequest) error
}

type ICustomerRequestService struct {
	customerRequestRepository repositories.ICustomerRequestRepository
}

func NewCustomerRequestService(customerRequestRepository repositories.ICustomerRequestRepository) *ICustomerRequestService {
	return &ICustomerRequestService{
		customerRequestRepository: customerRequestRepository,
	}
}

func (s *ICustomerRequestService) AddCustomerRequest(userID string, body entities.CustomerRequest) error {
	return s.customerRequestRepository.AddCustomerRequest(userID, body)
}
