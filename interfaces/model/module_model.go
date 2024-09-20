package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/sayyidinside/gofiber-clean-fresh/domain/entity"
)

type (
	ModuleDetail struct {
		ID          uint                `json:"id"`
		UUID        uuid.UUID           `json:"uuid"`
		Name        string              `json:"name"`
		Permissions []entity.Permission `json:"permissions"`
		CreatedAt   time.Time           `json:"created_at"`
		UpdatedAt   time.Time           `json:"updated_at"`
	}
)
