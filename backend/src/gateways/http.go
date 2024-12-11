package gateways

import (
	service "recycle-waste-management-backend/src/services"

	"github.com/gofiber/fiber/v2"
)

type HTTPGateway struct {
	UserService    service.IUsersService
	RecycleService service.IRecycleWasteService
	AuthService    service.IAuthService
}

func NewHTTPGateway(app *fiber.App, users service.IUsersService, recWasteSV service.IRecycleWasteService, authService service.IAuthService) {
	gateway := &HTTPGateway{
		UserService:    users,
		RecycleService: recWasteSV,
		AuthService:    authService,
	}

	RouteUsers(*gateway, app)
	RouteRecycle(*gateway, app)
	RouteCategoryWaste(*gateway, app)
	RouteAuth(*gateway, app)
}
