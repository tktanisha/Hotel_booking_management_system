package middlewares

import (
	"context"
	"net/http"
	"strings"

	"github.com/tktanisha/booking_system/internal/constants"
	"github.com/tktanisha/booking_system/internal/models"
	"github.com/tktanisha/booking_system/internal/utils"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			utils.WriteErrorResponse(w, http.StatusUnauthorized, "missing or invalid authorization header", "authorization header must be in format 'Bearer <token>'")
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")

		claims, err := utils.ValidateJWT(token)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusUnauthorized, "invalid or expired token", "token is either invalid or has expired")
			return
		}

		userCtx := &models.UserContext{
			Id:   claims.UserID,
			Role: claims.Role,
		}

		ctx := context.WithValue(r.Context(), constants.UserContextKey, userCtx)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
