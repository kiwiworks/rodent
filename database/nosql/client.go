package nosql

import (
	"context"

	"github.com/surrealdb/surrealdb.go"

	"github.com/kiwiworks/rodent/errors"
)

type Client struct {
	db     *surrealdb.DB
	ns     string
	dbName string
}

func NewClient(cfg *SurrealDBConfig) (*Client, error) {
	endpoint := cfg.ConnectionURI()
	db, err := surrealdb.New(endpoint)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to connect to %s", endpoint)
	}

	_, err = db.Signin(map[string]any{
		"user": cfg.Username,
		"pass": cfg.Password,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to sign in %s", endpoint)
	}
	return &Client{
		db:     db,
		ns:     cfg.Namespace,
		dbName: cfg.Namespace,
	}, nil
}

func (c *Client) OnStart(context.Context) error {
	if _, err := c.db.Use(c.ns, c.dbName); err != nil {
		return errors.Wrapf(err, "failed to use namespace/database %s/%s", c.ns, c.dbName)
	}
	return nil
}

func (c *Client) OnStop(context.Context) error {
	c.db.Close()
	return nil
}
