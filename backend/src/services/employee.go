package services

import (
	"crypto/rand"
	"encoding/base32"
	"fmt"
	"math"

	"golang.org/x/crypto/bcrypt"

	"recycle-waste-management-backend/src/domain/entities"
	"recycle-waste-management-backend/src/domain/models"
	"recycle-waste-management-backend/src/repositories"
)

type employeeService struct {
	EmployeeRepository repositories.IEmployeeRepository
}

type IEmployeeService interface {
	CreateEmployee(req *entities.CreateEmployeeRequest) (*entities.EmployeeResponse, error)
	GetEmployeeByID(employeeID string) (*entities.EmployeeResponse, error)
	GetEmployeesByShopID(shopID string, page, pageSize int) (*entities.EmployeeListResponse, error)
	UpdateEmployee(employeeID string, req *entities.UpdateEmployeeRequest) (*entities.EmployeeResponse, error)
	DeleteEmployee(employeeID string) error
	ValidateEmployeePassword(username, password string) (*models.Employee, error)
}

func NewEmployeeService(repo repositories.IEmployeeRepository) IEmployeeService {
	return &employeeService{
		EmployeeRepository: repo,
	}
}

func (s *employeeService) CreateEmployee(req *entities.CreateEmployeeRequest) (*entities.EmployeeResponse, error) {
	// Check if username already exists
	existingEmployee, err := s.EmployeeRepository.GetEmployeeByUsername(req.Username)
	if err != nil {
		return nil, fmt.Errorf("error checking existing username: %v", err)
	}
	if existingEmployee != nil {
		return nil, fmt.Errorf("username already exists")
	}

	// Generate employee ID
	secret := make([]byte, 8)
	_, err = rand.Reader.Read(secret)
	if err != nil {
		return nil, fmt.Errorf("error generating employee ID: %v", err)
	}
	var b32NoPadding = base32.StdEncoding.WithPadding(base32.NoPadding)
	employeeID := "EMP_" + b32NoPadding.EncodeToString(secret)

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("error hashing password: %v", err)
	}

	// Create employee model
	employee := &models.Employee{
		EmployeeID: employeeID,
		ShopID:     req.ShopID,
		FirstName:  req.FirstName,
		LastName:   req.LastName,
		Username:   req.Username,
		Password:   string(hashedPassword),
	}

	// Save to database
	if err := s.EmployeeRepository.CreateEmployee(employee); err != nil {
		return nil, fmt.Errorf("error creating employee: %v", err)
	}

	// Return response (without password)
	response := &entities.EmployeeResponse{
		EmployeeID: employee.EmployeeID,
		ShopID:     employee.ShopID,
		FirstName:  employee.FirstName,
		LastName:   employee.LastName,
		Username:   employee.Username,
		CreatedAt:  employee.CreatedAt,
		UpdatedAt:  employee.UpdatedAt,
	}

	return response, nil
}

func (s *employeeService) GetEmployeeByID(employeeID string) (*entities.EmployeeResponse, error) {
	employee, err := s.EmployeeRepository.GetEmployeeByID(employeeID)
	if err != nil {
		return nil, err
	}

	response := &entities.EmployeeResponse{
		EmployeeID: employee.EmployeeID,
		ShopID:     employee.ShopID,
		FirstName:  employee.FirstName,
		LastName:   employee.LastName,
		Username:   employee.Username,
		CreatedAt:  employee.CreatedAt,
		UpdatedAt:  employee.UpdatedAt,
	}

	return response, nil
}

func (s *employeeService) GetEmployeesByShopID(shopID string, page, pageSize int) (*entities.EmployeeListResponse, error) {
	employees, total, err := s.EmployeeRepository.GetEmployeesByShopID(shopID, page, pageSize)
	if err != nil {
		return nil, err
	}

	// Convert to response format (without passwords)
	var employeeResponses []entities.EmployeeResponse
	for _, emp := range employees {
		employeeResponses = append(employeeResponses, entities.EmployeeResponse{
			EmployeeID: emp.EmployeeID,
			ShopID:     emp.ShopID,
			FirstName:  emp.FirstName,
			LastName:   emp.LastName,
			Username:   emp.Username,
			CreatedAt:  emp.CreatedAt,
			UpdatedAt:  emp.UpdatedAt,
		})
	}

	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))

	response := &entities.EmployeeListResponse{
		Employees:  employeeResponses,
		TotalPages: totalPages,
		TotalItems: total,
		Page:       page,
		PageSize:   pageSize,
	}

	return response, nil
}

func (s *employeeService) UpdateEmployee(employeeID string, req *entities.UpdateEmployeeRequest) (*entities.EmployeeResponse, error) {
	// Get existing employee
	existingEmployee, err := s.EmployeeRepository.GetEmployeeByID(employeeID)
	if err != nil {
		return nil, err
	}

	// Check if username is being changed and if it already exists
	if req.Username != "" && req.Username != existingEmployee.Username {
		usernameExists, err := s.EmployeeRepository.GetEmployeeByUsername(req.Username)
		if err != nil {
			return nil, fmt.Errorf("error checking existing username: %v", err)
		}
		if usernameExists != nil {
			return nil, fmt.Errorf("username already exists")
		}
	}

	// Update fields
	updateData := &models.Employee{
		EmployeeID: employeeID,
		ShopID:     existingEmployee.ShopID,
		FirstName:  existingEmployee.FirstName,
		LastName:   existingEmployee.LastName,
		Username:   existingEmployee.Username,
		Password:   existingEmployee.Password,
		CreatedAt:  existingEmployee.CreatedAt,
	}

	if req.FirstName != "" {
		updateData.FirstName = req.FirstName
	}
	if req.LastName != "" {
		updateData.LastName = req.LastName
	}
	if req.Username != "" {
		updateData.Username = req.Username
	}

	// Hash new password if provided
	if req.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, fmt.Errorf("error hashing password: %v", err)
		}
		updateData.Password = string(hashedPassword)
	}

	// Update in database
	if err := s.EmployeeRepository.UpdateEmployee(employeeID, updateData); err != nil {
		return nil, err
	}

	// Get updated employee
	updatedEmployee, err := s.EmployeeRepository.GetEmployeeByID(employeeID)
	if err != nil {
		return nil, err
	}

	response := &entities.EmployeeResponse{
		EmployeeID: updatedEmployee.EmployeeID,
		ShopID:     updatedEmployee.ShopID,
		FirstName:  updatedEmployee.FirstName,
		LastName:   updatedEmployee.LastName,
		Username:   updatedEmployee.Username,
		CreatedAt:  updatedEmployee.CreatedAt,
		UpdatedAt:  updatedEmployee.UpdatedAt,
	}

	return response, nil
}

func (s *employeeService) DeleteEmployee(employeeID string) error {
	return s.EmployeeRepository.DeleteEmployee(employeeID)
}

func (s *employeeService) ValidateEmployeePassword(username, password string) (*models.Employee, error) {
	employee, err := s.EmployeeRepository.GetEmployeeByUsername(username)
	if err != nil {
		return nil, fmt.Errorf("error finding employee: %v", err)
	}
	if employee == nil {
		return nil, fmt.Errorf("employee not found")
	}

	// Compare hashed password
	err = bcrypt.CompareHashAndPassword([]byte(employee.Password), []byte(password))
	if err != nil {
		return nil, fmt.Errorf("invalid password")
	}

	return employee, nil
}
