package gateways

import (
	service "recycle-waste-management-backend/src/services"

	"github.com/gofiber/fiber/v2"
)

type HTTPGateway struct {
	UserService    service.IUsersService
	RecycleService service.IRecycleWasteService
	AuthService    service.IAuthService
	ImageService   service.IImageService
	ShopService    service.IShopService
}

func NewHTTPGateway(app *fiber.App, users service.IUsersService, recWasteSV service.IRecycleWasteService, authService service.IAuthService, imageService service.IImageService, shopService service.IShopService) {
	gateway := &HTTPGateway{
		UserService:    users,
		RecycleService: recWasteSV,
		AuthService:    authService,
		ImageService:   imageService,
		ShopService:    shopService,
	}

	RouteUsers(*gateway, app)
	RouteRecycle(*gateway, app)
	RouteCategoryWaste(*gateway, app)
	RouteAuth(*gateway, app)
	RouteShop(*gateway, app)
}
