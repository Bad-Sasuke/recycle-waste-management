package gateways

import (
	"fmt"
	"io"
	"recycle-waste-management-backend/src/domain/entities"
	"recycle-waste-management-backend/src/middlewares"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func (h *HTTPGateway) CreateShop(ctx *fiber.Ctx) error {
	// Decode JWT token to get user ID
	tokenDetails, err := middlewares.DecodeJWTToken(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(entities.ResponseMessage{Message: "Unauthorized access"})
	}

	name := ctx.FormValue("name")
	description := ctx.FormValue("description")
	address := ctx.FormValue("address")
	phone := ctx.FormValue("phone")
	email := ctx.FormValue("email")
	openingTime := ctx.FormValue("opening_time")
	closingTime := ctx.FormValue("closing_time")
	latitudeStr := ctx.FormValue("latitude")
	longitudeStr := ctx.FormValue("longitude")

	var latitude, longitude float64
	if latitudeStr != "" {
		fmt.Sscanf(latitudeStr, "%f", &latitude)
	}
	if longitudeStr != "" {
		fmt.Sscanf(longitudeStr, "%f", &longitude)
	}

	imageFile, err := ctx.FormFile("image")
	if err != nil {
		fmt.Println("Image file error:", err)
	}

	fileBytes := []byte{}
	if imageFile != nil {
		file, err := imageFile.Open()
		if err != nil {
			fmt.Println("Error opening file:", err)
		} else {
			defer file.Close()
			fileBytes, err = io.ReadAll(file)
			if err != nil {
				fmt.Println("Error reading file:", err)
			}
		}
	}

	shopRequest := entities.CreateShopRequest{
		Name:        name,
		Description: description,
		Address:     address,
		Phone:       phone,
		Email:       email,
		OpeningTime: openingTime,
		ClosingTime: closingTime,
		Latitude:    latitude,
		Longitude:   longitude,
	}

	if err := h.ShopService.CreateShop(tokenDetails.UserID, shopRequest, fileBytes); err != nil {
		return ctx.Status(fiber.StatusForbidden).JSON(entities.ResponseMessage{Message: err.Error()})
	}

	return ctx.Status(fiber.StatusOK).JSON(entities.ResponseModel{Message: "Shop created successfully"})
}

func (h *HTTPGateway) GetShopByShopID(ctx *fiber.Ctx) error {
	shopID := ctx.Params("shop_id")
	if shopID == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(entities.ResponseMessage{Message: "Shop ID is required"})
	}

	shop, err := h.ShopService.GetShopByShopID(shopID)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(entities.ResponseMessage{Message: err.Error()})
	}

	return ctx.Status(fiber.StatusOK).JSON(entities.ResponseModel{Message: "success", Data: shop})
}

func (h *HTTPGateway) GetShopByUserID(ctx *fiber.Ctx) error {
	// Decode JWT token to get user ID
	tokenDetails, err := middlewares.DecodeJWTToken(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(entities.ResponseMessage{Message: "Unauthorized access"})
	}

	shop, err := h.ShopService.GetShopByUserID(tokenDetails.UserID)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(entities.ResponseMessage{Message: err.Error()})
	}

	return ctx.Status(fiber.StatusOK).JSON(entities.ResponseModel{Message: "success", Data: shop})
}

func (h *HTTPGateway) GetAllShops(ctx *fiber.Ctx) error {
	page := 1
	limit := 12

	pageQuery := ctx.Query("page")
	limitQuery := ctx.Query("limit")

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

	shops, totalCount, err := h.ShopService.GetAllShops(page, limit)
	if err != nil {
		return ctx.Status(fiber.StatusForbidden).JSON(entities.ResponseModel{Message: err.Error()})
	}

	totalPages := int((totalCount + int64(limit) - 1) / int64(limit)) // Ceiling division

	response := entities.ResponsePaginationModel{
		Message:    "success",
		Data:       shops,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
		TotalItems: totalCount,
	}

	return ctx.Status(fiber.StatusOK).JSON(response)
}

