package logstuff

// CtxKeyer is a minimal interface that helps fork a logger and add ctx key capture to it.
type CtxKeyer interface {
	WithCtxKeys(...CtxKey) LogSink
}
