package logger

import (
	"os"
	"github.com/rs/zerolog"
)

// Logger wraps zerolog.Logger for dependency injection
type Logger struct {
	*zerolog.Logger
}

// NewLogger creates a new logger instance with console output in plain text
func NewLogger() *Logger {
	// Configure zerolog for plain text console output
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	// Create a new logger with plain text output (no JSON)
	output := zerolog.ConsoleWriter{Out: os.Stdout}
	log := zerolog.New(output).Level(zerolog.InfoLevel)

	return &Logger{Logger: &log}
}

// Info logs an info level message
func (l *Logger) Info(msg string) {
	l.Logger.Info().Msg(msg)
}

// Infof logs a formatted info message
func (l *Logger) Infof(msg string, fields ...interface{}) {
	l.Logger.Info().Msgf(msg, fields...)
}

// Warn logs a warning level message
func (l *Logger) Warn(msg string) {
	l.Logger.Warn().Msg(msg)
}

// Warnf logs a formatted warning message
func (l *Logger) Warnf(msg string, fields ...interface{}) {
	l.Logger.Warn().Msgf(msg, fields...)
}

// Error logs an error level message
func (l *Logger) Error(msg string) {
	l.Logger.Error().Msg(msg)
}

// Errorf logs a formatted error message
func (l *Logger) Errorf(msg string, fields ...interface{}) {
	l.Logger.Error().Msgf(msg, fields...)
}

// Debug logs a debug level message
func (l *Logger) Debug(msg string) {
	l.Logger.Debug().Msg(msg)
}

// Debugf logs a formatted debug message
func (l *Logger) Debugf(msg string, fields ...interface{}) {
	l.Logger.Debug().Msgf(msg, fields...)
}

// WithContext returns a logger with context fields
func (l *Logger) WithContext(key string, value interface{}) *Logger {
	log := l.Logger.With().Interface(key, value).Logger()
	return &Logger{Logger: &log}
}
