package logger

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
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
						zap.String("proto", r.Proto),
						zap.String("method", r.Method),
						zap.String("path", r.URL.Path),
						zap.String("id", middleware.GetReqID(ctx)),
						zap.Duration("elapsedTime", time.Since(now)),
						zap.Int("status", responseWriter.Status()),
						zap.Int("size", responseWriter.BytesWritten()),
					).
					Sugar()
				switch status := responseWriter.Status(); status {
				case 200, 201, 203, 301, 304:
					log.Infof("%s %s: %d", r.Method, r.URL.Path, status)
				default:
					log.Warnf("%s %s: %d", r.Method, r.URL.Path, status)
				}
			}()
			next.ServeHTTP(responseWriter, r.WithContext(ctx))
		})
	}
}
