package props

import (
	"time"

	"go.uber.org/zap"
)

func AppName(name string) zap.Field {
	return zap.String("app.name", name)
}

func AppVersion(version string) zap.Field {
	return zap.String("app.version", version)
}

func AppStartTimeout(timeout time.Duration) zap.Field {
	return zap.Duration("app.timeouts.start", timeout)
}

func AppStopTimeout(timeout time.Duration) zap.Field {
	return zap.Duration("app.timeouts.stop", timeout)
}
