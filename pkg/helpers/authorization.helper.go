package helpers

import (
	"context"

	"github.com/sayyidinside/gofiber-clean-fresh/pkg/utils/constant"
)

func SelfOrAdminOnly(ctx context.Context, user_id uint) bool {
	session_user_id := ctx.Value(constant.CtxKeyUserID).(float64)
	is_admin := ctx.Value(constant.CtxKeyIsAdmin).(bool)
	if user_id != uint(session_user_id) && !is_admin {
		return false
	}

	return true
}
