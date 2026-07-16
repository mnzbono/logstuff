package logstuff

// slogSink is a wrapper around SlogCtx that satisfies LogSink interface.
type slogSink struct {
	SlogCtx
}

// NewSlogSink returns a wrapper around SlogCtx that satisfies LogSink interface.
func NewSlogSink(l LogLevel, opts ...SlogHandler) LogSink {
	return &slogSink{*NewSlogCtx(l, opts...)}
}

// --- LogSink interface

// Log() is promoted from SlogCtx embedding and already satisfies LogSink.

func (l *slogSink) With(args ...any) LogSink {
	if len(args) == 0 {
		return l
	}
	return &slogSink{
		*l.SlogCtx.With(args...)}
}

// --- optional interfaces

func (l *slogSink) WithCtxKeys(keys ...CtxKey) LogSink {
	if len(keys) == 0 {
		return l
	}
	return &slogSink{
		*l.SlogCtx.WithCtxKeys(keys...)}
}

// Enabled() is promoted from SlogCtx embedding and already satisfies LogSink.
