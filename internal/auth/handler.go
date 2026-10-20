package auth

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/deeerain/nebula-stub/internal"
)

func AuthHandler(botToken string) http.HandlerFunc {
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

		expIn := 24 * time.Hour

		var userData TGUserData

		data, err := internal.ValidateTelegramInitData(botToken, req.InitData, expIn)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": "Invalid Telegram Data"})
			slog.Error("Failed to validate init data", "error", err)
			return
		}

		if err := json.Unmarshal([]byte(data.Get("user")), &userData); err != nil {
			http.Error(w, err.Error(), 400)
			return
		}

		sessionToken, err := GenerateToken(TokenData{
			ID: userData.ID,
		})

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(AuthResponse{Token: sessionToken})
	}
}
