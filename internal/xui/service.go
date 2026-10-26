package xui

import (
	"fmt"
	"sync"

	"github.com/deeerain/nebula-stub/internal/client"
)

type XuiService struct {
	apiClient *ApiClient
	lock      sync.Mutex
}

func (s *XuiService) GetClientByTgID(id int64) (client.ProxyClient, error) {
	var proxyClient client.ProxyClient

	cd, err := s.apiClient.GetClinetByTgID(id)
	if err != nil {
		return proxyClient, fmt.Errorf("api error: %w", err)
	}

	if len(cd) <= 0 {
		return proxyClient, nil
	}

	c := cd[0]
	proxyClient = client.ProxyClient{
		TgID:        c.Client.TgID,
		SubID:       c.Client.SubID,
		Email:       c.Client.Email,
		Enable:      c.Client.Enable,
		LimitIP:     c.Client.LimitIP,
		LimitHwid:   c.Client.LimitHwid,
		ExpiryTime:  c.Client.ExpiryTime,
		UsedTraffic: c.UsedTraffic,
	}

	return proxyClient, nil
}

// GetLinksBySubID implements [subscription.SubscriptionService].
func (s *XuiService) GetLinksBySubID(subID string) ([]string, error) {
	links, err := s.apiClient.GetSubLinks(subID)
	if err != nil {
		return nil, fmt.Errorf("failed to get subscription links: %w", err)
	}

	return links, nil
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
