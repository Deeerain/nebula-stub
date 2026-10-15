package xui

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"

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

type Connection struct {
}

type Subscription struct {
}

type ApiClient struct {
	baseURL    *url.URL
	httpClient *http.Client
	apiToken   string
}

func (s *ApiClient) GetSubscriptionsByTgID(telegramID int64) ([]string, error) {
	//TODO

	return nil, nil
}

func (s *ApiClient) GetClinetByTgID(telegramID int64) ([]ClientData, error) {
	url := s.baseURL.JoinPath("clients/get/tgId/").JoinPath(fmt.Sprint(telegramID))

	resp, err := s.get(url.String())
	if err != nil {
		return nil, fmt.Errorf("failed to execute: %w", err)
	}

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

	data, err := types.DecodeResponse[[]string](resp.Body)

	return data.Obj, err
}

func (s *ApiClient) GetLinks(email string) ([]string, error) {
	url := s.baseURL.JoinPath("clients/links").JoinPath(email)

	resp, err := s.get(url.String())
	if err != nil {
		return nil, fmt.Errorf("failed to get response: %w", err)
	}

	data, err := types.DecodeResponse[[]string](resp.Body)

	return data.Obj, err
}

func (s *ApiClient) NewSubscription(telegramID int64) (string, error) {
	return "", nil
}

func (s *ApiClient) do(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.apiToken))
	req.Header.Set("Content-Type", "application/json")

	slog.Info("api request", "url", req.URL, "auth", req.Header.Get("Authorization"), "ct", req.Header.Get("Content-type"))

	return s.httpClient.Do(req)
}

func (s *ApiClient) get(url string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	resp, err := s.do(req)
	return resp, err
}

func NewApiClient(baseURL *url.URL, apiToken string) *ApiClient {
	return &ApiClient{
		baseURL:    baseURL,
		httpClient: http.DefaultClient,
		apiToken:   apiToken,
	}
}
