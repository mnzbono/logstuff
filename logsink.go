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

const (
	levelBase = LevelTrace
	levelStep = 4
)

var levelNames = [...]string{
	"TRACE",
	"DEBUG",
	"INFO",
	"WARN",
	"ERROR",
}

func (l LogLevel) String() string {
	// a little ugly but inlines
	n := uint((l - levelBase) / levelStep)
	if n < uint(len(levelNames)) {
		return levelNames[n]
	}
	return l.string()
}

//go:noinline
func (l LogLevel) string() string {
	return "LVL<" + strconv.Itoa(int(l)) + ">"
}
