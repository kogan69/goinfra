package logging_test

import (
	"errors"
	"testing"

	"github.com/kogan69/goinfra/logging"
)

func Test_Log_Config(t *testing.T) {
	logging.Init("app1", "1.2.3", "dev", "debug", nil)
	logging.Debug("debug message")
	logging.Info("info message")
	logging.Warn("warning message")
	logging.Error(errors.New("error message"))
}
