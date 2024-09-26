package errors

import (
	"errors"
	"fmt"
	"runtime"
)

func As[T error](err error) *T {
	var t T
	if errors.As(err, &t) {
		return &t
	}
	return nil
}

func Is(err, target error) bool {
	return errors.Is(err, target)
}

func Newf(format string, values ...any) error {
	return errors.New(fmt.Sprintf(format, values...))
}

func Wrapf(err error, format string, values ...any) error {
	if err == nil {
		return Newf(fmt.Sprintf(format, values...))
	}
	caller, file, line, ok := runtime.Caller(1)
	if !ok {
		return errors.Join(err, errors.New(fmt.Sprintf(format, values...)))
	}
	fn := runtime.FuncForPC(caller)
	msg := fmt.Sprintf(format, values...)
	final := fmt.Sprintf("%s %s:%d\n -%s", fn.Name(), file, line, msg)
	return errors.Join(err, errors.New(final))
}

func Must[T any](ok T, err error) T {
	if err != nil {
		panic(err)
	}
	return ok
}
