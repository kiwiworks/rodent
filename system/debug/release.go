//go:build !debug

package debug

import "context"

//go:inline
func Do(func()) {}

//go:inline
func DoIf(bool, func()) {}

//go:inline
func DoWith(context.Context, func(context.Context)) {}
