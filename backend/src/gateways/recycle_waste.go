package gateways

import (
	"fmt"
	"io"
	"recycle-waste-management-backend/src/domain/entities"
	"recycle-waste-management-backend/src/middlewares"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func (h *HTTPGateway) GetRecycleWaste(ctx *fiber.Ctx) error {
	page := 1
	limit := 12

	pageQuery := ctx.Query("page")
	limitQuery := ctx.Query("limit")
	shopIDQuery := ctx.Query("shop_id") // Get shop_id from query params

	if pageQuery != "" {
		if p, err := strconv.Atoi(pageQuery); err == nil && p > 0 {
			page = p
		}
	}

	if limitQuery != "" {
		if l, err := strconv.Atoi(limitQuery); err == nil && l > 0 {
			limit = l
		}
	}

	var data *[]entities.RecyclableItemsModel
	var totalCount int64
	var err error

	if shopIDQuery != "" {
		// Get recyclable items for a specific shop
		data, totalCount, err = h.RecycleService.GetRecyclableItemsByShopIDPaginated(shopIDQuery, page, limit)
	} else {
		// Get all recyclable items regardless of shop
		data, totalCount, err = h.RecycleService.GetRecyclableItemsPaginated(page, limit)
	}

	if err != nil {
		return ctx.Status(fiber.StatusForbidden).JSON(entities.ResponseModel{Message: err.Error()})
	}

	totalPages := int((totalCount + int64(limit) - 1) / int64(limit)) // Ceiling division

	response := entities.ResponsePaginationModel{
		Message:    "success",
		Data:       data,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
		TotalItems: totalCount,
	}

	return ctx.Status(fiber.StatusOK).JSON(response)
}

func (h *HTTPGateway) AddRecycleWaste(ctx *fiber.Ctx) error {
	// Decode JWT token to get user ID
	tokenDetails, err := middlewares.DecodeJWTToken(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(entities.ResponseMessage{Message: "Unauthorized access"})
	}

	// Get user data to check role
	userData, err := h.UserService.GetUser(tokenDetails.UserID)
	if err != nil {
		return ctx.Status(fiber.StatusForbidden).JSON(entities.ResponseMessage{Message: "User not found"})
	}

	name := ctx.FormValue("name")
	price := ctx.FormValue("price")
	category := ctx.FormValue("category")
	// Don't get shop_id from the request anymore - we'll determine it automatically
	imageFile, err := ctx.FormFile("image_file")
	if err != nil {
		fmt.Println(err)
	}
	fileBytes := []byte{}
	if imageFile != nil {
		file, err := imageFile.Open()
		if err != nil {
			fmt.Println(err)
		}
		defer file.Close()

		fileBytes, err = io.ReadAll(file)
		if err != nil {
			fmt.Println(err)
		}
	}
	priceFloat, err := strconv.ParseFloat(price, 64)
	if err != nil {
		fmt.Println(err)
		return ctx.Status(fiber.StatusBadRequest).JSON(entities.ResponseMessage{Message: err.Error()})
	}

	// Automatically determine the shop_id based on the user_id
	shopID := ""
	if userData.Role != string(entities.UserRoleAdmin) && userData.Role != string(entities.UserRoleModerator) {
		// For regular users, get their shop automatically
		userShop, err := h.ShopService.GetShopByUserID(tokenDetails.UserID)
		if err != nil || userShop == nil {
			return ctx.Status(fiber.StatusForbidden).JSON(entities.ResponseMessage{Message: "Access denied: You must have a shop to add items to"})
		}
		shopID = userShop.ShopID
	} else {
		// For admin and moderator, allow them to specify a shop_id in the request (or get from request for backward compatibility)
		// But normally, they would have some way to specify the target shop in a real implementation
		requestedShopID := ctx.FormValue("shop_id") // Still allow shop_id in the request for admin/moderator
		if requestedShopID != "" {
			shopID = requestedShopID
		} else {
			// If no shop_id in request, get the admin/moderator's own shop if they have one
			userShop, err := h.ShopService.GetShopByUserID(tokenDetails.UserID)
			if err != nil || userShop == nil {
				return ctx.Status(fiber.StatusForbidden).JSON(entities.ResponseMessage{Message: "Access denied: Admin/Moderator must specify a shop_id or have a shop"})
			}
			shopID = userShop.ShopID
		}
	}

	bodyData := entities.RecyclableItemsModel{
		Name:     name,
		Price:    priceFloat,
		Category: category,
		ShopID:   shopID, // Include the auto-determined shop_id in the data
	}
	if err := h.RecycleService.AddRecycleWaste(bodyData, fileBytes); err != nil {
		return ctx.Status(fiber.StatusForbidden).JSON(entities.ResponseMessage{Message: err.Error()})
	}
	return ctx.Status(fiber.StatusOK).JSON(entities.ResponseModel{Message: "success"})
}

func (h *HTTPGateway) DeleteRecycleWaste(ctx *fiber.Ctx) error {
	// Decode JWT token to get user ID
	tokenDetails, err := middlewares.DecodeJWTToken(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(entities.ResponseMessage{Message: "Unauthorized access"})
	}

	// Get user data to check role
	userData, err := h.UserService.GetUser(tokenDetails.UserID)
	if err != nil {
		return ctx.Status(fiber.StatusForbidden).JSON(entities.ResponseMessage{Message: "User not found"})
	}

	id := ctx.Params("waste_id")
	if id == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(entities.ResponseMessage{Message: "invalid query params"})
	}

	// First, get the existing waste item to verify ownership
	foundItem, err := h.RecycleService.GetRecyclableItemByWasteID(id)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(entities.ResponseMessage{Message: "Waste item not found"})
	}

	// Check permissions: admin and moderator can delete items for any shop
	// Regular users can only delete items from their own shop
	if userData.Role != string(entities.UserRoleAdmin) && userData.Role != string(entities.UserRoleModerator) {
		// Check if user owns the shop that has this item
		if foundItem.ShopID != "" {
			shop, err := h.ShopService.GetShopByShopID(foundItem.ShopID)
			if err != nil || shop == nil {
				return ctx.Status(fiber.StatusForbidden).JSON(entities.ResponseMessage{Message: "Access denied: Shop not found"})
			}
			if shop.UserID != tokenDetails.UserID {
				return ctx.Status(fiber.StatusForbidden).JSON(entities.ResponseMessage{Message: "Access denied: You don't own this shop"})
			}
		} else {
			// If the item doesn't have a shop_id, user needs to be admin/moderator
			return ctx.Status(fiber.StatusForbidden).JSON(entities.ResponseMessage{Message: "Access denied: You don't have permission to delete this item"})
		}
	}

	err = h.RecycleService.DeleteWasteItem(id)
	if err != nil {
		return ctx.Status(fiber.StatusForbidden).JSON(entities.ResponseMessage{Message: err.Error()})
	}
	return ctx.Status(fiber.StatusOK).JSON(entities.ResponseModel{Message: "success"})
}

