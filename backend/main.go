package main

import (
	"fmt"
	"os"
	"recycle-waste-management-backend/src/configuration"
	ds "recycle-waste-management-backend/src/domain/datasources"
	"recycle-waste-management-backend/src/gateways"
	"recycle-waste-management-backend/src/middlewares"
	repo "recycle-waste-management-backend/src/repositories"
	sv "recycle-waste-management-backend/src/services"
	ws "recycle-waste-management-backend/src/websocket"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load(".env.local")
	if err != nil {
		fmt.Println("Error loading .env file")
	}
}

func main() {
	app := fiber.New(configuration.NewFiberConfiguration())

	app.Use(middlewares.ScalarMiddleware(middlewares.Config{
		PathURL: "/api/docs",
		SpecURL: "./src/docs/swagger.yaml",
	}))
	app.Use(middlewares.MonitorMiddleware("/api/monitor"))
	app.Use(recover.New())
	app.Use(cors.New())
	app.Use(middlewares.Logger())

	mongodb := ds.NewMongoDB(10)

	// Initialize WebSocket Chat Hub
	ws.InitChatHub()

	userMongo := repo.NewUsersRepository(mongodb)
	recycleWastes := repo.NewRecyclableItemsRepository(mongodb)
	categoryWasteRepo := repo.NewCategoryWasteRepository(mongodb)
	shopRepo := repo.NewShopRepository(mongodb)
	settingsRepo := repo.NewSettingsRepository(mongodb)
	customerRequestRepo := repo.NewCustomerRequestRepository(mongodb)
	reviewRepo := repo.NewReviewRepository(mongodb)
	stockRepo := repo.NewStockRepository(mongodb) // Moved up

	userSV := sv.NewUsersService(userMongo)
	stockSV := sv.NewStockService(stockRepo, recycleWastes)                                // Pass recycleWastes repo
	recycleWasteSV := sv.NewRecycleWasteService(recycleWastes, categoryWasteRepo, stockSV) // Updated
	authSV := sv.NewAuthService(userMongo)
	imageSV := sv.NewImageService()
	shopSV := sv.NewShopService(shopRepo, reviewRepo)
	settingsSV := sv.NewSettingsService(settingsRepo)
	customerRequestSV := sv.NewCustomerRequestService(customerRequestRepo, shopRepo)
	reviewSV := sv.NewReviewService(reviewRepo, customerRequestRepo)

	gateways.NewHTTPGateway(app, userSV, recycleWasteSV, authSV, imageSV, shopSV, settingsSV, customerRequestSV)

	// Initialize Review Gateway
	reviewGateway := gateways.NewReviewGateway(reviewSV)
	gateways.RouteReview(reviewGateway, app)

	// Initialize Receipt Gateway
	receiptRepo := repo.NewReceiptRepository(mongodb)
	receiptItemRepo := repo.NewReceiptItemRepository(mongodb)
	// stockRepo and stockSV are already initialized above
	receiptSV := sv.NewReceiptService(receiptRepo, receiptItemRepo, recycleWastes, stockSV, customerRequestRepo, shopRepo, userMongo)
	receiptGateway := gateways.NewReceiptGateway(receiptSV)
	gateways.RouteReceipt(receiptGateway, app)

	// Initialize Stock Gateway
	stockGateway := gateways.NewStockGateway(stockSV)
	gateways.RouteStock(stockGateway, app)

	// Initialize Employee Gateway
	employeeRepo := repo.NewEmployeeRepository(mongodb)
	employeeSV := sv.NewEmployeeService(employeeRepo)
	employeeGateway := gateways.NewEmployeeGateway(employeeSV, shopRepo)
	gateways.RouteEmployee(employeeGateway, app)

	PORT := os.Getenv("PORT")
	if PORT == "" {
		PORT = "8080"
	}

	app.Listen(":" + PORT)
}
