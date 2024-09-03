package loggo

import (
	"context"
	"io"
	"time"
)

// Option is a function that configures a Logger.
type Option func(*Logger)

// TimeProvider is a function that returns the current time.
type TimeProvider func() time.Time

// CallerProvider is a function that returns the path of the caller, the file name, and the line number, and a
// boolean indicating if the information is available.
type CallerProvider func() (pc uintptr, file string, line int, ok bool)

// Hook is a function that is executed before or after logging a message.
type Hook func(l *Logger, message *string)

// WithOutput configures the output destination of a Logger. The default output is os.Stdout.
//
// Parameters:
//   - output: The io.Writer to use as the output destination.
//
// Example:
//
//	logger := loggo.New(loggo.LevelInfo, loggo.WithOutput(os.Stderr))
func WithOutput(output io.Writer) Option {
	return func(l *Logger) {
		l.output = output
	}
}

// WithTemplate configures the log message template of a Logger. The default template is
// "{{.Time}} [{{printf \"%5s\" .Level}}]: {{.Message}}".
//
// Parameters:
//   - template: The template string for log messages.
//
// Example:
//
//	logger := loggo.New(loggo.LevelInfo, loggo.WithTemplate("{{.Time}}: {{.Message}}"))
func WithTemplate(template string) Option {
	return func(l *Logger) {
		l.template = template
	}
}

// WithTimeProvider configures the time provider function of a Logger. The default time provider is time.Now.
//
// Parameters:
//   - provider: The TimeProvider function to use.
//
// Example:
//
//	logger := loggo.New(loggo.LevelInfo, loggo.WithTimeProvider(func() time.Time { return time.Unix(0, 0) }))
func WithTimeProvider(provider TimeProvider) Option {
	return func(l *Logger) {
		l.now = provider
	}
}

// WithTimeFormat configures the time format of a Logger. The default time format is "2006-01-02 15:04:05".
//
// Parameters:
//   - format: The format string for the time in the log message.
//
// Example:
//
//	logger := loggo.New(loggo.LevelInfo, loggo.WithTimeFormat("2006-01-02 15:04:05"))
func WithTimeFormat(format string) Option {
	return func(l *Logger) {
		l.timeFormat = format
	}
}

// WithMaxSize configures the maximum size of a log message. The default maximum size is 1000.
//
// Parameters:
//   - size: The maximum size of the log message.
//
// Example:
//
//	logger := loggo.New(loggo.LevelInfo, loggo.WithMaxSize(1000))
func WithMaxSize(size int) Option {
	return func(l *Logger) {
		l.maxSize = size
	}
}

// WithCallerProvider configures the caller provider function of a Logger. The default caller provider is runtime.Caller.
//
// Parameters:
//   - provider: The CallerProvider function to use.
//
// Example:
//
//	logger := loggo.New(loggo.LevelInfo, loggo.WithCallerProvider(func(skip int) (pc uintptr, file string, line int, ok bool) {
//		return runtime.Caller(skip)
//	}))
func WithCallerProvider(provider CallerProvider) Option {
	return func(l *Logger) {
		l.callerProvider = provider
	}
}

// WithContext configures the context of a Logger. The default context is context.Background.
//
// Parameters:
//   - Context: The context to use.
//
// Example:
//
//	logger := loggo.New(loggo.LevelInfo, loggo.WithContext(context.Background()))
func WithContext(ctx context.Context) Option {
	return func(l *Logger) {
		l.Context = ctx
	}
}

// WithPreHook adds a pre-hook to a Logger. Pre-hooks are executed before logging a message.
//
// Parameters:
//   - hook: The pre-hook function to add.
//
// Example:
//
//	logger := loggo.New(loggo.LevelInfo, loggo.WithPreHook(func(Context context.Context, level loggo.Level, message string) {
//		// Do something before logging the message
//	}))
func WithPreHook(hook Hook) Option {
	return func(l *Logger) {
		l.preHooks = append(l.preHooks, hook)
	}
}

// WithPostHook adds a post-hook to a Logger. Post-hooks are executed after logging a message.
//
// Parameters:
//   - hook: The post-hook function to add.
//
// Example:
//
//	logger := loggo.New(loggo.LevelInfo, loggo.WithPostHook(func(Context context.Context, level loggo.Level, message string) {
//		// Do something after logging the message
//	}))
func WithPostHook(hook Hook) Option {
	return func(l *Logger) {
		l.postHooks = append(l.postHooks, hook)
	}
}
