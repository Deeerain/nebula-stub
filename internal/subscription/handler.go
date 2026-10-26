package subscription

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"

	"github.com/deeerain/nebula-stub/internal/auth"
	"github.com/deeerain/nebula-stub/internal/client"
)

var deepLinks = map[string]string{
	"incy": "incy://add/%s",
	"happ": "happ://add/%s",
}

type ClientService interface {
	GetClientByTgID(tgID int64) (client.ProxyClient, error)
}

type SubscriptionService interface {
	GetLinksBySubID(string) ([]string, error)
}

func GetLinks(ctx context.Context, service SubscriptionService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		subID := r.PathValue("id")
		if subID == "" {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		links, err := service.GetLinksBySubID(subID)
		if err != nil {
			slog.Error("failed to get subscription links", "error", err, "sub_id", subID)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(links); err != nil {
			slog.Error("failed to write subscription links", "error", err)
		}
	}
}

func ConnectHandler(subscriptionURL url.URL, service ClientService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		app := r.PathValue("app")
		appDeepLink, ok := deepLinks[app]
		if !ok {
			http.Error(w, "Unknown app", http.StatusBadRequest)
			return
		}

		tgID, ok := auth.UserIDFromContext(r.Context())
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		proxyClient, err := service.GetClientByTgID(tgID)
		if err != nil || proxyClient.SubID == "" {
			slog.Warn("failed to get client for deep link", "error", err, "tg_id", tgID)
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		subLink := subscriptionURL.JoinPath(proxyClient.SubID)

		http.Redirect(w, r, fmt.Sprintf(appDeepLink, subLink.String()), http.StatusFound)
	}
}
