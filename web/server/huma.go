package server

import (
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
	"go.uber.org/fx"

	"github.com/kiwiworks/rodent/system/manifest"
	"github.com/kiwiworks/rodent/web/auth"
)

type HumaParams struct {
	fx.In
	Manifest      *manifest.Manifest
	Mux           *chi.Mux
	AuthProviders *auth.Providers
}

func NewHuma(params HumaParams) huma.API {
	api := humachi.New(params.Mux, huma.DefaultConfig(params.Manifest.Application, params.Manifest.Version.String()))
	doc := api.OpenAPI()
	doc.Components.SecuritySchemes = make(map[string]*huma.SecurityScheme)
	params.AuthProviders.HydrateOas3(doc)
	return api
}
