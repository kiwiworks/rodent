package logger

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/sanity-io/litter"
	"go.uber.org/atomic"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/kiwiworks/rodent/logger/props"
)

func init() {
	litter.Config.HideZeroValues = false
	litter.Config.HidePrivateFields = false
}

var (
	legacy   = atomic.NewInt32(int32(LDebug))
	logLevel = zap.NewAtomicLevelAt(LDebug)
)

func SetLevel(level Level) {
	legacy.Store(int32(level))
}

type (
	Level = zapcore.Level
)

const (
	LDebug   = zapcore.DebugLevel
	LInfo    = zapcore.InfoLevel
	LWarning = zapcore.WarnLevel
	LError   = zapcore.ErrorLevel
	LFatal   = zapcore.FatalLevel
)

// do not migrate this to application OnStart hooks
// as we need a correctly configured logger ASAP
func init() {
	zap.ReplaceGlobals(New())
}

func New(opts ...Option) *zap.Logger {
	options := newLoggerOptions()
	options.apply(opts...)
	logMode := strings.ToUpper(os.Getenv("LOG_MODE"))
	var logger *zap.Logger
	var err error

	switch logMode {
	case "PROD":
		cfg := zap.NewProductionConfig()
		cfg.Level = logLevel
		logger, err = cfg.Build()
		cfg.DisableStacktrace = true
	case "DEV":
		cfg := zap.NewDevelopmentConfig()
		cfg.Level = logLevel
		cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		cfg.DisableStacktrace = true
		logger, err = cfg.Build()
	default:
		cfg := zap.NewDevelopmentConfig()
		cfg.Level = logLevel
		cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		cfg.DisableStacktrace = true
		logger, err = cfg.Build()
	}
	if err != nil {
		// we want to panic if no *zap.Logger instance can be acquired
		// because if you can't even open a TTY, running a full application is completely out of the question
		// so yeah, panic it is
		panic(fmt.Sprintf("logger could not be initialised: %+v", err))
	}
	return logger
}

func ctxDecorator(ctx context.Context) []Option {
	opts := make([]Option, 0)
	deadline, ok := ctx.Deadline()
	if ok {
		opts = append(opts,
			Decorate(func(s string) string {
				return fmt.Sprintf("%s!", s)
			}),
			Fields(zap.Time("context.deadline", deadline)),
		)
	}

	opts = append(opts, SkipCallFrame(2))
	return opts
}

func FromContext(ctx context.Context, opts ...Option) *zap.Logger {
	opts = append(opts, ctxDecorator(ctx)...)
	return New(opts...)
}

func FromRequest(req *http.Request, opts ...Option) *zap.Logger {
	opts = append(opts, Fields(
		props.HttpMethod(req.Method),
		props.HttpRequestUrl(req.URL),
	))
	return FromContext(req.Context(), opts...)
}

func Flush() {
	_ = zap.L().Sync()
}
