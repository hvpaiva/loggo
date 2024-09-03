package loggo

import (
	"fmt"
)

// templateData is a structure that holds the data for a log message template.
type templateData struct {
	Level   string
	Time    string
	Message string
	Caller  string
}

// getTemplateData returns the data for a log message template.
func getTemplateData(level Level, message string, logger *Logger) templateData {
	data := templateData{
		Level:   level.String(),
		Time:    logger.now().Format(logger.timeFormat),
		Message: truncateString(message, logger.maxSize),
		Caller:  getCaller(logger.callerProvider),
	}

	return data
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
