package assert

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFuncHasReturn(t *testing.T) {
	r := require.New(t)

	r.NoError(FuncHasReturn[bool](func() bool { return true }))
	r.Error(FuncHasReturn[bool](func() {}))
	r.Error(FuncHasReturn[bool](func() (bool, error) { return true, nil }))

	r.NoError(FuncHasReturnWithErr[int](func() (int, error) { return 1, nil }))
	r.Error(FuncHasReturnWithErr[int](func() int { return 1 }))
	r.Error(FuncHasReturnWithErr[int](func() error { return nil }))
}
