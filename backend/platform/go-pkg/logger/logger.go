package logger

import (
	"io"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

// Level level log
type Level string

const (
	DebugLevel Level = "debug"
	InfoLevel  Level = "info"
	WarnLevel  Level = "warn"
	ErrorLevel Level = "error"
	FatalLevel Level = "fatal"
)

// Logger wrapper logrus
type Logger struct {
	entry *logrus.Entry
}

// Fields untuk context tambahan
type Fields map[string]interface{}

// New buat instance logger
func New(serviceName string) *Logger {
	log := logrus.New()

	// Format JSON (untuk production, gampang diparse)
	log.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339,
	})

	// Output ke stdout
	log.SetOutput(os.Stdout)

	// Default level
	log.SetLevel(logrus.InfoLevel)

	// Tambah field default (service name)
	entry := log.WithField("service", serviceName)

	return &Logger{entry: entry}
}

// SetLevel ubah level log
func (l *Logger) SetLevel(level Level) {
	switch level {
	case DebugLevel:
		l.entry.Logger.SetLevel(logrus.DebugLevel)
	case InfoLevel:
		l.entry.Logger.SetLevel(logrus.InfoLevel)
	case WarnLevel:
		l.entry.Logger.SetLevel(logrus.WarnLevel)
	case ErrorLevel:
		l.entry.Logger.SetLevel(logrus.ErrorLevel)
	case FatalLevel:
		l.entry.Logger.SetLevel(logrus.FatalLevel)
	}
}

// SetOutput ubah output (untuk test)
func (l *Logger) SetOutput(w io.Writer) {
	l.entry.Logger.SetOutput(w)
}

// WithFields tambah context ke log
func (l *Logger) WithFields(fields Fields) *Logger {
	return &Logger{
		entry: l.entry.WithFields(logrus.Fields(fields)),
	}
}

// WithError tambah error ke log
func (l *Logger) WithError(err error) *Logger {
	return &Logger{
		entry: l.entry.WithError(err),
	}
}

// Debug log debug
func (l *Logger) Debug(msg string) {
	l.entry.Debug(msg)
}

// Debugf log debug dengan format
func (l *Logger) Debugf(format string, args ...interface{}) {
	l.entry.Debugf(format, args...)
}

// Info log info
func (l *Logger) Info(msg string) {
	l.entry.Info(msg)
}

// Infof log info dengan format
func (l *Logger) Infof(format string, args ...interface{}) {
	l.entry.Infof(format, args...)
}

// Warn log warning
func (l *Logger) Warn(msg string) {
	l.entry.Warn(msg)
}

// Warnf log warning dengan format
func (l *Logger) Warnf(format string, args ...interface{}) {
	l.entry.Warnf(format, args...)
}

// Error log error
func (l *Logger) Error(msg string) {
	l.entry.Error(msg)
}

// Errorf log error dengan format
func (l *Logger) Errorf(format string, args ...interface{}) {
	l.entry.Errorf(format, args...)
}

// Fatal log fatal + exit
func (l *Logger) Fatal(msg string) {
	l.entry.Fatal(msg)
}

// Fatalf log fatal dengan format + exit
func (l *Logger) Fatalf(format string, args ...interface{}) {
	l.entry.Fatalf(format, args...)
}

// Print log biasa (untuk compatibility)
func (l *Logger) Print(v ...interface{}) {
	l.entry.Print(v...)
}

// Printf log biasa dengan format
func (l *Logger) Printf(format string, v ...interface{}) {
	l.entry.Printf(format, v...)
}