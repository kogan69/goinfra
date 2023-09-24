package logging

import (
	"strings"

	"github.com/kogan69/goinfra/config"
)

type Reporter interface {
	ReportException(err error)
	ReportFatal(err error)
}

func New(env string, name string, dsn string) Reporter {
	if dsn != "" && strings.ToLower(env) == "production" {
		return NewSentryErrorReporter(env, name, dsn)
	}

	return NewTestReporter()
}

func NewFromEnv(env string, name string) Reporter {
	var sentryDsn string
	if strings.ToLower(env) == "production" {
		sentryDsn = config.GetEnvValueOrPanic("SENTRY_DSN")
	}
	return New(env, name, sentryDsn)
}
