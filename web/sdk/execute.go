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

	"github.com/kiwiworks/rodent/errors"
	"github.com/kiwiworks/rodent/lang"
	"github.com/kiwiworks/rodent/logger"
	"github.com/kiwiworks/rodent/logger/props"
	"github.com/kiwiworks/rodent/web/header"
)

func checkRequestErrors(request *Request) error {
	if len(request.Errors) != 0 {
		return errors.Wrapf(multierr.Combine(request.Errors...), "'%s': malformed request", request.Endpoint.String())
	}
	return nil
}

func buildURLWithQueryParams(request *Request) url.URL {
	query := request.Endpoint.Query()
	for k, v := range request.Values {
		query.Set(k, strings.Join(v, ","))
	}
	u := request.Endpoint
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

func checkResponseStatusCode(res *http.Response) error {
	if res.StatusCode >= 400 || res.StatusCode < 200 {
		return errors.Newf("server replied with '%d' status ('%s')", res.StatusCode, res.Status)
	}
	return nil
}

func deserializeResponse[Response any](responseBytes []byte, contentType, uri string) (*Response, error) {
	var result lang.Either[Response, string]
	var err error
	switch contentType {
	case header.ContentTypeForm:
		err = errors.Newf("deserializing response from '%s' for '%s' is not implemented yet", header.ContentTypeForm, uri)
	case header.ContentTypeJson:
		err = json.Unmarshal(responseBytes, &result)
	default:
		err = json.Unmarshal(responseBytes, &result)
	}
	if err != nil {
		return nil, errors.Wrapf(err, "'%s': could not deserialize response", uri)
	}
	if response, isResponse := result.Left(); isResponse {
		return &response, nil
	}
	if other, isOther := result.Right(); isOther {
		return nil, errors.Newf("'%s': unexpected response: '%s'", uri, other)
	}
	return nil, errors.Newf("'%s': unexpected response", uri)
}

func Execute[Response any](
	ctx context.Context,
	endpoint Client,
	request *Request,
) (*Response, error) {
	log := logger.FromContext(ctx).With(props.HttpMethod(request.Method))
	cfg := endpoint.cfg

	if err := checkRequestErrors(request); err != nil {
		return nil, err
	}

	urlWithQueryParams := buildURLWithQueryParams(request)

	requestCtx, cancel := createContextWithTimeout(ctx, request.Timeout)
	defer cancel()

	req, err := buildHTTPRequest(requestCtx, request.Method, urlWithQueryParams.String(), request.Body, request.Headers)
	if err != nil {
		return nil, errors.Wrapf(err, "'%s': could not build a valid request", request.Endpoint.String())
	}
	log = log.With(props.HttpPath(req.URL.String()))

	if err = cfg.InterceptRequest(ctx, req); err != nil {
		log.Error("failed to intercept request", zap.Error(err))
		return nil, errors.Wrapf(err, "failed to intercept request")
	}

	log.Debug("executing request")

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

	var finalErr error

	err = checkResponseStatusCode(res)
	if err != nil {
		log.Error("invalid http response", zap.Error(err))
		finalErr = err
	}
	responseBytes, err := io.ReadAll(res.Body)
	if err != nil {
		log.Error("could not read response", zap.Error(err))
		finalErr = multierr.Combine(finalErr, errors.Wrapf(err, "'%s': could not read response", req.RequestURI))
		return nil, finalErr
	}

	response, err := deserializeResponse[Response](responseBytes, res.Header.Get(header.ContentType), req.RequestURI)
	if err != nil {
		log.Error("could not deserialize response", zap.Error(err))
		finalErr = multierr.Combine(finalErr, err)
	}
	return response, finalErr
}
