package errors

import "fmt"

type NotFoundError struct {
	Entity              string
	RequestedIdentifier string
}

func (e NotFoundError) Error() string {
	return fmt.Sprintf("Entity '%s' with identifier '%s' not found", e.Entity, e.RequestedIdentifier)
}

func NotFound(entity string, requestedIdentifier fmt.Stringer) NotFoundError {
	return NotFoundError{
		Entity:              entity,
		RequestedIdentifier: requestedIdentifier.String(),
	}
}

func IsNotFound(err error) bool {
	return As[NotFoundError](err) != nil
}
