package gateways

import (
	"fmt"
	"io"
	"recycle-waste-management-backend/src/domain/entities"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func (h *HTTPGateway) GetRecycleWaste(ctx *fiber.Ctx) error {
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

	data, totalCount, err := h.RecycleService.GetRecyclableItemsPaginated(page, limit)
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
	name := ctx.FormValue("name")
	price := ctx.FormValue("price")
	category := ctx.FormValue("category")
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
	bodyData := entities.RecyclableItemsModel{
		Name:     name,
		Price:    priceFloat,
		Category: category,
	}
	if err := h.RecycleService.AddRecycleWaste(bodyData, fileBytes); err != nil {
		return ctx.Status(fiber.StatusForbidden).JSON(entities.ResponseMessage{Message: err.Error()})
	}
	return ctx.Status(fiber.StatusOK).JSON(entities.ResponseModel{Message: "success"})
}

func (h *HTTPGateway) DeleteRecycleWaste(ctx *fiber.Ctx) error {
	id := ctx.Params("waste_id")
	if id == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(entities.ResponseMessage{Message: "invalid query params"})
	}
	err := h.RecycleService.DeleteWasteItem(id)
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
	wasteID := ctx.Params("waste_id")
	name := ctx.FormValue("name")
	price := ctx.FormValue("price")
	category := ctx.FormValue("category")
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
	bodyData := entities.RecyclableItemsModel{
		Name:     name,
		Price:    priceFloat,
		Category: category,
	}
	if err := h.RecycleService.EditWasteItem(wasteID, bodyData, fileBytes); err != nil {
		return ctx.Status(fiber.StatusForbidden).JSON(entities.ResponseMessage{Message: err.Error()})
	}
	return ctx.Status(fiber.StatusOK).JSON(entities.ResponseModel{Message: "success"})
}
