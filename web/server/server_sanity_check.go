package server

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
	"go.uber.org/zap"

	"github.com/kiwiworks/rodent/logger"
	"github.com/kiwiworks/rodent/logger/props"
	"github.com/kiwiworks/rodent/web/http"
)

func securitySchemesForOperation(operation *huma.Operation) []map[string][]string {
	if operation == nil || operation.Security == nil {
		return nil
	}
	return operation.Security
}

func securitySchemesForPathItem(pathItem *huma.PathItem) map[http.Method][]map[string][]string {
	return map[http.Method][]map[string][]string{
		http.GET:     securitySchemesForOperation(pathItem.Get),
		http.DELETE:  securitySchemesForOperation(pathItem.Delete),
		http.HEAD:    securitySchemesForOperation(pathItem.Head),
		http.OPTIONS: securitySchemesForOperation(pathItem.Options),
		http.PATCH:   securitySchemesForOperation(pathItem.Patch),
		http.POST:    securitySchemesForOperation(pathItem.Post),
		http.PUT:     securitySchemesForOperation(pathItem.Put),
		http.TRACE:   securitySchemesForOperation(pathItem.Trace),
	}
}

func (s *Server) sanityCheck(ctx context.Context) {
	log := logger.FromContext(ctx)

	doc := s.router.api.OpenAPI()
	securitySchemes := doc.Components.SecuritySchemes
	for path, pathItem := range doc.Paths {
		schemes := securitySchemesForPathItem(pathItem)
		for method, securities := range schemes {
			if securities == nil {
				continue
			}
			for _, security := range securities {
				for provider, scopes := range security {
					_, ok := securitySchemes[provider]
					if !ok {
						log.Warn(
							"no provider found for endpoint, but one has been specified by the api.Handler",
							props.HttpMethod(string(method)),
							props.HttpPath(path),
							zap.String("auth.provider", provider),
							zap.Strings("auth.scopes", scopes),
						)
					}
					if scopes == nil {
						continue
					}
				}
			}
		}
	}
}
