package slices

// ToMap transforms a slice into a map by applying a provided transformation function to each element.
//
// The transformation function `fn` is called for each element in the `slice`, with its index and value.
// It returns a key-value pair that will be used to populate the resulting map.
//
// The generic type parameters are:
// - `T`: The type of elements in the input slice.
// - `K`: The type of keys in the resulting map. Must be comparable.
// - `V`: The type of values in the resulting map.
//
// Parameters:
// - `slice`: A slice of elements to transform.
// - `fn`: A function that takes an element's index and the element itself, and returns a key-value pair.
//
// Returns:
// A map of type `map[K]V` where each key-value pair is obtained by applying `fn` to each element of `slice`.
//
// Example:
//
//	// Example of transforming a slice of strings into a map where the key is the string and the value is its length:
//	words := []string{"apple", "banana", "cherry"}
//	result := ToMap(words, func(idx int, elem string) (string, int) {
//		return elem, len(elem)
//	})
//	// result is map[string]int{"apple": 5, "banana": 6, "cherry": 6}
func ToMap[T any, K comparable, V any](
	slice []T,
	fn func(idx int, elem T) (K, V),
) map[K]V {
	out := make(map[K]V, len(slice))
	for i, elem := range slice {
		k, v := fn(i, elem)
		out[k] = v
	}
	return out
}
