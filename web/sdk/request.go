package sdk

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/google/go-querystring/query"

	"github.com/kiwiworks/rodent/errors"
	"github.com/kiwiworks/rodent/maps"
	"github.com/kiwiworks/rodent/system/opt"
	"github.com/kiwiworks/rodent/web/header"
)

type Request struct {
	method   string
	endpoint url.URL
	values   url.Values
	headers  http.Header
	body     io.ReadCloser
	errors   []error
	timeout  time.Duration
}

func WithTimeout(ttl time.Duration) opt.Option[Request] {
	return func(opt *Request) {
		opt.timeout = ttl
	}
}

func WithJsonBody[Body any](value Body) opt.Option[Request] {
	return func(request *Request) {
		WithHeader(header.ContentType, header.ContentTypeJson)(request)
		buf, err := json.Marshal(value)
		if err != nil {
			request.errors = append(request.errors, errors.Wrapf(err, "JSON serialization of body '%T' failed", value))
			return
		}

		// why this ? because it allows us to use file upload or more complex types (ie: multipart/streaming)
		// down the line without changing the ease of use of the API
		request.body = io.NopCloser(bytes.NewReader(buf))
	}
}

func WithMultipartFormBody(formWriter multipart.Writer, value *bytes.Buffer) opt.Option[Request] {
	return func(request *Request) {
		WithHeader(header.ContentType, formWriter.FormDataContentType())(request)
		request.body = io.NopCloser(value)
	}
}

func WithFormBody[Body any](value Body) opt.Option[Request] {
	return func(request *Request) {
		WithHeader(header.ContentType, header.ContentTypeForm)(request)
		buf, err := query.Values(value)
		if err != nil {
			request.errors = append(request.errors, errors.Wrapf(err, "Form serialization of body '%T' failed", value))
			return
		}

		// why this ? because it allows us to use file upload or more complex types (ie: multipart/streaming)
		// down the line without changing the ease of use of the API
		request.body = io.NopCloser(strings.NewReader(buf.Encode()))
	}
}

func WithHeader(name string, value string) opt.Option[Request] {
	return func(request *Request) {
		// yeah, go headers can have multiple values 🤷
		request.headers[name] = []string{value}
	}
}

func WithQueryParam(name string, value string) opt.Option[Request] {
	return func(request *Request) {
		if _, ok := request.values[name]; !ok {
			request.values[name] = []string{value}
		} else {
			request.values[name] = append(request.values[name], value)
		}
	}
}

func WithQueryInt(name string, value any) opt.Option[Request] {
	return func(request *Request) {
		var upcast int64
		switch v := value.(type) {
		case int:
			upcast = int64(v)
		case int8:
			upcast = int64(v)
		case int16:
			upcast = int64(v)
		case int32:
			upcast = int64(v)
		case int64:
			upcast = v
		default:
			request.errors = append(request.errors, errors.Newf("invalid type '%T' for query param '%s' wanted int|int8|int16|int32|int64", value, name))
			return
		}
		str := strconv.FormatInt(upcast, 10)
		WithQueryParam(name, str)(request)
	}
}

func WithQueryFloat(name string, value any) opt.Option[Request] {
	return func(request *Request) {
		var upcast float64
		switch v := value.(type) {
		case float32:
			upcast = float64(v)
		case float64:
			upcast = v
		default:
			request.errors = append(request.errors, errors.Newf("invalid type '%T' for query param '%s' wanted float32|float64", value, name))
			return
		}
		str := strconv.FormatFloat(upcast, 'f', -1, 64) // 'f' for decimal notation
		WithQueryParam(name, str)(request)
	}
}

func WithQueryStrings(name string, values ...string) opt.Option[Request] {
	return func(request *Request) {
		for _, value := range values {
			WithQueryParam(name, value)(request)
		}
	}
}

func WithQueryStruct[T any](value T) opt.Option[Request] {
	return func(request *Request) {
		values, err := query.Values(value)
		if err != nil {
			request.errors = append(request.errors, errors.Wrapf(err, "Query serialization of body '%T' failed", value))
			return
		}
		request.values = maps.Merged(request.values, values)
	}
}
