package gateways

import (
	"recycle-waste-management-backend/src/domain/entities"
	"recycle-waste-management-backend/src/middlewares"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

func (h *HTTPGateway) GetAllUserData(ctx *fiber.Ctx) error {

	data, err := h.UserService.GetAllUser()
	if err != nil {
		return ctx.Status(fiber.StatusForbidden).JSON(entities.ResponseModel{Message: "cannot get all users data"})
	}
	return ctx.Status(fiber.StatusOK).JSON(entities.ResponseModel{Message: "success", Data: data})
}

func (h *HTTPGateway) CreateNewUserAccount(ctx *fiber.Ctx) error {

	bodyData := new(entities.NewUserBody)
	if err := ctx.BodyParser(&bodyData); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(entities.ResponseMessage{Message: "invalid json body"})
	}

	err := h.UserService.InsertNewAccount(bodyData)
	if err != nil {
		return ctx.Status(fiber.StatusForbidden).JSON(entities.ResponseMessage{Message: "cannot insert new user data"})
	}
	log.Info("new user created successfully")
	return ctx.Status(fiber.StatusCreated).JSON(entities.ResponseModel{Message: "created successfully"})
}

func (h *HTTPGateway) UpdateUserData(ctx *fiber.Ctx) error {

	bodyData := new(entities.NewUserBody)
	if err := ctx.BodyParser(&bodyData); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(entities.ResponseMessage{Message: "invalid json body"})
	}
	params := ctx.Queries()
	userID := params["user_id"]

	if userID == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(entities.ResponseMessage{Message: "invalid query params"})
	}

	err := h.UserService.UpdateUser(userID, bodyData)
	if err != nil {
		return ctx.Status(fiber.StatusForbidden).JSON(entities.ResponseMessage{Message: "cannot update user data"})
	}
	return ctx.Status(fiber.StatusOK).JSON(entities.ResponseModel{Message: "updated successfully"})
}

func (h *HTTPGateway) DeleteUser(ctx *fiber.Ctx) error {
	userID := ctx.Params("user_id")
	if userID == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(entities.ResponseMessage{Message: "invalid query params"})
	}
	err := h.UserService.DeleteUser(userID)
	if err != nil {
		return ctx.Status(fiber.StatusForbidden).JSON(entities.ResponseMessage{Message: "cannot delete user data"})
	}
	return ctx.Status(fiber.StatusOK).JSON(entities.ResponseModel{Message: "deleted successfully"})
}

func (h *HTTPGateway) GetUser(ctx *fiber.Ctx) error {
	userID := ctx.Query("user_id")
	if userID == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(entities.ResponseMessage{Message: "invalid query params"})
	}
	data, err := h.UserService.GetUser(userID)
	if err != nil {
		return ctx.Status(fiber.StatusForbidden).JSON(entities.ResponseMessage{Message: "cannot get user data"})
	}
	return ctx.Status(fiber.StatusOK).JSON(entities.ResponseModel{Message: "success", Data: data})
}

func (h *HTTPGateway) GetCurrentUser(ctx *fiber.Ctx) error {
	tokenDetails, err := middlewares.DecodeJWTToken(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(entities.ResponseMessage{Message: "unauthorized"})
	}

	data, err := h.UserService.GetUser(tokenDetails.UserID)
	if err != nil {
		return ctx.Status(fiber.StatusForbidden).JSON(entities.ResponseMessage{Message: "cannot get user data"})
	}
	return ctx.Status(fiber.StatusOK).JSON(entities.ResponseModel{Message: "success", Data: data})
}

func (h *HTTPGateway) UpdateCurrentUser(ctx *fiber.Ctx) error {
	tokenDetails, err := middlewares.DecodeJWTToken(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(entities.ResponseMessage{Message: "unauthorized"})
	}

	bodyData := new(entities.NewUserBody)
	if err := ctx.BodyParser(&bodyData); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(entities.ResponseMessage{Message: "invalid json body"})
	}

	err = h.UserService.UpdateUser(tokenDetails.UserID, bodyData)
	if err != nil {
		return ctx.Status(fiber.StatusForbidden).JSON(entities.ResponseMessage{Message: "cannot update user data"})
	}
	return ctx.Status(fiber.StatusOK).JSON(entities.ResponseModel{Message: "updated successfully"})
}

