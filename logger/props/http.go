package props

import (
	"net/url"

	"go.uber.org/zap"
)

func HttpMethod(method string) zap.Field {
	return zap.String("http.method", method)
}

func HttpPath(path string) zap.Field {
	return zap.String("http.path", path)
}

func HttpUserAgent(userAgent string) zap.Field {
	return zap.String("http.user_agent", userAgent)
}

func HttpStatusCode(statusCode int) zap.Field {
	return zap.Int("http.status_code", statusCode)
}

func HttpProtocol(protocol string) zap.Field {
	return zap.String("http.protocol", protocol)
}

func HttpRequestID(requestID string) zap.Field {
	return zap.String("http.request.id", requestID)
}

func HttpContentLength(requestSize int64) zap.Field {
	return zap.Int64("http.request.size", requestSize)
}

func HttpResponseSize(responseSize int) zap.Field {
	return zap.Int("http.response.size", responseSize)
}

func HttpRequestUrl(url *url.URL) zap.Field {
	return zap.String("http.request.url", url.String())
}
