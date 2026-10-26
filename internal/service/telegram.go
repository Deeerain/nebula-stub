package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"time"
)

type Curency int

const (
	Stars Curency = iota
)

var currencyName = map[Curency]string{
	Stars: "XTR",
}

func (c Curency) String() string {
	return currencyName[c]
}

type CreateInvoiceLink struct {
	Title         string                   `json:"title"`
	Description   string                   `json:"description"`
	Payload       string                   `json:"payload"`
	ProviderToken string                   `json:"provider_token"`
	Currency      string                   `json:"currency"`
	Prices        []CreateInvoiceLinkPrice `json:"prices"`
}

type CreateInvoiceLinkPrice struct {
	Label  string `json:"label"`
	Amount int64  `json:"amount"`
}

type Telegram interface {
	SetWebhook(url string) error
	CreateInvoiceLink(data CreateInvoiceLink) (string, error)
}

type telegramServie struct {
	baseURL    *url.URL
	httpClient *http.Client
}

func (t *telegramServie) SetWebhook(webhookURL string) error {
	u := t.baseURL.JoinPath("setWebhook")
	query := u.Query()
	query.Set("url", webhookURL)
	u.RawQuery = query.Encode()

	resp, err := t.request(http.MethodGet, u.String(), nil)
	if err != nil {
		return fmt.Errorf("failed to set webhook: %w", err)
	}
	defer resp.Body.Close()

	var data struct {
		Ok          bool   `json:"ok"`
		Result      bool   `json:"result"`
		Description string `json:"description"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return fmt.Errorf("failed to decode object: %w", err)
	}

	if !data.Ok {
		return fmt.Errorf("telegram error: %s", data.Description)
	}

	slog.Debug("Webhook result", "data", data)

	return nil
}

func (t *telegramServie) request(httpMethod string, url string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(httpMethod, url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	return t.httpClient.Do(req)
}

func (t *telegramServie) CreateInvoiceLink(data CreateInvoiceLink) (string, error) {
	u := t.baseURL.JoinPath("createInvoiceLink")

	payload, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("failed to marshal invoice data: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, u.String(), bytes.NewReader(payload))
	if err != nil {
		return "", fmt.Errorf("failed to create invoice request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := t.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to create invoice link: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("failed to create invoice link: %s", string(body))
	}

	var respData struct {
		Ok          bool   `json:"ok"`
		Result      string `json:"result"`
		Description string `json:"description"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&respData); err != nil {
		return "", fmt.Errorf("failed to decode invoice link response: %w", err)
	}
	if !respData.Ok {
		return "", fmt.Errorf("telegram error: %s", respData.Description)
	}

	return respData.Result, nil
}

func NewTelegramService(token string) (Telegram, error) {
	if token == "" {
		return nil, errors.New("telegram bot token is empty")
	}

	baseURL, err := url.Parse(fmt.Sprintf("https://api.telegram.org/bot%s", token))
	if err != nil {
		return nil, fmt.Errorf("parse base url: %w", err)
	}

	return &telegramServie{
		httpClient: &http.Client{Timeout: 15 * time.Second},
		baseURL:    baseURL,
	}, nil
}
