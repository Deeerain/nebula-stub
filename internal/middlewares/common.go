package middlewares

import "net/http"

func CommonMiddleware() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r.Header.Set("Content-type", "application/json")
			next.ServeHTTP(w, r)
		})
	}
}
