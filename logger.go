package loggo

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"runtime"
	"sync"
	"text/template"
	"time"
)

// Logger is the structure that holds the logger information.
// It includes the log level Threshold, output destination, message template, and time provider.
type Logger struct {
	Context        context.Context // Context for the logger
	Threshold      Level           // Minimum log level to output
	mu             sync.RWMutex    // Ensures thread-safe access to the logger
	output         io.Writer       // Destination for log output
	template       string          // Template for log messages
	now            TimeProvider    // Function to get the current time
	timeFormat     string          // Format for the time in the log message
	maxSize        int             // Maximum size of the log message
	callerProvider CallerProvider  // Function to get the caller information
	preHooks       []Hook          // Pre-hooks to run before logging
	postHooks      []Hook          // Post-hooks to run after logging
}

// New creates a new Logger with the given Threshold and options.
// The default output is os.Stdout, the default template is "%s [%5s]: %s", and the default time provider is time.Now.
//
// Parameters:
//   - Threshold: Minimum log level to output.
//   - options: Variadic options to configure the Logger.
//
// Returns:
//   - A pointer to the newly created Logger.
//
// Example:
//
//	logger := loggo.New(loggo.LevelInfo, loggo.WithOutput(os.Stderr))
//	logger.Info("This is an info message")
func New(threshold Level, options ...Option) *Logger {
	defaultCaller := func() (pc uintptr, file string, line int, ok bool) {
		pc, file, line, ok = runtime.Caller(5)

		return
	}
	log := &Logger{
		Threshold:      threshold,
		Context:        context.Background(),
		output:         os.Stdout,
		template:       "{{.Time}} [{{printf \"%5s\" .Level}}]: {{.Message}}",
		now:            time.Now,
		timeFormat:     "2006-01-02 15:04:05",
		maxSize:        1000,
		callerProvider: defaultCaller,
		preHooks:       []Hook{},
		postHooks:      []Hook{},
	}

	for _, option := range options {
		option(log)
	}

	return log
}

// Log logs a message at the given log level.
// If the log level is below the Threshold, the message is not logged. If an error occurs while logging the message, it is ignored.
//
// Parameters:
//   - level: The log level of the message.
//   - message: The message to log.
//
// Example:
//
//	logger := loggo.New(loggo.LevelInfo)
//	logger.Log(loggo.LevelInfo, "This is an info message")
func (l *Logger) Log(level Level, message string) {
	_ = l.LogE(level, message)
}

