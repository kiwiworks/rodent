package sdk

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"go.uber.org/multierr"
	"go.uber.org/zap"

	"github.com/kiwiworks/rodent/system/errors"
	"github.com/kiwiworks/rodent/system/logger"
	"github.com/kiwiworks/rodent/web/header"
)

func checkRequestErrors(request *Request) error {
	if len(request.errors) != 0 {
		return errors.Wrapf(multierr.Combine(request.errors...), "'%s': malformed request", request.endpoint.String())
	}
	return nil
}

func buildURLWithQueryParams(request *Request) url.URL {
	query := request.endpoint.Query()
	for k, v := range request.values {
		query.Set(k, strings.Join(v, ","))
	}
	u := request.endpoint
	u.RawQuery = query.Encode()
	return u
}

func createContextWithTimeout(ctx context.Context, timeout time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(ctx, timeout)
}

func buildHTTPRequest(ctx context.Context, method, url string, body io.Reader, headers http.Header) (
	*http.Request,
	error,
) {
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}
	for k, v := range headers {
		req.Header.Set(k, strings.Join(v, ","))
	}
	return req, nil
}

func closeResponseBody(body io.ReadCloser, log *zap.Logger) {
	if err := body.Close(); err != nil {
		log.Error("cleanup of response failed", zap.Error(err))
	}
}

func checkResponseStatusCode(res *http.Response, uri string) error {
	if res.StatusCode >= 400 || res.StatusCode < 200 {
		return errors.Newf("'%s': server replied with '%d' status ('%s')", uri, res.StatusCode, res.Status)
	}
	return nil
}

func deserializeResponse[Response any](responseBytes []byte, contentType, uri string) (*Response, error) {
	var response Response
	var err error
	switch contentType {
	case header.ContentTypeForm:
		err = errors.Newf("deserializing response from '%s' for '%s' is not implemented yet", header.ContentTypeForm, uri)
	case header.ContentTypeJson:
		err = json.Unmarshal(responseBytes, &response)
	default:
		err = json.Unmarshal(responseBytes, &response)
	}
	if err != nil {
		return nil, errors.Wrapf(err, "'%s': could not deserialize response", uri)
	}
	return &response, nil
}

func Execute[Response any](
	ctx context.Context,
	endpoint Client,
	request *Request,
) (*Response, error) {
	log := logger.FromContext(ctx).With(zap.String("urlWithQueryParams", request.endpoint.String()))
	cfg := endpoint.cfg

	if err := checkRequestErrors(request); err != nil {
		return nil, err
	}

	urlWithQueryParams := buildURLWithQueryParams(request)

	requestCtx, cancel := createContextWithTimeout(ctx, request.timeout)
	defer cancel()

	req, err := buildHTTPRequest(requestCtx, string(request.method), urlWithQueryParams.String(), request.body, request.headers)
	if err != nil {
		return nil, errors.Wrapf(err, "'%s': could not build a valid request", request.endpoint.String())
	}

	if err = cfg.InterceptRequest(ctx, req); err != nil {
		log.Error("failed to intercept request", zap.Error(err))
		return nil, errors.Wrapf(err, "failed to intercept request")
	}

	log.Debug("executing request",
		zap.String("uri", req.RequestURI),
		zap.String("method", req.Method))

	client := http.Client{}
	res, err := client.Do(req)
	if err != nil || res == nil {
		return nil, errors.Wrapf(err, "'%s': request failed", req.RequestURI)
	}

	defer closeResponseBody(res.Body, log)

	if err := cfg.InterceptResponse(ctx, res); err != nil {
		log.Error("failed to intercept response", zap.Error(err))
		return nil, errors.Wrapf(err, "failed to intercept response")
	}

	if err := checkResponseStatusCode(res, req.RequestURI); err != nil {
		return nil, err
	}

	responseBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrapf(err, "'%s': could not read from response body", req.RequestURI)
	}

	return deserializeResponse[Response](responseBytes, res.Header.Get(header.ContentType), req.RequestURI)
}
