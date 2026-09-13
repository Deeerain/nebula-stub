package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type ApiResponse[T any] struct {
	Success bool   `json:"success"`
	Msg     string `json:"msg"`
	Obj     T      `json:"obj"`
}

type Client struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	SubID    string `json:"subId"`
	UUID     string `json:"uuid"`
	TgID     int64  `json:"telegramId"`
	Inbounds []int  `json:"inbounds"`
}

type ClientData struct {
	Client Client `json:"client"`
}

type Connection struct {
}

type Subscription struct {
}

type XrayService interface {
	NewSubscription(telegramID int64) (string, error)
	GetClinetByTgID(telegramID int64) ([]ClientData, error)
	GetLinksByEmail(email string) ([]string, error)
}

type xuiService struct {
	baseURL    *url.URL
	httpClient *http.Client
	apiToken   string
	clients    []ClientData
}

func (s *xuiService) GetClinetByTgID(telegramID int64) ([]ClientData, error) {
	url := s.baseURL.JoinPath("clients/get/tgId/").JoinPath(fmt.Sprint(telegramID))

	resp, err := s.get(url.String())
	if err != nil {
		return nil, fmt.Errorf("failed to execute: %w", err)
	}

	var data ApiResponse[[]ClientData]

	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, fmt.Errorf("failed to decode response data: %w", err)
	}

	return data.Obj, nil
}

func (s *xuiService) GetLinksByEmail(email string) ([]string, error) {
	url := s.baseURL.JoinPath("clients/get/links").JoinPath(email)

	resp, err := s.get(url.String())
	if err != nil {
		return nil, fmt.Errorf("failed to get response: %w", err)
	}

	var data ApiResponse[[]string]

	err = json.NewDecoder(resp.Body).Decode(data)
	if err != nil {
		return nil, fmt.Errorf("failed to decode object: %w", err)
	}

	return data.Obj, nil
}

func (s *xuiService) GetSubLinks(subId string) ([]string, error) {
	url := s.baseURL.JoinPath("clients/subLinks").JoinPath(subId)

	resp, err := s.get(url.String())
	if err != nil {
		return nil, fmt.Errorf("failed to get response: %w", err)
	}

	var data ApiResponse[[]string]

	err = json.NewDecoder(resp.Body).Decode(data)
	if err != nil {
		return nil, fmt.Errorf("failed to decode object: %w", err)
	}

	return data.Obj, nil
}

func (s *xuiService) NewSubscription(telegramID int64) (string, error) {
	return "", nil
}

func (s *xuiService) do(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.apiToken))
	req.Header.Set("Content-type", "application/json")
	return s.httpClient.Do(req)
}

func (s *xuiService) get(url string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	resp, err := s.do(req)
	return resp, err
}

func NewXrayService(baseURL *url.URL, apiToken string) XrayService {
	return &xuiService{
		baseURL:    baseURL,
		httpClient: http.DefaultClient,
		apiToken:   apiToken,
	}
}
