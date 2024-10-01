package lang

func Default[T comparable](value T, defaultValue T) T {
	var zero T
	if value == zero {
		return defaultValue
	}
	return value
}
