package helpers

import (
	"github.com/gofiber/fiber/v2"
)

type BaseResponse struct {
	Status  int          `json:"status"`
	Success bool         `json:"success"`
	Message string       `json:"message"`
	Data    *interface{} `json:"data,omitempty"`
	Errors  *interface{} `json:"errors,omitempty"`
	Meta    *Meta        `json:"meta,omitempty"`
}

type SuccessResponse struct {
	Status  int          `json:"status"`
	Success bool         `json:"success"`
	Message string       `json:"message"`
	Data    *interface{} `json:"data,omitempty"`
	Meta    *Meta        `json:"meta"`
}

type Meta struct {
	RequestID  string      `json:"request_id"`
	Pagination *Pagination `json:"pagination,omitempty"`
}

type Pagination struct {
	CurrentPage int     `json:"current_page"`
	TotalItems  int     `json:"total_items"`
	TotalPages  int     `json:"total_pages"`
	ItemPerPage int     `json:"item_per_page"`
	FromRow     int     `json:"from_row"`
	ToRow       int     `json:"to_row"`
	Self        string  `json:"self"`
	Next        *string `json:"next"`
	Prev        *string `json:"prev"`
}

type ErrorResponse struct {
	Status  int          `json:"status"`
	Success bool         `json:"success"`
	Message string       `json:"message"`
	Errors  *interface{} `json:"errors,omitempty"`
}

func ResponseFormatter(c *fiber.Ctx, res BaseResponse) error {
	return c.Status(res.Status).JSON(res)
}
