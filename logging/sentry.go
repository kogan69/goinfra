package logging

import (
	"log"
	"sync"

	"github.com/getsentry/sentry-go"
)

type SentryErrorReporter struct {
	serviceName string
	envName     string
}

var initOnce sync.Once

func NewSentryErrorReporter(envName, serviceName, dsn string) *SentryErrorReporter {
	initOnce.Do(func() {
		err := sentry.Init(sentry.ClientOptions{
			Dsn:              dsn,
			TracesSampleRate: 1.0,
			Environment:      envName,
		})
		if err != nil {
			log.Fatalln(err)
		}
	})
	return &SentryErrorReporter{serviceName: serviceName, envName: envName}
}
func (s *SentryErrorReporter) ReportException(err error) {
	localHub := sentry.CurrentHub().Clone()
	localHub.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetTag("service", s.serviceName)
		scope.SetLevel(sentry.LevelError)
	})
	localHub.CaptureException(err)
}

func (s *SentryErrorReporter) ReportFatal(err error) {
	localHub := sentry.CurrentHub().Clone()
	localHub.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetTag("service", s.serviceName)
		scope.SetLevel(sentry.LevelFatal)
	})
	localHub.CaptureException(err)
}
