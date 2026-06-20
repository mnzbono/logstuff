package logstuff

import (
	"context"
)

type Logger struct {
	sink LogSink
}

func NewLogger(l LogSink) *Logger {
	return &Logger{sink: l}
}

func (l *Logger) With(args ...any) *Logger {
	return &Logger{
		sink: l.sink.With(args...),
	}
}
func (l *Logger) Log(ctx context.Context, lvl LogLevel, msg string, args ...any) {
	l.sink.Log(ctx, lvl, msg, args...)
}

func (l *Logger) Enabled(ctx context.Context, lvl LogLevel) bool {
	if e, ok := l.sink.(levelEnabler); ok {
		return e.Enabled(ctx, lvl)
	}
	return true // sink doesn't support checking, assume enabled
}

func (l *Logger) Trace(msg string, args ...any) { l.TraceContext(context.Background(), msg, args...) }
func (l *Logger) TraceContext(ctx context.Context, msg string, args ...any) {
	l.sink.Log(ctx, LevelTrace, msg, args...)
}
func (l *Logger) Debug(msg string, args ...any) { l.DebugContext(context.Background(), msg, args...) }
func (l *Logger) DebugContext(ctx context.Context, msg string, args ...any) {
	l.sink.Log(ctx, LevelDebug, msg, args...)
}
func (l *Logger) Info(msg string, args ...any) { l.InfoContext(context.Background(), msg, args...) }
func (l *Logger) InfoContext(ctx context.Context, msg string, args ...any) {
	l.sink.Log(ctx, LevelInfo, msg, args...)
}
func (l *Logger) Warn(msg string, args ...any) { l.WarnContext(context.Background(), msg, args...) }
func (l *Logger) WarnContext(ctx context.Context, msg string, args ...any) {
	l.sink.Log(ctx, LevelWarn, msg, args...)
}
func (l *Logger) Error(msg string, args ...any) { l.ErrorContext(context.Background(), msg, args...) }
func (l *Logger) ErrorContext(ctx context.Context, msg string, args ...any) {
	l.sink.Log(ctx, LevelError, msg, args...)
}
