package slices

import "fmt"

func Of[T any](items ...T) []T {
	return items
}

func OfStringers(items ...fmt.Stringer) []string {
	return Map(items, func(in fmt.Stringer) string {
		return in.String()
	})
}
