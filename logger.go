package library

import (
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

func Logging(LOG *zap.Logger, rp RespParams) {
	switch rp.Severity {
	case DEBUG:
		LOG.Debug(rp.Section,
			zap.String("connection", rp.URL),
			zap.Any("parameters", rp.Input),
			zap.Error(rp.Error))
	case INFO:
		LOG.Info(rp.Section,
			zap.String("connection", rp.URL),
			zap.Any("parameters", rp.Input))
	case WARN:
		LOG.Warn(rp.Section,
			zap.String("connection", rp.URL),
			zap.Any("parameters", rp.Input),
			zap.Error(rp.Error))
	case ERROR:
		LOG.Error(rp.Section,
			zap.String("connection", rp.URL),
			zap.Any("parameters", rp.Input),
			zap.Error(rp.Error))
	}
}

func NewLogger(path string, debug bool) (*zap.Logger, error) {
	w := zapcore.AddSync(&lumberjack.Logger{
		Filename:   path + time.Now().Format("20060102") + ".log",
		MaxSize:    100, // megabytes
		MaxBackups: 30,
		MaxAge:     30, // days
	})

	pe := zap.NewProductionEncoderConfig()
	if debug {
		pe = zap.NewDevelopmentEncoderConfig()
	}

	fileEncoder := zapcore.NewJSONEncoder(pe)
	pe.EncodeTime = zapcore.ISO8601TimeEncoder

	consoleEncoder := zapcore.NewConsoleEncoder(pe)

	level := zap.InfoLevel
	if debug {
		level = zap.DebugLevel
	}

	core := zapcore.NewTee(
		zapcore.NewCore(fileEncoder, zapcore.AddSync(w), level),
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), level),
	)

	l := zap.New(core)

	return l, nil
}
