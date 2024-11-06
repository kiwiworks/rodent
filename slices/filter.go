package slices

// Filter returns a new slice containing elements from the input slice that satisfy
// the given predicate function. The function checks each element in the input slice
// and appends it to the output slice if the predicate returns true for that element.
// The input slice remains unchanged.
//
// The type parameter T represents the type of elements in the input and output slices.
func Filter[T any](vs []T, f func(T) bool) []T {
	out := make([]T, 0)
	for _, v := range vs {
		if f(v) {
			out = append(out, v)
		}
	}
	return out
}

// FilterMap filters and transforms a slice of elements of type `From` to a slice of elements of type `To`.
// The transformation and filtering are specified by the provided function `f`.
//
// The function `f` takes an element of type `From` and returns a tuple `(To, bool)`.
// If the boolean value is `true`, the transformed value is included in the result slice.
// If the boolean value is `false`, the transformed value is excluded from the result slice.
//
// BlockType Parameters:
// - From: the type of the elements in the input slice.
// - To: the type of the elements in the output slice.
//
// Parameters:
// - vs: the input slice containing elements of type `From`.
// - f: the transformation and filtering function.
//
// Returns:
// - A slice of elements of type `To` that satisfies the filter condition.
//
// Example:
//
//	func isNonEmptyString(s string) (string, bool) {
//	    if s != "" {
//	        return s, true
//	    }
//	    return "", false
//	}
//
// input := []string{"a", "", "b", "c", ""}
// result := FilterMap(input, isNonEmptyString) // result: []string{"a", "b", "c"}
func FilterMap[From any, To any](vs []From, f func(From) (To, bool)) []To {
	out := make([]To, 0)
	for _, v := range vs {
		if val, ok := f(v); ok {
			out = append(out, val)
		}
	}
	return out
}

// UniqueBy returns a new slice that contains only the unique elements from the input slice `vs`
// based on a specified key selected by the `selector` function. It uses the provided `selector`
// function to determine the key for comparison.
//
// This function leverages the `Filter` function to iterate through the input slice `vs`
// and include elements in the output slice only if their selector key is unique within the slice.
//
// BlockType Parameters:
//   - T: The type of elements in the input and output slices.
//   - V: The type of the key returned by the `selector` function for comparison; it must be comparable.
//
// Parameters:
//   - vs: A slice of elements of type T to be filtered for uniqueness.
//   - selector: A function that takes an element of type T and returns a value of type V,
//     which is used for determining uniqueness.
//
// Returns:
//
//	A new slice containing only the unique elements from the input slice `vs`,
//	based on the keys provided by the `selector` function.
//
// Example usage:
//
//	items := []Item{{ID: 1}, {ID: 2}, {ID: 1}}
//	uniqueItems := UniqueBy(items, func(i Item) int { return i.ID })
//	// uniqueItems will contain [{ID: 1}, {ID: 2}]
func UniqueBy[T any, V comparable](vs []T, selector func(v T) V) []T {
	seen := make(map[V]struct{})
	result := make([]T, 0, len(vs))

	for _, v := range vs {
		key := selector(v)
		if _, exists := seen[key]; !exists {
			seen[key] = struct{}{}
			result = append(result, v)
		}
	}
	return result
}

// Ident returns the input value without modifications.
// It is a generic identity function that works for any type.
func Ident[T any](t T) T {
	return t
}

// Unique returns a new slice that contains only the unique elements from the input slice `vs`.
// It uses the Ident function as the selector to determine uniqueness.
func Unique[T comparable](vs []T) []T {
	return UniqueBy[T](vs, Ident[T])
}
