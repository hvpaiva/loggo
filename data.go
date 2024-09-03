package loggo

import (
	"fmt"
)

// logData is a structure that holds the data for a log message.
type logData struct {
	Level   string
	Time    string
	Message string
	Caller  string
}

// getData returns the data for a log message.
func getData(level Level, message string, logger *Logger) logData {
	return logData{
		Level:   level.String(),
		Time:    logger.now().Format(logger.timeFormat),
		Message: truncateString(message, logger.maxSize),
		Caller:  getCaller(logger.callerProvider),
	}
}

// getCaller returns the file and line number of the caller.
func getCaller(cp CallerProvider) string {
	_, file, line, ok := cp()
	if !ok {
		return "unknown"
	}

	return fmt.Sprintf("%s:%d", file, line)
}

// truncateString truncates the input string to the specified maxSize.
func truncateString(input string, maxSize int) string {
	if len(input) > maxSize {
		return input[:maxSize]
	}

	return input
}
