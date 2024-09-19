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

//go:inline
func Do(fn func()) {
	printLocation()
	fn()
}

//go:inline
func DoIf(cond bool, fn func()) {
	if cond {
		printLocation()
		fn()
	}
}

//go:inline
func DoWith(ctx context.Context, fn func(ctx context.Context)) {
	printLocation()
	fn(ctx)
}
