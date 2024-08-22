package logging

import (
	"os"

	log "github.com/inconshreveable/log15"
)

const (
	logLevelEnv = "LOG_LEVEL"
	logFilePath = "/var/log/app.log"
)

var (
	localLogger log.Logger
)

type config struct {
	LogLevel   log.Lvl
	AppVersion string
	Service    string
}

func NewLogger() log.Logger {
	return localLogger
}

func init() {
	// logLevel := utils.GetEnv(logLevelEnv, "debug")
	logLevel := "debug"
	// appVersion := utils.GetEnv(utils.AppVersionEnv, "unknown")
	appVersion := "v1.0.0"
	// serviceName := utils.GetEnv(utils.ServiceNameEnv, "unknown")
	serviceName := "log15-sample-app"

	configure(config{
		LogLevel:   getLog15LogLevelFromEnv(logLevel),
		AppVersion: appVersion,
		Service:    serviceName,
	})
}

func configure(config config) {
	srvlog := log.New(
		"app_version", config.AppVersion,
		"service", config.Service,
	)

	srvlog.SetHandler(log.MultiHandler(
		log.StreamHandler(os.Stdout, log.JsonFormat()),
		log.Must.FileHandler(logFilePath, log.JsonFormat()),
	))

	localLogger = srvlog
}

func getLog15LogLevelFromEnv(level string) log.Lvl {
	switch level {
	case "debug":
		return log.LvlDebug
	case "info":
		return log.LvlInfo
	case "warn":
		return log.LvlWarn
	case "error":
		return log.LvlError
	default:
		return log.LvlDebug
	}
}
