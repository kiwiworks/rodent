package auth

import (
	"github.com/danielgtaylor/huma/v2"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/kiwiworks/rodent/logger"
)

type (
	Middleware struct {
		api       huma.API
		providers map[string]Provider
	}
	MiddlewareParams struct {
		fx.In
		API       huma.API
		Providers *Providers
	}
)

func NewMiddleware(params MiddlewareParams) *Middleware {
	return &Middleware{
		api:       params.API,
		providers: params.Providers.Registry,
	}
}

func (m *Middleware) writeError(ctx huma.Context, status int, msg string, err error) {
	op := ctx.Operation()
	log := logger.New().With(zap.String("http.method", ctx.Method()), zap.String("http.path", op.Path), zap.Int("http.status", status))

	if err = huma.WriteErr(m.api, ctx, status, msg, err); err != nil {
		log.Error("failed to write error response", zap.Error(err))
	}
}

func (m *Middleware) Middleware(ctx huma.Context, next func(ctx huma.Context)) {
	op := ctx.Operation()
	log := logger.New().With(zap.String("http.method", ctx.Method()), zap.String("http.path", op.Path))

	if op.Security == nil {
		next(ctx)
		return
	}

	for _, security := range op.Security {
		for name, scopes := range security {
			//todo handle scopes
			_ = scopes
			provider, exists := m.providers[name]
			if !exists {
				log.Warn("operation defines a security scheme that is not supported: (missing provider)", zap.String("provider.name", name))
				continue
			}
			user, err := provider.UserResolver(ctx)
			if err != nil {
				m.writeError(ctx, 401, "could not resolve user from credentials", err)
				return
			}
			next(huma.WithContext(ctx, injectUser(ctx.Context(), user)))
			return
		}
	}
}
