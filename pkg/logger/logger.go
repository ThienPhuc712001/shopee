// Package logger provides structured logging using Logrus
package logger

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

// ============================================================================
// LOG LEVELS
// ============================================================================

// Level represents a logging level
type Level string

const (
	DebugLevel   Level = "debug"
	InfoLevel    Level = "info"
	WarnLevel    Level = "warn"
	ErrorLevel   Level = "error"
	FatalLevel   Level = "fatal"
	PanicLevel   Level = "panic"
)

// ParseLevel converts a string to logrus Level
func ParseLevel(level string) logrus.Level {
	switch level {
	case "debug":
		return logrus.DebugLevel
	case "info":
		return logrus.InfoLevel
	case "warn", "warning":
		return logrus.WarnLevel
	case "error":
		return logrus.ErrorLevel
	case "fatal":
		return logrus.FatalLevel
	case "panic":
		return logrus.PanicLevel
	default:
		return logrus.InfoLevel
	}
}

// ============================================================================
// FIELDS
// ============================================================================

// Fields is a map of key-value pairs for structured logging
type Fields map[string]interface{}

// ============================================================================
// LOGGER
// ============================================================================

// Logger wraps logrus.Logger with additional functionality
type Logger struct {
	*logrus.Logger
}

// New creates a new Logger with the given configuration
func New(config *Config) (*Logger, error) {
	log := logrus.New()

	// Set log level
	level, err := logrus.ParseLevel(config.Level)
	if err != nil {
		level = logrus.InfoLevel
	}
	log.SetLevel(level)

	// Set log format
	if config.Format == "json" {
		log.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: time.RFC3339,
		})
	} else {
		log.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02 15:04:05",
			ForceColors:     config.Color,
		})
	}

	// Set output (file, stdout, or both)
	var writers []io.Writer

	// Add file writer if configured
	if config.FilePath != "" {
		fileWriter, err := createFileWriter(config)
		if err != nil {
			return nil, fmt.Errorf("failed to create file writer: %w", err)
		}
		writers = append(writers, fileWriter)
	}

	// Add stdout if configured or if no file path
	if config.Stdout || config.FilePath == "" {
		writers = append(writers, os.Stdout)
	}

	// If error logging to file is enabled, create separate error log
	if config.ErrorFilePath != "" {
		errorWriter, err := createErrorFileWriter(config)
		if err != nil {
			return nil, fmt.Errorf("failed to create error file writer: %w", err)
		}
		
		// Create a hook for error-level logs
		log.AddHook(&ErrorHook{
			ErrorWriter: errorWriter,
			LogLevels: []logrus.Level{
				logrus.ErrorLevel,
				logrus.FatalLevel,
				logrus.PanicLevel,
			},
		})
	}

	// Combine writers
	log.SetOutput(io.MultiWriter(writers...))

	// Add default fields
	if config.Service != "" {
		log.WithField("service", config.Service)
	}
	if config.Environment != "" {
		log.WithField("environment", config.Environment)
	}

	return &Logger{Logger: log}, nil
}

// createFileWriter creates a writer with log rotation
func createFileWriter(config *Config) (io.Writer, error) {
	// Ensure directory exists
	dir := filepath.Dir(config.FilePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}

	return &lumberjack.Logger{
		Filename:   config.FilePath,
		MaxSize:    config.MaxSize,    // megabytes
		MaxBackups: config.MaxBackups, // number of files
		MaxAge:     config.MaxAge,     // days
		Compress:   config.Compress,
		LocalTime:  true,
	}, nil
}

// createErrorFileWriter creates a writer for error logs with rotation
func createErrorFileWriter(config *Config) (io.Writer, error) {
	// Ensure directory exists
	dir := filepath.Dir(config.ErrorFilePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}

	return &lumberjack.Logger{
		Filename:   config.ErrorFilePath,
		MaxSize:    config.MaxSize,
		MaxBackups: config.MaxBackups,
		MaxAge:     config.MaxAge,
		Compress:   config.Compress,
		LocalTime:  true,
	}, nil
}

// ============================================================================
// LOGGING METHODS
// ============================================================================

// WithFields creates a new logger with the given fields
func (l *Logger) WithFields(fields Fields) *logrus.Entry {
	return l.Logger.WithFields(logrus.Fields(fields))
}

// WithField creates a new logger with the given field
func (l *Logger) WithField(key string, value interface{}) *logrus.Entry {
	return l.Logger.WithField(key, value)
}

// WithError creates a new logger with the given error
func (l *Logger) WithError(err error) *logrus.Entry {
	return l.Logger.WithError(err)
}

// Debug logs a debug message
func (l *Logger) Debug(args ...interface{}) {
	l.Logger.Debug(args...)
}

// Info logs an info message
func (l *Logger) Info(args ...interface{}) {
	l.Logger.Info(args...)
}

