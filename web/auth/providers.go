package auth

import (
	"github.com/danielgtaylor/huma/v2"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/kiwiworks/rodent/logger"
)

type (
	Providers struct {
		Registry map[string]Provider
	}
	ProvidersParams struct {
		fx.In
		Providers []Provider `group:"auth.provider"`
	}
)

func NewProviders(params ProvidersParams) *Providers {
	log := logger.New()

	registry := make(map[string]Provider)
	for _, provider := range params.Providers {
		scheme := provider.SecurityScheme()
		log.Info("registering auth provider",
			zap.String("auth.provider.name", scheme.Name),
			zap.String("auth.provider.type", scheme.Type),
		)
		registry[scheme.Name] = provider
	}
	return &Providers{
		Registry: registry,
	}
}

func (p *Providers) HydrateOas3(doc *huma.OpenAPI) {
	for _, provider := range p.Registry {
		scheme := provider.SecurityScheme()
		doc.Components.SecuritySchemes[scheme.Name] = scheme
	}
}
