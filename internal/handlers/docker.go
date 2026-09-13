package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/deeerain/nebula-stub/internal/service"
)

func GetContainerStatuses(ctx context.Context, dockerService service.DockerService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if dockerService == nil {
			http.NotFound(w, r)
		}

		containers, err := dockerService.GetContainerList(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusServiceUnavailable)
		}

		json.NewEncoder(w).Encode(containers)
	}
}
