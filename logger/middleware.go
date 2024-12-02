package logger

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"

	"github.com/kiwiworks/rodent/logger/props"
)

func ChiMiddleware() func(next http.Handler) http.Handler {
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
						props.HttpUserAgent(r.UserAgent()),
						zap.Duration("elapsedTime", time.Since(now)),
					)
				switch status := responseWriter.Status(); status {
				case 200, 201, 203, 301, 304:
					log.Info("ok")
				default:
					log.Warn("error")
				}
			}()
			next.ServeHTTP(responseWriter, r.WithContext(ctx))
		})
	}
}
