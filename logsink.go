// Package logstuff provides a few minimal logging abstractions
//   - LogSink, minimal log interface for logging swapability.
//   - SlogCtx, slog wrapper that adds ctx key propagation and trace level.
//   - slogSink, the LogSink adapter for SlogCtx.
//   - Logger, adds logging cenvenience methods to a LogSink.
//   - CtxKeyer, LevelEnbaler and more.
package logstuff

import (
	"context"
	"strconv"
)

// LogSink is a minimal interface that helps logger swapping.
type LogSink interface {
	Log(context.Context, LogLevel, string, ...any)
	With(...any) LogSink
}

type LogLevel int

const (
	LevelTrace LogLevel = -8
	LevelDebug LogLevel = -4
	LevelInfo  LogLevel = 0
	LevelWarn  LogLevel = 4
	LevelError LogLevel = 8
)

func (l LogLevel) String() string {
	switch l {
	case LevelTrace:
		return "TRACE"
	case LevelDebug:
		return "DEBUG"
	case LevelInfo:
		return "INFO"
	case LevelWarn:
		return "WARN"
	case LevelError:
		return "ERROR"
	default:
		return "UNKNOWNLOGLVL<" + strconv.Itoa(int(l)) + ">"
	}
}
