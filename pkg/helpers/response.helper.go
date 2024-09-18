package helpers

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

type BaseResponse struct {
	Status  int          `json:"status"`
	Success bool         `json:"success"`
	Message string       `json:"message"`
	Data    *interface{} `json:"data,omitempty"`
	Errors  *interface{} `json:"errors,omitempty"`
	Meta    *Meta        `json:"meta,omitempty"`
	Log     *Log         `json:"log,omitempty"`
}

type SuccessResponse struct {
	Status  int          `json:"status"`
	Success bool         `json:"success"`
	Message string       `json:"message"`
	Data    *interface{} `json:"data,omitempty"`
	Meta    *Meta        `json:"meta,omitempty"`
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

type Log struct {
	Location  string
	StartTime time.Time
}

func ResponseFormatter(c *fiber.Ctx, res BaseResponse) error {
	// Insert log
	var username string
	if sessionUsername := c.Locals("username"); sessionUsername != nil {
		username = sessionUsername.(string)
	} else {
		username = ""
	}

	LogSystem(LogSystemParam{
		Identifier: c.GetRespHeader(fiber.HeaderXRequestID),
		StatusCode: res.Status,
		Location:   res.Log.Location,
		Message:    res.Message,
		StartTime:  res.Log.StartTime,
		EndTime:    time.Now(),
		Err:        res.Errors,
		Username:   username,
	})

	res.Log = nil
	return c.Status(res.Status).JSON(res)
}
