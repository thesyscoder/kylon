/**
 * @File: logger.go
 * @Title: Centralized Logrus Logger Configuration
 * @Description: Configures and provides a singleton Logrus logger instance for the application,
 * @Description: ensuring consistent logging behavior, format, and level management across all components.
 * @Author: thesyscoder (github.com/thesyscoder)
 */

package logger

import "github.com/sirupsen/logrus"

// log is the package-level singleton Logrus logger instance.
var log *logrus.Logger

// init function is automatically called once when the package is imported.
// It initializes the logger with a default configuration.
func init() {
	log = logrus.New()
	// Sets the default logging level to Info, meaning logs at Info, Warn, Error, Fatal, and Panic levels will be output.
	log.SetLevel(logrus.InfoLevel)
	// Configures the formatter to produce human-readable text output with full timestamps.
	log.SetFormatter(&logrus.JSONFormatter{
		PrettyPrint: true,
	})
	// Enables reporting of the file and line number where the log call originated, aiding in debugging.
	log.SetReportCaller(true)
	// Directs logger output to the standard logger's output, which is typically os.Stderr by default.
	log.SetOutput(logrus.StandardLogger().Out)

	// Logs an initial message to confirm logger setup, including a service identifier for context.
	log.WithFields(logrus.Fields{"service": "kylon"}).Info("Logger initialized with default settings")
}

// GetLogger returns the global singleton Logrus logger instance.
// This function provides a consistent way for other packages to obtain and use the pre-configured logger.
func GetLogger() *logrus.Logger {
	return log
}

// SetLogger dynamically configures the global logger's minimum output level.
// It parses the provided `level` string (e.g., "info", "debug", "error").
// If the `level` string is invalid or unparseable, it defaults the logger's level to InfoLevel
// and logs an error message.
func SetLogger(level string) {
	parsedLevel, err := logrus.ParseLevel(level)
	if err != nil {
		// Logs an error if the provided level string is invalid, and falls back to InfoLevel.
		log.WithFields(logrus.Fields{"error": err.Error(), "provided_level": level}).
			Error("Failed to parse log level, using default InfoLevel")
		parsedLevel = logrus.InfoLevel // Fallback to InfoLevel
	}
	// Applies the determined log level to the global logger.
	log.SetLevel(parsedLevel)
	// Logs a confirmation message indicating the new active log level.
	log.WithFields(logrus.Fields{"new_level": parsedLevel.String()}).Info("Log level set")
}
