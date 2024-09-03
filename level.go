package loggo

// Level represents an available log level.
//
// The log levels are ordered by severity, with LevelDebug being the lowest and LevelFatal being the highest.
// The levels are:
// - LevelDebug: Used for debugging purposes.
// - LevelInfo: Used to log general information about the application.
// - LevelWarn: Used to log warnings about potential issues.
// - LevelError: Used to log errors that do not cause the application to stop.
// - LevelFatal: Used to log fatal errors that cause the application to stop.
type Level byte

// Available log levels.
const (
	// LevelDebug is the lowest level and is mostly used for debugging purposes.
	LevelDebug Level = iota
	// LevelInfo is used to log general information about the application.
	LevelInfo
	// LevelWarn is used to log warnings about potential issues.
	LevelWarn
	// LevelError is used to log errors that do not cause the application to stop.
	LevelError
	// LevelFatal is used to log fatal errors that cause the application to stop.
	LevelFatal
)

// String returns the string representation of the log level.
func (l Level) String() string {
	return [...]string{"DEBUG", "INFO", "WARN", "ERROR", "FATAL"}[l]
}
