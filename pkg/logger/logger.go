package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// ZapAdapter is the Zap logger adapter that implements Logger interface
type ZapAdapter struct {
	logger *zap.Logger
}

var defaultLogger Logger

// Init initializes the default logger
func Init(level string) error {
	var zapLevel zapcore.Level
	switch level {
	case "debug":
		zapLevel = zapcore.DebugLevel
	case "info":
		zapLevel = zapcore.InfoLevel
	case "warn":
		zapLevel = zapcore.WarnLevel
	case "error":
		zapLevel = zapcore.ErrorLevel
	default:
		zapLevel = zapcore.InfoLevel
	}

	config := zap.Config{
		Level:            zap.NewAtomicLevelAt(zapLevel),
		Development:      false,
		Encoding:         "json",
		EncoderConfig:    zap.NewProductionEncoderConfig(),
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}

	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	zapLogger, err := config.Build()
	if err != nil {
		return err
	}

	defaultLogger = &ZapAdapter{logger: zapLogger}
	return nil
}

// InitDevelopment initializes the logger for development
func InitDevelopment() error {
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder

	zapLogger, err := config.Build()
	if err != nil {
		return err
	}

	defaultLogger = &ZapAdapter{logger: zapLogger}
	return nil
}

// New creates a new ZapAdapter with custom zap logger
func New(zapLogger *zap.Logger) Logger {
	return &ZapAdapter{logger: zapLogger}
}

// GetDefault returns the default logger instance
func GetDefault() Logger {
	if defaultLogger == nil {
		// Fallback to development logger if not initialized
		_ = InitDevelopment()
	}
	return defaultLogger
}

// Debug logs a debug message
func (l *ZapAdapter) Debug(msg string, fields ...Field) {
	// For now, just log the message without fields
	// In real implementation, we'd convert fields properly
	l.logger.Debug(msg)
}

// Info logs an info message
func (l *ZapAdapter) Info(msg string, fields ...Field) {
	l.logger.Info(msg)
}

// Warn logs a warning message
func (l *ZapAdapter) Warn(msg string, fields ...Field) {
	l.logger.Warn(msg)
}

// Error logs an error message
func (l *ZapAdapter) Error(msg string, fields ...Field) {
	l.logger.Error(msg)
}

// Fatal logs a fatal message and exits
func (l *ZapAdapter) Fatal(msg string, fields ...Field) {
	l.logger.Fatal(msg)
	os.Exit(1)
}

// Sync flushes any buffered log entries
func (l *ZapAdapter) Sync() error {
	return l.logger.Sync()
}

// Global convenience functions that use the default logger
// These accept zap.Field directly for backward compatibility

// Debug logs a debug message using the default logger
func Debug(msg string, fields ...zap.Field) {
	if defaultLogger == nil {
		GetDefault()
	}
	if adapter, ok := defaultLogger.(*ZapAdapter); ok {
		adapter.logger.Debug(msg, fields...)
	}
}

// Info logs an info message using the default logger
func Info(msg string, fields ...zap.Field) {
	if defaultLogger == nil {
		GetDefault()
	}
	if adapter, ok := defaultLogger.(*ZapAdapter); ok {
		adapter.logger.Info(msg, fields...)
	}
}

// Warn logs a warning message using the default logger
func Warn(msg string, fields ...zap.Field) {
	if defaultLogger == nil {
		GetDefault()
	}
	if adapter, ok := defaultLogger.(*ZapAdapter); ok {
		adapter.logger.Warn(msg, fields...)
	}
}

// Error logs an error message using the default logger
func Error(msg string, fields ...zap.Field) {
	if defaultLogger == nil {
		GetDefault()
	}
	if adapter, ok := defaultLogger.(*ZapAdapter); ok {
		adapter.logger.Error(msg, fields...)
	}
}

// Fatal logs a fatal message and exits using the default logger
func Fatal(msg string, fields ...zap.Field) {
	if defaultLogger == nil {
		GetDefault()
	}
	if adapter, ok := defaultLogger.(*ZapAdapter); ok {
		adapter.logger.Fatal(msg, fields...)
	}
	os.Exit(1)
}

// Sync flushes any buffered log entries from the default logger
func Sync() error {
	if defaultLogger != nil {
		return defaultLogger.Sync()
	}
	return nil
}
