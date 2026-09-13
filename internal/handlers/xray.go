package handlers

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/deeerain/nebula-stub/internal"
	"github.com/deeerain/nebula-stub/internal/service"
)

func GetClinet(ctx context.Context, xrayService service.XrayService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authorizationValue := r.Header.Get("Authorization")
		initData := strings.Split(authorizationValue, " ")[1]

		botToken := os.Getenv("BOT_TOKEN")
		data, err := internal.ValidateTelegramInitData(botToken, initData, 24*time.Hour)
		if err != nil {
			http.Error(w, err.Error(), 400)
		}

		slog.Info(data.Get("user"))

		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
		}
	}
}
