package main

import (
	"flag"
	"github.com/gorilla/securecookie"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Debugging                    bool
	Port                         int
	DbFilePath                   string
	DbSaveInterval               string
	SessionKey                   []byte
	CookieKey                    []byte
	AdminEmails                  []string
	SessionExpiryIntervalSeconds int
}

var site = Site{}

func (config *Config) initialise() {
	debugPtr := flag.Bool("debug", false, "Enable debug logging")
	hostPtr := flag.String("host", "http://localhost", "Protocol and host")
	portPtr := flag.Int("port", 80, "Port number")
	dbPathPtr := flag.String("db", "ldn_inbox.sqlite", "Path to to Database file")
	adminEmailPtr := flag.String("admin", "paul@antleaf.com", "Admin email address")
	flag.Parse()
	config.Debugging = *debugPtr
	config.Debugging = true
	zapLogger, _ = configureZapLogger(config.Debugging)
	if config.Debugging == true {
		zapLogger.Info("Debugging enabled")
	}
	config.Port = *portPtr
	zapLogger.Debug("Set port to ", zap.Int("port", config.Port))
	if config.Port == 80 {
		site.BaseUrl = *hostPtr
	} else {
		site.BaseUrl = *hostPtr + ":" + strconv.Itoa(config.Port)
	}
	zapLogger.Debug("Set base url to ", zap.String("host", site.BaseUrl))
	config.DbFilePath = *dbPathPtr
	zapLogger.Debug("Set db path to ", zap.String("db path", config.DbFilePath))
	config.SessionKey = []byte(os.Getenv("COAR_NOTIFY_INBOX_SESSION_KEY"))
	if len(config.SessionKey) == 0 {
		zapLogger.Fatal("Session Key is nil - please set this in your environment with key 'COAR_NOTIFY_INBOX_SESSION_KEY'")
	}
	config.CookieKey = []byte(os.Getenv("COAR_NOTIFY_INBOX_COOKIE_STORE_KEY"))
	if len(config.CookieKey) == 0 {
		zapLogger.Fatal("Cookie Key is nil - please set this in your environment with key 'COAR_NOTIFY_INBOX_COOKIE_STORE_KEY; For now, using random key'")
		config.CookieKey = securecookie.GenerateRandomKey(16)
	}
	config.AdminEmails = strings.Split(*adminEmailPtr, ",")
	config.SessionExpiryIntervalSeconds = 86400
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
