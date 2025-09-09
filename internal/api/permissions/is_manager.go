package permissions

import (
	user_role "github.com/tktanisha/booking_system/internal/enums/user"
	"github.com/tktanisha/booking_system/internal/models"
)

func IsManager(userCtx *models.UserContext) bool {
	return userCtx != nil && userCtx.Role == user_role.RoleManager
}
