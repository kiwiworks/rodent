package logger

import (
	"context"
	"fmt"
	"net/http"

	"github.com/sanity-io/litter"
	"go.uber.org/atomic"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func init() {
	litter.Config.HideZeroValues = false
	litter.Config.HidePrivateFields = false
}

var (
	globalLevel = atomic.NewInt32(int32(LDebug))
)

func SetLevel(level Level) {
	globalLevel.Store(int32(level))
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
	logMode := "DEV"
	var logger *zap.Logger
	var err error

	switch logMode {
	case "PROD":
		cfg := zap.NewProductionConfig()
		logger, err = cfg.Build()
		cfg.DisableStacktrace = true
	case "DEV":
		cfg := zap.NewDevelopmentConfig()
		cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		cfg.DisableStacktrace = true
		logger, err = cfg.Build()
	default:
		cfg := zap.NewDevelopmentConfig()
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
		zap.String("http.method", req.Method),
		zap.String("http.url", req.URL.String()),
	))
	return FromContext(req.Context(), opts...)
}

func Flush() {
	_ = zap.L().Sync()
}
