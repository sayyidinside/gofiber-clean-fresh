package handlers

import (
	"github.com/sayyidinside/gofiber-clean-fresh/interfaces/http/user"
)

type UserManagementHandler struct {
	UserHandler user.UserHandler
}

type Handlers struct {
	UserManagementHandler *UserManagementHandler
}
