package server

import (
	"github.com/gorilla/mux"
	"go.uber.org/fx"
)

func AsMiddleware(name string, middlewareFunc mux.MiddlewareFunc) any {
	return fx.Annotate(func() *Middleware {
		return NewMiddleware(name, middlewareFunc)
	}, fx.ResultTags(`group:"mux.middleware"`))
}
