package payments

import (
	"os"
	"testing"

	"github.com/deeerain/nebula-stub/internal/service"
)

func TestCreateInvoiceLink(t *testing.T) {
	token := os.Getenv("BOT_TOKEN")
	service, err := service.NewTelegramService(token)
	if err != nil {
		t.Fatal("failed to init service", err.Error())
	}

	payments := NewTelegramPaymentProvider(service)
	link, err := payments.CreateInvoiceLink(PlanSupport, 0)

	if err != nil {
		t.Fatal("failed to create invoice link", err.Error())
	}

	if link == "" {
		t.Fatal("link is empty")
	}
}
