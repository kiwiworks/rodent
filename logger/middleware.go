package logger

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"

	"github.com/kiwiworks/rodent/logger/props"
)

func Middleware() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			responseWriter := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			now := time.Now()
			defer func() {
				log := FromContext(ctx).
					With(
						props.HttpProtocol(r.Proto),
						props.HttpMethod(r.Method),
						props.HttpPath(r.URL.Path),
						props.HttpRequestID(middleware.GetReqID(ctx)),
						props.HttpContentLength(r.ContentLength),
						props.HttpResponseSize(responseWriter.BytesWritten()),
						props.HttpStatusCode(responseWriter.Status()),
						zap.Duration("elapsedTime", time.Since(now)),
					)
				switch status := responseWriter.Status(); status {
				case 200, 201, 203, 301, 304:
				default:
					log.Warn("http request error")
				}
			}()
			next.ServeHTTP(responseWriter, r.WithContext(ctx))
		})
	}
}
