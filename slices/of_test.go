package slices

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type testStringer struct {
	value string
}

func (ts *testStringer) String() string {
	return ts.value
}

func TestOf(t *testing.T) {
	tt := []struct {
		name   string
		input  []int
		output []int
	}{
		{
			name:   "Empty",
			input:  []int{},
			output: []int{},
		},
		{
			name:   "SingleElement",
			input:  []int{1},
			output: []int{1},
		},
		{
			name:   "MultipleElements",
			input:  []int{1, 2, 3, 4, 5},
			output: []int{1, 2, 3, 4, 5},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			res := Of(tc.input...)
			require.Equal(t, tc.output, res)
		})
	}
}

func TestOfStringers(t *testing.T) {
	tt := []struct {
		name   string
		input  []fmt.Stringer
		output []string
	}{
		{
			name:   "Empty",
			input:  []fmt.Stringer{},
			output: []string{},
		},
		{
			name:   "SingleStringer",
			input:  []fmt.Stringer{&testStringer{value: "test"}},
			output: []string{"test"},
		},
		{
			name:   "MultipleStringers",
			input:  []fmt.Stringer{&testStringer{value: "one"}, &testStringer{value: "two"}},
			output: []string{"one", "two"},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			res := OfStringers(tc.input...)
			require.Equal(t, tc.output, res)
		})
	}
}
