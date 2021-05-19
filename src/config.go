package main

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Config struct {
	Debugging  bool
	Port       int
	DbFilePath string
}

var site = Site{}

func (config *Config) initialise(debug *bool, port *int, dbPathPtr, baseUrlPtr *string) {
	config.Debugging = *debug
	if config.Debugging == true {
		zapLogger, _ = configureZapLogger(true)
		EnableDebugging()
		zapLogger.Info("Debugging enabled")
	}
	config.Port = *port
	config.DbFilePath = *dbPathPtr
	site.BaseUrl = *baseUrlPtr
}

func EnableDebugging() {
	zapLogger, _ = configureZapLogger(true)
}

func configureZapLogger(debugging bool) (*zap.Logger, error) {
	level := zapcore.InfoLevel
	encoderConfig := zapcore.EncoderConfig{
		MessageKey:  "message",
		LevelKey:    "level",
		TimeKey:     "",
		EncodeLevel: zapcore.CapitalColorLevelEncoder,
	}
	if debugging == true {
		level = zapcore.DebugLevel
		encoderConfig = zapcore.EncoderConfig{
			MessageKey:   "message",
			LevelKey:     "level",
			TimeKey:      "",
			EncodeLevel:  zapcore.CapitalColorLevelEncoder,
			CallerKey:    "caller",
			EncodeCaller: zapcore.ShortCallerEncoder,
		}
	}
	zapConfig := zap.Config{
		Encoding:      "console",
		Level:         zap.NewAtomicLevelAt(level),
		OutputPaths:   []string{"stdout"},
		EncoderConfig: encoderConfig,
	}
	return zapConfig.Build()
}
