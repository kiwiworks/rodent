package module

import (
	"context"

	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/kiwiworks/rodent/errors"
	"github.com/kiwiworks/rodent/logger"
	"github.com/kiwiworks/rodent/system/manifest"
)

type ServiceRegistry struct {
	services map[string]HealthCheck
	timeouts *manifest.Timeouts
}

type ServiceRegistryParams struct {
	fx.In
	Services []HealthCheck
	Manifest *manifest.Manifest
}

func NewServiceRegistry(params ServiceRegistryParams) (*ServiceRegistry, error) {
	services := make(map[string]HealthCheck)
	for _, service := range params.Services {
		inspect := service.Inspect()
		if _, alreadyExists := services[inspect.Name]; alreadyExists {
			return nil, errors.Newf("duplicate service name: %s", inspect.Name)
		}
		services[inspect.Name] = service
	}
	return &ServiceRegistry{
		services: services,
		timeouts: &params.Manifest.Timeouts,
	}, nil
}

func (s *ServiceRegistry) stopService(ctx context.Context, healthCheck HealthCheck) error {
	shutdownCtx, cancel := context.WithTimeout(ctx, s.timeouts.Start)
	defer cancel()
	if err := healthCheck.OnStop(shutdownCtx); err != nil {
		return errors.Wrapf(err, "failed to stop service '%s'", healthCheck.Inspect().Name)
	}
	return nil
}

func (s *ServiceRegistry) startService(ctx context.Context, healthCheck HealthCheck) error {
	startupCtx, cancel := context.WithTimeout(ctx, s.timeouts.Start)
	defer cancel()
	if err := healthCheck.OnStart(startupCtx); err != nil {
		return errors.Wrapf(err, "failed to start service '%s'", healthCheck.Inspect().Name)
	}
	return nil
}

func (s *ServiceRegistry) RestartService(ctx context.Context, name string) error {
	log := logger.FromContext(ctx)
	service, ok := s.services[name]
	if !ok {
		return errors.Newf("service '%s' not found", name)
	}
	if err := s.stopService(ctx, service); err != nil {
		log.Error("failed to stop service", zap.Error(err))
	}
	if err := s.startService(ctx, service); err != nil {
		log.Error("failed to start service", zap.Error(err))
		return err
	}
	return nil
}

func (s *ServiceRegistry) OnStart(context.Context) error {
	return nil
}

func (s *ServiceRegistry) OnStop(context.Context) error {
	return nil
}
