package logger

import (
	"log"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func Init(logLevelRaw string) *zap.Logger {
	prodConfig := zap.NewProductionConfig()
	prodConfig.Encoding = "json"
	prodConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	prodConfig.EncoderConfig.EncodeDuration = zapcore.StringDurationEncoder

	// set logger level
	var logLevel zapcore.Level
	err := logLevel.UnmarshalText([]byte(logLevelRaw))
	if err != nil { // if err, set level debug as a default
		logLevel = zapcore.DebugLevel
	}

	prodConfig.Level = zap.NewAtomicLevelAt(logLevel)
	logger, err := prodConfig.Build()
	if err != nil {
		log.Fatalln("build logger from config failed")
	}
	zap.ReplaceGlobals(logger) // don't remove this, otherwise logger won't work

	zap.L().Info("logger has been initiated", zap.String("level", logLevel.String()))

	return logger
}
