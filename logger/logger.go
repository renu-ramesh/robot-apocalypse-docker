package logger

import (
	"github.com/renu-ramesh/robot-apocalypse-docker/models"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func InitLogger(env models.EnvVariables) (*zap.Logger, error) {

	var logger *zap.Logger
	var err error
	var level zap.AtomicLevel
	var encoding string

	switch env.LogLevel {
	case "debug":
		level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	case "warn":
		level = zap.NewAtomicLevelAt(zapcore.WarnLevel)
	case "error":
		level = zap.NewAtomicLevelAt(zapcore.ErrorLevel)
	case "fatal":
		level = zap.NewAtomicLevelAt(zapcore.FatalLevel)
	case "panic":
		level = zap.NewAtomicLevelAt(zapcore.PanicLevel)
	default:
		level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	}

	switch env.LogFormat {
	case "json":
		encoding = "json"
	default:
		encoding = "console"
	}

	logger, err = zap.Config{
		Level:            level,
		Encoding:         encoding,
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey: "message",

			LevelKey:    "level",
			EncodeLevel: zapcore.CapitalColorLevelEncoder,

			TimeKey:    "timestamp",
			EncodeTime: zapcore.ISO8601TimeEncoder,

			CallerKey:    "caller",
			EncodeCaller: zapcore.ShortCallerEncoder,
		},
	}.Build()

	if err != nil {
		return nil, err
	}

	return logger, nil

}
