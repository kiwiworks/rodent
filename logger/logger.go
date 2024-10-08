package logger

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/kiwiworks/rodent/logger/props"
)

var (
	once     sync.Once
	logMode  string
	logLevel zap.AtomicLevel
)

func SetLevel(level Level) {
	logLevel.SetLevel(level)
}

type (
	Level = zapcore.Level
)

const (
	DebugLevel   = zapcore.DebugLevel
	InfoLevel    = zapcore.InfoLevel
	WarningLevel = zapcore.WarnLevel
	ErrorLevel   = zapcore.ErrorLevel
	FatalLevel   = zapcore.FatalLevel
)

// do not migrate this to application OnStart hooks
// as we need a correctly configured logger ASAP
func init() {
	once.Do(func() {
		logMode = strings.ToUpper(os.Getenv("LOG_MODE"))
		if logMode == "" {
			logMode = "DEV"
		}
		logLevel = zap.NewAtomicLevelAt(zapcore.InfoLevel)
		if level := os.Getenv("LOG_LEVEL"); level != "" {
			if err := logLevel.UnmarshalText([]byte(level)); err != nil {
				panic(fmt.Sprintf("invalid log level: %s", level))
			}
		}
	})
	zap.ReplaceGlobals(New())
}

func New(opts ...Option) *zap.Logger {
	options := newLoggerOptions()
	options.apply(opts...)
	var logger *zap.Logger
	var err error

	switch logMode {
	case "PROD", "PRODUCTION":
		cfg := zap.NewProductionConfig()
		cfg.Level = logLevel
		logger, err = cfg.Build()
		cfg.DisableStacktrace = true
	case "DEV", "DEVELOPMENT":
		cfg := zap.NewDevelopmentConfig()
		cfg.Level = logLevel
		cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		cfg.DisableStacktrace = true
		logger, err = cfg.Build()
	default:
		panic(fmt.Sprintf("invalid log mode: %s, expected one of [DEV, DEVELOPMENT, PROD, PRODUCTION]", logMode))
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
