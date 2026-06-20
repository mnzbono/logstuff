package logstuff

import (
	"context"
)

type LogSink interface {
	With(...any) LogSink
	Log(context.Context, LogLevel, string, ...any)
}

type LogLevel int

const (
	LevelTrace LogLevel = -8
	LevelDebug LogLevel = -4
	LevelInfo  LogLevel = 0
	LevelWarn  LogLevel = 4
	LevelError LogLevel = 8
)

