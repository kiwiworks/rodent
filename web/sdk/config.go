package sdk

import (
	"context"
	"net/http"
	"time"

	"rodent/system/opt"
)

type (
	Config struct {
		Timeout              time.Duration
		RequestInterceptors  []RequestInterceptor
		ResponseInterceptors []ResponseInterceptor
	}
	RequestInterceptor  func(ctx context.Context, req *http.Request) error
	ResponseInterceptor func(ctx context.Context, resp *http.Response) error
)

func NewConfig(opts ...opt.Option[Config]) Config {
	cfg := Config{
		Timeout:              5 * time.Second,
		RequestInterceptors:  []RequestInterceptor{},
		ResponseInterceptors: []ResponseInterceptor{},
	}
	opt.Apply(&cfg, opts...)
	return cfg
}

func (c *Config) InterceptRequest(ctx context.Context, req *http.Request) error {
	for _, interceptor := range c.RequestInterceptors {
		if err := interceptor(ctx, req); err != nil {
			return err
		}
	}
	return nil
}

func (c *Config) InterceptResponse(ctx context.Context, resp *http.Response) error {
	for _, interceptor := range c.ResponseInterceptors {
		if err := interceptor(ctx, resp); err != nil {
			return err
		}
	}
	return nil
}
