package model

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type (
	UserDetail struct {
		ID          uint         `json:"id"`
		UUID        uuid.UUID    `json:"uuid"`
		RoleID      uint         `json:"role_id"`
		Role        string       `json:"role"`
		Name        string       `json:"name"`
		Username    string       `json:"username"`
		Email       string       `json:"email"`
		ValidatedAt sql.NullTime `json:"validated_at"`
		CreatedAt   time.Time    `json:"created_at"`
		UpdatedAt   time.Time    `json:"updated_at"`
	}
)
