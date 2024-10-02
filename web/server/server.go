package server

import (
	"context"
	"net/http"

	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/kiwiworks/rodent/errors"
	"github.com/kiwiworks/rodent/logger"
	"github.com/kiwiworks/rodent/web/api"
)

type (
	Addr   string
	Server struct {
		server http.Server
		router *Router
	}
	Params struct {
		fx.In
		Addr      Addr `optional:"true"`
		Router    *Router
		ApiConfig *api.Config    `optional:"true"`
		Handlers  []*api.Handler `group:"api.handler"`
	}
)

func New(params Params) *Server {
	if params.ApiConfig == nil {
		params.ApiConfig = api.DefaultConfig()
	}
	for _, handler := range params.Handlers {
		handler.Mount(params.Router.api, *params.ApiConfig)
	}
	addr := params.Addr
	if addr == "" {
		addr = "[::1]:8080"
	}
	server := &Server{
		server: http.Server{
			Addr:    string(addr),
			Handler: params.Router.mux,
		},
		router: params.Router,
	}

	return server
}

func (s *Server) OnStart(ctx context.Context) error {
	log := logger.FromContext(ctx)
	s.sanityCheck(ctx)
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
