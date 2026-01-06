package constant

type contextKey string

const (
	TABLE_MODULE          string = "modules"
	TABLE_PERMISSION      string = "permissions"
	TABLE_ROLE            string = "roles"
	TABLE_USER            string = "users"
	TABLE_PROFILE         string = "profiles"
	TABLE_REFRESH_TOKEN   string = "refresh_tokens"
	TABLE_ROLE_PERMISSION string = "role_permissions"

	// CONTEXT KEY
	CtxKeyIdentifier contextKey = "identifier"
	CtxKeyUsername   contextKey = "username"
	CtxKeyUserID     contextKey = "user_id"
	CtxKeyIsAdmin    contextKey = "is_admin"
	CtxKeyFunction   contextKey = "function"
)
