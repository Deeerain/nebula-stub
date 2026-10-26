package auth

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/deeerain/nebula-stub/internal"
)

const (
	devInitData = "dev"
	devUserID   = int64(1)
)

func AuthHandler(botToken string, devMode bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req AuthRequestBody

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		var userData TGUserData

		if devMode && req.InitData == devInitData {
			slog.Warn("Dev auth is used", "user_id", devUserID)
			userData.ID = devUserID
		} else {
			data, err := internal.ValidateTelegramInitData(botToken, req.InitData, 24*time.Hour)
			if err != nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(map[string]string{"error": "Invalid Telegram Data"})
				slog.Error("Failed to validate init data", "error", err)
				return
			}

			if err := json.Unmarshal([]byte(data.Get("user")), &userData); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		}

		sessionToken, err := GenerateToken(TokenData{
			ID: userData.ID,
		})

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(AuthResponse{Token: sessionToken}); err != nil {
			slog.Error("failed to write auth response", "error", err)
		}
	}
}
