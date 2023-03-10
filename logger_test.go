package library

import (
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest"
)

// use this for init logger on ServiceContext , for local test
func MockLogger(t *testing.T) *zap.Logger {
	logger := zaptest.NewLogger(t, zaptest.WrapOptions(
		zap.Hooks(func(e zapcore.Entry) error {
			if e.Level == zap.ErrorLevel {
				t.Fatal("error should not happen")
			}
			return nil
		})))
	return logger
}

func Test_Logger(t *testing.T) {
	logger := MockLogger(t)
	logger.Info("success")
}
