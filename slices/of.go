package slices

import "fmt"

// Of takes a variadic number of items of any type and returns them as a slice.
func Of[T any](items ...T) []T {
	return items
}

// OfStringers converts a variadic list of fmt.Stringer items to a slice of strings by calling String() on each item.
func OfStringers(items ...fmt.Stringer) []string {
	return Map(items, func(in fmt.Stringer) string {
		return in.String()
	})
}
