package handler

type UserManagementHandler struct {
	UserHandler       UserHandler
	PermissionHandler PermissionHandler
	ModuleHandler     ModuleHandler
}

type Handlers struct {
	UserManagementHandler *UserManagementHandler
}
