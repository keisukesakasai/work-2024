package logging

import (
	"context"
	"datadog-log-sample-app/internal/utils"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	logLevelEnv = "LOG_LEVEL"
	logFilePath = "/var/log/app.log"
)

type contextKeyLoggerKey int

const (
	contextKeyLogger contextKeyLoggerKey = iota
)

var (
	localLogger *zap.SugaredLogger
)

type config struct {
	LogLevel   zapcore.Level
	AppVersion string
	Service    string
}

func NewLogger() *zap.SugaredLogger {
	return localLogger
}

func init() {
	logLevel := utils.GetEnv(logLevelEnv, "debug")
	appVersion := utils.GetEnv(utils.AppVersionEnv, "unknown")
	serviceName := utils.GetEnv(utils.ServiceNameEnv, "unknown")

	configure(config{
		LogLevel:   getZapLogLevelFromEnv(logLevel),
		AppVersion: appVersion,
		Service:    serviceName,
	})
}

func configure(config config) {
	zapConfig := defaultZapConfig()

	logger, err := zapConfig.Build()
	if err != nil {
		panic(err)
	}
	fields := zap.Fields([]zap.Field{
		zap.String(utils.AppVersionKey, config.AppVersion),
		zap.String(utils.ServiceNameKey, config.Service),
	}...)

	localLogger = logger.WithOptions(fields).Sugar()
}

func defaultZapConfig() zap.Config {
	return zap.Config{
		Level:    zap.NewAtomicLevelAt(zapcore.DebugLevel),
		Encoding: "json",
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:     "message",
			LevelKey:       "severity",
			TimeKey:        "time",
			NameKey:        "logger",
			CallerKey:      "caller",
			StacktraceKey:  "stacktrace",
			EncodeLevel:    zapcore.CapitalLevelEncoder,
			EncodeTime:     zapcore.RFC3339NanoTimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		OutputPaths:      []string{"stdout", logFilePath}, // 標準出力とファイルの両方に出力
		ErrorOutputPaths: []string{"stderr"},
	}
}

func GetLoggerFromCtx(ctx context.Context) *zap.SugaredLogger {
	logger, ok := ctx.Value(contextKeyLogger).(*zap.SugaredLogger)
	if ok {
		return logger
	}

	return NewLogger()
}

func SetLoggerToCtx(ctx context.Context, logger *zap.SugaredLogger) context.Context {
	return context.WithValue(ctx, contextKeyLogger, logger)
}
