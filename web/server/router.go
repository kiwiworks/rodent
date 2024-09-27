package server

import (
	"github.com/danielgtaylor/huma/v2"
	"github.com/go-chi/chi/v5"
	"go.uber.org/fx"

	"github.com/kiwiworks/rodent/web/auth"
)

type Router struct {
	mux *chi.Mux
	api huma.API
}

type RouterConfig struct {
	fx.In
	Mux            *chi.Mux
	Api            huma.API
	AuthMiddleware *auth.Middleware
}

func NewRouter(cfg RouterConfig) *Router {
	return &Router{
		mux: cfg.Mux,
		api: cfg.Api,
	}
}
