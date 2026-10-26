package client_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/deeerain/nebula-stub/internal/auth"
	"github.com/deeerain/nebula-stub/internal/client"
)

type stubService struct {
	proxyClient client.ProxyClient
}

func (s stubService) GetClientByTgID(id int64) (client.ProxyClient, error) {
	return s.proxyClient, nil
}

func TestGetClientReturnsExpectedContract(t *testing.T) {
	svc := stubService{proxyClient: client.ProxyClient{
		TgID:        7,
		SubID:       "sub-1",
		Email:       "user7",
		Enable:      true,
		LimitIP:     2,
		LimitHwid:   3,
		ExpiryTime:  1793094589000,
		UsedTraffic: 1024,
	}}

	token, err := auth.GenerateToken(auth.TokenData{ID: 7})
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	handler := auth.AuthMiddleware()(client.GetClient(context.Background(), svc))
	req := httptest.NewRequest(http.MethodGet, "/api/client", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("got status %d, want %d", rec.Code, http.StatusOK)
	}

	var body map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	for key, want := range map[string]any{
		"tg_id":        float64(7),
		"sub_id":       "sub-1",
		"email":        "user7",
		"enable":       true,
		"limit_ip":     float64(2),
		"limit_hwid":   float64(3),
		"expiry_time":  float64(1793094589000),
		"used_traffic": float64(1024),
	} {
		if got := body[key]; got != want {
			t.Fatalf("field %q: got %v, want %v", key, got, want)
		}
	}
}

func TestGetClientReturnsNotFoundWithoutSubscription(t *testing.T) {
	token, err := auth.GenerateToken(auth.TokenData{ID: 7})
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	handler := auth.AuthMiddleware()(client.GetClient(context.Background(), stubService{}))
	req := httptest.NewRequest(http.MethodGet, "/api/client", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("got status %d, want %d", rec.Code, http.StatusNotFound)
	}
}
