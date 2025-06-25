package auth

import (
	"github.com/danielgtaylor/huma/v2"
)

type Manifest struct {
	Name           string
	SecurityScheme *huma.SecurityScheme
}

type Provider interface {
	Manifest() *Manifest
	UserResolver(ctx huma.Context) (*ResolvedUser, error)
	AuthMiddleware(ctx huma.Context, next func(ctx huma.Context))
}
