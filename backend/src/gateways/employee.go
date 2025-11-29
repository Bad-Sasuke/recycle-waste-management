package gateways

import (
	"strconv"

	"recycle-waste-management-backend/src/domain/entities"
	"recycle-waste-management-backend/src/middlewares"
	"recycle-waste-management-backend/src/repositories"
	"recycle-waste-management-backend/src/services"

	"github.com/gofiber/fiber/v2"
)

type EmployeeGateway struct {
	EmployeeService services.IEmployeeService
	ShopRepository  repositories.IShopRepository
}

func NewEmployeeGateway(service services.IEmployeeService, shopRepo repositories.IShopRepository) *EmployeeGateway {
	return &EmployeeGateway{
		EmployeeService: service,
		ShopRepository:  shopRepo,
	}
}

func (g *EmployeeGateway) SetupRoutes(app *fiber.App) {
	employee := app.Group("/api/employees")

	// Public route - login
	employee.Post("/login", g.EmployeeLogin)

	// Protected routes - require JWT authentication
	employee.Post("/", middlewares.SetJWtHeaderHandler(), g.CreateEmployee)
	employee.Get("/shop/:shop_id", middlewares.SetJWtHeaderHandler(), g.GetEmployeesByShopID)
	employee.Get("/:employee_id", middlewares.SetJWtHeaderHandler(), g.GetEmployeeByID)
	employee.Put("/:employee_id", middlewares.SetJWtHeaderHandler(), g.UpdateEmployee)
	employee.Delete("/:employee_id", middlewares.SetJWtHeaderHandler(), g.DeleteEmployee)
}

// EmployeeLogin handles employee authentication
func (g *EmployeeGateway) EmployeeLogin(c *fiber.Ctx) error {
	var req struct {
		ShopCode string `json:"shop_code"`
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(entities.ResponseModel{
			Message: "Invalid request body",
			Status:  fiber.StatusBadRequest,
		})
	}

	// Validate required fields
	if req.ShopCode == "" || req.Username == "" || req.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(entities.ResponseModel{
			Message: "Shop code, username, and password are required",
			Status:  fiber.StatusBadRequest,
		})
	}

	// Check if shop exists
	shop, err := g.ShopRepository.GetByShopCode(req.ShopCode)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(entities.ResponseModel{
			Message: "Invalid shop code",
			Status:  fiber.StatusUnauthorized,
		})
	}

	// Validate employee credentials
	employee, err := g.EmployeeService.ValidateEmployeePassword(req.Username, req.Password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(entities.ResponseModel{
			Message: "Invalid username or password",
			Status:  fiber.StatusUnauthorized,
		})
	}

	// Verify employee belongs to this shop
	if employee.ShopID != shop.ShopID {
		return c.Status(fiber.StatusUnauthorized).JSON(entities.ResponseModel{
			Message: "Employee does not belong to this shop",
			Status:  fiber.StatusUnauthorized,
		})
	}

	// Generate JWT token for employee
	tokenDetails, err := middlewares.GenerateJWTToken(employee.EmployeeID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(entities.ResponseModel{
			Message: "Error generating token",
			Status:  fiber.StatusInternalServerError,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Login successful",
		"token":   *tokenDetails.Token,
		"employee": map[string]interface{}{
			"employee_id": employee.EmployeeID,
			"shop_id":     employee.ShopID,
			"first_name":  employee.FirstName,
			"last_name":   employee.LastName,
			"username":    employee.Username,
		},
	})
}

func (g *EmployeeGateway) CreateEmployee(c *fiber.Ctx) error {
	var req entities.CreateEmployeeRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(entities.ResponseModel{
			Message: "Invalid request body",
			Status:  fiber.StatusBadRequest,
		})
	}

	// Validate required fields
	if req.FirstName == "" || req.LastName == "" || req.Username == "" || req.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(entities.ResponseModel{
			Message: "First name, last name, username, and password are required",
			Status:  fiber.StatusBadRequest,
		})
	}

	if len(req.Password) < 6 {
		return c.Status(fiber.StatusBadRequest).JSON(entities.ResponseModel{
			Message: "Password must be at least 6 characters",
			Status:  fiber.StatusBadRequest,
		})
	}

	employee, err := g.EmployeeService.CreateEmployee(&req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(entities.ResponseModel{
			Message: err.Error(),
			Status:  fiber.StatusInternalServerError,
		})
	}

	return c.Status(fiber.StatusCreated).JSON(entities.ResponseModel{
		Message: "Employee created successfully",
		Data:    employee,
		Status:  fiber.StatusCreated,
	})
}

