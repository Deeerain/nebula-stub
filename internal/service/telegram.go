package service

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
)

type Telegram interface {
	SetWebhook(url string) error
}

type telegramServie struct {
	baseURL    *url.URL
	httpClient *http.Client
}

func (t *telegramServie) SetWebhook(url string) error {
	u := t.baseURL.JoinPath("setWebhook")
	u.Query().Add("url", url)

	resp, err := t.request(http.MethodGet, u.String(), nil)
	if err != nil {
		return fmt.Errorf("failed to set webhook: %w", err)
	}

	var data struct {
		Ok          bool   `json:"ok"`
		Result      bool   `json:"result"`
		Description string `json:"description"`
	}

	if err := json.NewDecoder(resp.Body).Decode(data); err != nil {
		return fmt.Errorf("failed to decode object: %w", err)
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

func NewTelegramService(token string) (Telegram, error) {
	baseUrl, err := url.Parse(fmt.Sprintf("https://api.telegram.org/bot%s", token))
	if err != nil {
		fmt.Errorf("Failed to parse base url", err)
	}

	return &telegramServie{
		httpClient: http.DefaultClient,
		baseURL:    baseUrl,
	}, nil
}
