package payments

import (
	"encoding/base64"
	"fmt"

	"github.com/deeerain/nebula-stub/internal/service"
)

var invoices = []string{}

var planTitles = map[Plan]string{
	PlanSupport: "Поддержать",
	Plan30:      "Подписка на 30 дней",
}

var planDescriptions = map[Plan]string{
	PlanSupport: "Поддержать сервис",
	Plan30:      "Подписка на сервис:\nVPN + PROXY на 30 дней",
}

var planPrices = map[Plan][]service.CreateInvoiceLinkPrice{
	PlanSupport: {
		{
			Label:  "Поддержать",
			Amount: 50,
		},
	},
	Plan30: {
		{
			Label:  "Подписка 30 дней",
			Amount: 30,
		},
	},
}

type TelegramPaymentService struct {
	telegramService service.Telegram
}

func (t *TelegramPaymentService) CreateInvoiceLink(plan Plan, userID int64) (string, error) {
	planBase := base64.StdEncoding.EncodeToString(fmt.Appendf(nil, "%s_%v", plan.String(), userID))
	invoices = append(invoices, planBase)

	link, err := t.telegramService.CreateInvoiceLink(service.CreateInvoiceLink{
		Title:       planTitles[plan],
		Description: planDescriptions[plan],
		Payload:     planBase,
		Currency:    service.Stars.String(),
		Prices:      planPrices[plan],
	})

	if err != nil {
		return "", fmt.Errorf("telegram error: %w", err)
	}

	return link, nil
}

func NewTelegramPaymentProvider(telegram service.Telegram) PaymentService {
	return &TelegramPaymentService{
		telegramService: telegram,
	}
}
