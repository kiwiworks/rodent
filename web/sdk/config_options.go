package sdk

import (
	"context"
	"fmt"
	"net/http"

	"rodent/system/opt"
)

func ApiKeyAuth(headerName, headerValue string) opt.Option[Config] {
	return func(cfg *Config) {
		cfg.RequestInterceptors = append(cfg.RequestInterceptors, func(ctx context.Context, req *http.Request) error {
			req.Header.Add(headerName, headerValue)
			return nil
		})
	}
}

func BearerAuth(token string) opt.Option[Config] {
	bearerToken := fmt.Sprintf("Bearer %s", token)
	return func(cfg *Config) {
		cfg.RequestInterceptors = append(cfg.RequestInterceptors, func(ctx context.Context, req *http.Request) error {
			req.Header.Add("Authorization", bearerToken)
			return nil
		})
	}
}
