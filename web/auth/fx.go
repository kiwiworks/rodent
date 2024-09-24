package auth

import (
	"fmt"

	"go.uber.org/fx"
)

func AsProvider(provider any) any {
	return fx.Annotate(provider, fx.As(new(Provider)))
}

func AsNamedProvider(name string, provider any) any {
	return fx.Annotate(provider, fx.As(new(Provider)), fx.ResultTags(fmt.Sprintf("`name=\"%s\"`", name)))
}
