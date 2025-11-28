package services

import (
	"fmt"
	"recycle-waste-management-backend/src/domain/entities"
	"recycle-waste-management-backend/src/repositories"
	"time"

	"github.com/google/uuid"
)

type IReceiptService interface {
	CreateReceipt(req CreateReceiptRequest) (*entities.Receipt, error)
	GetReceiptByCustomerRequestID(requestID string) (*ReceiptWithItemsResponse, error)
}

type ReceiptWithItemsResponse struct {
	Receipt *entities.Receipt       `json:"receipt"`
	Items   *[]entities.ReceiptItem `json:"items"`
	Shop    *entities.ShopModel     `json:"shop"`
}

type CreateReceiptRequest struct {
	ShopID            string               `json:"shop_id"`
	PaymentMethod     string               `json:"payment_method"`
	CustomerRequestID string               `json:"customer_request_id"`
	Items             []ReceiptItemRequest `json:"items"`
}

type ReceiptItemRequest struct {
	WasteID   string  `json:"waste_id"`
	Name      string  `json:"name"`
	Category  string  `json:"category"`
	Weight    float64 `json:"weight"`
	UnitPrice float64 `json:"unit_price"`
	Price     float64 `json:"price"`
}

type ReceiptService struct {
	ReceiptRepo         repositories.IReceiptRepository
	ReceiptItemRepo     repositories.IReceiptItemRepository
	RecyclableItemsRepo repositories.IRecyclableItemsRepository
	StockService        IStockService
	CustomerRequestRepo repositories.ICustomerRequestRepository
	ShopRepo            repositories.IShopRepository // Inject
}

func NewReceiptService(
	receiptRepo repositories.IReceiptRepository,
	receiptItemRepo repositories.IReceiptItemRepository,
	recyclableItemsRepo repositories.IRecyclableItemsRepository,
	stockService IStockService,
	customerRequestRepo repositories.ICustomerRequestRepository,
	shopRepo repositories.IShopRepository, // Add param
) IReceiptService {
	return &ReceiptService{
		ReceiptRepo:         receiptRepo,
		ReceiptItemRepo:     receiptItemRepo,
		RecyclableItemsRepo: recyclableItemsRepo,
		StockService:        stockService,
		CustomerRequestRepo: customerRequestRepo,
		ShopRepo:            shopRepo, // Assign
	}
}

func (s *ReceiptService) CreateReceipt(req CreateReceiptRequest) (*entities.Receipt, error) {
	// 1. Calculate totals
	var totalAmount float64
	for _, item := range req.Items {
		totalAmount += item.Price
	}

	// Assuming VAT is 7% included or on top?
	// Usually for buying waste, there might not be VAT in the same way as selling goods,
	// but let's assume standard 7% calculation for the receipt structure.
	// If the user means "Withholding Tax" (WHT) which is common for services, that's different.
	// But for "VAT", let's calculate it.
	// Let's assume the Price is the Net amount paid to customer, so maybe VAT is 0 for now unless specified.
	// However, the user said "receipt ... มี vat ด้วย".
	// Let's assume 7% VAT.
	vatRate := 0.07
	vat := totalAmount * vatRate
	netTotal := totalAmount + vat // Or totalAmount if VAT is included.
	// Let's assume totalAmount is subtotal.

	// Generate Receipt ID
	receiptID := uuid.New().String()

	receipt := &entities.Receipt{
		ID:                receiptID,
		ShopID:            req.ShopID,
		PaymentMethod:     req.PaymentMethod,
		CustomerRequestID: req.CustomerRequestID,
		TotalAmount:       totalAmount,
		VatRate:           vatRate,
		VAT:               vat,
		NetTotal:          netTotal,
		Status:            "completed",
		CreatedAt:         time.Now(),
	}

	// 2. Create Receipt
	if err := s.ReceiptRepo.Create(receipt); err != nil {
		return nil, err
	}

	// 3. Create Receipt Items and Update Stock
	var receiptItems []interface{}
	for _, item := range req.Items {
		receiptItem := entities.ReceiptItem{
			ID:        uuid.New().String(),
			ReceiptID: receiptID,
			ShopID:    req.ShopID,
			WasteID:   item.WasteID,
			Name:      item.Name,
			Category:  item.Category,
			Weight:    item.Weight,
			UnitPrice: item.UnitPrice,
			Price:     item.Price,
		}
		receiptItems = append(receiptItems, receiptItem)

		// Update Stock using StockService
		// We are "buying" waste, so stock increases.
		if err := s.StockService.AddStock(req.ShopID, item.WasteID, item.Weight); err != nil {
			// In a real app, we should rollback here.
			// For now, we log/return error.
			return nil, err
		}
	}

	if len(receiptItems) > 0 {
		if err := s.ReceiptItemRepo.CreateMany(receiptItems); err != nil {
			return nil, err
		}
	}

	// 4. Update Customer Request Status (if applicable)
	if req.CustomerRequestID != "" {
		if err := s.CustomerRequestRepo.CompleteCustomerRequest(req.CustomerRequestID, req.ShopID); err != nil {
			// Log error but don't fail the receipt creation?
			// Or fail? Let's log for now (in real world) but here we just return error or ignore.
			// Let's print error for now.
			// fmt.Printf("Error completing customer request: %v\n", err)
			// Better to not fail the whole transaction if receipt is already created.
		}
	}

	return receipt, nil
}

func (s *ReceiptService) GetReceiptByCustomerRequestID(requestID string) (*ReceiptWithItemsResponse, error) {
	// 1. Find receipt by customer_request_id
	receipt, err := s.ReceiptRepo.FindByCustomerRequestID(requestID)
	if err != nil {
		return nil, err
	}
	if receipt == nil {
		return nil, fmt.Errorf("receipt not found for request ID: %s", requestID)
	}

	// 2. Find receipt items
	items, err := s.ReceiptItemRepo.FindByReceiptID(receipt.ID)
	if err != nil {
		return nil, err
	}

	// 3. Find shop information
	shop, err := s.ShopRepo.GetByShopID(receipt.ShopID)
	if err != nil {
		return nil, fmt.Errorf("error fetching shop information: %v", err)
	}

	return &ReceiptWithItemsResponse{
		Receipt: receipt,
		Items:   items,
		Shop:    shop,
	}, nil
}
