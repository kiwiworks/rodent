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
	"go.opentelemetry.io/otel/trace/noop"
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

	// Set up telemetry only if OTEL endpoint is configured
	if hasOtelEndpoint() {
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

		if err = runtime.Start(runtime.WithMinimumReadMemStatsInterval(time.Second)); err != nil {
			return nil, errors.Wrapf(err, "failed to start runtime")
		}
	} else {
		// Set up no-op providers when no OTEL endpoint is configured
		noopProvider := noop.NewTracerProvider()
		t.traceProvider = nil // We'll use the noop provider via the Tracer method
		otel.SetTracerProvider(noopProvider)
		t.meterProvider = metric.NewMeterProvider()
		otel.SetMeterProvider(t.meterProvider)
	}

	return t, nil
}

func (t *Telemetry) Tracer(name string) trace2.Tracer {
	if t.traceProvider == nil {
		// If no trace provider exists, return a no-op tracer
		return noop.NewTracerProvider().Tracer(name)
	}
	return t.traceProvider.Tracer(name)
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
		opts := []otlptracegrpc.Option{}

		// Use explicitly defined endpoint if available
		if endpoint := os.Getenv("OTEL_EXPORTER_OTLP_TRACES_ENDPOINT"); endpoint != "" {
			opts = append(opts, otlptracegrpc.WithEndpoint(endpoint))
		} else if endpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT"); endpoint != "" {
			opts = append(opts, otlptracegrpc.WithEndpoint(endpoint))
		}

		client = otlptracegrpc.NewClient(opts...)
	case "http":
		opts := []otlptracehttp.Option{}

		// Use explicitly defined endpoint if available
		if endpoint := os.Getenv("OTEL_EXPORTER_OTLP_TRACES_ENDPOINT"); endpoint != "" {
			opts = append(opts, otlptracehttp.WithEndpoint(endpoint))
		} else if endpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT"); endpoint != "" {
			opts = append(opts, otlptracehttp.WithEndpoint(endpoint))
		}

		client = otlptracehttp.NewClient(opts...)
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
	// Configure the metric exporter with explicit options
	opts := []otlpmetrichttp.Option{}

	// Use explicitly defined endpoint if available
	if endpoint := os.Getenv("OTEL_EXPORTER_OTLP_METRICS_ENDPOINT"); endpoint != "" {
		opts = append(opts, otlpmetrichttp.WithEndpoint(endpoint))
	} else if endpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT"); endpoint != "" {
		opts = append(opts, otlpmetrichttp.WithEndpoint(endpoint))
	}

	exporter, err := otlpmetrichttp.New(ctx, opts...)
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

// hasOtelEndpoint checks if any OTEL endpoint environment variables are set
func hasOtelEndpoint() bool {
	// Check for standard OTEL endpoint variables
	if os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT") != "" {
		return true
	}
	if os.Getenv("OTEL_EXPORTER_OTLP_TRACES_ENDPOINT") != "" {
		return true
	}
	if os.Getenv("OTEL_EXPORTER_OTLP_METRICS_ENDPOINT") != "" {
		return true
	}

	// Check HTTP specific endpoints
	if os.Getenv("OTEL_EXPORTER_OTLP_HTTP_ENDPOINT") != "" {
		return true
	}
	if os.Getenv("OTEL_EXPORTER_OTLP_TRACES_HTTP_ENDPOINT") != "" {
		return true
	}
	if os.Getenv("OTEL_EXPORTER_OTLP_METRICS_HTTP_ENDPOINT") != "" {
		return true
	}

	// Check GRPC specific endpoints
	if os.Getenv("OTEL_EXPORTER_OTLP_GRPC_ENDPOINT") != "" {
		return true
	}
	if os.Getenv("OTEL_EXPORTER_OTLP_TRACES_GRPC_ENDPOINT") != "" {
		return true
	}
	if os.Getenv("OTEL_EXPORTER_OTLP_METRICS_GRPC_ENDPOINT") != "" {
		return true
	}

	// If protocol is set but no endpoints, don't enable telemetry by default
	// as that leads to connection errors with default localhost endpoints
	return false
}

func (t *Telemetry) OnStop(ctx context.Context) error {
	var err error
	if hasOtelEndpoint() && t.traceExporter != nil {
		if t.traceProvider != nil {
			if err = t.traceProvider.Shutdown(ctx); err != nil {
				err = multierr.Combine(err, errors.Wrapf(err, "failed to shutdown tracer provider"))
			}
		}
		if err = t.traceExporter.Shutdown(ctx); err != nil {
			err = multierr.Combine(err, errors.Wrapf(err, "failed to shutdown trace exporter"))
		}
		if t.meterProvider != nil {
			if err = t.meterProvider.Shutdown(ctx); err != nil {
				err = multierr.Combine(err, errors.Wrapf(err, "failed to shutdown meter provider"))
			}
		}
		t.cancel()
	}
	return err
}
