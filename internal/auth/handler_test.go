package auth_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/deeerain/nebula-stub/internal/auth"
)

func TestAuthHandlerDevMode(t *testing.T) {
	authenticate := func(t *testing.T, devMode bool) *httptest.ResponseRecorder {
		t.Helper()

		handler := auth.AuthHandler("", devMode)
		req := httptest.NewRequest(http.MethodPost, "/api/auth", strings.NewReader(`{"initData":"dev"}`))
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		return rec
	}

	t.Run("enabled", func(t *testing.T) {
		rec := authenticate(t, true)

		if rec.Code != http.StatusOK {
			t.Fatalf("got status %d, want %d", rec.Code, http.StatusOK)
		}

		var resp auth.AuthResponse
		if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		if resp.Token == "" {
			t.Fatal("token is empty")
		}

		claims, err := auth.ValidateToken(resp.Token)
		if err != nil {
			t.Fatalf("failed to validate token: %v", err)
		}

		if claims.ID != 1 {
			t.Fatalf("got user id %d, want 1", claims.ID)
		}
	})

	t.Run("disabled", func(t *testing.T) {
		rec := authenticate(t, false)

		if rec.Code != http.StatusUnauthorized {
			t.Fatalf("got status %d, want %d", rec.Code, http.StatusUnauthorized)
		}
	})
}
