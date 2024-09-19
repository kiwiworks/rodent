package server

import (
	"github.com/danielgtaylor/huma/v2"
)

type Security interface {
	SecurityConfig() *SecurityConfig
	SecurityMiddleware(ctx huma.Context, next func(ctx huma.Context))
}

type OpenApiConfig struct {
	Securities map[string]Security
}
