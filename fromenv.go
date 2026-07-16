package logstuff

import (
	"os"
	"strings"
)

// EnvLogLevel defines the env var to use when setting log level, default "LOGLEVEL"
var EnvLogLevel = "LOGLEVEL"

// LevelFromEnv reads log level form env. Reads exported var EnvLogLevel value, default "LOGLEVEL"
func LevelFromEnv() LogLevel {
	envlevel := strings.ToLower(strings.TrimSpace(os.Getenv(EnvLogLevel)))
	switch envlevel {
	case "trace":
		return LevelTrace
	case "debug":
		return LevelDebug
	case "warn":
		return LevelWarn
	case "error":
		return LevelError
	default:
		return LevelInfo
	}
}
