package handlers

import (
	"log/slog"
	"net/http"

	"github.com/deeerain/nebula-stub/internal/service"
)

func WebHook(telegramService service.Telegram) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if telegramService == nil {
			http.NotFound(w, r)
			return
		}

		slog.Info("New webhook request", "req", r.Body)
	}
}
