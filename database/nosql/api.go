package nosql

import (
	"context"

	"rodent/system/errors"
)

func (c *Client) Create(ctx context.Context, id string, data any) (any, error) {
	res, err := c.db.Create(id, data)
	if err != nil {
		return nil, errors.Wrapf(err, "could not create '%s'", id)
	}
	return res, nil
}
