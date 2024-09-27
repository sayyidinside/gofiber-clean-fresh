package model

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/sayyidinside/gofiber-clean-fresh/domain/entity"
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

	LogUserInfo struct {
		ID       string `json:"id"`
		Username string `json:"username"`
	}

	UserList struct {
		ID       uint      `json:"id"`
		UUID     uuid.UUID `json:"uuid"`
		Name     string    `json:"name"`
		Username string    `json:"username"`
		Email    string    `json:"email"`
		Role     string    `json:"role"`
	}
)

func UserToDetailModel(user *entity.User) *UserDetail {
	return &UserDetail{
		ID:          user.ID,
		UUID:        user.UUID,
		RoleID:      user.RoleID,
		Role:        user.Role.Name,
		Name:        user.Name,
		Username:    user.Username,
		Email:       user.Email,
		ValidatedAt: user.ValidatedAt,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
	}
}

func UserToModel(user *entity.User) *UserList {
	return &UserList{
		ID:       user.ID,
		UUID:     user.UUID,
		Name:     user.Name,
		Username: user.Username,
		Email:    user.Email,
		Role:     user.Role.Name,
	}
}

func UserToListModel(users *[]entity.User) *[]UserList {
	listUsers := []UserList{}

	for _, user := range *users {
		listUsers = append(listUsers, *UserToModel(&user))
	}

	return &listUsers
}
