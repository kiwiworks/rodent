package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"

	"github.com/kiwiworks/rodent/logger"
)

func NewMux() *chi.Mux {
	mux := chi.NewRouter()
	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "X-Refresh-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))
	mux.Use(logger.Middleware())
	return mux
}
