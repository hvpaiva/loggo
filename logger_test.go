package loggo_test

import (
	"strings"
	"testing"
	"time"

	"github.com/hvpaiva/loggo"
)

var fakeNow = func() time.Time {
	return time.Date(2022, 1, 25, 0, 0, 0, 0, time.UTC)
}
var fakeNowString = "2022-01-25 00:00:00"

var okCallerProvider = func() (pc uintptr, file string, line int, ok bool) {
	return 0, "file", 1, true
}

var errorCallerProvider = func() (pc uintptr, file string, line int, ok bool) {
	return 0, "", 0, false
}

func TestLogger_Log(t *testing.T) {
	type testCase struct {
		name      string
		threshold loggo.Level
		level     loggo.Level
		message   string
		want      string
	}

	testCases := []testCase{
		{
			name:      "info",
			threshold: loggo.LevelInfo,
			level:     loggo.LevelInfo,
			message:   "This is an info log message",
			want:      fakeNowString + " [ INFO]: This is an info log message\n",
		},
		{
			name:      "info, but in debug threshold",
			threshold: loggo.LevelDebug,
			level:     loggo.LevelInfo,
			message:   "This is an info log message",
			want:      fakeNowString + " [ INFO]: This is an info log message\n",
		},
		{
			name:      "fatal, but in debug threshold",
			threshold: loggo.LevelDebug,
			level:     loggo.LevelFatal,
			message:   "This is an fatal log message",
			want:      fakeNowString + " [FATAL]: This is an fatal log message\n",
		},
		{
			name:      "error, but in fatal threshold",
			threshold: loggo.LevelFatal,
			level:     loggo.LevelError,
			message:   "This is an error log message",
			want:      "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			w := &strings.Builder{}
			logger := loggo.New(tc.threshold, loggo.WithOutput(w), loggo.WithTimeProvider(fakeNow))
			logger.Log(tc.level, tc.message)

			if w.String() != tc.want {
				t.Errorf("Logger.Log() = %q, want %q", w.String(), tc.want)
			}
		})
	}
}

func TestLogger_Logf(t *testing.T) {
	type testCase struct {
		name      string
		threshold loggo.Level
		level     loggo.Level
		message   string
		args      []any
		want      string
	}

	testCases := []testCase{
		{
			name:      "info",
			threshold: loggo.LevelInfo,
			level:     loggo.LevelInfo,
			message:   "This is an info %s log message",
			args:      []any{"format"},
			want:      fakeNowString + " [ INFO]: This is an info format log message\n",
		},
		{
			name:      "info, but in debug threshold",
			threshold: loggo.LevelDebug,
			level:     loggo.LevelInfo,
			message:   "This is an info %s log message",
			args:      []any{"format"},
			want:      fakeNowString + " [ INFO]: This is an info format log message\n",
		},
		{
			name:      "fatal, but in debug threshold",
			threshold: loggo.LevelDebug,
			level:     loggo.LevelFatal,
			message:   "This is an fatal %s log message",
			args:      []any{"format"},
			want:      fakeNowString + " [FATAL]: This is an fatal format log message\n",
		},
		{
			name:      "error, but in fatal threshold",
			threshold: loggo.LevelFatal,
			level:     loggo.LevelError,
			message:   "This is an error %s log message",
			args:      []any{"format"},
			want:      "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			w := &strings.Builder{}
			logger := loggo.New(tc.threshold, loggo.WithOutput(w), loggo.WithTimeProvider(fakeNow))
			logger.Logf(tc.level, tc.message, tc.args...)

			if w.String() != tc.want {
				t.Errorf("Logger.Logf() = %q, want %q", w.String(), tc.want)
			}
		})
	}
}

func TestLogger(t *testing.T) {
	type testCase struct {
		name    string
		logger  func(message string)
		message string
		want    string
	}

	w := &strings.Builder{}
	log := loggo.New(loggo.LevelInfo, loggo.WithOutput(w), loggo.WithTimeProvider(fakeNow))

	testCases := []testCase{
		{
			name:    "debug",
			logger:  log.Debug,
			message: "This is an debug log message",
			want:    "",
		},
		{
			name:    "info",
			logger:  log.Info,
			message: "This is an info log message",
			want:    fakeNowString + " [ INFO]: This is an info log message\n",
		},
		{
			name:    "warn",
			logger:  log.Warn,
			message: "This is an warn log message",
			want:    fakeNowString + " [ WARN]: This is an warn log message\n",
		},
		{
			name:    "error",
			logger:  log.Error,
			message: "This is an error log message",
			want:    fakeNowString + " [ERROR]: This is an error log message\n",
		},
		{
			name:    "fatal",
			logger:  log.Fatal,
			message: "This is an fatal log message",
			want:    fakeNowString + " [FATAL]: This is an fatal log message\n",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.logger(tc.message)

			if w.String() != tc.want {
				w.Reset()
				t.Errorf("Logger.Log() = %q, want %q", w.String(), tc.want)
			}

			w.Reset()
		})
	}
}

