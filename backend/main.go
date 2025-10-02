package main

import (
	"fmt"
	"os"
	"recycle-waste-management-backend/src/configuration"
	ds "recycle-waste-management-backend/src/domain/datasources"
	repo "recycle-waste-management-backend/src/domain/repositories"
	"recycle-waste-management-backend/src/gateways"
	"recycle-waste-management-backend/src/middlewares"
	sv "recycle-waste-management-backend/src/services"

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

	userMongo := repo.NewUsersRepository(mongodb)
	recycleWastes := repo.NewRecyclableItemsRepository(mongodb)
	categoryWasteRepo := repo.NewCategoryWasteRepository(mongodb)
	shopRepo := repo.NewShopRepository(mongodb)
	settingsRepo := repo.NewSettingsRepository(mongodb)

	userSV := sv.NewUsersService(userMongo)
	recycleWasteSV := sv.NewRecycleWasteService(recycleWastes, categoryWasteRepo)
	authSV := sv.NewAuthService(userMongo)
	imageSV := sv.NewImageService()
	shopSV := sv.NewShopService(shopRepo)
	settingsSV := sv.NewSettingsService(settingsRepo)
	gateways.NewHTTPGateway(app, userSV, recycleWasteSV, authSV, imageSV, shopSV, settingsSV)

	PORT := os.Getenv("PORT")
	if PORT == "" {
		PORT = "8080"
	}

	app.Listen(":" + PORT)
}
