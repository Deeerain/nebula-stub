package client

import (
	"fmt"
	"maps"
	"slices"
)

type Client interface {
	GetID() int64
	GetEmail() string
}

type ClientStore[T Client] struct {
	items map[int64]T
}

func (s *ClientStore[T]) GetClientByID(id int64) (T, error) {
	ids := slices.Collect(maps.Keys(s.items))
	if slices.Contains(ids, id) {
		return s.items[id], nil
	}

	var zero T
	return zero, fmt.Errorf("client not found")
}

func (s *ClientStore[T]) AddClient(id int64, data T) {
	s.items[id] = data
}

func NewClientStore[T Client]() *ClientStore[T] {
	return &ClientStore[T]{
		items: make(map[int64]T),
	}
}
