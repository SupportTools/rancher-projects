package logging

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/sirupsen/logrus"
)

var (
	logger *logrus.Logger
	once   sync.Once
)

// SetupLogging initializes the logger with the appropriate settings.
func SetupLogging() *logrus.Logger {
	once.Do(func() {
		logger = logrus.New()
		logger.SetOutput(os.Stdout)
		logger.SetReportCaller(true)

		// Set log level
		logLevel := getLogLevel(os.Getenv("LOG_LEVEL"))
		logger.SetLevel(logLevel)

		// Set log formatter
		logFormat := strings.ToLower(os.Getenv("LOG_FORMAT"))
		switch logFormat {
		case "json":
			logger.SetFormatter(&logrus.JSONFormatter{
				DisableTimestamp: false,
				PrettyPrint:      false,
			})
		default:
			logger.SetFormatter(&CustomTextFormatter{
				DisableTimestamp: false,
			})
		}

		logger.Debugf("Logger initialized with level: %s and format: %s", logLevel.String(), logFormat)
	})

	return logger
}

// CustomTextFormatter formats log entries in a more concise way.
type CustomTextFormatter struct {
	DisableTimestamp bool
}

func (f *CustomTextFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	caller := ""
	if entry.HasCaller() {
		file := filepath.Base(entry.Caller.File)
		caller = file + ":" + strconv.Itoa(entry.Caller.Line)
	}

	// Example format: level=info msg="Incoming request" func=handleRequest file=server.go:42
	logMessage := strings.Builder{}
	if !f.DisableTimestamp {
		logMessage.WriteString(entry.Time.Format("2006-01-02 15:04:05") + " ")
	}
	logMessage.WriteString("level=" + entry.Level.String())
	logMessage.WriteString(" msg=\"" + entry.Message + "\"")
	if caller != "" {
		logMessage.WriteString(" caller=" + caller)
	}
	logMessage.WriteString("\n")

	return []byte(logMessage.String()), nil
}

// getLogLevel returns the logrus log level based on the input string.
func getLogLevel(level string) logrus.Level {
	switch strings.ToLower(level) {
	case "debug":
		return logrus.DebugLevel
	case "warn", "warning":
		return logrus.WarnLevel
	case "error":
		return logrus.ErrorLevel
	case "fatal":
		return logrus.FatalLevel
	case "panic":
		return logrus.PanicLevel
	default:
		logrus.Warnf("Invalid LOG_LEVEL '%s'. Defaulting to InfoLevel.", level)
		return logrus.InfoLevel
	}
}
