package logstuff

import (
	"context"
	"log/slog"
)

type slogSink struct {
	*slogCtx
}

func NewSlogSink(l LogLevel) LogSink {
	return &slogSink{NewSlogCtx(l)}
}

func (l *slogSink) With(args ...any) LogSink {
	return &slogSink{
		&slogCtx{
			Logger: l.Logger.With(args...),
			keys:   l.keys,
		},
	}
}
func (l *slogSink) Log(ctx context.Context, lvl LogLevel, msg string, args ...any) {
	l.Logger.Log(ctx, slog.Level(lvl), msg, l.appendCtx(ctx, args)...)
}
