package server

import (
	"context"
	"net/http"

	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/kiwiworks/rodent/system/errors"
	"github.com/kiwiworks/rodent/system/logger"
	"github.com/kiwiworks/rodent/web/api"
)

type (
	Server struct {
		server http.Server
	}
	Config struct {
		fx.In
		Router    *Router
		ApiConfig *api.Config    `optional:"true"`
		Handlers  []*api.Handler `group:"api.handler"`
	}
)

func New(cfg Config) *Server {
	if cfg.ApiConfig == nil {
		cfg.ApiConfig = api.DefaultConfig()
	}
	for _, handler := range cfg.Handlers {
		handler.Mount(cfg.Router.api, *cfg.ApiConfig)
	}
	server := &Server{
		server: http.Server{
			Addr:    "[::1]:8080",
			Handler: cfg.Router.mux,
		},
	}

	return server
}

func (s *Server) OnStart(ctx context.Context) error {
	log := logger.FromContext(ctx)
	go func() {
		log.Info("starting server", zap.String("address", s.server.Addr))
		if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("server crashed!", zap.Error(err))
		}
	}()
	return nil
}

func (s *Server) OnStop(ctx context.Context) error {
	log := logger.FromContext(ctx)
	log.Info("stopping server", zap.String("address", s.server.Addr))
	return s.server.Shutdown(ctx)
}
