package xui

import (
	"fmt"
	"sync"
)

type XuiService struct {
	apiClient *ApiClient
	lock      sync.Mutex
}

func (s *XuiService) GetClient(tgId TgID) (any, error) {
	cd, err := s.apiClient.GetClinetByTgID(int64(tgId))
	if err != nil {
		return nil, fmt.Errorf("api error: %w", err)
	}

	return cd, nil
}

func (s *XuiService) GetSubscription(subId SubscriptionID) (any, error) {
	subLinks, err := s.apiClient.GetSubLinks(string(subId))
	if err != nil {
		return nil, fmt.Errorf("failed to get subscriptions: %w", err)
	}

	return subLinks, nil
}

func NewXuiService(client *ApiClient) *XuiService {
	return &XuiService{
		apiClient: client,
	}
}
