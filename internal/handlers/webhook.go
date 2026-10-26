package handlers

import (
	"io"
	"log/slog"
	"net/http"

	"github.com/deeerain/nebula-stub/internal/service"
)

const maxWebhookBodySize = 1 << 20

func WebHook(telegramService service.Telegram) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if telegramService == nil {
			http.NotFound(w, r)
			return
		}

		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		body, err := io.ReadAll(io.LimitReader(r.Body, maxWebhookBodySize))
		if err != nil {
			slog.Error("failed to read telegram webhook body", "error", err)
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		slog.Info("telegram webhook received", "body", string(body))
		w.WriteHeader(http.StatusOK)
	}
}
