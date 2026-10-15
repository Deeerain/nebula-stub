package types

import (
	"encoding/json"
	"io"
)

type ApiResponse[T any] struct {
	Success bool   `json:"success"`
	Msg     string `json:"msg"`
	Obj     T      `json:"obj"`
}

func DecodeResponse[T any](r io.Reader) (ApiResponse[T], error) {
	var obj ApiResponse[T]

	err := json.NewDecoder(r).Decode(&obj)

	return obj, err
}
