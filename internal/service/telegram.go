package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"time"
)

type Telegram interface {
	SetWebhook(url string) error
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
