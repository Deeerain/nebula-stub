package model

type PageData[T any] struct {
	Title string
	Data  *T
}
