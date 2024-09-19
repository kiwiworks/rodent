package slices

// Map applies a given mapper function to all elements of the input slice
// and returns a new slice containing the results.
//
// Parameters:
//   - items: A slice of elements of type From, representing the input elements.
//   - mapper: A function that takes an element of type From and returns an element of type To,
//     representing the transformation to be applied to each element of the input slice.
//
// Returns:
//   - A new slice of elements of type To, containing the results of applying the mapper function
//     to each element of the input slice.
//
// Example:
//
//	transformed := Map([]int{1, 2, 3}, func(in int) string {
//	    return fmt.Sprintf("Number %d", in)
//	})
//	// transformed would be []string{"Number 1", "Number 2", "Number 3"}
func Map[From any, To any](items []From, mapper func(in From) To) []To {
	out := make([]To, len(items))
	for idx, item := range items {
		out[idx] = mapper(item)
	}
	return out
}

// TryMap applies a given mapper function to all elements of the input slice and returns a new slice containing the results.
// Similar to Map function, TryMap allows the mapper function to return a pointer to the transformed element along with an error.
// If any error occurs during the mapping process, TryMap stops and returns the current result along with the error.
//
// Parameters:
//   - items: A slice of elements of type From, representing the input elements.
//   - mapper: A function that takes an element of type From and returns a pointer to an element of type To,
//     along with an error. This function represents the transformation to be applied to each element of the input slice.
//
// Returns:
//   - A new slice of elements of type To, containing the results of applying the mapper function to each element of the input slice.
//   - If any error occurs during the mapping process, TryMap returns the current result along with the error.
//
// Example:
//
//	transformed, err := TryMap([]int{1, 2, 3}, func(in int) (*string, error) {
//	    str := fmt.Sprintf("Number %d", in)
//	    return &str, nil
//	})
//	// transformed would be []string{"Number 1", "Number 2", "Number 3"}, and err would be nil
func TryMap[From any, To any](items []From, mapper func(in From) (*To, error)) ([]To, error) {
	out := make([]To, len(items))
	for idx, item := range items {
		result, err := mapper(item)
		if err != nil {
			return out, err
		}
		out[idx] = *result
	}
	return out, nil
}
