package server

import (
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"go.uber.org/fx"

	"github.com/kiwiworks/rodent/cas"
	"github.com/kiwiworks/rodent/logger"
	"github.com/kiwiworks/rodent/system/manifest"
)

type Router struct {
	mux *chi.Mux
	api huma.API
}

func NewMux() *chi.Mux {
	mux := chi.NewRouter()
	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "X-Refresh-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))
	mux.Use(logger.Middleware())
	return mux
}

func NewHuma(mux *chi.Mux, manifest *manifest.Manifest) huma.API {
	api := humachi.New(mux, huma.DefaultConfig(manifest.Application, manifest.Version.String()))
	oas := api.OpenAPI()
	oas.Components.SecuritySchemes = make(map[string]*huma.SecurityScheme)
	oas.Components.SecuritySchemes["protected"] = &huma.SecurityScheme{
		Type:         "oauth2",
		Description:  "Casdoor managed authentication",
		Name:         "protected",
		BearerFormat: "Bearer",
	}
	return api
}

type RouterConfig struct {
	fx.In
	Mux           *chi.Mux
	Api           huma.API
	CasMiddleware *cas.Middleware
}

func NewRouter(cfg RouterConfig) *Router {
	return &Router{
		mux: cfg.Mux,
		api: cfg.Api,
	}
}
