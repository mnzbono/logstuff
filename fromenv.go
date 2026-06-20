package logstuff

import (
	"os"
	"strings"
)

// "LOGLEVEL"
var EnvLogLevel = "LOGLEVEL"

// reads var EnvLogLevel, default "LOGLEVEL"
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
