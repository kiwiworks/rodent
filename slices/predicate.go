package slices

// Contains checks if a given value exists within a slice.
// The function is generic and can operate on slices of any type
// that supports comparison.
//
// Parameters:
//   - slice: The slice of elements to be searched.
//   - value: The value to search for within the slice.
//
// Returns:
//   - bool: Returns true if the value is found within the slice,
//     otherwise returns false.
//
// Example usage:
//   - Contains([]int{1, 2, 3, 4}, 3) // returns true
//   - Contains([]string{"apple", "banana", "cherry"}, "banana") // returns true
//   - Contains([]string{"apple", "banana", "cherry"}, "date") // returns false
//
// Type constraints:
//   - T: The type of elements in the slice, which must be comparable.
func Contains[T comparable](slice []T, value T) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}
