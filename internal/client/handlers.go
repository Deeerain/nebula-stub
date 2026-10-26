package client

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/deeerain/nebula-stub/internal/auth"
	"github.com/deeerain/nebula-stub/internal/server"
)

type ClientService interface {
	GetClientByTgID(id int64) (ProxyClient, error)
}

func GetClient(ctx context.Context, service ClientService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userTgID, ok := auth.UserIDFromContext(r.Context())
		if !ok {
			server.WriteError(w, http.StatusUnauthorized, "unauthorized")
			return
		}

		proxyClient, err := service.GetClientByTgID(userTgID)
		if err != nil {
			server.WriteError(w, http.StatusBadRequest, "failed to get proxy client")
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		if proxyClient.SubID == "" {
			server.WriteError(w, http.StatusInternalServerError, "sub id is null")
			http.Error(w, "Client not found", http.StatusNotFound)
			return
		}

		if err := server.WriteJSON(w, proxyClient); err != nil {
			slog.Error("failed to write client response", "error", err)
		}
	}
}
