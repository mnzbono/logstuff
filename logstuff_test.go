package logstuff_test

import (
	"bytes"
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/mnzbono/logstuff"
)

type testCall struct {
	ctx  context.Context
	lvl  logstuff.LogLevel
	msg  string
	args []any
}
type testSink struct {
	calls   []testCall
	with    [][]any
	initLvl logstuff.LogLevel
}

func (t *testSink) With(args ...any) logstuff.LogSink { t.with = append(t.with, args); return t }
func (t *testSink) Log(ctx context.Context, lvl logstuff.LogLevel, msg string, args ...any) {
	t.calls = append(t.calls, testCall{ctx: ctx, lvl: lvl, msg: msg, args: args})
}
func (t *testSink) Enabled(_ context.Context, lvl logstuff.LogLevel) bool { return lvl >= t.initLvl }

type notEnabled struct{}

func (t *notEnabled) With(...any) logstuff.LogSink                                            { return t }
func (t *notEnabled) Log(ctx context.Context, lvl logstuff.LogLevel, msg string, args ...any) {}

func TestSlogSink_Log(t *testing.T) {
	var buf bytes.Buffer
	sink := logstuff.NewSlogSink(logstuff.LevelInfo, logstuff.WithTextHandler(&buf))
	sink.Log(context.Background(), logstuff.LevelInfo, "hello", "key", "val")
	_, out, _ := strings.Cut(buf.String(), " ")
	if strings.TrimSpace(out) != `level=INFO msg=hello key=val` {
		t.Errorf("unexpected output: %q", out)
	}
}

func TestSlogSink_With(t *testing.T) {
	var extraKey, extraVal = "extra", "xxl"
	var buf bytes.Buffer
	sink := logstuff.NewSlogSink(logstuff.LevelInfo, logstuff.WithTextHandler(&buf)).With(extraKey, extraVal)
	sink.Log(context.Background(), logstuff.LevelInfo, "hello", "key", "val")
	_, out, _ := strings.Cut(buf.String(), " ")
	if !strings.Contains(buf.String(), extraKey+"="+extraVal) {
		t.Errorf("unexpected output: %q", out)
	}
}

func TestSlogSinkLevels(t *testing.T) {
	cases := map[string]logstuff.LogLevel{
		"TRACE": logstuff.LevelTrace, "DEBUG": logstuff.LevelDebug,
		"INFO": logstuff.LevelInfo,
		"WARN": logstuff.LevelWarn, "ERROR": logstuff.LevelError,
	}
	for k, v := range cases {
		var b bytes.Buffer
		logger := logstuff.NewSlogSink(v, logstuff.WithTextHandler(&b))
		logger.Log(context.Background(), v, k)
		out := b.String()
		if !strings.Contains(out, "level="+k) {
			t.Fatalf("unexpected output: %q", out)
		}
	}
}

func TestLevelFromEnv(t *testing.T) {
	cases := map[string]logstuff.LogLevel{
		"trace": logstuff.LevelTrace, "debug": logstuff.LevelDebug, "": logstuff.LevelInfo,
		"warn": logstuff.LevelWarn, "error": logstuff.LevelError, "garbage": logstuff.LevelInfo,
	}
	for env, want := range cases {
		t.Setenv(logstuff.EnvLogLevel, env)
		if got := logstuff.LevelFromEnv(); got != want {
			t.Errorf("env=%q: got %v, want %v", env, got, want)
		}
	}
}

func TestSlogCtx_PropagatesCtxKeys(t *testing.T) {
	var buf bytes.Buffer
	slogCtx := logstuff.NewSlogSink(logstuff.LevelInfo, logstuff.WithTextHandler(&buf)).(logstuff.CtxKeyer).WithCtxKeys("request_id", "and")
	ctx := context.WithValue(context.Background(), logstuff.CtxKey("request_id"), "abc-123")
	slogCtx.Log(ctx, logstuff.LevelInfo, "hi")
	out := buf.String()
	if !strings.Contains(out, "request_id=abc-123") {
		t.Errorf("ctx key not propagated: %q", out)
	}
	if strings.Contains(out, "and=") {
		t.Fatalf("tracked ctx keys that are not set in the ctx should not be logged")
	}
}

func TestSlogCtxWithCtxReturnsEqual(t *testing.T) {
	var b1, b2 bytes.Buffer
	logger := logstuff.NewSlogSink(logstuff.LevelDebug, logstuff.WithTextHandler(&b1))
	keyed := logger.(logstuff.CtxKeyer).WithCtxKeys()
	if logger != keyed {
		t.Fatalf("Empty keyed must return caller")
	}
	keyed = logstuff.NewSlogSink(logstuff.LevelDebug, logstuff.WithTextHandler(&b2)).(logstuff.CtxKeyer).WithCtxKeys()
	logger.Log(context.Background(), logstuff.LevelTrace, "yolo")
	keyed.Log(context.Background(), logstuff.LevelTrace, "yolo")
	if b1.String() != b2.String() {
		t.Fatalf("log must be identical")
	}
}

func TestSlogCtxDuplicateKeys(t *testing.T) {
	var buf bytes.Buffer
	ctx := context.WithValue(context.Background(), logstuff.CtxKey("key1"), "v1")
	ctx = context.WithValue(ctx, logstuff.CtxKey("key2"), "v2")
	logger := logstuff.NewSlogCtx(logstuff.LevelDebug, logstuff.WithTextHandler(&buf)).WithCtxKeys("key1")
	loggerKeyed := logger.WithCtxKeys("key1", "key2", "key2")
	loggerKeyed.Log(ctx, logstuff.LevelWarn, "msg")
	_, out, _ := strings.Cut(buf.String(), " ")
	if strings.TrimSpace(out) != `level=WARN msg=msg key1=v1 key2=v2` {
		t.Errorf("unexpected output: %q", out)
	}
	same := loggerKeyed.WithCtxKeys()
	if same != loggerKeyed {
		t.Fatal()
	}
}