func (h *HTTPGateway) GetCategoryWaste(ctx *fiber.Ctx) error {
	data, err := h.RecycleService.GetCategoryWaste()
	if err != nil {
		return ctx.Status(fiber.StatusForbidden).JSON(entities.ResponseModel{Message: err.Error()})
	}
	return ctx.Status(fiber.StatusOK).JSON(entities.ResponseModel{Message: "success", Data: data})
}

func (h *HTTPGateway) EditRecycleWaste(ctx *fiber.Ctx) error {
	// Decode JWT token to get user ID
	tokenDetails, err := middlewares.DecodeJWTToken(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(entities.ResponseMessage{Message: "Unauthorized access"})
	}

	// Get user data to check role
	userData, err := h.UserService.GetUser(tokenDetails.UserID)
	if err != nil {
		return ctx.Status(fiber.StatusForbidden).JSON(entities.ResponseMessage{Message: "User not found"})
	}

	wasteID := ctx.Params("waste_id")
	name := ctx.FormValue("name")
	price := ctx.FormValue("price")
	category := ctx.FormValue("category")
	// Don't get shop_id from the request anymore - we'll determine it automatically
	imageFile, err := ctx.FormFile("image_file")
	if err != nil {
		fmt.Println(err)
	}
	fileBytes := []byte{}
	if imageFile != nil {
		file, err := imageFile.Open()
		if err != nil {
			fmt.Println(err)
		}
		defer file.Close()

		fileBytes, err = io.ReadAll(file)
		if err != nil {
			fmt.Println(err)
		}
	}
	priceFloat, err := strconv.ParseFloat(price, 64)
	if err != nil {
		fmt.Println(err)
		return ctx.Status(fiber.StatusBadRequest).JSON(entities.ResponseMessage{Message: err.Error()})
	}

	// First, get the existing waste item to verify ownership
	foundItem, err := h.RecycleService.GetRecyclableItemByWasteID(wasteID)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(entities.ResponseMessage{Message: "Waste item not found"})
	}

	// Check permissions: admin and moderator can edit items for any shop
	// Regular users can only edit items from their own shop
	if userData.Role != string(entities.UserRoleAdmin) && userData.Role != string(entities.UserRoleModerator) {
		// Check if user owns the shop that has this item
		if foundItem.ShopID != "" {
			shop, err := h.ShopService.GetShopByShopID(foundItem.ShopID)
			if err != nil || shop == nil {
				return ctx.Status(fiber.StatusForbidden).JSON(entities.ResponseMessage{Message: "Access denied: Shop not found"})
			}
			if shop.UserID != tokenDetails.UserID {
				return ctx.Status(fiber.StatusForbidden).JSON(entities.ResponseMessage{Message: "Access denied: You don't own this shop"})
			}
		} else {
			// If the item doesn't have a shop_id, user needs to be admin/moderator
			return ctx.Status(fiber.StatusForbidden).JSON(entities.ResponseMessage{Message: "Access denied: You don't have permission to edit this item"})
		}
	}

	// Automatically determine the shop_id based on the user_id
	shopID := ""
	if userData.Role != string(entities.UserRoleAdmin) && userData.Role != string(entities.UserRoleModerator) {
		// For regular users, get their shop automatically
		userShop, err := h.ShopService.GetShopByUserID(tokenDetails.UserID)
		if err != nil || userShop == nil {
			return ctx.Status(fiber.StatusForbidden).JSON(entities.ResponseMessage{Message: "Access denied: You must have a shop to edit items in"})
		}
		shopID = userShop.ShopID
	} else {
		// For admin and moderator, allow them to specify a shop_id in the request (or get from request for backward compatibility)
		// But normally, they would have some way to specify the target shop in a real implementation
		requestedShopID := ctx.FormValue("shop_id") // Still allow shop_id in the request for admin/moderator
		if requestedShopID != "" {
			shopID = requestedShopID
		} else {
			// If no shop_id in request, use the existing shop_id of the item if it's theirs
			// Or get the admin/moderator's own shop if they have one
			userShop, err := h.ShopService.GetShopByUserID(tokenDetails.UserID)
			if err != nil || userShop == nil {
				return ctx.Status(fiber.StatusForbidden).JSON(entities.ResponseMessage{Message: "Access denied: Admin/Moderator must specify a shop_id or have a shop"})
			}
			shopID = userShop.ShopID
		}
	}

	bodyData := entities.RecyclableItemsModel{
		Name:     name,
		Price:    priceFloat,
		Category: category,
		ShopID:   shopID, // Include the auto-determined shop_id in the data
	}
	if err := h.RecycleService.EditWasteItem(wasteID, bodyData, fileBytes); err != nil {
		return ctx.Status(fiber.StatusForbidden).JSON(entities.ResponseMessage{Message: err.Error()})
	}
	return ctx.Status(fiber.StatusOK).JSON(entities.ResponseModel{Message: "success"})
}