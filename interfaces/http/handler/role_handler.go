package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/sayyidinside/gofiber-clean-fresh/domain/service"
	"github.com/sayyidinside/gofiber-clean-fresh/interfaces/model"
	"github.com/sayyidinside/gofiber-clean-fresh/pkg/helpers"
)

type RoleHandler interface {
	GetRole(c *fiber.Ctx) error
	GetAllRole(c *fiber.Ctx) error
	CreateRole(c *fiber.Ctx) error
	UpdateRole(c *fiber.Ctx) error
	DeleteRole(c *fiber.Ctx) error
}

type roleHandler struct {
	service service.RoleService
}

func NewRoleHandler(service service.RoleService) RoleHandler {
	return &roleHandler{
		service: service,
	}
}

func (h *roleHandler) GetRole(c *fiber.Ctx) error {
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

func (h *roleHandler) GetAllRole(c *fiber.Ctx) error {
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

func (h *roleHandler) CreateRole(c *fiber.Ctx) error {
	log := helpers.CreateLog(c)
	var input model.RoleInput

	if err := c.BodyParser(&input); err != nil {
		return helpers.ResponseFormatter(c, helpers.BaseResponse{
			Status:  fiber.StatusBadRequest,
			Success: false,
			Message: "Invalid or malformed request body",
			Log:     &log,
		})
	}

	model.SanitizeRoleInput(&input)

	if err := helpers.ValidateInput(input); err != nil {
		iError := interface{}(err)
		return helpers.ResponseFormatter(c, helpers.BaseResponse{
			Status:  fiber.StatusBadRequest,
			Success: false,
			Message: "Invalid or malformed request body",
			Errors:  &iError,
			Log:     &log,
		})
	}

	response := h.service.Create(c.Context(), &input)
	response.Log = &log

	return helpers.ResponseFormatter(c, response)
}

func (h *roleHandler) UpdateRole(c *fiber.Ctx) error {
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

	var input model.RoleInput

	if err := c.BodyParser(&input); err != nil {
		return helpers.ResponseFormatter(c, helpers.BaseResponse{
			Status:  fiber.StatusBadRequest,
			Success: false,
			Message: "Invalid or malformed request body",
			Log:     &log,
		})
	}

	model.SanitizeRoleInput(&input)

	if err := helpers.ValidateInput(input); err != nil {
		return helpers.ResponseFormatter(c, helpers.BaseResponse{
			Status:  fiber.StatusBadRequest,
			Success: false,
			Message: "Invalid or malformed request body",
			Errors:  err,
			Log:     &log,
		})
	}

	response := h.service.UpdateByID(c.Context(), &input, uint(id))
	response.Log = &log

	return helpers.ResponseFormatter(c, response)
}

func (h *roleHandler) DeleteRole(c *fiber.Ctx) error {
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
