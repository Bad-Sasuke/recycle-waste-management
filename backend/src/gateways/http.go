package gateways

import (
	service "recycle-waste-management-backend/src/services"

	"github.com/gofiber/fiber/v2"
)

type HTTPGateway struct {
	UserService    service.IUsersService
	RecycleService service.IRecycleWasteService
}

func NewHTTPGateway(app *fiber.App, users service.IUsersService, recWasteSV service.IRecycleWasteService) {
	gateway := &HTTPGateway{
		UserService:    users,
		RecycleService: recWasteSV,
	}

	RouteUsers(*gateway, app)
	RouteRecycle(*gateway, app)
}