// LogE logs a message at the given log level and returns an error if the message could not be logged.
// If the log level is below the Threshold, the message is not logged.
//
// Parameters:
//   - level: The log level of the message.
//   - message: The message to log.
//
// Returns:
//   - An error if the message could not be logged, nil otherwise.
//
// Example:
//
//	logger := loggo.New(loggo.LevelInfo)
//	err := logger.LogE(loggo.LevelInfo, "This is an info message")
//	if err != nil {
//		log.Fatal(err)
//	}
func (l *Logger) LogE(level Level, message string) error {
	for _, hook := range l.preHooks {
		hook(l, &message)
	}

	if l.Threshold > level {
		return nil
	}

	data := getTemplateData(level, message, l)

	tmpl, err := template.New("log").Parse(l.template + "\n")
	if err != nil {
		return errors.New("error parsing template: " + err.Error())
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	if err = tmpl.Execute(l.output, data); err != nil {
		return errors.New("error executing template: " + err.Error())
	}

	for _, hook := range l.postHooks {
		hook(l, &message)
	}

	return nil
}

// Logf logs a formatted message at the given log level.
// If the log level is below the Threshold, the message is not logged. If an error occurs while logging the message, it is ignored.
//
// Parameters:
//   - level: The log level of the message.
//   - format: The format string for the message.
//   - args: The arguments for the format string.
//
// Example:
//
//	logger := loggo.New(loggo.LevelInfo)
//	logger.Logf(loggo.LevelInfo, "This is an info message with a %s", "format")
func (l *Logger) Logf(level Level, format string, args ...any) {
	l.Log(level, fmt.Sprintf(format, args...))
}

// LogfE logs a formatted message at the given log level and returns an error if the message could not be logged.
// If the log level is below the Threshold, the message is not logged.
//
// Parameters:
//   - level: The log level of the message.
//   - format: The format string for the message.
//   - args: The arguments for the format string.
//
// Returns:
//   - An error if the message could not be logged, nil otherwise.
//
// Example:
//
//	logger := loggo.New(loggo.LevelInfo)
//	err := logger.LogfE(loggo.LevelInfo, "This is an info message with a %s", "format")
//	if err != nil {
//		log.Fatal(err)
//	}
func (l *Logger) LogfE(level Level, format string, args ...any) error {
	return l.LogE(level, fmt.Sprintf(format, args...))
}

// Debug logs a message at the LevelDebug. If an error occurs while logging the message, it is ignored.
//
// Parameters:
//   - message: The debug message to log.
//
// Example:
//
//	logger := loggo.New(loggo.LevelDebug)
//	logger.Debug("This is a debug message")
func (l *Logger) Debug(message string) {
	l.Log(LevelDebug, message)
}

// Debugf logs a formatted message at the LevelDebug. If an error occurs while logging the message, it is ignored.
//
// Parameters:
//   - format: The format string for the debug message.
//   - args: The arguments for the format string.
//
// Example:
//
//	logger := loggo.New(loggo.LevelDebug)
//	logger.Debugf("This is a debug message with a %s", "format")
func (l *Logger) Debugf(format string, args ...any) {
	l.Logf(LevelDebug, format, args...)
}

// Info logs a message at the LevelInfo. If an error occurs while logging the message, it is ignored.
//
// Parameters:
//   - message: The info message to log.
//
// Example:
//
//	logger := loggo.New(loggo.LevelInfo)
//	logger.Info("This is an info message")
func (l *Logger) Info(message string) {
	l.Log(LevelInfo, message)
}

// Infof logs a formatted message at the LevelInfo. If an error occurs while logging the message, it is ignored.
//
// Parameters:
//   - format: The format string for the info message.
//   - args: The arguments for the format string.
//
// Example:
//
//	logger := loggo.New(loggo.LevelInfo)
//	logger.Infof("This is an info message with a %s", "format")
func (l *Logger) Infof(format string, args ...any) {
	l.Logf(LevelInfo, format, args...)
}

// Warn logs a message at the LevelWarn. If an error occurs while logging the message, it is ignored.
//
// Parameters:
//   - message: The warn message to log.
//
// Example:
//
//	logger := loggo.New(loggo.LevelWarn)
//	logger.Warn("This is a warn message")
func (l *Logger) Warn(message string) {
	l.Log(LevelWarn, message)
}

// Warnf logs a formatted message at the LevelWarn. If an error occurs while logging the message, it is ignored.
//
// Parameters:
//   - format: The format string for the warn message.
//   - args: The arguments for the format string.
//
// Example:
//
//	logger := loggo.New(loggo.LevelWarn)
//	logger.Warnf("This is a warn message with a %s", "format")
func (l *Logger) Warnf(format string, args ...any) {
	l.Logf(LevelWarn, format, args...)
}

// Error logs a message at the LevelError. If an error occurs while logging the message, it is ignored.
//
// Parameters:
//   - message: The error message to log.
//
// Example:
//
//	logger := loggo.New(loggo.LevelError)
//	logger.Error("This is an error message")
func (l *Logger) Error(message string) {
	l.Log(LevelError, message)
}

// Errorf logs a formatted message at the LevelError. If an error occurs while logging the message, it is ignored.
//
// Parameters:
//   - format: The format string for the error message.
//   - args: The arguments for the format string.
//
// Example:
//
//	logger := loggo.New(loggo.LevelError)
//	logger.Errorf("This is an error message with a %s", "format")
func (l *Logger) Errorf(format string, args ...any) {
	l.Logf(LevelError, format, args...)
}

// Fatal logs a message at the LevelFatal. If an error occurs while logging the message, it is ignored.
//
// Parameters:
//   - message: The fatal message to log.
//
// Example:
//
//	logger := loggo.New(loggo.LevelFatal)
//	logger.Fatal("This is a fatal message")
func (l *Logger) Fatal(message string) {
	l.Log(LevelFatal, message)
}

// Fatalf logs a formatted message at the LevelFatal. If an error occurs while logging the message, it is ignored.
//
// Parameters:
//   - format: The format string for the fatal message.
//   - args: The arguments for the format string.
//
// Example:
//
//	logger := loggo.New(loggo.LevelFatal)
//	logger.Fatalf("This is a fatal message with a %s", "format")
func (l *Logger) Fatalf(format string, args ...any) {
	l.Logf(LevelFatal, format, args...)
}
