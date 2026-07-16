[![tests/vet](https://github.com/mnzbono/logstuff/actions/workflows/ci.yml/badge.svg)](https://github.com/mnzbono/logstuff/actions/workflows/ci.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Reference](https://pkg.go.dev/badge/github.com/mnzbono/logstuff.svg)](https://pkg.go.dev/github.com/mnzbono/logstuff)

# logstuff

A few small logging abstractions I use across my own Go projects: 
- a minimal swappable logging interface, 
- a slog wrapper with context-key propagation and trace level, 
- and a couple of optional capability interfaces (LevelEnabler, CtxKeyer).

## LogSink
`LogSink` is a minimal logger interface. It only has two methods, `Log` and `With`, and uses exported `LogLevel` constants for the level argument. This keeps call sites light and behind interface, making it easy to modify logger setup without touching call sites.

```go
type LogSink interface {
	Log(context.Context, LogLevel, string, ...any)
	With(...any) LogSink
}
```
#### Why `LogSink` at all?
Short answer: refactor pain insurance.

This started as a direct slog dependency, then got refactored into a custom logger, and the coupling between call sites and the concrete logger caused friction during that move.

LogSink exists so call sites depend on a two-method shape, not a specific implementation, so a future swap (by me, or by anyone forking/collaborating who already has their own logging setup) doesn't mean touching every call site again.

## A few more things in this package, just convenience stuff:
- `SlogCtx`
    - It's a `slog` wrapper that adds context keys (useful for request-id, traces or similar).
    - It also adds `trace` level.
- `slogSink`
    - It's just a ready-made `SlogCtx` wrapper that satisfies the LogSink interface.
- `Logger`
    - LogSink wrapper.
    - Provides all the convenience methods on top of a LogSink, `Debug`, `Info`, `Error`, etc.
    - It can be used with any LogSink but it doesn't satisfy LogSink interface itself.
- `discardSink`
    - is just a ready-made implementation of a discard logger (a logger that doesn't do anything) that satisfies the LogSink interface.
- `LevelFromEnv`
    - Package function for reading log level from env, returns directly a `logstuff.LogLevel`.

Building-block interfaces:
- `CtxKeyer`
    - minimal interface to add ctx keys for the logger to track. `WithCtxKeys` method.
- `LevelEnabler`
    - minimal interface to add finer control over what's logged. `Enabled` method.

## Install
```sh
# Requires Go 1.26+ (uses the stdlib slog.MultiHandler)
go get github.com/mnzbono/logstuff
```

## Usage example
```go
// ...
type MyThingy struct {
    logger logstuff.LogSink // this is the interface that allows swapping
    // ...
}
// ...
func NewMyThingy() *MyThingy {
    logger := logstuff.NewSlogSink(logstuff.LevelInfo) // you can use the ready-made adapters or quickly roll your own.
    return &MyThingy{logger:logger}
}
// ...
func (m *MyThingy) Work() {
    m.logger.Log(ctx, logstuff.LevelInfo, "just logging around...")
    // ...
}
// ...
func TestMyThingy() {
    m := &MyThingy{logger: logstuff.NewDiscard() } // convenience discard logger while avoiding nil
    // ...
}
```

## `SlogCtx` & `slogSink`
`SlogCtx` wraps slog and adds ability to track context keys.  
If it isn't tracking any key it's basically slog, but without any methods, only Log and With.  
It can use any `io.Writer`, `text`/`json`, or prebuilt handler (slog capabilities). If you don't provide any handler argument it builds a texthandler that outputs to `os.Stdout`. If you provide more than one handler it creates the new slog `MultiHandler`.

`slogSink` just wraps it so it can be used as a LogSink.

example:
```go
// this example will create a slog.MultiHandler as a LogSink
logger := logstuff.NewSlogSink(
    logstuff.LevelDebug, 
    logstuff.WithJSONHandler(w), // you can write json to a buffer/file 
    logstuff.WithTextHandler(os.Stderr, logstuff.LevelError), // error only to stderr
    logstuff.WithHandler(myCustomdHandler) // or pass slog.Handler you already built.
)
// logger (slogSink) satisfies LogSink interface
logger.Log(ctx, logstuff.LevelInfo, "just logging around...")
```

context keys:
```go
// good hygiene to set the keys in a single place
var ridKey logstuff.CtxKey = "request-id"

// we create a slogSink and set it to track the ctx key "request-id"
// we use interface assertion in this example, but you can compose an interface with CtxKeyer
baseLogger := logstuff.NewSlogSink(logstuff.LevelTrace)
ctxLogger := baseLogger.(logstuff.CtxKeyer).WithCtxKeys(
    ridKey,
    logstuff.CtxKey("user"), // variadic — pass any number of keys, deduped before saving
    )
// ...
func (h *Handler) HandleRequest(ctx context.Context) {
    ctx = context.WithValue(ctx, ridKey, newRequestId()) // request-id=123
    h.service.Work(ctx)
}
// if s.logger is set up to track ctx keys it captures them.
func (s *Service) Work(ctx context.Context) {
    s.logger.Log(ctx, logstuff.LevelDebug, "working...") // request-id=123
}
```
## Wait, `Logger` doesn't satisfy `LogSink`?
No. LogSink is meant to be minimal, so you can quickly write a wrapper.  
Logger is meant for when you want the convenience methods back. It wraps a LogSink (which only has `Log` method) and provides all the `Debug`, `TraceCtx`, `ErrorCtx`, etc, methods.
easier to see some code:
```go
// ...
type MyThingy struct {
    logger *logstuff.Logger // this will wrap the LogSink
    // ...
}
// ...
func NewMyThingy() *MyThingy {
    // imagine you build a logsink adapter to some kind of logger
    myAdapter := NewMyLoggerAdapter() 
    // you just wrap it with Logger and it gives you the convenience methods instead of just Log()
    logger := logstuff.NewLogger(myAdapter) 
    return &MyThingy{logger:logger}
}
// ...
func (m *MyThingy) Work() {
    m.logger.Debug("just debugging around...")
    m.logger.TraceContext(ctx, "just tracing around...")
    // etc ...
}
```
---
*AI disclaimer: AI was used for code review, documentation review and english phrasing, and wrote a handful of tests. Design, implementation, all code, and most tests are human-made.*