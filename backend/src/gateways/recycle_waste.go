package gateways

import (
	"fmt"
	"io"
	"recycle-waste-management-backend/src/domain/entities"
	"recycle-waste-management-backend/src/middlewares"
	"sort"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func (h *HTTPGateway) GetRecycleWaste(ctx *fiber.Ctx) error {
	page := 1
	limit := 12

	pageQuery := ctx.Query("page")
	limitQuery := ctx.Query("limit")
	shopIDQuery := ctx.Query("shop_id")    // Get shop_id from query params
	searchQuery := ctx.Query("search")     // Get search query from params
	categoryQuery := ctx.Query("category") // Get category query from params

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
		// Check if category filter is also provided
		if categoryQuery != "" {
			// Get recyclable items for a specific shop and category
			data, totalCount, err = h.RecycleService.GetRecyclableItemsByShopIDAndCategoryPaginated(shopIDQuery, categoryQuery, page, limit)
		} else {
			// Get recyclable items for a specific shop
			data, totalCount, err = h.RecycleService.GetRecyclableItemsByShopIDPaginated(shopIDQuery, page, limit)
		}
	} else {
		// Get all recyclable items regardless of shop
		data, totalCount, err = h.RecycleService.GetRecyclableItemsPaginated(page, limit)
	}

	if err != nil {
		return ctx.Status(fiber.StatusForbidden).JSON(entities.ResponseModel{Message: err.Error()})
	}

	// Apply search filter if provided
	if searchQuery != "" {
		var filteredData []entities.RecyclableItemsModel
		searchLower := strings.ToLower(strings.TrimSpace(searchQuery))

		for _, item := range *data {
			// Check if the item's name contains the search query (case-insensitive)
			if item.Name != "" && strings.Contains(strings.ToLower(item.Name), searchLower) {
				filteredData = append(filteredData, item)
			} else if item.Category != "" && strings.Contains(strings.ToLower(item.Category), searchLower) {
				// Also check category for search matches
				filteredData = append(filteredData, item)
			}
		}

		// Update data to be the filtered results
		filteredDataPtr := &filteredData
		data = filteredDataPtr

		// Update total count to reflect filtered results
		totalCount = int64(len(filteredData))
	}

	totalPages := int((totalCount + int64(limit) - 1) / int64(limit)) // Ceiling division

	// Group items by name if not filtering by shop_id
	var responseData interface{}
	if shopIDQuery == "" {
		// Group the recyclable items by name
		groupedData := groupRecyclableItemsByName(*data, h)
		responseData = groupedData
	} else {
		// Return original data if filtering by specific shop
		responseData = data
	}

	response := entities.ResponsePaginationModel{
		Message:    "success",
		Data:       responseData,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
		TotalItems: totalCount,
	}

	return ctx.Status(fiber.StatusOK).JSON(response)
}

// getShopInfo retrieves the shop info (name and image) for a given shop_id
func getShopInfo(shopID string, gateway *HTTPGateway) entities.ShopInfo {
	if shopID == "" {
		return entities.ShopInfo{
			ShopID:       "",
			ShopName:     "",
			ShopImageURL: "",
		}
	}
	shop, err := gateway.ShopService.GetShopByShopID(shopID)
	if err != nil || shop == nil {
		return entities.ShopInfo{
			ShopID:       shopID,
			ShopName:     "",
			ShopImageURL: "",
		}
	}
	return entities.ShopInfo{
		ShopID:       shopID,
		ShopName:     shop.Name,
		ShopImageURL: shop.ImageURL,
	}
}

// groupRecyclableItemsByName groups recyclable items by their name and collects shop info
func groupRecyclableItemsByName(items []entities.RecyclableItemsModel, gateway *HTTPGateway) []entities.GroupedRecyclableItem {
	groupMap := make(map[string]*entities.GroupedRecyclableItem)

	for _, item := range items {
		key := item.Name
		shopInfo := getShopInfo(item.ShopID, gateway)

		if existingItem, exists := groupMap[key]; exists {
			// If item with same name exists, add shop info and waste_id to arrays
			existingItem.Shops = append(existingItem.Shops, shopInfo)
			existingItem.WasteIDs = append(existingItem.WasteIDs, item.WasteID)
			// Update the LastUpdate to the most recent one
			if item.LastUpdate.After(existingItem.LastUpdate) {
				existingItem.LastUpdate = item.LastUpdate
			}
		} else {
			// Create new grouped item
			groupMap[key] = &entities.GroupedRecyclableItem{
				Name:       item.Name,
				Category:   item.Category,
				Price:      item.Price,
				LastUpdate: item.LastUpdate,
				Hours:      item.Hours,
				URL:        item.URL,
				Shops:      []entities.ShopInfo{shopInfo},
				WasteIDs:   []string{item.WasteID},
			}
		}
	}

	// Filter out shop entries with empty shop_id from each grouped item
	for key, item := range groupMap {
		filteredShops := []entities.ShopInfo{}
		for _, shop := range item.Shops {
			if shop.ShopID != "" {
				filteredShops = append(filteredShops, shop)
			}
		}
		// Update the shops array to only include entries with non-empty shop_id
		groupMap[key].Shops = filteredShops
	}

	// Convert map to slice
	groupedItems := make([]entities.GroupedRecyclableItem, 0, len(groupMap))
	for _, item := range groupMap {
		groupedItems = append(groupedItems, *item)
	}

	// Sort the grouped items by LastUpdate (most recent first)
	sort.Slice(groupedItems, func(i, j int) bool {
		return groupedItems[i].LastUpdate.After(groupedItems[j].LastUpdate)
	})

	return groupedItems
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
