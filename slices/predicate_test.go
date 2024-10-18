package slices

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIntersects(t *testing.T) {
	type testCase struct {
		slice  []int
		values []int
		want   bool
	}

	testCases := []testCase{
		{slice: []int{1, 2, 3, 4}, values: []int{5, 6, 7, 8}, want: false},
		{slice: []int{1, 2, 3, 4}, values: []int{3, 6, 7, 8}, want: true},
		{slice: []int{1, 2, 3, 4}, values: []int{}, want: false},
		{slice: []int{}, values: []int{3, 6, 7, 8}, want: false},
		{slice: []int{}, values: []int{}, want: false},
	}

	for _, tc := range testCases {
		got := Intersects(tc.slice, tc.values)
		require.Equal(t, tc.want, got)
	}
}

func TestContains(t *testing.T) {
	type testCase struct {
		slice []int
		value int
		want  bool
	}

	testCases := []testCase{
		{slice: []int{1, 2, 3, 4}, value: 1, want: true},
		{slice: []int{1, 2, 3, 4}, value: 5, want: false},
		{slice: []int{}, value: 3, want: false},
		{slice: []int{2, 4}, value: 4, want: true},
		{slice: []int{7}, value: 7, want: true},
	}

	for _, tc := range testCases {
		got := Contains(tc.slice, tc.value)
		require.Equal(t, tc.want, got)
	}
}
