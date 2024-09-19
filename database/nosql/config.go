package nosql

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/kiwiworks/rodent/slices"
	"github.com/kiwiworks/rodent/system/errors"
)

type SurrealDBConfig struct {
	Username  string
	Password  string
	Host      string
	Port      uint16
	Secure    bool
	Namespace string
	Database  string
}

var validSchemes = []string{"ws", "wss", "http", "https"}

func ConfigFromUrl(uri url.URL) (*SurrealDBConfig, error) {
	if !slices.Contains(validSchemes, uri.Scheme) {
		return nil, errors.Newf("invalid database scheme '%s' expected one of '%s'", uri.Scheme, strings.Join(validSchemes, ","))
	}
	user := uri.User
	if user == nil {
		return nil, errors.Newf("invalid database user")
	}
	username := user.Username()
	password, exists := user.Password()
	if !exists {
		return nil, errors.Newf("invalid database password")
	}
	secure := uri.Scheme == "wss"
	portStr := uri.Port()
	if portStr == "" {
		return nil, errors.Newf("missing port")
	}
	port, err := strconv.ParseInt(portStr, 10, 16)
	if err != nil {
		return nil, errors.Wrapf(err, "invalid port '%s'", portStr)
	}

	query := uri.Query()
	database := query.Get("db")
	namespace := query.Get("ns")

	return &SurrealDBConfig{
		Username:  username,
		Password:  password,
		Host:      uri.Hostname(),
		Port:      uint16(port),
		Secure:    secure,
		Namespace: namespace,
		Database:  database,
	}, nil
}

func (s SurrealDBConfig) ConnectionURI() string {
	scheme := "wss"
	if !s.Secure {
		scheme = "ws"
	}
	return (&url.URL{
		Scheme: scheme,
		Host:   fmt.Sprintf("%s:%d", s.Host, s.Port),
		Path:   "rpc",
	}).String()
}