func (h *HTTPGateway) UpdateShop(ctx *fiber.Ctx) error {
	// Decode JWT token to get user ID
	tokenDetails, err := middlewares.DecodeJWTToken(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(entities.ResponseMessage{Message: "Unauthorized access"})
	}

	shopID := ctx.Params("shop_id")
	if shopID == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(entities.ResponseMessage{Message: "Shop ID is required"})
	}

	// Get existing shop to verify ownership
	existingShop, err := h.ShopService.GetShopByShopID(shopID)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(entities.ResponseMessage{Message: "Shop not found"})
	}

	// Check if the current user owns this shop
	if existingShop.UserID != tokenDetails.UserID {
		return ctx.Status(fiber.StatusForbidden).JSON(entities.ResponseMessage{Message: "Access denied: You don't own this shop"})
	}

	name := ctx.FormValue("name")
	description := ctx.FormValue("description")
	address := ctx.FormValue("address")
	phone := ctx.FormValue("phone")
	email := ctx.FormValue("email")
	openingTime := ctx.FormValue("opening_time")
	closingTime := ctx.FormValue("closing_time")
	latitudeStr := ctx.FormValue("latitude")
	longitudeStr := ctx.FormValue("longitude")

	var updateRequest entities.UpdateShopRequest
	if name != "" {
		updateRequest.Name = &name
	}
	if description != "" {
		updateRequest.Description = &description
	}
	if address != "" {
		updateRequest.Address = &address
	}
	if phone != "" {
		updateRequest.Phone = &phone
	}
	if email != "" {
		updateRequest.Email = &email
	}
	if openingTime != "" {
		updateRequest.OpeningTime = &openingTime
	}
	if closingTime != "" {
		updateRequest.ClosingTime = &closingTime
	}
	if latitudeStr != "" {
		var latitude float64
		fmt.Sscanf(latitudeStr, "%f", &latitude)
		updateRequest.Latitude = &latitude
	}
	if longitudeStr != "" {
		var longitude float64
		fmt.Sscanf(longitudeStr, "%f", &longitude)
		updateRequest.Longitude = &longitude
	}

	imageFile, err := ctx.FormFile("image")
	if err != nil {
		fmt.Println("Image file error:", err)
	}

	fileBytes := []byte{}
	if imageFile != nil {
		file, err := imageFile.Open()
		if err != nil {
			fmt.Println("Error opening file:", err)
		} else {
			defer file.Close()
			fileBytes, err = io.ReadAll(file)
			if err != nil {
				fmt.Println("Error reading file:", err)
			}
		}
	}

	if err := h.ShopService.UpdateShop(shopID, updateRequest, fileBytes); err != nil {
		return ctx.Status(fiber.StatusForbidden).JSON(entities.ResponseMessage{Message: err.Error()})
	}

	return ctx.Status(fiber.StatusOK).JSON(entities.ResponseModel{Message: "Shop updated successfully"})
}

func (h *HTTPGateway) DeleteShop(ctx *fiber.Ctx) error {
	// Decode JWT token to get user ID
	tokenDetails, err := middlewares.DecodeJWTToken(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(entities.ResponseMessage{Message: "Unauthorized access"})
	}

	shopID := ctx.Params("shop_id")
	if shopID == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(entities.ResponseMessage{Message: "Shop ID is required"})
	}

	// Get existing shop to verify ownership
	existingShop, err := h.ShopService.GetShopByShopID(shopID)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(entities.ResponseMessage{Message: "Shop not found"})
	}

	// Check if the current user owns this shop
	if existingShop.UserID != tokenDetails.UserID {
		return ctx.Status(fiber.StatusForbidden).JSON(entities.ResponseMessage{Message: "Access denied: You don't own this shop"})
	}

	if err := h.ShopService.DeleteShop(shopID); err != nil {
		return ctx.Status(fiber.StatusForbidden).JSON(entities.ResponseMessage{Message: err.Error()})
	}

	return ctx.Status(fiber.StatusOK).JSON(entities.ResponseModel{Message: "Shop deleted successfully"})
}
