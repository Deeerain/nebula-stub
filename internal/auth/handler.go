package auth

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/deeerain/nebula-stub/internal"
	"github.com/deeerain/nebula-stub/internal/server"
)

const (
	devInitData = "dev"
	devUserID   = int64(1)
)

func AuthHandler(botToken string, devMode bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			server.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}

		var req AuthRequestBody

		if err := server.ReadJSON(r, &req); err != nil {
			server.WriteError(w, http.StatusBadRequest, "bad request")
			return
		}

		var userData TGUserData

		if devMode && req.InitData == devInitData {
			slog.Warn("dev auth is used", "user_id", devUserID)
			userData.ID = devUserID
		} else {
			data, err := internal.ValidateTelegramInitData(botToken, req.InitData, 24*time.Hour)
			if err != nil {
				slog.Error("failed to validate init data", "error", err)
				server.WriteError(w, http.StatusUnauthorized, "invalid telegram data")
				return
			}

			if err := json.Unmarshal([]byte(data.Get("user")), &userData); err != nil {
				slog.Error("failed to parse telegram user data", "error", err)
				server.WriteError(w, http.StatusBadRequest, "invalid user data")
				return
			}
		}

		sessionToken, err := GenerateToken(TokenData{
			ID: userData.ID,
		})
		if err != nil {
			slog.Error("failed to generate token", "error", err)
			server.WriteError(w, http.StatusInternalServerError, "internal server error")
			return
		}

		if err := server.WriteJSON(w, AuthResponse{Token: sessionToken}); err != nil {
			slog.Error("failed to write auth response", "error", err)
		}
	}
}
