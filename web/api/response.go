package api

import (
	"fmt"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
)

type Created struct {
	Id uuid.UUID `json:"id"`
}

type Response[T any] struct {
	Body T
}

func Ok[T any](body T) (*Response[T], error) {
	return &Response[T]{Body: body}, nil
}

func Forbidden[T any](msg string, err error) (*Response[T], error) {
	return (*Response[T])(nil), huma.Error403Forbidden(msg, err)
}

func Unauthorized[T any](msg string, err error) (*Response[T], error) {
	return (*Response[T])(nil), huma.Error401Unauthorized(msg, err)
}

func NotFound[T any](msg string) (*Response[T], error) {
	return (*Response[T])(nil), huma.Error404NotFound(msg)
}

func OkCreated(id uuid.UUID) (*Response[Created], error) {
	return &Response[Created]{Body: Created{Id: id}}, nil
}

func NotImplemented[T any](format string, args ...any) (*Response[T], error) {
	return (*Response[T])(nil), huma.Error501NotImplemented(fmt.Sprintf(format, args...))
}
