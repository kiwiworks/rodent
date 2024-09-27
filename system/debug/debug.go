//go:build debug

package debug

import (
	"context"
	"fmt"
	"runtime"
)

func printLocation() {
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		panic("Unable to get caller information")
	}
	fmt.Printf("%s:%d\n", file, line)
}

// Do execute the provided function `fn` and prints the caller's file and line number.
// This is only available in debug builds.
//
//go:inline
func Do(fn func()) {
	printLocation()
	fn()
}

// DoIf executes the provided function if the condition is true and prints the call location.
// This is only available in debug builds.
//
//go:inline
func DoIf(cond bool, fn func()) {
	if cond {
		printLocation()
		fn()
	}
}

// DoWith calls printLocation to output the caller's file and line number, then executes the provided function with the context.
// This is only available in debug builds.
//
//go:inline
func DoWith(ctx context.Context, fn func(ctx context.Context)) {
	printLocation()
	fn(ctx)
}
