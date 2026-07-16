package logstuff

import (
	"context"
	"io"
	"log/slog"
	"os"
)

// CtxKey is the type of context keys that SlogCtx can track.
type CtxKey string

// SlogCtx is a wrapper around slog.Logger that adds support for context keys.
type SlogCtx struct {
	logger *slog.Logger // slog.Logger private so its methods aren't promoted (breaks the ctx thing)
	keys   []CtxKey
}

// SlogHandler builds a slog.Handler from the given options.
type SlogHandler func(*slog.HandlerOptions) slog.Handler

// WithJSONHandler creates a JSON slog.Handler. If logLevel is passed it overwrites the base logLevel of constructor.
func WithJSONHandler(w io.Writer, logLevel ...LogLevel) SlogHandler {
	return func(opts *slog.HandlerOptions) slog.Handler {
		myOpts := *opts
		if len(logLevel) > 0 {
			myOpts.Level = slog.Level(logLevel[0])
		}
		return slog.NewJSONHandler(w, &myOpts)
	}
}

// WithTextHandler creates a TEXT slog.Handler. If logLevel is passed it overwrites the base logLevel of constructor.
func WithTextHandler(w io.Writer, logLevel ...LogLevel) SlogHandler {
	return func(opts *slog.HandlerOptions) slog.Handler {
		myOpts := *opts
		if len(logLevel) > 0 {
			myOpts.Level = slog.Level(logLevel[0])
		}
		return slog.NewTextHandler(w, &myOpts)
	}
}

// WithHandler allows to passing directly a pre-built slog.Handler.
// NOTE: Handlers passed in WithHandler won't do ReplaceAttr (trace level).
func WithHandler(h slog.Handler) SlogHandler {
	return func(opts *slog.HandlerOptions) slog.Handler { return h }
}

// NewSlogCtx returns a *SlogCtx. variadic SlogHandler allows to pass 0, 1 or more build options.
//   - If no SlogHandler is passed, the default TEXT slog.Handler that outputs to stdout will be used.
//   - If more than one is passed, then a slog.Multihandler will be created.
func NewSlogCtx(l LogLevel, opts ...SlogHandler) *SlogCtx {
	var h slog.Handler
	handlerOpts := &slog.HandlerOptions{Level: slog.Level(l), ReplaceAttr: replaceLevelAttr}
	switch len(opts) {
	case 0:
		h = slog.NewTextHandler(os.Stdout, handlerOpts)
	case 1:
		h = opts[0](handlerOpts)
	default:
		handlers := make([]slog.Handler, 0, len(opts))
		for _, v := range opts {
			handlers = append(handlers, v(handlerOpts))
		}
		h = slog.NewMultiHandler(handlers...)
	}
	return &SlogCtx{logger: slog.New(h)}
}

func replaceLevelAttr(groups []string, a slog.Attr) slog.Attr {
	if a.Key == slog.LevelKey {
		level := a.Value.Any().(slog.Level)
		a.Value = slog.StringValue(
			LogLevel(level).String(),
		)
	}
	return a
}

func (l *SlogCtx) With(args ...any) *SlogCtx {
	return &SlogCtx{
		logger: l.logger.With(args...),
		keys:   l.keys,
	}
}
func (l *SlogCtx) Log(ctx context.Context, lvl LogLevel, msg string, args ...any) {
	l.logger.Log(ctx, slog.Level(lvl), msg, l.appendCtx(ctx, args)...)
}

func (l *SlogCtx) Enabled(ctx context.Context, lvl LogLevel) bool {
	return l.logger.Enabled(ctx, slog.Level(lvl))
}

func (l *SlogCtx) WithCtxKeys(args ...CtxKey) *SlogCtx {
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
	return &SlogCtx{
		logger: l.logger,
		keys:   newKeys,
	}
}

// appendCtx looks for tracked CtxKeys and adds them to the logger args.
func (l *SlogCtx) appendCtx(ctx context.Context, args []any) []any {
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
