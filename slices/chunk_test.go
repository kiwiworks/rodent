package slices

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestChunk(t *testing.T) {
	type testCase struct {
		input     []int
		chunkSize int
		expected  [][]int
	}

	testCases := []testCase{
		{
			input:     []int{1, 2, 3, 4, 5, 6},
			chunkSize: 2,
			expected:  [][]int{{1, 2}, {3, 4}, {5, 6}},
		},
		{
			input:     []int{1, 2, 3, 4, 5, 6},
			chunkSize: 3,
			expected:  [][]int{{1, 2, 3}, {4, 5, 6}},
		},
		{
			input:     []int{1, 2},
			chunkSize: 5,
			expected:  [][]int{{1, 2}},
		},
		{
			input:     []int{},
			chunkSize: 3,
			expected:  [][]int{},
		},
		{
			input:     []int{1, 2, 3, 4, 5, 6},
			chunkSize: 0,
			expected:  [][]int{},
		},
	}

	r := require.New(t)
	for _, tc := range testCases {
		result := Chunk(tc.input, tc.chunkSize)
		r.Equal(tc.expected, result, "Chunk(%v, %d)", tc.input, tc.chunkSize)
	}
}
