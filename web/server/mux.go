package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"go.uber.org/fx"

	"github.com/kiwiworks/rodent/logger"
)

//todo replace by cors config

type UseCors bool

type MuxParams struct {
	fx.In
	UseCors UseCors `optional:"true"`
}

func NewMux(params MuxParams) *chi.Mux {
	log := logger.New()

	mux := chi.NewRouter()
	if params.UseCors {
		log.Info("Using CORS")
		mux.Use(cors.Handler(cors.Options{
			AllowedOrigins:   []string{"*"},
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "X-Refresh-Token"},
			ExposedHeaders:   []string{"Link"},
			AllowCredentials: false,
			MaxAge:           300,
		}))
	}
	mux.Use(logger.ChiMiddleware())
	return mux
}