func (g *EmployeeGateway) GetEmployeesByShopID(c *fiber.Ctx) error {
	shopID := c.Params("shop_id")
	if shopID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(entities.ResponseModel{
			Message: "Shop ID is required",
			Status:  fiber.StatusBadRequest,
		})
	}

	// Get pagination parameters
	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(c.Query("page_size", "10"))
	if err != nil || pageSize < 1 {
		pageSize = 10
	}

	employees, err := g.EmployeeService.GetEmployeesByShopID(shopID, page, pageSize)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(entities.ResponseModel{
			Message: err.Error(),
			Status:  fiber.StatusInternalServerError,
		})
	}

	return c.Status(fiber.StatusOK).JSON(entities.ResponsePaginationModel{
		Message:    "Employees retrieved successfully",
		Data:       employees.Employees,
		TotalPages: employees.TotalPages,
		TotalItems: int64(employees.TotalItems),
		Page:       employees.Page,
		Limit:      employees.PageSize,
		Status:     fiber.StatusOK,
	})
}

func (g *EmployeeGateway) GetEmployeeByID(c *fiber.Ctx) error {
	employeeID := c.Params("employee_id")
	if employeeID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(entities.ResponseModel{
			Message: "Employee ID is required",
			Status:  fiber.StatusBadRequest,
		})
	}

	employee, err := g.EmployeeService.GetEmployeeByID(employeeID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(entities.ResponseModel{
			Message: err.Error(),
			Status:  fiber.StatusNotFound,
		})
	}

	return c.Status(fiber.StatusOK).JSON(entities.ResponseModel{
		Message: "Employee retrieved successfully",
		Data:    employee,
		Status:  fiber.StatusOK,
	})
}

func (g *EmployeeGateway) UpdateEmployee(c *fiber.Ctx) error {
	employeeID := c.Params("employee_id")
	if employeeID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(entities.ResponseModel{
			Message: "Employee ID is required",
			Status:  fiber.StatusBadRequest,
		})
	}

	var req entities.UpdateEmployeeRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(entities.ResponseModel{
			Message: "Invalid request body",
			Status:  fiber.StatusBadRequest,
		})
	}

	// Validate password length if provided
	if req.Password != "" && len(req.Password) < 6 {
		return c.Status(fiber.StatusBadRequest).JSON(entities.ResponseModel{
			Message: "Password must be at least 6 characters",
			Status:  fiber.StatusBadRequest,
		})
	}

	employee, err := g.EmployeeService.UpdateEmployee(employeeID, &req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(entities.ResponseModel{
			Message: err.Error(),
			Status:  fiber.StatusInternalServerError,
		})
	}

	return c.Status(fiber.StatusOK).JSON(entities.ResponseModel{
		Message: "Employee updated successfully",
		Data:    employee,
		Status:  fiber.StatusOK,
	})
}

func (g *EmployeeGateway) DeleteEmployee(c *fiber.Ctx) error {
	employeeID := c.Params("employee_id")
	if employeeID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(entities.ResponseModel{
			Message: "Employee ID is required",
			Status:  fiber.StatusBadRequest,
		})
	}

	if err := g.EmployeeService.DeleteEmployee(employeeID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(entities.ResponseModel{
			Message: err.Error(),
			Status:  fiber.StatusInternalServerError,
		})
	}

	return c.Status(fiber.StatusOK).JSON(entities.ResponseModel{
		Message: "Employee deleted successfully",
		Status:  fiber.StatusOK,
	})
}
