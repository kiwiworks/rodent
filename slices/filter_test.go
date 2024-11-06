package slices

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUniqueBy(t *testing.T) {
	tests := []struct {
		name     string
		vs       []int
		selector func(int) int
		want     []int
	}{
		{
			name: "empty slice",
			vs:   []int{},
			selector: func(n int) int {
				return n
			},
			want: []int{},
		},
		{
			name: "non-empty unique slice",
			vs:   []int{1, 2, 3},
			selector: func(n int) int {
				return n
			},
			want: []int{1, 2, 3},
		},
		{
			name: "non-empty slice with duplicates",
			vs:   []int{1, 2, 2, 3, 4, 3},
			selector: func(n int) int {
				return n
			},
			want: []int{1, 2, 3, 4},
		},
		{
			name: "non-empty slice with negative numbers",
			vs:   []int{1, -2, -2, 3, 4, -3},
			selector: func(n int) int {
				return n
			},
			want: []int{1, -2, 3, 4, -3},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := require.New(t)
			got := UniqueBy(tt.vs, tt.selector)
			r.Equal(tt.want, got)
		})
	}
}
