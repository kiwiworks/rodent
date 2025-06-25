package auth

import (
	"fmt"

	"github.com/danielgtaylor/huma/v2"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/kiwiworks/rodent/logger"
	"github.com/kiwiworks/rodent/slices"
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
		manifest := provider.Manifest()
		scheme := manifest.SecurityScheme
		log.Info("registering auth provider",
			zap.String("auth.provider.name", manifest.Name),
			zap.String("auth.provider.type", scheme.Type),
		)
		registry[manifest.Name] = provider
	}
	return &Providers{
		Registry: registry,
	}
}

func (p *Providers) HydrateOas3(doc *huma.OpenAPI) {
	for _, provider := range p.Registry {
		manifest := provider.Manifest()
		scheme := manifest.SecurityScheme
		doc.Components.SecuritySchemes[manifest.Name] = scheme
	}
}

func (p *Providers) ValidateSecuritySchemes(doc *huma.OpenAPI) error {
	var missingProviders []struct {
		Name string
		Path string
	}

	// Check all operations for security schemes
	for path, pathItem := range doc.Paths {
		operations := []*huma.Operation{
			pathItem.Get, pathItem.Post, pathItem.Put, pathItem.Patch, pathItem.Delete,
			pathItem.Head, pathItem.Options, pathItem.Trace,
		}

		for _, op := range operations {
			if op == nil || op.Security == nil {
				continue
			}

			for _, security := range op.Security {
				for schemeName := range security {
					if _, exists := p.Registry[schemeName]; !exists {
						if !slices.Any(missingProviders, func(p struct{ Name, Path string }) bool {
							return p.Name == schemeName
						}) {
							missingProviders = append(missingProviders, struct {
								Name string
								Path string
							}{Name: schemeName, Path: path})
						}
					}
				}
			}
		}
	}

	if len(missingProviders) > 0 {
		return fmt.Errorf("authentication providers not registered for security schemes: %v", missingProviders)
	}

	return nil
}
