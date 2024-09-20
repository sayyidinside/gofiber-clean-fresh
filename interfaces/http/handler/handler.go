package handler

type UserManagementHandler struct {
	UserHandler       UserHandler
	PermissionHandler PermissionHandler
}

type Handlers struct {
	UserManagementHandler *UserManagementHandler
}
