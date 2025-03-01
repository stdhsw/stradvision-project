package logger

import (
	"path/filepath"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	DefaultPath         string = "/var/log/app/"
	DefaultLogExtention string = ".log"
	DefaultAppName      string = "app"

	DefaultMaxSize    int  = 100
	DefaultMaxBackups int  = 3
	DefaultMaxAge     int  = 7
	DefaultLocalTime  bool = true
	DefaultCompress   bool = true
)

type config struct {
	appName string
	encoder zapcore.Encoder
	level   zapcore.Level
	logger  lumberjack.Logger
}

func defaultOption() config {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	return config{
		appName: DefaultAppName,
		encoder: zapcore.NewJSONEncoder(encoderConfig),
		level:   zapcore.InfoLevel,
		logger: lumberjack.Logger{
			Filename:   DefaultPath + DefaultAppName + DefaultLogExtention,
			MaxSize:    DefaultMaxSize,
			MaxBackups: DefaultMaxBackups,
			MaxAge:     DefaultMaxAge,
			LocalTime:  DefaultLocalTime,
			Compress:   DefaultCompress,
		},
	}
}

type Option func(*config)

func fromOptions(appName string, opts ...Option) *config {
	c := defaultOption()
	c.appName = appName

	for _, opt := range opts {
		opt(&c)
	}

	return &c
}

func WithPath(path string) Option {
	return func(c *config) {
		if path != "" {
			c.logger.Filename = filepath.Join(path, c.appName+DefaultLogExtention)
		}
	}
}

func WithLogMaxSize(size int) Option {
	return func(c *config) {
		if size > 0 {
			c.logger.MaxSize = size
		}
	}
}

func WithLogMaxBackups(backups int) Option {
	return func(c *config) {
		if backups > 0 {
			c.logger.MaxBackups = backups
		}
	}
}

func WithLogMaxAge(age int) Option {
	return func(c *config) {
		if age > 0 {
			c.logger.MaxAge = age
		}
	}
}

func WithLogLocalTime(localTime bool) Option {
	return func(c *config) {
		c.logger.LocalTime = localTime
	}
}

func WithLogCompress(compress bool) Option {
	return func(c *config) {
		c.logger.Compress = compress
	}
}

func WithLogLevel(level string) Option {
	return func(c *config) {
		switch strings.ToUpper(level) {
		case "DEBUG":
			c.level = zapcore.DebugLevel
		case "WARN", "WARNING":
			c.level = zapcore.WarnLevel
		case "ERROR", "ERR":
			c.level = zapcore.ErrorLevel
		case "DPANIC":
			c.level = zapcore.DPanicLevel
		case "PANIC":
			c.level = zapcore.PanicLevel
		case "FATAL":
			c.level = zapcore.FatalLevel
		case "INFO", "INF":
			fallthrough
		default:
			c.level = zapcore.InfoLevel
		}
	}
}

func WitchEncoder(encoder string) Option {
	return func(c *config) {
		encoderConfig := zap.NewProductionEncoderConfig()
		encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

		switch strings.ToUpper(encoder) {
		case "CONSOLE":
			c.encoder = zapcore.NewConsoleEncoder(encoderConfig)
		case "JSON":
			fallthrough
		default:
			c.encoder = zapcore.NewJSONEncoder(encoderConfig)
		}
	}
}
