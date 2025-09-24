package logger

import (
	"log/slog"
	"os"
)

var logger = slog.New(
	slog.NewTextHandler(
		os.Stdout,
		&slog.HandlerOptions{
			Level: slog.LevelInfo},
	))

// getEnvAndLog retrieves the value of an environment variable and logs if it is not set.
func GetEnvVarAndLog(key string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		logger.Warn("Environment variable is not set", "variable", key)
	}
	return value
}

func Error(msg string, args ...any) {
	logger.Error(msg, args...)
}

func Warn(msg string, args ...any) {
	logger.Error(msg, args...)
}

func Info(msg string, args ...any) {
	logger.Error(msg, args...)
}

func Fatal(msg string, args ...any) {
	logger.Error(msg, args...)
	os.Exit(1)
}
