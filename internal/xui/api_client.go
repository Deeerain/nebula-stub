package xui

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"time"

	"github.com/deeerain/nebula-stub/internal/xui/types"
)

type Client struct {
	ID         int    `json:"id"`
	Email      string `json:"email"`
	SubID      string `json:"subId"`
	UUID       string `json:"uuid"`
	TgID       int64  `json:"tgId"`
	Enable     bool   `json:"enable"`
	Reset      int    `json:"reset"`
	ResetDays  int    `json:"resetDay"`
	LimitIP    int    `json:"limitIP"`
	LimitHwid  int    `json:"limitHwid"`
	ExpiryTime int    `json:"expiryTime"`
}

type ClientData struct {
	Client      Client `json:"client"`
	Inbounds    []int  `json:"inboundIds"`
	UsedTraffic int    `json:"usedTraffic"`
}

type ApiClient struct {
	baseURL    *url.URL
	httpClient *http.Client
	apiToken   string
}

func (s *ApiClient) GetClinetByTgID(telegramID int64) ([]ClientData, error) {
	url := s.baseURL.JoinPath("clients/get/tgId/").JoinPath(fmt.Sprint(telegramID))

	resp, err := s.get(url.String())
	if err != nil {
		return nil, fmt.Errorf("failed to execute: %w", err)
	}
	defer resp.Body.Close()

	var data types.ApiResponse[[]ClientData]

	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, fmt.Errorf("failed to decode response data: %w", err)
	}

	slog.Debug("Response", "data", data, "status", resp.Status)

	return data.Obj, nil
}

func (s *ApiClient) GetLinksByEmail(email string) ([]string, error) {
	url := s.baseURL.JoinPath("clients/get/links").JoinPath(email)

	resp, err := s.get(url.String())
	if err != nil {
		return nil, fmt.Errorf("failed to get response: %w", err)
	}
	defer resp.Body.Close()

	var data types.ApiResponse[[]string]

	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, fmt.Errorf("failed to decode object: %w", err)
	}

	return data.Obj, nil
}

func (s *ApiClient) GetSubLinks(subId string) ([]string, error) {
	url := s.baseURL.JoinPath("clients/subLinks").JoinPath(subId)

	resp, err := s.get(url.String())
	if err != nil {
		return nil, fmt.Errorf("failed to get response: %w", err)
	}
	defer resp.Body.Close()

	data, err := types.DecodeResponse[[]string](resp.Body)

	return data.Obj, err
}

func (s *ApiClient) GetLinks(email string) ([]string, error) {
	url := s.baseURL.JoinPath("clients/links").JoinPath(email)

	resp, err := s.get(url.String())
	if err != nil {
		return nil, fmt.Errorf("failed to get response: %w", err)
	}
	defer resp.Body.Close()

	data, err := types.DecodeResponse[[]string](resp.Body)

	return data.Obj, err
}

func (s *ApiClient) do(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.apiToken))
	req.Header.Set("Content-Type", "application/json")

	slog.Debug("api request", "method", req.Method, "url", req.URL)

	return s.httpClient.Do(req)
}

func (s *ApiClient) get(rawURL string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, rawURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := s.do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		resp.Body.Close()
		return nil, fmt.Errorf("unexpected response status: %s", resp.Status)
	}

	return resp, nil
}

func NewApiClient(baseURL *url.URL, apiToken string) *ApiClient {
	if baseURL == nil {
		baseURL = &url.URL{}
	}

	return &ApiClient{
		baseURL:    baseURL,
		httpClient: &http.Client{Timeout: 15 * time.Second},
		apiToken:   apiToken,
	}
}
