package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/sayyidinside/gofiber-clean-fresh/domain/service"
	"github.com/sayyidinside/gofiber-clean-fresh/interfaces/model"
	"github.com/sayyidinside/gofiber-clean-fresh/pkg/helpers"
)

type PermissionHandler struct {
	service service.PermissionService
}

func NewPermissionHandler(service service.PermissionService) *PermissionHandler {
	return &PermissionHandler{
		service: service,
	}
}

func (h *PermissionHandler) GetPermission(c *fiber.Ctx) error {
	log := helpers.CreateLog(c)

	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return helpers.ResponseFormatter(c, helpers.BaseResponse{
			Status:  fiber.StatusBadRequest,
			Success: false,
			Message: "Invalid ID format",
			Log:     &log,
		})
	}

	response := h.service.GetByID(c.Context(), uint(id))
	response.Log = &log

	return helpers.ResponseFormatter(c, response)
}

func (h *PermissionHandler) GetAllPermission(c *fiber.Ctx) error {
	log := helpers.CreateLog(c)
	query := new(model.QueryGet)

	if err := c.QueryParser(query); err != nil {
		return helpers.ResponseFormatter(c, helpers.BaseResponse{
			Status:  fiber.StatusBadRequest,
			Success: false,
			Message: "Invalid or malformed request query",
			Log:     &log,
		})
	}

	model.SanitizeQueryGet(query)

	url := c.BaseURL() + c.OriginalURL()
	response := h.service.GetAll(c.Context(), query, url)
	response.Log = &log

	return helpers.ResponseFormatter(c, response)
}

func (h *PermissionHandler) CreatePermission(c *fiber.Ctx) error {
	log := helpers.CreateLog(c)
	var input model.PermissionInput

	if err := c.BodyParser(&input); err != nil {
		return helpers.ResponseFormatter(c, helpers.BaseResponse{
			Status:  fiber.StatusBadRequest,
			Success: false,
			Message: "Invalid or malformed request body",
			Log:     &log,
		})
	}

	model.SanitizePermissionInput(&input)

	if err := helpers.ValidateInput(input); err != nil {
		errorData := interface{}(err)

		return helpers.ResponseFormatter(c, helpers.BaseResponse{
			Status:  fiber.StatusBadRequest,
			Success: false,
			Message: "Invalid or malformed request body",
			Errors:  &errorData,
		})
	}

	response := h.service.Create(c.Context(), &input)
	response.Log = &log

	return helpers.ResponseFormatter(c, response)
}

func (h *PermissionHandler) UpdatePermission(c *fiber.Ctx) error {
	log := helpers.CreateLog(c)
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return helpers.ResponseFormatter(c, helpers.BaseResponse{
			Status:  fiber.StatusBadRequest,
			Success: false,
			Message: "Invalid ID format",
			Log:     &log,
		})
	}

	var input model.PermissionInput

	if err := c.BodyParser(&input); err != nil {
		return err
	}

	model.SanitizePermissionInput(&input)

	if err := helpers.ValidateInput(input); err != nil {
		errorData := interface{}(err)

		return helpers.ResponseFormatter(c, helpers.BaseResponse{
			Status:  fiber.StatusBadRequest,
			Success: false,
			Message: "Invalid or malformed request body",
			Errors:  &errorData,
		})
	}

	response := h.service.UpdateByID(c.Context(), &input, uint(id))
	response.Log = &log

	return helpers.ResponseFormatter(c, response)
}

func (h *PermissionHandler) DeletePermission(c *fiber.Ctx) error {
	log := helpers.CreateLog(c)
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return helpers.ResponseFormatter(c, helpers.BaseResponse{
			Status:  fiber.StatusBadRequest,
			Success: false,
			Message: "Invalid ID format",
			Log:     &log,
		})
	}

	response := h.service.DeleteByID(c.Context(), uint(id))
	response.Log = &log

	return helpers.ResponseFormatter(c, response)
}
