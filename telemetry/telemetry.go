package telemetry

import (
	"context"
	"os"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/runtime"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	trace2 "go.opentelemetry.io/otel/trace"
	"go.uber.org/multierr"

	"github.com/kiwiworks/rodent/errors"
	"github.com/kiwiworks/rodent/system/manifest"
)

type Telemetry struct {
	manifest      *manifest.Manifest
	cancel        context.CancelFunc
	traceExporter *otlptrace.Exporter
	traceProvider *trace.TracerProvider
	meterProvider *metric.MeterProvider
}

func New(manifest *manifest.Manifest) (*Telemetry, error) {
	t := &Telemetry{
		manifest: manifest,
	}
	ctx, cancel := context.WithCancel(context.Background())
	t.cancel = cancel

	propagator := newPropagator()
	otel.SetTextMapPropagator(propagator)

	traceExporter, err := newTraceExporter(ctx, os.Getenv("OTEL_EXPORTER_OTLP_PROTOCOL"))
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create trace exporter")
	}
	t.traceExporter = traceExporter
	traceProvider, err := newTracerProvider(traceExporter, t.manifest)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create trace provider")
	}
	t.traceProvider = traceProvider
	otel.SetTracerProvider(traceProvider)

	httpMetricExporter, err := newHttpMetricExporter(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create metric exporter")
	}
	meterProvider, err := newMeterProvider(httpMetricExporter)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create metric provider")
	}
	otel.SetMeterProvider(meterProvider)
	t.meterProvider = meterProvider

	if os.Getenv("OTEL_EXPORTER_OTLP_PROTOCOL") != "" {
		if err = runtime.Start(runtime.WithMinimumReadMemStatsInterval(time.Second)); err != nil {
			return nil, errors.Wrapf(err, "failed to start runtime")
		}
	}
	return t, nil
}

func (t *Telemetry) Tracer(name string) trace2.Tracer {
	tracer := t.traceProvider.Tracer(name)
	return tracer
}

func newPropagator() propagation.TextMapPropagator {
	return propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
}

func newResource(manifest *manifest.Manifest) (*resource.Resource, error) {
	return resource.Merge(resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(manifest.Application),
			semconv.ServiceVersion(manifest.Version.String()),
		),
	)
}

func newTracerProvider(
	traceExporter *otlptrace.Exporter,
	manifest *manifest.Manifest,
) (*trace.TracerProvider, error) {
	traceResource, err := newResource(manifest)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create resource")
	}
	traceProvider := trace.NewTracerProvider(
		trace.WithBatcher(traceExporter),
		trace.WithResource(traceResource),
	)
	return traceProvider, nil
}

func newTraceExporter(ctx context.Context, protocol string) (*otlptrace.Exporter, error) {
	var client otlptrace.Client
	if protocol == "" {
		protocol = "http"
	}
	switch protocol {
	case "grpc":
		client = otlptracegrpc.NewClient()
	case "http":
		client = otlptracehttp.NewClient()
	default:
		return nil, errors.Newf("unsupported protocol: %s", protocol)
	}

	traceExporter, err := otlptrace.New(ctx, client)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create trace exporter")
	}
	return traceExporter, nil
}

func newHttpMetricExporter(ctx context.Context) (*otlpmetrichttp.Exporter, error) {
	exporter, err := otlpmetrichttp.New(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create metric exporter")
	}
	return exporter, nil
}

func newMeterProvider(httpMetricExporter *otlpmetrichttp.Exporter) (*metric.MeterProvider, error) {
	meterProvider := metric.NewMeterProvider(metric.WithReader(
		metric.NewPeriodicReader(httpMetricExporter),
	))
	return meterProvider, nil
}

func (t *Telemetry) OnStop(ctx context.Context) error {
	var err error
	if os.Getenv("OTEL_EXPORTER_OTLP_PROTOCOL") != "" {
		if err = t.traceProvider.Shutdown(ctx); err != nil {
			err = multierr.Combine(err, errors.Wrapf(err, "failed to shutdown tracer provider"))
		}
		if err = t.traceExporter.Shutdown(ctx); err != nil {
			err = multierr.Combine(err, errors.Wrapf(err, "failed to shutdown trace exporter"))
		}
		if err = t.meterProvider.Shutdown(ctx); err != nil {
			err = multierr.Combine(err, errors.Wrapf(err, "failed to shutdown meter provider"))
		}
		t.cancel()
	}
	return err
}
