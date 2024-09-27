package sdk

import (
	"net/http"
	"net/url"

	"github.com/kiwiworks/rodent/errors"
	"github.com/kiwiworks/rodent/system/opt"
)

type Client struct {
	endpoint url.URL
	cfg      Config
}

func New(endpoint string, opts ...opt.Option[Config]) (*Client, error) {
	endpointURL, err := url.Parse(endpoint)
	if err != nil {
		return nil, errors.Wrapf(err, "invalid url for endpoint '%s'", endpoint)
	}
	cfg := NewConfig(opts...)
	return &Client{
		endpoint: *endpointURL,
		cfg:      cfg,
	}, nil
}

func (c Client) Request(method, path string, opts ...opt.Option[Request]) *Request {
	endpoint := c.endpoint.JoinPath(path)
	request := Request{
		method:   method,
		endpoint: *endpoint,
		values:   make(url.Values),
		headers:  make(http.Header),
		errors:   []error{},
		timeout:  c.cfg.Timeout,
	}
	opt.Apply(&request, opts...)
	return &request
}
