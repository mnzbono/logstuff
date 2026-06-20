package logstuff

import (
	"context"
	"log/slog"
	"os"
)

type CtxKey string

type slogCtx struct {
	*slog.Logger
	keys []CtxKey
}

func NewSlogCtx(l LogLevel) *slogCtx {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.Level(l),
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.LevelKey {
				level := a.Value.Any().(slog.Level)

				switch LogLevel(level) {
				case LevelTrace:
					a.Value = slog.StringValue("TRACE")
				case LevelDebug:
					a.Value = slog.StringValue("DEBUG")
				case LevelInfo:
					a.Value = slog.StringValue("INFO")
				case LevelWarn:
					a.Value = slog.StringValue("WARN")
				case LevelError:
					a.Value = slog.StringValue("ERROR")
				}
			}
			return a
		},
	}))
	return &slogCtx{Logger: logger}
}

func (l *slogCtx) With(args ...any) *slogCtx {
	return &slogCtx{
		Logger: l.Logger.With(args...),
		keys:   l.keys,
	}
}
func (l *slogCtx) Log(ctx context.Context, lvl LogLevel, msg string, args ...any) {
	l.Logger.Log(ctx, slog.Level(lvl), msg, l.appendCtx(ctx, args)...)
}

func (l *slogCtx) Enabled(ctx context.Context, lvl LogLevel) bool {
	return l.Logger.Enabled(ctx, slog.Level(lvl))
}

func (l *slogCtx) WithCtxKeys(args ...CtxKey) *slogCtx {
	lenArgs := len(args)
	if lenArgs < 1 {
		return l
	}
	totalLength := len(l.keys) + lenArgs
	newKeys := make([]CtxKey, 0, totalLength)
	seen := make(map[CtxKey]struct{}, totalLength)
	for _, v := range l.keys {
		if _, ok := seen[v]; ok {
			continue
		}
		newKeys = append(newKeys, v)
		seen[v] = struct{}{}
	}
	for _, v := range args {
		if _, ok := seen[v]; ok {
			continue
		}
		newKeys = append(newKeys, v)
		seen[v] = struct{}{}
	}
	return &slogCtx{
		Logger: l.Logger,
		keys:   newKeys,
	}
}

func (l *slogCtx) appendCtx(ctx context.Context, args []any) []any {
	if len(l.keys) == 0 {
		return args
	}
	for _, key := range l.keys {
		if inCtx := ctx.Value(key); inCtx != nil {
			args = append(args, string(key), inCtx)
		}
	}
	return args
}