func TestSlogSinkImplementsCtxKeyer(t *testing.T) {
	if _, ok := logstuff.NewSlogSink(logstuff.LevelInfo).(logstuff.CtxKeyer); !ok {
		t.Fatalf("SlogSink must implement CtxKeyer")
	}
}

func TestSlogSinkImplementsLevelEnabler(t *testing.T) {
	if _, ok := logstuff.NewSlogSink(logstuff.LevelInfo).(logstuff.LevelEnabler); !ok {
		t.Fatalf("SlogSink must implement LevelEnabler")
	}
}

func TestSlogSinkEnabled(t *testing.T) {
	sink := logstuff.NewSlogSink(logstuff.LevelInfo)
	if sink.(logstuff.LevelEnabler).Enabled(context.Background(), logstuff.LevelDebug) {
		t.Fatalf("should be disabled")
	}
}

func TestDiscard(t *testing.T) {
	d := logstuff.NewDiscard()
	if d.With("any", "thing") != d {
		t.Fatal()
	}
	d.Log(context.Background(), logstuff.LevelInfo, "msg")
	if d.(logstuff.LevelEnabler).Enabled(context.Background(), logstuff.LevelInfo) {
		t.Fatalf("discard must be always disabled")
	}
	d.Log(context.Background(), logstuff.LevelInfo, "msg")
}

func TestLogger(t *testing.T) {
	cases := make([]testCall, 0, 5)
	msg := make([]string, 5)
	theAnys := make([][]any, 5)
	for i, v := range "abcde" {
		msg[i] = string(v)
		theAnys[i] = []any{any(i), any(v)}
	}
	for i := range 5 {
		cases = append(cases, testCall{
			ctx:  context.Background(),
			lvl:  logstuff.LogLevel(-8 + 4*i),
			msg:  msg[i],
			args: theAnys[i],
		})
	}
	sink := testSink{}
	logger := logstuff.NewLogger(&sink)
	t.Run("Methods", func(t *testing.T) {
		for i, f := range []func(string, ...any){logger.Trace, logger.Debug, logger.Info, logger.Warn, logger.Error} {
			f(msg[i], theAnys[i]...)
		}
		for i, w := range cases {
			c := sink.calls[i]
			if c.ctx != w.ctx || c.lvl != w.lvl || c.msg != w.msg || c.args[1] != w.args[1] || c.args[0] != w.args[0] {
				t.Fatal()
			}
		}
	})
	t.Run("With", func(t *testing.T) {
		for i, v := range theAnys {
			logger.With(v...)
			if sink.with[i][0] != theAnys[i][0] || sink.with[i][1] != theAnys[i][1] {
				t.Fatal()
			}
		}
	})
	t.Run("With", func(t *testing.T) {
		sink := testSink{}
		logger := logstuff.NewLogger(&sink)
		w := cases[0]
		logger.Log(w.ctx, w.lvl, w.msg, w.args...)
		c := sink.calls[0]
		if c.ctx != w.ctx || c.lvl != w.lvl || c.msg != w.msg || c.args[1] != w.args[1] || c.args[0] != w.args[0] {
			t.Fatal()
		}
	})
	t.Run("Enabled", func(t *testing.T) {
		sink := testSink{initLvl: logstuff.LevelError}
		logger := logstuff.NewLogger(&sink)
		for _, v := range cases[:4] {
			if logger.Enabled(context.Background(), v.lvl) {
				t.Fatal()
			}
		}
		if !logger.Enabled(context.Background(), logstuff.LevelError) {
			t.Fatal()
		}
	})
	t.Run("notEnabled", func(t *testing.T) {
		sink := notEnabled{}
		logger := logstuff.NewLogger(&sink)
		if !logger.Enabled(context.Background(), logstuff.LevelTrace) {
			t.Fatalf("sink without enabled must always pass")
		}
	})
}

func TestSlogCtx_JSONHandler(t *testing.T) {
	var buf bytes.Buffer
	sink := logstuff.NewSlogSink(logstuff.LevelInfo, logstuff.WithJSONHandler(&buf))
	sink.Log(context.Background(), logstuff.LevelInfo, "hi", "key", "val")

	var got map[string]any
	if err := json.Unmarshal(buf.Bytes(), &got); err != nil {
		t.Fatalf("expected valid JSON, got %q: %v", buf.String(), err)
	}
	if got["key"] != "val" {
		t.Errorf("got %v", got)
	}
}

func TestNewSlogCtx_MultiHandler(t *testing.T) {
	var textBuf, jsonBuf bytes.Buffer
	sink := logstuff.NewSlogSink(logstuff.LevelInfo,
		logstuff.WithTextHandler(&textBuf),
		logstuff.WithJSONHandler(&jsonBuf),
	)
	sink.Log(context.Background(), logstuff.LevelInfo, "fanout test")

	if !strings.Contains(textBuf.String(), "fanout test") {
		t.Errorf("text handler missed the record: %q", textBuf.String())
	}
	if !strings.Contains(jsonBuf.String(), "fanout test") {
		t.Errorf("json handler missed the record: %q", jsonBuf.String())
	}
}