func TestLogger_format(t *testing.T) {
	type testCase struct {
		name    string
		logger  func(message string, args ...any)
		message string
		args    []any
		want    string
	}

	w := &strings.Builder{}
	log := loggo.New(loggo.LevelInfo, loggo.WithOutput(w), loggo.WithTimeProvider(fakeNow))

	testCases := []testCase{
		{
			name:    "debug",
			logger:  log.Debugf,
			message: "This is an debug log message %s",
			args:    []any{"format"},
			want:    "",
		},
		{
			name:    "info",
			logger:  log.Infof,
			message: "This is an info log message %s",
			args:    []any{"format"},
			want:    fakeNowString + " [ INFO]: This is an info log message format\n",
		},
		{
			name:    "warn",
			logger:  log.Warnf,
			message: "This is an warn log message %s",
			args:    []any{"format"},
			want:    fakeNowString + " [ WARN]: This is an warn log message format\n",
		},
		{
			name:    "error",
			logger:  log.Errorf,
			message: "This is an error log message %s",
			args:    []any{"format"},
			want:    fakeNowString + " [ERROR]: This is an error log message format\n",
		},
		{
			name:    "fatal",
			logger:  log.Fatalf,
			message: "This is an fatal log message %s",
			args:    []any{"format"},
			want:    fakeNowString + " [FATAL]: This is an fatal log message format\n",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.logger(tc.message, tc.args...)

			if w.String() != tc.want {
				w.Reset()
				t.Errorf("Logger.Log() = %q, want %q", w.String(), tc.want)
			}

			w.Reset()
		})
	}
}

func TestLogger_LogfE(t *testing.T) {
	type testCase struct {
		name     string
		message  string
		template string
		wantErr  string
	}

	testCases := []testCase{
		{
			name:     "error parsing template",
			message:  "This is an info log message",
			template: "{{.Level",
			wantErr:  "error parsing template: template: log:2: unclosed action started at log:1",
		},
		{
			name:     "error executing template",
			message:  "This is an info log message",
			template: "{{.SomeField}}",
			wantErr:  "error executing template: template: log:1:2: executing \"log\" at <.SomeField>: can't evaluate field SomeField in type loggo.logData",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			w := &strings.Builder{}
			logger := loggo.New(loggo.LevelInfo, loggo.WithOutput(w), loggo.WithTimeProvider(fakeNow), loggo.WithTemplate(tc.template))
			err := logger.LogfE(loggo.LevelInfo, tc.message)

			if err == nil {
				t.Errorf("Logger.LogfE() was nil, wantErr %q", tc.wantErr)
			}

			if err.Error() != tc.wantErr {
				t.Errorf("Logger.LogfE() error = %q, wantErr %q", err, tc.wantErr)
			}

		})
	}
}

func TestLogger_Log_unknownCaller(t *testing.T) {
	w := &strings.Builder{}
	logger := loggo.New(loggo.LevelInfo, loggo.WithOutput(w), loggo.WithTimeProvider(fakeNow), loggo.WithTemplate("{{.Caller}}"))
	logger.Log(loggo.LevelInfo, "This is an info log message")

}

func ExampleLogger_Log() {
	logger := loggo.New(loggo.LevelInfo, loggo.WithTimeProvider(fakeNow))
	logger.Log(loggo.LevelInfo, "This is an info log message")
	// Output: 2022-01-25 00:00:00 [ INFO]: This is an info log message
}

func ExampleLogger_Logf() {
	logger := loggo.New(loggo.LevelInfo, loggo.WithTimeProvider(fakeNow))
	logger.Logf(loggo.LevelInfo, "This is an info log message with a %q", "format")
	// Output: 2022-01-25 00:00:00 [ INFO]: This is an info log message with a "format"
}

func ExampleLogger_Log_threshold() {
	logger := loggo.New(loggo.LevelWarn, loggo.WithTimeProvider(fakeNow))
	logger.Log(loggo.LevelInfo, "This is an info log message")
	// Output:
}

func ExampleLogger_Logf_threshold() {
	logger := loggo.New(loggo.LevelWarn, loggo.WithTimeProvider(fakeNow))
	logger.Logf(loggo.LevelInfo, "This is an info log message with a %s", "format")
	// Output:
}

func ExampleLogger_Debug() {
	logger := loggo.New(loggo.LevelDebug, loggo.WithTimeProvider(fakeNow))
	logger.Debug("This is a debug log message")
	// Output: 2022-01-25 00:00:00 [DEBUG]: This is a debug log message
}

