package auth

import (
	"go.uber.org/fx"
)

func AsProvider(provider any) any {
	return fx.Annotate(provider, fx.As(new(Provider)), fx.ResultTags(`group:"auth.provider"`))
}
