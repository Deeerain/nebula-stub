package handlers

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/deeerain/nebula-stub/internal"
	"github.com/deeerain/nebula-stub/internal/xui"
)

func GetClinet(ctx context.Context, service *xui.XuiService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authorizationValue := r.Header.Get("Authorization")

		if authorizationValue == "" {
			http.Error(w, "not authorized", http.StatusUnauthorized)
			return
		}

		initData := strings.Split(authorizationValue, " ")[1]

		botToken := os.Getenv("BOT_TOKEN")
		dataRaw, err := internal.ValidateTelegramInitData(botToken, initData, 24*time.Hour)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}

		var userData struct {
			ID int64 `json:"id"`
		}

		slog.Info("user data", "data", dataRaw.Get("user"))

		if err := json.Unmarshal([]byte(dataRaw.Get("user")), &userData); err != nil {
			http.Error(w, err.Error(), 400)
			return
		}

		client, err := service.GetClient(xui.TgID(userData.ID))
		if err != nil {
			slog.Error("failed to get client", "msg", err, "client", client)
			http.Error(w, err.Error(), 500)
			return
		}
		json.NewEncoder(w).Encode(client)
	}
}

func GetLinks(ctx context.Context, service *xui.XuiService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authorizationValue := r.Header.Get("Authorization")

		if authorizationValue == "" {
			http.Error(w, "not authorized", http.StatusUnauthorized)
			return
		}

		initData := strings.Split(authorizationValue, " ")[1]

		botToken := os.Getenv("BOT_TOKEN")
		_, err := internal.ValidateTelegramInitData(botToken, initData, 24*time.Hour)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}

		subId := r.PathValue("id")

		subs, err := service.GetSubscription(xui.SubscriptionID(subId))
		if err != nil {
			http.Error(w, err.Error(), 500)
		}

		json.NewEncoder(w).Encode(subs)
	}
}