func ExampleLogger_Debugf() {
	logger := loggo.New(loggo.LevelDebug, loggo.WithTimeProvider(fakeNow))
	logger.Debugf("This is a debug log message with a %q", "format")
	// Output: 2022-01-25 00:00:00 [DEBUG]: This is a debug log message with a "format"
}

func ExampleLogger_Info() {
	logger := loggo.New(loggo.LevelInfo, loggo.WithTimeProvider(fakeNow))
	logger.Info("This is an info log message")
	// Output: 2022-01-25 00:00:00 [ INFO]: This is an info log message
}

func ExampleLogger_Infof() {
	logger := loggo.New(loggo.LevelInfo, loggo.WithTimeProvider(fakeNow))
	logger.Infof("This is an info log message with a %q", "format")
	// Output: 2022-01-25 00:00:00 [ INFO]: This is an info log message with a "format"
}

func ExampleLogger_Warn() {
	logger := loggo.New(loggo.LevelWarn, loggo.WithTimeProvider(fakeNow))
	logger.Warn("This is a warn log message")
	// Output: 2022-01-25 00:00:00 [ WARN]: This is a warn log message
}

func ExampleLogger_Warnf() {
	logger := loggo.New(loggo.LevelWarn, loggo.WithTimeProvider(fakeNow))
	logger.Warnf("This is a warn log message with a %q", "format")
	// Output: 2022-01-25 00:00:00 [ WARN]: This is a warn log message with a "format"
}

func ExampleLogger_Error() {
	logger := loggo.New(loggo.LevelError, loggo.WithTimeProvider(fakeNow))
	logger.Error("This is an error log message")
	// Output: 2022-01-25 00:00:00 [ERROR]: This is an error log message
}

func ExampleLogger_Errorf() {
	logger := loggo.New(loggo.LevelError, loggo.WithTimeProvider(fakeNow))
	logger.Errorf("This is an error log message with a %q", "format")
	// Output: 2022-01-25 00:00:00 [ERROR]: This is an error log message with a "format"
}

func ExampleLogger_Fatal() {
	logger := loggo.New(loggo.LevelFatal, loggo.WithTimeProvider(fakeNow))
	logger.Fatal("This is a fatal log message")
	// Output: 2022-01-25 00:00:00 [FATAL]: This is a fatal log message
}

func ExampleLogger_Fatalf() {
	logger := loggo.New(loggo.LevelFatal, loggo.WithTimeProvider(fakeNow))
	logger.Fatalf("This is a fatal log message with a %q", "format")
	// Output: 2022-01-25 00:00:00 [FATAL]: This is a fatal log message with a "format"
}

func ExampleLogger_Log_maxSize() {
	logger := loggo.New(loggo.LevelInfo, loggo.WithTimeProvider(fakeNow), loggo.WithMaxSize(10))
	logger.Log(loggo.LevelInfo, "This is an info log message")
	// Output: 2022-01-25 00:00:00 [ INFO]: This is an
}

func ExampleLogger_Log_template() {
	logger := loggo.New(loggo.LevelInfo, loggo.WithTimeProvider(fakeNow), loggo.WithTemplate("[{{.Level}}]: {{.Message}}"))
	logger.Log(loggo.LevelInfo, "This is an info log message")
	// Output: [INFO]: This is an info log message
}

func ExampleLogger_Log_templateComplete() {
	logger := loggo.New(
		loggo.LevelInfo,
		loggo.WithTimeProvider(fakeNow),
		loggo.WithTemplate("{{.Time}} {{.Caller}} | [{{.Level}}]: {{.Message}}"),
		loggo.WithCallerProvider(okCallerProvider),
	)
	logger.Log(loggo.LevelInfo, "This is an info log message")
	// Output: 2022-01-25 00:00:00 file:1 | [INFO]: This is an info log message
}

func ExampleLogger_Log_callerProviderErr() {
	logger := loggo.New(
		loggo.LevelInfo,
		loggo.WithTimeProvider(fakeNow),
		loggo.WithTemplate("{{.Caller}} [{{.Level}}]: {{.Message}}"),
		loggo.WithCallerProvider(errorCallerProvider),
	)
	logger.Log(loggo.LevelInfo, "This is an info log message")
	// Output: unknown [INFO]: This is an info log message
}

func ExampleLogger_Log_timeFormat() {
	logger := loggo.New(loggo.LevelInfo, loggo.WithTimeProvider(fakeNow), loggo.WithTimeFormat("02/01/06 00:00"))
	logger.Log(loggo.LevelInfo, "This is an info log message")
	// Output: 25/01/22 00:00 [ INFO]: This is an info log message
}
