package auth

import (
	"context"
	"log/slog"
	"net/http"
	"strings"

	"github.com/deeerain/nebula-stub/internal/middlewares"
)

type contextKey string

const userIDKey contextKey = "user_id"

// UserIDFromContext returns the Telegram user ID stored by AuthMiddleware.
func UserIDFromContext(ctx context.Context) (int64, bool) {
	id, ok := ctx.Value(userIDKey).(int64)
	return id, ok
}

func AuthMiddleware() middlewares.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authVar := r.Header.Get("Authorization")

			if !strings.HasPrefix(authVar, "Bearer ") {
				slog.Warn("Missing or malformed Authorization header")
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			tokenString := strings.TrimSpace(strings.TrimPrefix(authVar, "Bearer "))
			if tokenString == "" {
				slog.Warn("Empty access token")
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			td, err := ValidateToken(tokenString)
			if err != nil {
				slog.Warn("Invalid token", "error", err)
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			newCtx := context.WithValue(r.Context(), userIDKey, td.ID)
			next.ServeHTTP(w, r.WithContext(newCtx))
		})
	}
}
