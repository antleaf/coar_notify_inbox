package main

import (
	"flag"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Config struct {
	Debugging      bool
	Port           int
	DbFilePath     string
	DbSaveInterval string
}

var site = Site{}

func (config *Config) initialise() {
	zapLogger, _ = configureZapLogger(config.Debugging)
	if config.Debugging == true {
		zapLogger.Info("Debugging enabled")
	}
	debugPtr := flag.Bool("debug", false, "Enable debug logging")
	portPtr := flag.Int("port", 80, "Port number")
	dbPathPtr := flag.String("db", "ldn_inbox.sqlite", "Path to to Database file")
	flag.Parse()
	config.Debugging = *debugPtr
	config.Port = *portPtr
	config.DbFilePath = *dbPathPtr
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

func configure() error {
	var err error
	config.initialise()
	if config.DbFilePath != "" {
		err = InitialiseDb(config.DbFilePath)
		if err != nil {
			zapLogger.Error(err.Error())
			return err
		}
	}
	initialiseRendering()
	router = ConfigureRouter()
	return err
}
