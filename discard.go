package logstuff

import "context"

type discardSink struct{}

// NewDiscard returns discardSink, a logger that does nothing and satisfies LogSink interface.
func NewDiscard() LogSink { return discardSink{} }

func (discardSink) With(...any) LogSink                           { return discardSink{} }
func (discardSink) Log(context.Context, LogLevel, string, ...any) {}

func (discardSink) Enabled(context.Context, LogLevel) bool { return false }