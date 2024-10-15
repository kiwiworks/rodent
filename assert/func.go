package assert

import (
	"reflect"

	"github.com/kiwiworks/rodent/errors"
)

func FuncHasReturnWithErr[ReturnT any](f any) error {
	var r ReturnT
	fType := reflect.TypeOf(f)
	rType := reflect.TypeOf(r)
	if fType.Kind() != reflect.Func {
		return errors.Newf("expected a function, got '%s'", fType.Kind())
	}
	if fType.NumOut() != 2 {
		return errors.Newf("expected a function with two return values, got '%d'", fType.NumOut())
	}
	if fType.Out(0) != rType {
		return errors.Newf("expected a function with return type '%s', got '%s'", rType, fType.Out(0))
	}
	if fType.Out(1) != reflect.TypeOf((*error)(nil)).Elem() {
		return errors.Newf("expected a function with error return type, got '%s'", fType.Out(1))
	}
	return nil
}

func FuncHasReturn[ReturnT any](f any) error {
	var r ReturnT
	fType := reflect.TypeOf(f)
	rType := reflect.TypeOf(r)
	if fType.Kind() != reflect.Func {
		return errors.Newf("expected a function, got '%s'", fType.Kind())
	}
	if fType.NumOut() != 1 {
		return errors.Newf("expected a function with one return value, got '%d'", fType.NumOut())
	}
	if fType.Out(0) != rType {
		return errors.Newf("expected a function with return type '%s', got '%s'", rType, fType.Out(0))
	}
	return nil
}
