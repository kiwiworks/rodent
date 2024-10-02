package slices

// Chunk divides a slice of elements into multiple smaller slices, each of a specified size.
//
// The function takes a generic type parameter `T` which allows it to work with slices of any type.
//
// Parameters:
//   - elements: A slice of elements of type `T` to be divided into chunks.
//   - chunkSize: An integer representing the size of each chunk. If chunkSize is less than or equal to 0,
//     an empty slice is returned.
//
// Returns:
//   - A two-dimensional slice of type `T` where each inner slice is of length `chunkSize` (except possibly the last one).
//     If chunkSize is less than or equal to 0, an empty slice of slices is returned. If the input slice is empty,
//     an empty slice of slices is returned.
//
// Example usage:
//
//	chunks := Chunk([]int{1, 2, 3, 4, 5, 6}, 2)
//	// chunks will be [][]int{{1, 2}, {3, 4}, {5, 6}}
//
//	chunks = Chunk([]int{1, 2, 3, 4, 5, 6}, 3)
//	// chunks will be [][]int{{1, 2, 3}, {4, 5, 6}}
//
//	chunks = Chunk([]int{1, 2, 3, 4, 5, 6}, 0)
//	// chunks will be [][]int{}
//
//	chunks = Chunk([]int{}, 3)
//	// chunks will be [][]int{}
func Chunk[T any](elements []T, chunkSize int) [][]T {
	if chunkSize <= 0 {
		return [][]T{}
	}

	capacity := (len(elements) + chunkSize - 1) / chunkSize
	chunks := make([][]T, capacity)

	for i, chunkIndex := 0, 0; i < len(elements); i += chunkSize {
		end := i + chunkSize
		isEndExceeding := end > len(elements)
		if isEndExceeding {
			end = len(elements)
		}
		chunks[chunkIndex] = elements[i:end]
		chunkIndex++
	}
	return chunks
}
