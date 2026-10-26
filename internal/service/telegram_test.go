package service

import (
	"os"
	"testing"
)

func TestNewTelegramServiceRequiresToken(t *testing.T) {
	if _, err := NewTelegramService(""); err == nil {
		t.Fatal("NewTelegramService(\"\") returned nil error")
	}
}

func TestCreateInvoiceLink(t *testing.T) {
	token := os.Getenv("BOT_TOKEN")
	t.Log(token)
	if token == "" {
		t.Fatal("bot token not set")
	}

	ts, err := NewTelegramService(token)
	if err != nil {
		t.Fatal("failed to init telegram service")
	}

	link, err := ts.CreateInvoiceLink(CreateInvoiceLink{
		Title:       "Test invoice",
		Description: "test invoice",
		Payload:     "test_invoice",
		Currency:    Stars.String(),
		Prices: []CreateInvoiceLinkPrice{
			{
				Label:  "Test invoice",
				Amount: 20,
			},
		},
	})
	if err != nil {
		t.Fatal("failed to create link", err.Error())
	}

	if link == "" {
		t.Fail()
	}

	t.Log(link)
}
