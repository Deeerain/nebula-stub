package model

import (
	"encoding/json"
	"net/http"
)

type Response[T any] interface {
	Write(writer http.ResponseWriter) error
}

type response[T any] struct {
	Success bool    `json:"success"`
	Error   *string `json:"error,omitempty"`
	Obj     T       `json:"obj,omitempty"`
}

func (r *response[T]) Write(writer http.ResponseWriter) error {
	writer.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(writer).Encode(r)
}

func NewResponse[T any](success bool, obj T, errormsg *string) Response[T] {
	return &response[T]{
		Success: success,
		Obj:     obj,
		Error:   errormsg,
	}
}
