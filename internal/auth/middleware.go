package auth

import (
	"context"
	"log/slog"
	"net/http"
	"strings"

	"github.com/deeerain/nebula-stub/internal/middlewares"
)

func AuthMiddleware() middlewares.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authVar := r.Header.Get("Authorization")

			if authVar == "" || strings.Contains("Bearer", authVar) {
				w.WriteHeader(http.StatusUnauthorized)
				slog.Warn("Auth method is empty", "vars", authVar)
				return
			}

			vars := strings.Split(authVar, " ")
			if len(vars) <= 0 {
				http.Error(w, "", http.StatusUnauthorized)
				slog.Warn("Failed to parse access token")
				return
			}

			td, err := ValidateToken(vars[1])
			if err != nil {
				slog.Warn("Invalid token")
				http.Error(w, "Invalid token", http.StatusUnauthorized)
			}

			newCtx := context.WithValue(r.Context(), "user_id", td.ID)
			next.ServeHTTP(w, r.WithContext(newCtx))
		})
	}
}
