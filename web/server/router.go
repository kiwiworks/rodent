package server

import (
	"github.com/danielgtaylor/huma/v2"
	"github.com/go-chi/chi/v5"
	"go.uber.org/fx"

	"github.com/kiwiworks/rodent/logger"
	"github.com/kiwiworks/rodent/web/auth"
)

type Router struct {
	mux            *chi.Mux
	api            huma.API
	authMiddleware *auth.Middleware
}

type RouterConfig struct {
	fx.In
	Mux            *chi.Mux
	Api            huma.API
	AuthMiddleware *auth.Middleware `optional:"true"`
}

func NewRouter(cfg RouterConfig) *Router {
	log := logger.New()
	if cfg.AuthMiddleware != nil {
		log.Info("using auth middleware")
		cfg.Api.UseMiddleware(cfg.AuthMiddleware.Middleware)
	}
	return &Router{
		mux:            cfg.Mux,
		api:            cfg.Api,
		authMiddleware: cfg.AuthMiddleware,
	}
}
