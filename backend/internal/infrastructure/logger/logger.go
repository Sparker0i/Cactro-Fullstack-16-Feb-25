package logger

import (
	"os"

	"github.com/Sparker0i/cactro-polls/internal/infrastructure/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger interface {
	Debug(msg string, fields ...Field)
	Info(msg string, fields ...Field)
	Warn(msg string, fields ...Field)
	Error(msg string, fields ...Field)
	Fatal(msg string, fields ...Field)
	With(fields ...Field) Logger
}

type Field = zapcore.Field

type zapLogger struct {
	logger *zap.Logger
}

func NewLogger(cfg *config.LoggerConfig) (Logger, error) {
	// Configure logger
	config := zap.NewProductionConfig()

	// Set level
	level, err := zapcore.ParseLevel(cfg.Level)
	if err != nil {
		return nil, err
	}
	config.Level = zap.NewAtomicLevelAt(level)

	// Configure output
	switch cfg.Format {
	case "console":
		config.Encoding = "console"
	default:
		config.Encoding = "json"
	}

	// Configure output destination
	var output zapcore.WriteSyncer
	switch cfg.Output {
	case "stderr":
		output = zapcore.AddSync(os.Stderr)
	default:
		output = zapcore.AddSync(os.Stdout)
	}

	// Create core
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(config.EncoderConfig),
		output,
		config.Level,
	)

	// Create logger
	logger := zap.New(core)
	return &zapLogger{logger: logger}, nil
}

func (l *zapLogger) Debug(msg string, fields ...Field) {
	l.logger.Debug(msg, fields...)
}

func (l *zapLogger) Info(msg string, fields ...Field) {
	l.logger.Info(msg, fields...)
}

func (l *zapLogger) Warn(msg string, fields ...Field) {
	l.logger.Warn(msg, fields...)
}

func (l *zapLogger) Error(msg string, fields ...Field) {
	l.logger.Error(msg, fields...)
}

func (l *zapLogger) Fatal(msg string, fields ...Field) {
	l.logger.Fatal(msg, fields...)
}

func (l *zapLogger) With(fields ...Field) Logger {
	return &zapLogger{logger: l.logger.With(fields...)}
}

// Helper functions for creating fields
func String(key, value string) Field {
	return zap.String(key, value)
}

func Int(key string, value int) Field {
	return zap.Int(key, value)
}

func Error(err error) Field {
	return zap.Error(err)
}
