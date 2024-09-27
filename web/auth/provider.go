package auth

import (
	"github.com/danielgtaylor/huma/v2"
)

type Provider interface {
	SecurityScheme() *huma.SecurityScheme
	UserResolver(ctx huma.Context) (*ResolvedUser, error)
	AuthMiddleware(ctx huma.Context, next func(ctx huma.Context))
}
