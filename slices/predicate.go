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

// Intersects determines if there is any intersection between two slices.
// The function is generic and can operate on slices of any type that supports comparison.
//
// Parameters:
//   - slice: The first slice of elements.
//   - values: The second slice of elements to be compared against the first slice.
//
// Returns:
//   - bool: Returns true if any element in the `values` slice exists within the `slice`,
//     otherwise returns false.
//
// Example usage:
//   - Intersects([]int{1, 2, 3, 4}, []int{3, 5}) // returns true
//   - Intersects([]string{"apple", "banana", "cherry"}, []string{"banana", "date"}) // returns true
//   - Intersects([]string{"apple", "banana", "cherry"}, []string{"date", "fig"}) // returns false
//
// Type constraints:
//   - T: The type of elements in the slices, which must be comparable.
func Intersects[T comparable](slice []T, values []T) bool {
	for _, v := range values {
		if Contains(slice, v) {
			return true
		}
	}
	return false
}

// All checks if all elements in a slice satisfy a given predicate function.
// The function is generic and can operate on slices of any type.
//
// Parameters:
//   - slice: The slice of elements to be checked.
//   - predicate: A function that takes an element and returns a boolean.
//
// Returns:
//   - bool: Returns true if all elements satisfy the predicate,
//     otherwise returns false.
//
// Example usage:
//   - All([]int{2, 4, 6}, func(x int) bool { return x%2 == 0 }) // returns true
//   - All([]int{2, 3, 4}, func(x int) bool { return x%2 == 0 }) // returns false
//
// Type constraints:
//   - T: The type of elements in the slice.
func All[T any](slice []T, predicate func(T) bool) bool {
	for _, v := range slice {
		if !predicate(v) {
			return false
		}
	}
	return true
}

// Any checks if any element in a slice satisfies a given predicate function.
// The function is generic and can operate on slices of any type.
//
// Parameters:
//   - slice: The slice of elements to be checked.
//   - predicate: A function that takes an element and returns a boolean.
//
// Returns:
//   - bool: Returns true if any element satisfies the predicate,
//     otherwise returns false.
//
// Example usage:
//   - Any([]int{1, 2, 3}, func(x int) bool { return x == 2 }) // returns true
//   - Any([]int{1, 3, 5}, func(x int) bool { return x%2 == 0 }) // returns false
//
// Type constraints:
//   - T: The type of elements in the slice.
func Any[T any](slice []T, predicate func(T) bool) bool {
	for _, v := range slice {
		if predicate(v) {
			return true
		}
	}
	return false
}
