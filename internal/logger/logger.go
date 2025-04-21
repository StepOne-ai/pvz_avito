package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

// InitializeLogger initializes and configures the logger.
func InitializeLogger() *logrus.Logger {
	logger := logrus.New()

	// Set the log level (default: InfoLevel)
	logger.SetLevel(logrus.InfoLevel)

	// Use JSON formatter for structured logging
	logger.SetFormatter(&logrus.JSONFormatter{})

	// Output logs to stdout
	logger.SetOutput(os.Stdout)

	return logger
}
