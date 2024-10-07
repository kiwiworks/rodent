package sdk

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"go.uber.org/multierr"

	"github.com/kiwiworks/rodent/errors"
	"github.com/kiwiworks/rodent/system/opt"
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

func AddRequestInterceptor(interceptor func(ctx context.Context, req *http.Request) error) opt.Option[Config] {
	return func(cfg *Config) {
		cfg.RequestInterceptors = append(cfg.RequestInterceptors, interceptor)
	}
}

func AddResponseInterceptor(interceptor func(ctx context.Context, resp *http.Response) error) opt.Option[Config] {
	return func(cfg *Config) {
		cfg.ResponseInterceptors = append(cfg.ResponseInterceptors, interceptor)
	}
}

type ErrorResponse[T any] struct {
	Response T
}

func (e ErrorResponse[T]) Error() string {
	bytes, err := json.Marshal(e.Response)
	if err != nil {
		bytes = []byte(fmt.Sprintf("INVALID"))
	}
	return fmt.Sprintf("error response: %s", string(bytes))
}

func SetErrorResponseType[T any](handler ...func(response T) string) opt.Option[Config] {
	return func(cfg *Config) {
		cfg.ResponseInterceptors = append(cfg.ResponseInterceptors, func(ctx context.Context, resp *http.Response) error {
			var errorResponse T

			err := checkResponseStatusCode(resp)
			if err == nil {
				return nil
			}
			if decodingErr := json.NewDecoder(resp.Body).Decode(&errorResponse); decodingErr != nil {
				return multierr.Combine(err, decodingErr)
			}
			if len(handler) > 0 {
				msg := handler[0](errorResponse)
				return errors.Wrapf(err, "api returned some details about the error: %s", msg)
			}
			return multierr.Combine(err, ErrorResponse[T]{Response: errorResponse})
		})
	}
}
