package slices

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
