package logging

import (
	"log/slog"
	"os"
	"sync"
)

var loggerOnce sync.Once
var errorReporter Reporter
var logger *slog.Logger

func Configure(name, version, environment, logLevel string, reporter Reporter) {
	loggerOnce.Do(func() {
		initLogger(name, version, environment, logLevel, reporter)
	})
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
