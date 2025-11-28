package gateways

import (
	"recycle-waste-management-backend/src/middlewares"

	"github.com/gofiber/contrib/websocket"
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
	api.Put("/update-role-user", middlewares.SetJWtHeaderHandler(), gateway.UpdateUserRoleUser)
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

func RouteShop(gateway HTTPGateway, app *fiber.App) {
	api := app.Group("/api/shops")

	// Public routes
	api.Get("/get-shops", gateway.GetAllShops)
	api.Get("/get-shop/:shop_id", gateway.GetShopByShopID)
	api.Get("/:shop_id", gateway.GetShopByShopID) // Add this line for cleaner URL

	// Protected routes requiring JWT authentication
	protected := api.Group("", middlewares.SetJWtHeaderHandler())
	protected.Post("/create-shop", gateway.CreateShop)
	protected.Get("/my-shop", gateway.GetShopByUserID)
	protected.Put("/update-shop/:shop_id", gateway.UpdateShop)
	protected.Delete("/delete-shop/:shop_id", gateway.DeleteShop)
}

func RouteSettings(gateway HTTPGateway, app *fiber.App) {
	api := app.Group("/api/settings")

	// Protected routes requiring JWT authentication
	protected := api.Group("", middlewares.SetJWtHeaderHandler())
	protected.Get("/", gateway.GetUserSettings)
	protected.Put("/", gateway.UpdateUserSettings)
	protected.Delete("/", gateway.DeleteUserSettings)
}

func RouteCustomerRequest(gateway HTTPGateway, app *fiber.App) {
	api := app.Group("/api/customer-request", middlewares.SetJWtHeaderHandler())
	api.Get("/my-request", gateway.GetCustomerRequestByRequestID)
	api.Post("", gateway.AddCustomerRequest)
	api.Delete("", gateway.DeleteCustomerRequest)
	api.Put("", gateway.UpdateCustomerRequest)
	api.Get("/all", gateway.GetCustomerRequests)
	api.Put("/accept/:id", gateway.AcceptCustomerRequest)
	api.Put("/cancel/:id", gateway.CancelCustomerRequest)
	api.Put("/complete/:id", gateway.CompleteCustomerRequest)

}

func RouteWebSocket(gateway HTTPGateway, app *fiber.App) {
	// WebSocket chat endpoint
	app.Get("/ws/chat",
		gateway.WebSocketChatUpgrade,
		websocket.New(gateway.HandleWebSocketChat),
	)
}

func RouteReview(reviewGateway *ReviewGateway, app *fiber.App) {
	api := app.Group("/api/reviews")

	// Public routes
	api.Get("/shop/:shop_id", reviewGateway.GetReviewsByShopID)
	api.Get("/check/:customer_request_id", reviewGateway.CheckReviewExists)

	// Protected routes requiring JWT authentication
	protected := api.Group("", middlewares.SetJWtHeaderHandler())
	protected.Post("", reviewGateway.CreateReview)
	protected.Post("/skip", reviewGateway.SkipReview)
}

func RouteReceipt(receiptGateway *ReceiptGateway, app *fiber.App) {
	api := app.Group("/api/receipts")

	// Protected routes requiring JWT authentication
	protected := api.Group("", middlewares.SetJWtHeaderHandler())
	protected.Post("", receiptGateway.CreateReceipt)
	protected.Get("/by-request/:request_id", receiptGateway.GetReceiptByRequestID)
}
