package middlewares

import (
	"net/http"
	"time"

	"github.com/deeerain/nebula-stub/internal/logger"
)

func LoggingMiddleware(logger logger.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			next.ServeHTTP(w, r)

			logger.Info(
				"Request", "ip", r.RemoteAddr,
				"method", r.Method,
				"path", r.URL.Path,
				"duration",
				time.Since(start).String(),
			)
		})
	}
}
