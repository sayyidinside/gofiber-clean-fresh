package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/sayyidinside/gofiber-clean-fresh/domain/service"
	"github.com/sayyidinside/gofiber-clean-fresh/interfaces/model"
	"github.com/sayyidinside/gofiber-clean-fresh/pkg/helpers"
)

type ModuleHandler interface {
	GetModule(c *fiber.Ctx) error
	GetAllModule(c *fiber.Ctx) error
	CreateModule(c *fiber.Ctx) error
	UpdateModule(c *fiber.Ctx) error
	DeleteModule(c *fiber.Ctx) error
}

type moduleHandler struct {
	service service.ModuleService
}

func NewModuleHandler(service service.ModuleService) ModuleHandler {
	return &moduleHandler{
		service: service,
	}
}

func (h *moduleHandler) GetModule(c *fiber.Ctx) error {
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

func (h *moduleHandler) GetAllModule(c *fiber.Ctx) error {
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

func (h *moduleHandler) CreateModule(c *fiber.Ctx) error {
	log := helpers.CreateLog(c)
	var input model.ModuleInput

	if err := c.BodyParser(&input); err != nil {
		return helpers.ResponseFormatter(c, helpers.BaseResponse{
			Status:  fiber.StatusBadRequest,
			Success: false,
			Message: "Invalid or malformed request body",
			Log:     &log,
		})
	}

	model.SanitizeModuleInput(&input)

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

func (h *moduleHandler) UpdateModule(c *fiber.Ctx) error {
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

	var input model.ModuleInput

	if err := c.BodyParser(&input); err != nil {
		return helpers.ResponseFormatter(c, helpers.BaseResponse{
			Status:  fiber.StatusBadRequest,
			Success: false,
			Message: "Invalid or malformed request body",
			Log:     &log,
		})
	}

	model.SanitizeModuleInput(&input)

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

	response := h.service.UpdateByID(c.Context(), &input, uint(id))
	response.Log = &log

	return helpers.ResponseFormatter(c, response)
}

func (h *moduleHandler) DeleteModule(c *fiber.Ctx) error {
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
