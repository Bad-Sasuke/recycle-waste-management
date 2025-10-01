package gateways

import (
	"recycle-waste-management-backend/src/middlewares"

	"github.com/gofiber/fiber/v2"
)

func RouteUsers(gateway HTTPGateway, app *fiber.App) {
	api := app.Group("/api/user")

	api.Post("/add_user", gateway.CreateNewUserAccount)
	api.Get("/users", gateway.GetAllUserData)
	api.Put("/update_user", gateway.UpdateUserData)
	api.Delete("/delete_user/:user_id", gateway.DeleteUser)
	api.Get("/get_user", gateway.GetUser)

	// Protected routes for current user profile
	api.Get("/profile", middlewares.SetJWtHeaderHandler(), gateway.GetCurrentUser)
	api.Put("/profile", middlewares.SetJWtHeaderHandler(), gateway.UpdateCurrentUser)
	api.Post("/update-image-profile", middlewares.SetJWtHeaderHandler(), gateway.UpdateProfileImage)
	// Role management (admin only)
	api.Put("/update-role", middlewares.SetJWtHeaderHandler(), gateway.UpdateUserRole)
	api.Get("/my-role", middlewares.SetJWtHeaderHandler(), gateway.GetCurrentUserRole)
}

func RouteRecycle(gateway HTTPGateway, app *fiber.App) {
	api := app.Group("/api/recycle-waste")
	api.Get("/get-wastes", gateway.GetRecycleWaste)

	// Protected routes requiring JWT authentication
	protected := api.Group("", middlewares.SetJWtHeaderHandler())
	protected.Post("/add-waste", gateway.AddRecycleWaste)
	protected.Delete("/delete-waste/:waste_id", gateway.DeleteRecycleWaste)
	protected.Put("/edit-waste/:waste_id", gateway.EditRecycleWaste)
}

func RouteCategoryWaste(gateway HTTPGateway, app *fiber.App) {
	api := app.Group("/api/category-waste")
	api.Get("/get-category", gateway.GetCategoryWaste)
}

func RouteAuth(gateway HTTPGateway, app *fiber.App) {
	api := app.Group("/api/auth")
	api.Get("github/callback", gateway.AuthGithubCallback)
	api.Get("google/callback", gateway.AuthGoogleCallback)
}
