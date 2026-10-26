package auth_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/deeerain/nebula-stub/internal/auth"
)

func TestMiddlewareRejectsBadTokensWithoutPanic(t *testing.T) {
	handler := auth.AuthMiddleware()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, ok := auth.UserIDFromContext(r.Context())
		if !ok {
			t.Fatal("user id is missing in context")
		}
		if id != 42 {
			t.Fatalf("unexpected user id: %d", id)
		}
		w.WriteHeader(http.StatusOK)
	}))

	valid, err := auth.GenerateToken(auth.TokenData{ID: 42})
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	expired, err := auth.GenerateToken(auth.TokenData{
		ID:        42,
		ExpiresAt: time.Now().Add(-time.Hour).Unix(),
	})
	if err != nil {
		t.Fatalf("failed to generate expired token: %v", err)
	}

	cases := []struct {
		name   string
		header string
		want   int
	}{
		{"missing header", "", http.StatusUnauthorized},
		{"no bearer prefix", "abcdef", http.StatusUnauthorized},
		{"empty bearer", "Bearer ", http.StatusUnauthorized},
		{"invalid token", "Bearer abc", http.StatusUnauthorized},
		{"expired token", "Bearer " + expired, http.StatusUnauthorized},
		{"valid token", "Bearer " + valid, http.StatusOK},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			if tc.header != "" {
				req.Header.Set("Authorization", tc.header)
			}

			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, req)

			if rec.Code != tc.want {
				t.Fatalf("got status %d, want %d", rec.Code, tc.want)
			}
		})
	}
}
