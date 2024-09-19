package module

import (
	"fmt"

	"go.uber.org/fx"
	"go.uber.org/zap"

	"rodent/system/logger"
	"rodent/system/opt"
)

func registerLifecycle[T OnStartStop](impl T, lifecycle fx.Lifecycle) {
	lifecycle.Append(fx.Hook{
		OnStart: impl.OnStart,
		OnStop:  impl.OnStop,
	})
}

func registerLifecycleStart[T OnStart](impl T, lifecycle fx.Lifecycle) {
	lifecycle.Append(fx.Hook{
		OnStart: impl.OnStart,
	})
}

func registerLifecycleStop[T OnStop](impl T, lifecycle fx.Lifecycle) {
	lifecycle.Append(fx.Hook{
		OnStop: impl.OnStop,
	})
}

func invoke[T any](*T) {
	log := logger.New()
	log.Info("service available", zap.String("service", fmt.Sprintf("%T", (*T)(nil))))
}

func Service[T any]() opt.Option[Module] {
	concrete := new(T)
	asAny := any(concrete)
	switch asAny.(type) {
	case OnStartStop:
		return func(opt *Module) {
			opt.Decorators = append(opt.Decorators, func(impl *T, lifecycle fx.Lifecycle) *T {
				asAny := any(impl)
				asService := asAny.(OnStartStop)
				registerLifecycle(asService, lifecycle)
				return impl
			})
			opt.Invokers = append(opt.Invokers, invoke[T])
		}
	case OnStart:
		return func(opt *Module) {
			opt.Decorators = append(opt.Decorators, func(impl *T, lifecycle fx.Lifecycle) *T {
				asAny := any(impl)
				asService := asAny.(OnStart)
				registerLifecycleStart(asService, lifecycle)
				return impl
			})
			opt.Invokers = append(opt.Invokers, invoke[T])
		}
	case OnStop:
		return func(opt *Module) {
			opt.Decorators = append(opt.Decorators, func(impl *T, lifecycle fx.Lifecycle) *T {
				asAny := any(impl)
				asService := asAny.(OnStop)
				registerLifecycleStop(asService, lifecycle)
				return impl
			})
			opt.Invokers = append(opt.Invokers, invoke[T])
		}
	default:
		panic(fmt.Sprintf("%T is not a valid Service, it should at least implement one of `app.OnStart` or `app.OnStop`", concrete))
	}
}
