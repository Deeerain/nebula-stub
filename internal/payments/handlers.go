package payments

import (
	"log/slog"
	"net/http"

	"github.com/deeerain/nebula-stub/internal/auth"
	"github.com/deeerain/nebula-stub/internal/server"
)

type PaymentService interface {
	CreateInvoiceLink() (string, error)
}

func InvoiceHandler(paymentService PaymentService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := auth.UserIDFromContext(r.Context())
		if !ok {
			server.WriteError(w, http.StatusUnauthorized, "unauthorized")
			return
		}

		var invoiceRequest InvoiceRequest
		if err := server.ReadJSON(r, &invoiceRequest); err != nil {
			server.WriteError(w, http.StatusBadRequest, "failed to parse request")
			return
		}

		link, _ := paymentService.CreateInvoiceLink()
		slog.Debug("Created invoice link", "user", userID, "link", link)

		if err := server.WriteJSON(w, link); err != nil {
			slog.Error("failed to write response", "error", err.Error())
		}
	}
}