func (h *HTTPGateway) UpdateProfileImage(ctx *fiber.Ctx) error {
	tokenDetails, err := middlewares.DecodeJWTToken(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(entities.ResponseMessage{Message: "unauthorized"})
	}

	// Get uploaded file
	file, err := ctx.FormFile("image")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(entities.ResponseMessage{Message: "no image file uploaded"})
	}

	// Process and upload image
	imageURL, err := h.ImageService.ProcessAndUploadProfileImage(file, tokenDetails.UserID)
	if err != nil {
		log.Error("Failed to process and upload image:", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(entities.ResponseMessage{Message: err.Error()})
	}

	// Get current user data to delete old image
	currentUser, err := h.UserService.GetUser(tokenDetails.UserID)
	if err == nil && currentUser != nil && currentUser.ImageURL != "" {
		// Delete old profile image (ignore errors)
		_ = h.ImageService.DeleteProfileImage(currentUser.ImageURL)
	}

	// Update user's image URL in database
	err = h.UserService.UpdateUserImage(tokenDetails.UserID, imageURL)
	if err != nil {
		// If database update fails, try to clean up uploaded image
		_ = h.ImageService.DeleteProfileImage(imageURL)
		return ctx.Status(fiber.StatusInternalServerError).JSON(entities.ResponseMessage{Message: "failed to update user image"})
	}

	return ctx.Status(fiber.StatusOK).JSON(entities.ResponseModel{
		Message: "profile image updated successfully",
		Data:    map[string]string{"image_url": imageURL},
	})
}

func (h *HTTPGateway) UpdateUserRole(ctx *fiber.Ctx) error {
	tokenDetails, err := middlewares.DecodeJWTToken(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(entities.ResponseMessage{Message: "unauthorized"})
	}

	// First, get the current user to check their role
	currentUser, err := h.UserService.GetUser(tokenDetails.UserID)
	if err != nil {
		return ctx.Status(fiber.StatusForbidden).JSON(entities.ResponseMessage{Message: "cannot get current user data"})
	}

	// Only allow admin users to update roles
	if currentUser.Role != string(entities.UserRoleAdmin) {
		return ctx.Status(fiber.StatusForbidden).JSON(entities.ResponseMessage{Message: "access denied - admin role required"})
	}

	// Get the target user ID and new role from the request
	var req struct {
		UserID string `json:"user_id"`
		Role   string `json:"role"`
	}
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(entities.ResponseMessage{Message: "invalid json body"})
	}

	if req.UserID == "" || req.Role == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(entities.ResponseMessage{Message: "user_id and role are required"})
	}

	// Validate the role
	validRoles := map[string]bool{
		string(entities.UserRoleUser):      true,
		string(entities.UserRoleAdmin):     true,
		string(entities.UserRoleModerator): true,
	}

	if !validRoles[req.Role] {
		return ctx.Status(fiber.StatusBadRequest).JSON(entities.ResponseMessage{Message: "invalid role specified"})
	}

	// Update the target user's role
	userRole := entities.UserRole(req.Role)
	err = h.UserService.UpdateUserRole(req.UserID, userRole)
	if err != nil {
		return ctx.Status(fiber.StatusForbidden).JSON(entities.ResponseMessage{Message: "cannot update user role"})
	}

	return ctx.Status(fiber.StatusOK).JSON(entities.ResponseModel{Message: "user role updated successfully", Data: map[string]string{"user_id": req.UserID, "role": req.Role}})
}

func (h *HTTPGateway) GetCurrentUserRole(ctx *fiber.Ctx) error {
	tokenDetails, err := middlewares.DecodeJWTToken(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(entities.ResponseMessage{Message: "unauthorized"})
	}

	currentUser, err := h.UserService.GetUser(tokenDetails.UserID)
	if err != nil {
		return ctx.Status(fiber.StatusForbidden).JSON(entities.ResponseMessage{Message: "cannot get current user data"})
	}

	// Return only the role information
	roleInfo := map[string]interface{}{
		"user_id":  currentUser.UserID,
		"username": currentUser.Username,
		"role":     currentUser.Role,
	}

	return ctx.Status(fiber.StatusOK).JSON(entities.ResponseModel{Message: "success", Data: roleInfo})
}
