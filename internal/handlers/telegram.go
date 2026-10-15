package handlers

import (
	"net/http"

	"github.com/deeerain/nebula-stub/internal/model"
)

type GetSubscriptionUseCase interface {
	Execute(telegramID int64) ([]string, error)
}

func GetSubscriptionHandler(useCase GetSubscriptionUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		subscriptions, err := useCase.Execute(0)
		if err != nil {
			http.Error(w, "", http.StatusBadRequest)
		}

		model.NewResponse(true, subscriptions, nil).
			Write(w)
	}
}
