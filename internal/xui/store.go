package xui

import (
	"fmt"
	"maps"
	"slices"
)

type TgID int64

type ClientStore map[TgID]any

func (cs ClientStore) addClient(tgId TgID, data any) {
	cs[tgId] = data
}

func (cs ClientStore) getClient(tgId TgID) (any, error) {
	data := cs[tgId]

	if data == nil {
		return nil, fmt.Errorf("clinet not found")
	}

	return data, nil
}

func (cs ClientStore) clear() {
	clear(cs)
}

func (cs ClientStore) has(tgId TgID) bool {
	keysitr := maps.Keys(cs)
	keys := slices.Collect(keysitr)

	return slices.Contains(keys, tgId)
}

type SubscriptionID string
type SubscriptionStore map[SubscriptionID]any

func (ss SubscriptionStore) has(id SubscriptionID) bool {
	return slices.Contains(slices.Collect(maps.Keys(ss)), id)
}