// Warn logs a warning message
func (l *Logger) Warn(args ...interface{}) {
	l.Logger.Warn(args...)
}

// Error logs an error message
func (l *Logger) Error(args ...interface{}) {
	l.Logger.Error(args...)
}

// Fatal logs a fatal message and exits
func (l *Logger) Fatal(args ...interface{}) {
	l.Logger.Fatal(args...)
}

// Panic logs a panic message and panics
func (l *Logger) Panic(args ...interface{}) {
	l.Logger.Panic(args...)
}

// Debugf logs a formatted debug message
func (l *Logger) Debugf(format string, args ...interface{}) {
	l.Logger.Debugf(format, args...)
}

// Infof logs a formatted info message
func (l *Logger) Infof(format string, args ...interface{}) {
	l.Logger.Infof(format, args...)
}

// Warnf logs a formatted warning message
func (l *Logger) Warnf(format string, args ...interface{}) {
	l.Logger.Warnf(format, args...)
}

// Errorf logs a formatted error message
func (l *Logger) Errorf(format string, args ...interface{}) {
	l.Logger.Errorf(format, args...)
}

// Fatalf logs a formatted fatal message and exits
func (l *Logger) Fatalf(format string, args ...interface{}) {
	l.Logger.Fatalf(format, args...)
}

// Panicf logs a formatted panic message and panics
func (l *Logger) Panicf(format string, args ...interface{}) {
	l.Logger.Panicf(format, args...)
}

// ============================================================================
// ERROR HOOK
// ============================================================================

// ErrorHook is a logrus hook that writes error logs to a separate file
type ErrorHook struct {
	ErrorWriter io.Writer
	LogLevels   []logrus.Level
}

// Fire is called when a log event is fired
func (hook *ErrorHook) Fire(entry *logrus.Entry) error {
	// Create a copy of the entry without the hook to avoid infinite loop
	entry.Logger.Level = logrus.PanicLevel // Prevent the hook from firing again
	
	// Format the log entry
	formatted, err := entry.Logger.Formatter.Format(entry)
	if err != nil {
		return err
	}
	
	// Write to error file
	_, err = hook.ErrorWriter.Write(formatted)
	return err
}

// Levels returns the log levels this hook handles
func (hook *ErrorHook) Levels() []logrus.Level {
	return hook.LogLevels
}

// ============================================================================
// CONFIGURATION
// ============================================================================

// Config holds logger configuration
type Config struct {
	// Level is the minimum log level (debug, info, warn, error, fatal, panic)
	Level string `mapstructure:"level" default:"info"`

	// Format is the log format (text, json)
	Format string `mapstructure:"format" default:"text"`

	// FilePath is the path to the log file
	FilePath string `mapstructure:"file_path" default:"logs/app.log"`

	// ErrorFilePath is the path to the error log file
	ErrorFilePath string `mapstructure:"error_file_path" default:"logs/error.log"`

	// Stdout enables logging to stdout
	Stdout bool `mapstructure:"stdout" default:"true"`

	// Color enables colored output in text format
	Color bool `mapstructure:"color" default:"true"`

	// Service is the service name to include in logs
	Service string `mapstructure:"service"`

	// Environment is the environment (development, staging, production)
	Environment string `mapstructure:"environment"`

	// MaxSize is the maximum size in megabytes before rotation
	MaxSize int `mapstructure:"max_size" default:"100"`

	// MaxBackups is the maximum number of old log files to retain
	MaxBackups int `mapstructure:"max_backups" default:"30"`

	// MaxAge is the maximum number of days to retain old log files
	MaxAge int `mapstructure:"max_age" default:"90"`

	// Compress enables compression of rotated log files
	Compress bool `mapstructure:"compress" default:"true"`
}

// DefaultConfig returns a default configuration
func DefaultConfig() *Config {
	return &Config{
		Level:         "info",
		Format:        "text",
		FilePath:      "logs/app.log",
		ErrorFilePath: "logs/error.log",
		Stdout:        true,
		Color:         true,
		MaxSize:       100,
		MaxBackups:    30,
		MaxAge:        90,
		Compress:      true,
	}
}

// ============================================================================
// GLOBAL LOGGER
// ============================================================================

var defaultLogger *Logger

// Init initializes the global logger with the given config
func Init(config *Config) error {
	var err error
	defaultLogger, err = New(config)
	return err
}

// Get returns the global logger
func Get() *Logger {
	if defaultLogger == nil {
		// Create a default logger if not initialized
		defaultLogger, _ = New(DefaultConfig())
	}
	return defaultLogger
}

// SetDefault sets the global logger
func SetDefault(log *Logger) {
	defaultLogger = log
}

// StandardLogger returns the global logger for package-level logging
func StandardLogger() *Logger {
	return Get()
}
