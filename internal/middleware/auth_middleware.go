package middleware

import (
	"context"
	"net/http"
	"olx_api/internal/httpx"
	"olx_api/internal/utils"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		parts := strings.Fields(authHeader)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			httpx.Error(w, "Access denied", httpx.ErrUnauthorized, http.StatusUnauthorized)
			return
		}
		tokenString := parts[1]
		token, err := utils.VerifyToken(tokenString, os.Getenv("JWT_SECRET"))

		if err != nil || token == nil || !token.Valid {
			httpx.Error(w, "Invalid or expired token", httpx.ErrUnauthorized, http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || claims == nil {
			httpx.Error(w, "Invalid token claims", httpx.ErrUnauthorized, http.StatusUnauthorized)
			return
		}
		userID, ok := claims["user_id"].(float64)
		if !ok || userID < 1 || userID != float64(uint(userID)) {
			http.Error(w, "Invalid user ID", http.StatusUnauthorized)
			return
		}
		
		ctx := context.WithValue(r.Context(), "user_id", uint(userID))
		next.ServeHTTP(w,  r.WithContext(ctx))

	})

}
