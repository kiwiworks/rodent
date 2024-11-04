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

type UseCors bool

type MuxParams struct {
	fx.In
	UseCors UseCors `optional:"true"`
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
