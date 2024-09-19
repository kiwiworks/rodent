package functional

import (
	"encoding/json"

	"github.com/kiwiworks/rodent/system/errors"
)

type Either[L, R any] struct {
	isLeft bool
	left   L
	right  R
}

func Left[L, R any](value L) Either[L, R] {
	return Either[L, R]{
		isLeft: true,
		left:   value,
	}
}

func Right[L, R any](value R) Either[L, R] {
	return Either[L, R]{
		right: value,
	}
}

func (e *Either[L, R]) IsLeft() bool {
	return e.isLeft
}

func (e *Either[L, R]) IsRight() bool {
	return !e.isLeft
}

func (e *Either[L, R]) Left() (L, bool) {
	return e.left, e.isLeft
}

func (e *Either[L, R]) Right() (R, bool) {
	return e.right, !e.isLeft
}

func (e *Either[L, R]) UnmarshalJSON(data []byte) error {
	var l L
	if err := json.Unmarshal(data, &l); err == nil {
		e.isLeft = true
		e.left = l
		return nil
	}

	var r R
	if err := json.Unmarshal(data, &r); err == nil {
		e.isLeft = false
		e.right = r
		return nil
	}

	return errors.Newf("data does not match either type, wanted either %T or %T", l, r)
}
