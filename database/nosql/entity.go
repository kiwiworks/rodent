package nosql

import (
	"fmt"
	"sync"

	"github.com/google/uuid"
	"github.com/surrealdb/surrealdb.go"

	"github.com/kiwiworks/rodent/system/errors"
)

type IEntity[T any] interface {
	AsEntity() Entity[T]
}

func IdOf(resource string, id uuid.UUID) string {
	return fmt.Sprintf("%s:%s", resource, id.String())
}

type Entity[T any] struct {
	lock  sync.RWMutex
	Id    string `json:"id"`
	Value *T     `json:"value"`
}

func EntityById[T any](resource string, id uuid.UUID, value T) Entity[T] {
	return Entity[T]{
		Id:    IdOf(resource, id),
		Value: &value,
	}
}

func (e *Entity[T]) id() string {
	e.lock.RLock()
	defer e.lock.RUnlock()
	return e.Id
}

func (e *Entity[T]) Load(client *Client) error {
	id := e.id()
	data, err := surrealdb.SmartUnmarshal[T](client.db.Select(id))
	if err != nil {
		return errors.Wrapf(err, "could not load '%s'", id)
	}
	e.lock.Lock()
	defer e.lock.Unlock()
	e.Value = &data
	return nil
}

func (e *Entity[T]) Save(client *Client) error {
	e.lock.RLock()
	_, err := surrealdb.SmartMarshal(client.db.Update, e.Value)
	e.lock.RUnlock()
	if err != nil {
		return errors.Wrapf(err, "could not load '%s'", e.id())
	}
	e.lock.Lock()
	defer e.lock.Unlock()
	return nil
}
