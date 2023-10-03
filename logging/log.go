package logging

import (
	"log/slog"
	"os"
	"sync"

	"github.com/kogan69/goinfra/config"
)

var loggerOnce sync.Once
var errorReporter Reporter
var logger *slog.Logger

func Init(name, version, environment, logLevel string, reporter Reporter) {
	loggerOnce.Do(func() {
		initLogger(name, version, environment, logLevel, reporter)
	})
}

func InitFromEnv() {
	name := config.GetEnvValue("APP_NAME", "unknown")
	version := config.GetEnvValue("APP_VERSION", "unknown")
	env := config.GetEnvValue("ENV", "development")
	logLevel := config.GetEnvValue("LOG_LEVEL", "info")
	reporter := NewFromEnv(env, name)
	Init(name, version, env, logLevel, reporter)
}
func initLogger(
	name string,
	version string,
	environment string,
	level string,
	reporter Reporter,
) {
	levelVar := &slog.LevelVar{}
	err := levelVar.UnmarshalText([]byte(level))
	if err != nil {
		panic(err)
	}

	handler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: levelVar}).WithAttrs([]slog.Attr{
		slog.String("env", environment),
		slog.String("name", name),
		slog.String("version", version),
	})

	logger = slog.New(handler)

	errorReporter = reporter
}

func Debug(message string, args ...any) {
	logger.Debug(message, args...)
}

func Info(message string, args ...any) {
	logger.Info(message, args...)
}
func Warn(message string, args ...any) {
	logger.Warn(message, args...)
}
func Error(err error) {
	logger.Error(err.Error())
	if errorReporter != nil {
		errorReporter.ReportFatal(err)
	}
}
