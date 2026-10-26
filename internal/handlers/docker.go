package handlers

import (
	"log/slog"
	"net/http"

	"github.com/deeerain/nebula-stub/internal/server"
	"github.com/deeerain/nebula-stub/internal/service"
)

func GetContainerStatuses(dockerService service.DockerService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if dockerService == nil {
			http.Error(w, "Docker service is unavailable", http.StatusServiceUnavailable)
			return
		}

		containers, err := dockerService.GetContainerList(r.Context())
		if err != nil {
			slog.Error("failed to get container list", "error", err)
			http.Error(w, "Docker service is unavailable", http.StatusServiceUnavailable)
			return
		}

		if err := server.WriteJSON(w, containers); err != nil {
			slog.Error("failed to write container list", "error", err)
		}
	}
}
