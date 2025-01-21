package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"

	"github.com/kiwiworks/rodent/logger"
)

//todo(mrkiwi): replace by cors config

// UseCors define if the server should use CORS
// Deprecated: use CorsConfig instead
type UseCors bool
type CorsConfig struct {
	AllowedOrigins   []string
	AllowedMethods   []string
	AllowedHeaders   []string
	ExposedHeaders   []string
	AllowCredentials bool
	MaxAge           int
}

func WithCors(config CorsConfig) *CorsConfig {
	return &config
}

type MuxParams struct {
	fx.In
	UseCors UseCors     `optional:"true"`
	Cors    *CorsConfig `optional:"true"`
}

func NewMux(params MuxParams) *chi.Mux {
	log := logger.New()

	mux := chi.NewRouter()
	mux.Use(middleware.RequestID)
	mux.Use(middleware.RealIP)
	mux.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := chi.NewRouteContext()
			routePattern := ""
			if mux.Match(ctx, r.Method, r.URL.Path) {
				routePattern = ctx.RoutePattern()
			}
			instrumentedHandler := otelhttp.NewHandler(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
				span := trace.SpanFromContext(request.Context())
				span.SetAttributes(
					semconv.HTTPRoute(routePattern),
				)
				next.ServeHTTP(writer, request)
			}), routePattern)
			instrumentedHandler.ServeHTTP(w, r)
		})
	})
	if params.UseCors && params.Cors == nil {
		log.Info("Using CORS with default configuration (deprecated)")
		mux.Use(cors.Handler(cors.Options{
			AllowedOrigins:   []string{"*"},
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "X-Refresh-Token", "X-Request-Id"},
			ExposedHeaders:   []string{"Link"},
			AllowCredentials: false,
			MaxAge:           300,
		}))
	}
	if params.Cors != nil {
		log.Info("Using CORS")
		mux.Use(cors.Handler(cors.Options{
			AllowedOrigins:   params.Cors.AllowedOrigins,
			AllowedMethods:   params.Cors.AllowedMethods,
			AllowedHeaders:   params.Cors.AllowedHeaders,
			ExposedHeaders:   params.Cors.ExposedHeaders,
			AllowCredentials: params.Cors.AllowCredentials,
			MaxAge:           params.Cors.MaxAge,
		}))
	}

	mux.Use(logger.ChiMiddleware())

	return mux
}
