package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/sayyidinside/gofiber-clean-fresh/domain/service"
	"github.com/sayyidinside/gofiber-clean-fresh/interfaces/model"
	"github.com/sayyidinside/gofiber-clean-fresh/pkg/helpers"
)

type UserHandler interface {
	GetUser(c *fiber.Ctx) error
	// GetUser(c *fiber.Ctx, ctx context.Context, log *helpers.Log) error
	GetAllUser(c *fiber.Ctx) error
	CreateUser(c *fiber.Ctx) error
	UpdateUser(c *fiber.Ctx) error
	ResetPassword(c *fiber.Ctx) error
	DeleteUser(c *fiber.Ctx) error
}

type userHandler struct {
	service service.UserService
}

func NewUserHandler(service service.UserService) UserHandler {
	return &userHandler{
		service: service,
	}
}

func (h *userHandler) GetUser(c *fiber.Ctx) error {
	ctx := helpers.ExtractIdentifierAndUsername(c)
	logData := helpers.CreateLog(h)

	defer helpers.LogSystemWithDefer(ctx, &logData)

	// Simulate delay
	// time.Sleep(100 * time.Millisecond)

	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		logData.Message = "Invalid ID Format"
		logData.Err = err
		// helpers.CreateLogSystem23(ctx, logData)
		return helpers.ResponseFormatter(c, helpers.BaseResponse{
			Status:  fiber.StatusBadRequest,
			Success: false,
			Message: logData.Message,
			Log:     &logData,
		})
	}

	response := h.service.GetByID(ctx, uint(id))
	response.Log = &logData

	// helpers.CreateLogSystem23(ctx, logData)

	return helpers.ResponseFormatter(c, response)
}

func (h *userHandler) GetAllUser(c *fiber.Ctx) error {
	logData := helpers.CreateLog(c)
	query := new(model.QueryGet)

	if err := c.QueryParser(query); err != nil {
		return helpers.ResponseFormatter(c, helpers.BaseResponse{
			Status:  fiber.StatusBadRequest,
			Success: false,
			Message: "Invalid or malformed request query",
			Log:     &logData,
		})
	}

	model.SanitizeQueryGet(query)

	url := c.BaseURL() + c.OriginalURL()
	response := h.service.GetAll(c.Context(), query, url)
	response.Log = &logData

	return helpers.ResponseFormatter(c, response)
}

func (h *userHandler) CreateUser(c *fiber.Ctx) error {
	logData := helpers.CreateLog(c)
	var input model.UserInput

	if err := c.BodyParser(&input); err != nil {
		return helpers.ResponseFormatter(c, helpers.BaseResponse{
			Status:  fiber.StatusBadRequest,
			Success: false,
			Message: "Invalid or malformed request body",
			Log:     &logData,
		})
	}

	model.SanitizeUserInput(&input)

	if err := helpers.ValidateInput(input); err != nil {
		return helpers.ResponseFormatter(c, helpers.BaseResponse{
			Status:  fiber.StatusBadRequest,
			Success: false,
			Message: "Invalid or malformed request body",
			Errors:  err,
			Log:     &logData,
		})
	}

	response := h.service.Create(c.Context(), &input)
	response.Log = &logData

	return helpers.ResponseFormatter(c, response)
}

func (h *userHandler) UpdateUser(c *fiber.Ctx) error {
	logData := helpers.CreateLog(c)

	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return helpers.ResponseFormatter(c, helpers.BaseResponse{
			Status:  fiber.StatusBadRequest,
			Success: false,
			Message: "Invalid ID format",
			Log:     &logData,
		})
	}

	var input model.UserUpdateInput
	if err := c.BodyParser(&input); err != nil {
		return helpers.ResponseFormatter(c, helpers.BaseResponse{
			Status:  fiber.StatusBadRequest,
			Success: false,
			Message: "Invalid or malformed request body",
			Log:     &logData,
		})
	}

	model.SanitizeUserUpdateInput(&input)

	if err := helpers.ValidateInput(input); err != nil {
		return helpers.ResponseFormatter(c, helpers.BaseResponse{
			Status:  fiber.StatusBadRequest,
			Success: false,
			Message: "Invalid or malformed request body",
			Errors:  err,
			Log:     &logData,
		})
	}

	response := h.service.UpdateByID(c.Context(), &input, uint(id))
	response.Log = &logData

	return helpers.ResponseFormatter(c, response)
}

func (h *userHandler) ResetPassword(c *fiber.Ctx) error {
	logData := helpers.CreateLog(c)

	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return helpers.ResponseFormatter(c, helpers.BaseResponse{
			Status:  fiber.StatusBadRequest,
			Success: false,
			Message: "Invalid ID format",
			Log:     &logData,
		})
	}

	var input model.ChangePasswordInput
	if err := c.BodyParser(&input); err != nil {
		return helpers.ResponseFormatter(c, helpers.BaseResponse{
			Status:  fiber.StatusBadRequest,
			Success: false,
			Message: "Invalid or malformed request body",
			Log:     &logData,
		})
	}

	model.SanitizeChangePasswordInput(&input)

	if err := helpers.ValidateInput(input); err != nil {
		return helpers.ResponseFormatter(c, helpers.BaseResponse{
			Status:  fiber.StatusBadRequest,
			Success: false,
			Message: "Invalid or malformed request body",
			Errors:  err,
			Log:     &logData,
		})
	}

	response := h.service.ChangePassByID(c.Context(), &input, uint(id))
	response.Log = &logData

	return helpers.ResponseFormatter(c, response)
}

func (h *userHandler) DeleteUser(c *fiber.Ctx) error {
	logData := helpers.CreateLog(c)

	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return helpers.ResponseFormatter(c, helpers.BaseResponse{
			Status:  fiber.StatusBadRequest,
			Success: false,
			Message: "Invalid ID format",
		})
	}

	response := h.service.DeleteByID(c.Context(), uint(id))
	response.Log = &logData

	return helpers.ResponseFormatter(c, response)
}
