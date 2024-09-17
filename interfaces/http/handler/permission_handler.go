package handler

import (
	"encoding/json"
	"log"
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
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID format"})
	}

	response := h.service.GetByID(c.Context(), uint(id))

	return helpers.ResponseFormatter(c, response)
}

func (h *PermissionHandler) GetAllPermission(c *fiber.Ctx) error {
	query := model.QueryGet{}

	if err := c.QueryParser(&query); err != nil {
		return err
	}

	jsonData, _ := json.MarshalIndent(query, "", "  ")
	log.Println(string(jsonData))

	response := h.service.GetAll(c.Context(), &query)

	return helpers.ResponseFormatter(c, response)
}
