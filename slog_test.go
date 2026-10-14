package zerolog

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"strings"
	"testing"
	"testing/slogtest"
	"time"
)

func newSlogLogger(buf *bytes.Buffer) *slog.Logger {
	zl := New(buf)
	return slog.New(NewSlogHandler(zl))
}

func decodeJSON(t *testing.T, buf *bytes.Buffer) map[string]interface{} {
	t.Helper()
	var m map[string]interface{}
	s := decodeIfBinaryToString(buf.Bytes())
	if err := json.Unmarshal([]byte(s), &m); err != nil {
		t.Fatalf("failed to decode JSON %q: %v", s, err)
	}
	return m
}

func TestSlogHandler_BasicInfo(t *testing.T) {
	var buf bytes.Buffer
	logger := newSlogLogger(&buf)

	logger.Info("hello world")

	m := decodeJSON(t, &buf)
	if m["level"] != "info" {
		t.Errorf("expected level info, got %v", m["level"])
	}
	if m["message"] != "hello world" {
		t.Errorf("expected message 'hello world', got %v", m["message"])
	}
}

func TestSlogHandler_Debug(t *testing.T) {
	var buf bytes.Buffer
	zl := New(&buf).Level(DebugLevel)
	logger := slog.New(NewSlogHandler(zl))

	logger.Debug("debug msg")

	m := decodeJSON(t, &buf)
	if m["level"] != "debug" {
		t.Errorf("expected level debug, got %v", m["level"])
	}
	if m["message"] != "debug msg" {
		t.Errorf("expected message 'debug msg', got %v", m["message"])
	}
}

func TestSlogHandler_Warn(t *testing.T) {
	var buf bytes.Buffer
	logger := newSlogLogger(&buf)

	logger.Warn("warn msg")

	m := decodeJSON(t, &buf)
	if m["level"] != "warn" {
		t.Errorf("expected level warn, got %v", m["level"])
	}
}

func TestSlogHandler_Error(t *testing.T) {
	var buf bytes.Buffer
	logger := newSlogLogger(&buf)

	logger.Error("error msg")

	m := decodeJSON(t, &buf)
	if m["level"] != "error" {
		t.Errorf("expected level error, got %v", m["level"])
	}
}

func TestSlogHandler_WithStringAttr(t *testing.T) {
	var buf bytes.Buffer
	logger := newSlogLogger(&buf)

	logger.Info("test", "key", "value")

	m := decodeJSON(t, &buf)
	if m["key"] != "value" {
		t.Errorf("expected key=value, got %v", m["key"])
	}
}

func TestSlogHandler_WithIntAttr(t *testing.T) {
	var buf bytes.Buffer
	logger := newSlogLogger(&buf)

	logger.Info("test", slog.Int("count", 42))

	m := decodeJSON(t, &buf)
	if m["count"] != float64(42) {
		t.Errorf("expected count=42, got %v", m["count"])
	}
}

func TestSlogHandler_WithBoolAttr(t *testing.T) {
	var buf bytes.Buffer
	logger := newSlogLogger(&buf)

	logger.Info("test", slog.Bool("flag", true))

	m := decodeJSON(t, &buf)
	if m["flag"] != true {
		t.Errorf("expected flag=true, got %v", m["flag"])
	}
}

func TestSlogHandler_WithFloat64Attr(t *testing.T) {
	var buf bytes.Buffer
	logger := newSlogLogger(&buf)

	logger.Info("test", slog.Float64("pi", 3.14))

	m := decodeJSON(t, &buf)
	if m["pi"] != 3.14 {
		t.Errorf("expected pi=3.14, got %v", m["pi"])
	}
}

func TestSlogHandler_WithTimeAttr(t *testing.T) {
	var buf bytes.Buffer
	logger := newSlogLogger(&buf)

	ts := time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)
	logger.Info("test", slog.Time("created", ts))

	m := decodeJSON(t, &buf)
	if m["created"] == nil {
		t.Error("expected created field to be present")
	}
}

func TestSlogHandler_WithDurationAttr(t *testing.T) {
	var buf bytes.Buffer
	logger := newSlogLogger(&buf)

	logger.Info("test", slog.Duration("elapsed", 5*time.Second))

	m := decodeJSON(t, &buf)
	if m["elapsed"] == nil {
		t.Error("expected elapsed field to be present")
	}
}

func TestSlogHandler_WithErrorAttr(t *testing.T) {
	var buf bytes.Buffer
	logger := newSlogLogger(&buf)

	logger.Info("test", slog.Any("err", errors.New("something failed")))

	m := decodeJSON(t, &buf)
	if m["err"] != "something failed" {
		t.Errorf("expected err='something failed', got %v", m["err"])
	}
}

func TestSlogHandler_WithAttrs(t *testing.T) {
	var buf bytes.Buffer
	zl := New(&buf)
	handler := NewSlogHandler(zl)

	child := handler.WithAttrs([]slog.Attr{
		slog.String("component", "auth"),
		slog.Int("version", 2),
	})
	logger := slog.New(child)

	logger.Info("request handled")

	m := decodeJSON(t, &buf)
	if m["component"] != "auth" {
		t.Errorf("expected component=auth, got %v", m["component"])
	}
	if m["version"] != float64(2) {
		t.Errorf("expected version=2, got %v", m["version"])
	}
	if m["message"] != "request handled" {
		t.Errorf("expected message 'request handled', got %v", m["message"])
	}
}

func TestSlogHandler_WithAttrsEmpty(t *testing.T) {
	var buf bytes.Buffer
	zl := New(&buf)
	handler := NewSlogHandler(zl)

	child := handler.WithAttrs(nil)
	if child != handler {
		t.Error("expected WithAttrs(nil) to return same handler")
	}
}

func TestSlogHandler_WithGroup(t *testing.T) {
	var buf bytes.Buffer
	zl := New(&buf)
	handler := NewSlogHandler(zl)

	child := handler.WithGroup("request")
	logger := slog.New(child)

	logger.Info("handled", "method", "GET", "status", 200)

	m := decodeJSON(t, &buf)
	req, ok := m["request"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected request group to be an object, got %T (%v)", m["request"], m["request"])
	}
	if req["method"] != "GET" {
		t.Errorf("expected request.method=GET, got %v", req["method"])
	}
	if req["status"] != float64(200) {
		t.Errorf("expected request.status=200, got %v", req["status"])
	}
}

func TestSlogHandler_WithGroupEmpty(t *testing.T) {
	var buf bytes.Buffer
	zl := New(&buf)
	handler := NewSlogHandler(zl)

	child := handler.WithGroup("")
	if child != handler {
		t.Error("expected WithGroup('') to return same handler")
	}
}

func TestSlogHandler_WithNestedGroups(t *testing.T) {
	var buf bytes.Buffer
	zl := New(&buf)
	handler := NewSlogHandler(zl)

	child := handler.WithGroup("http").WithGroup("request")
	logger := slog.New(child)

	logger.Info("handled", "method", "POST")

	m := decodeJSON(t, &buf)
	httpObj, ok := m["http"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected http group to be an object, got %T (%v)", m["http"], m["http"])
	}
	reqObj, ok := httpObj["request"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected http.request group to be an object, got %T (%v)", httpObj["request"], httpObj["request"])
	}
	if reqObj["method"] != "POST" {
		t.Errorf("expected http.request.method=POST, got %v", reqObj["method"])
	}
}

func TestSlogHandler_WithGroupAndAttrs(t *testing.T) {
	var buf bytes.Buffer
	zl := New(&buf)
	handler := NewSlogHandler(zl)

	child := handler.WithGroup("server").WithAttrs([]slog.Attr{
		slog.String("host", "localhost"),
	})
	logger := slog.New(child)

	logger.Info("started", "port", 8080)

	m := decodeJSON(t, &buf)
	srv, ok := m["server"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected server group to be an object, got %T (%v)", m["server"], m["server"])
	}
	if srv["host"] != "localhost" {
		t.Errorf("expected server.host=localhost, got %v", srv["host"])
	}
	if srv["port"] != float64(8080) {
		t.Errorf("expected server.port=8080, got %v", srv["port"])
	}
}

func TestSlogHandler_GroupAttrInRecord(t *testing.T) {
	var buf bytes.Buffer
	logger := newSlogLogger(&buf)

	logger.Info("test", slog.Group("user",
		slog.String("name", "alice"),
		slog.Int("age", 30),
	))

	m := decodeJSON(t, &buf)
	user, ok := m["user"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected user group to be an object, got %T (%v)", m["user"], m["user"])
	}
	if user["name"] != "alice" {
		t.Errorf("expected user.name=alice, got %v", user["name"])
	}
	if user["age"] != float64(30) {
		t.Errorf("expected user.age=30, got %v", user["age"])
	}
}

func TestSlogHandler_LevelFiltering(t *testing.T) {
	var buf bytes.Buffer
	zl := New(&buf).Level(WarnLevel)
	handler := NewSlogHandler(zl)

	if handler.Enabled(nil, slog.LevelDebug) {
		t.Error("expected debug to be filtered at warn level")
	}
	if handler.Enabled(nil, slog.LevelInfo) {
		t.Error("expected info to be filtered at warn level")
	}
	if !handler.Enabled(nil, slog.LevelWarn) {
		t.Error("expected warn to be enabled at warn level")
	}
	if !handler.Enabled(nil, slog.LevelError) {
		t.Error("expected error to be enabled at warn level")
	}
}

func TestSlogHandler_FilteredMessageNotWritten(t *testing.T) {
	var buf bytes.Buffer
	zl := New(&buf).Level(ErrorLevel)
	logger := slog.New(NewSlogHandler(zl))

	logger.Info("should not appear")

	if buf.Len() != 0 {
		t.Errorf("expected no output for filtered message, got %q", buf.String())
	}
}

func TestSlogHandler_MultipleAttrs(t *testing.T) {
	var buf bytes.Buffer
	logger := newSlogLogger(&buf)

	logger.Info("multi",
		slog.String("a", "1"),
		slog.Int("b", 2),
		slog.Bool("c", true),
		slog.Float64("d", 3.5),
	)

	m := decodeJSON(t, &buf)
	if m["a"] != "1" {
		t.Errorf("expected a=1, got %v", m["a"])
	}
	if m["b"] != float64(2) {
		t.Errorf("expected b=2, got %v", m["b"])
	}
	if m["c"] != true {
		t.Errorf("expected c=true, got %v", m["c"])
	}
	if m["d"] != 3.5 {
		t.Errorf("expected d=3.5, got %v", m["d"])
	}
}

func TestSlogHandler_LogValuer(t *testing.T) {
	var buf bytes.Buffer
	logger := newSlogLogger(&buf)

	logger.Info("test", "addr", testLogValuer{host: "example.com", port: 443})

	m := decodeJSON(t, &buf)
	addr, ok := m["addr"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected addr to resolve to a map object, got %T (%v)", m["addr"], m["addr"])
	}
	if addr["host"] != "example.com" {
		t.Errorf("expected addr.host=example.com, got %v", addr["host"])
	}
	if addr["port"] != float64(443) {
		t.Errorf("expected addr.port=443, got %v", addr["port"])
	}
}

type testLogValuer struct {
	host string
	port int
}

func (v testLogValuer) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("host", v.host),
		slog.Int("port", v.port),
	)
}

func TestSlogHandler_WithAttrsImmutability(t *testing.T) {
	var buf bytes.Buffer
	zl := New(&buf)
	h := NewSlogHandler(zl)

	child1 := h.WithAttrs([]slog.Attr{slog.String("from", "child1")})
	child2 := child1.WithAttrs([]slog.Attr{slog.String("from2", "child2")})

	buf.Reset()
	slog.New(child1).Info("test")
	m1 := decodeJSON(t, &buf)
	if m1["from"] != "child1" {
		t.Errorf("expected from=child1, got %v", m1["from"])
	}
	if m1["from2"] != nil {
		t.Errorf("expected child1 to not have from2, got %v", m1["from2"])
	}

	buf.Reset()
	slog.New(child2).Info("test")
	m2 := decodeJSON(t, &buf)
	if m2["from"] != "child1" {
		t.Errorf("expected child2 to have from=child1, got %v", m2["from"])
	}
	if m2["from2"] != "child2" {
		t.Errorf("expected child2 to have from2=child2, got %v", m2["from2"])
	}

	buf.Reset()
	slog.New(h).Info("test")
	m0 := decodeJSON(t, &buf)
	if m0["from"] != nil || m0["from2"] != nil {
		t.Errorf("expected root handler to have no child attrs, got %v", m0)
	}
}

func TestSlogHandler_LevelMapping(t *testing.T) {
	tests := []struct {
		slogLevel slog.Level
		wantLevel string
	}{
		{slog.LevelDebug - 4, "trace"},
		{slog.LevelDebug, "debug"},
		{slog.LevelInfo, "info"},
		{slog.LevelWarn, "warn"},
		{slog.LevelError, "error"},
	}

	for _, tt := range tests {
		var buf bytes.Buffer
		zl := New(&buf).Level(TraceLevel)
		logger := slog.New(NewSlogHandler(zl))

		logger.Log(nil, tt.slogLevel, "test")

		m := decodeJSON(t, &buf)
		if m["level"] != tt.wantLevel {
			t.Errorf("slog level %d: expected zerolog level %q, got %q",
				tt.slogLevel, tt.wantLevel, m["level"])
		}
		buf.Reset()
	}
}

func TestSlogHandler_EmptyMessage(t *testing.T) {
	var buf bytes.Buffer
	logger := newSlogLogger(&buf)

	logger.Info("", "key", "val")

	m := decodeJSON(t, &buf)
	if m["key"] != "val" {
		t.Errorf("expected key=val, got %v", m["key"])
	}
	if _, ok := m[MessageFieldName]; ok {
		t.Errorf("expected empty message to be omitted, got %v", m[MessageFieldName])
	}
}

func TestSlogHandler_WithContext(t *testing.T) {
	var buf bytes.Buffer
	zl := New(&buf).With().Str("service", "api").Logger()
	logger := slog.New(NewSlogHandler(zl))

	logger.Info("request")

	m := decodeJSON(t, &buf)
	if m["service"] != "api" {
		t.Errorf("expected service=api, got %v", m["service"])
	}
	if m["message"] != "request" {
		t.Errorf("expected message 'request', got %v", m["message"])
	}
}

func TestSlogHandler_EnabledRespectsGlobalLevel(t *testing.T) {
	var buf bytes.Buffer
	zl := New(&buf).Level(DebugLevel)
	handler := NewSlogHandler(zl)

	if !handler.Enabled(nil, slog.LevelInfo) {
		t.Fatal("expected info to be enabled before setting global level")
	}

	SetGlobalLevel(ErrorLevel)
	defer SetGlobalLevel(TraceLevel)

	if handler.Enabled(nil, slog.LevelInfo) {
		t.Error("expected info to be disabled when GlobalLevel is error")
	}
	if !handler.Enabled(nil, slog.LevelError) {
		t.Error("expected error to be enabled when GlobalLevel is error")
	}
}

func TestSlogHandler_EnabledNilWriter(t *testing.T) {
	var zl Logger
	handler := NewSlogHandler(zl)

	if handler.Enabled(nil, slog.LevelError) {
		t.Error("expected disabled when logger writer is nil")
	}
}

func TestSlogHandler_HandlePropagatesContext(t *testing.T) {
	var buf bytes.Buffer
	type ctxKey struct{}
	ctx := context.WithValue(context.Background(), ctxKey{}, "test-value")

	var gotCtx context.Context
	hook := HookFunc(func(e *Event, level Level, msg string) {
		gotCtx = e.GetCtx()
	})

	zl := New(&buf).Hook(hook)
	handler := NewSlogHandler(zl)

	record := slog.NewRecord(time.Now(), slog.LevelInfo, "test", 0)
	_ = handler.Handle(ctx, record)

	if gotCtx == nil {
		t.Fatal("expected context to be propagated to event")
	}
	if gotCtx.Value(ctxKey{}) != "test-value" {
		t.Error("expected context value to be preserved")
	}
}

func TestSlogHandler_NoDuplicateTimestamp(t *testing.T) {
	var buf bytes.Buffer
	zl := New(&buf).With().Timestamp().Logger()
	handler := NewSlogHandler(zl)

	record := slog.NewRecord(time.Now(), slog.LevelInfo, "test", 0)
	_ = handler.Handle(context.Background(), record)

	output := decodeIfBinaryToString(buf.Bytes())
	count := strings.Count(output, `"`+TimestampFieldName+`":`)
	if count != 1 {
		t.Errorf("expected exactly 1 timestamp field, got %d in output: %s", count, output)
	}
}

func TestSlogHandler_TimestampWithoutHook(t *testing.T) {
	var buf bytes.Buffer
	zl := New(&buf)
	handler := NewSlogHandler(zl)

	ts := time.Date(2024, 6, 15, 12, 0, 0, 0, time.UTC)
	record := slog.NewRecord(ts, slog.LevelInfo, "test", 0)
	_ = handler.Handle(context.Background(), record)

	m := decodeJSON(t, &buf)
	if m[TimestampFieldName] == nil {
		t.Error("expected timestamp field when logger has no timestamp hook")
	}
}

// Issue 787 Reproduction and Contract Compliance Tests

func TestSlogHandler_Issue787_NestedGroups(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(NewSlogHandler(New(&buf).With().Timestamp().Logger()))
	logger.WithGroup("g").Info("msg", "k", "v")

	m := decodeJSON(t, &buf)
	g, ok := m["g"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected nested group 'g' object, got %T (%v)", m["g"], m["g"])
	}
	if g["k"] != "v" {
		t.Errorf("expected g.k=v, got %v", g["k"])
	}
}

func TestSlogHandler_Issue787_SequentialWithGroups(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(NewSlogHandler(New(&buf).With().Timestamp().Logger()))
	logger.With("l", 1).WithGroup("g").With("l", 2).WithGroup("G").Info("msg", "l", 3)

	m := decodeJSON(t, &buf)
	if m["l"] != float64(1) {
		t.Errorf("expected top-level l=1, got %v", m["l"])
	}
	g, ok := m["g"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected group 'g' to be an object, got %T (%v)", m["g"], m["g"])
	}
	if g["l"] != float64(2) {
		t.Errorf("expected g.l=2, got %v", g["l"])
	}
	bigG, ok := g["G"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected group 'G' to be an object inside 'g', got %T (%v)", g["G"], g["G"])
	}
	if bigG["l"] != float64(3) {
		t.Errorf("expected g.G.l=3, got %v", bigG["l"])
	}
}

func TestSlogHandler_Issue787_RecordTimeRespected(t *testing.T) {
	var buf bytes.Buffer
	h := NewSlogHandler(New(&buf).With().Timestamp().Logger())
	targetTime := time.Date(2001, time.April, 1, 0, 0, 0, 0, time.UTC)
	record := slog.NewRecord(targetTime, slog.LevelInfo, "msg", 0)
	if err := h.Handle(context.Background(), record); err != nil {
		t.Fatalf("Handle failed: %v", err)
	}

	m := decodeJSON(t, &buf)
	timeVal, ok := m[TimestampFieldName].(string)
	if !ok {
		t.Fatalf("expected string timestamp, got %T (%v)", m[TimestampFieldName], m[TimestampFieldName])
	}
	parsed, err := time.Parse(time.RFC3339Nano, timeVal)
	if err != nil {
		parsed, err = time.Parse(time.RFC3339, timeVal)
	}
	if err != nil {
		t.Fatalf("failed to parse output timestamp %q: %v", timeVal, err)
	}
	if !parsed.Equal(targetTime) {
		t.Errorf("expected timestamp %v, got %v", targetTime, parsed)
	}
}

func TestSlogHandler_Issue787_ZeroRecordTimeDropped(t *testing.T) {
	var buf bytes.Buffer
	h := NewSlogHandler(New(&buf).With().Timestamp().Logger())
	record := slog.NewRecord(time.Time{}, slog.LevelInfo, "msg", 0)
	if err := h.Handle(context.Background(), record); err != nil {
		t.Fatalf("Handle failed: %v", err)
	}

	m := decodeJSON(t, &buf)
	if _, ok := m[TimestampFieldName]; ok {
		t.Errorf("expected zero timestamp to be omitted, got %v", m[TimestampFieldName])
	}
}

func TestSlogHandler_Issue787_CustomTimestampHookNoDuplicate(t *testing.T) {
	var buf bytes.Buffer
	hook := HookFunc(func(e *Event, level Level, message string) {
		e.Time(TimestampFieldName, time.Now())
	})

	logger := slog.New(NewSlogHandler(New(&buf).Hook(hook)))
	logger.Info("msg")

	output := decodeIfBinaryToString(buf.Bytes())
	count := strings.Count(output, `"`+TimestampFieldName+`":`)
	if count != 1 {
		t.Fatalf("expected exactly 1 timestamp field in output, got %d: %s", count, output)
	}
}

func TestSlogHandler_Issue787_CallerRespected(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(NewSlogHandler(New(&buf).With().Timestamp().Caller().Logger()))
	logger.Info("msg")

	m := decodeJSON(t, &buf)
	callerVal, ok := m[CallerFieldName].(string)
	if !ok {
		t.Fatalf("expected caller field, got %T (%v)", m[CallerFieldName], m[CallerFieldName])
	}
	if !strings.Contains(callerVal, "slog_test.go") {
		t.Errorf("expected caller to reference slog_test.go call site, got %q", callerVal)
	}
	if strings.Contains(callerVal, "slog.go") {
		t.Errorf("expected caller not to reference slog.go handler internals, got %q", callerVal)
	}
}

func TestSlogHandler_Issue787_ZeroPCDropped(t *testing.T) {
	var buf bytes.Buffer
	h := NewSlogHandler(New(&buf).With().Timestamp().Caller().Logger())
	record := slog.NewRecord(time.Now(), slog.LevelInfo, "msg", 0)
	if err := h.Handle(context.Background(), record); err != nil {
		t.Fatalf("Handle failed: %v", err)
	}

	m := decodeJSON(t, &buf)
	if _, ok := m[CallerFieldName]; ok {
		t.Errorf("expected zero PC caller to be omitted, got %v", m[CallerFieldName])
	}
}

func TestSlogHandler_Issue787_EmptyKeyLogged(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(NewSlogHandler(New(&buf).With().Timestamp().Logger()))
	logger.Info("msg", "", "v")

	m := decodeJSON(t, &buf)
	if m[""] != "v" {
		t.Errorf("expected empty key with value 'v', got %v", m[""])
	}
}

func TestSlogHandler_Issue787_ZeroAttrDropped(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(NewSlogHandler(New(&buf)))
	logger.Info("msg", slog.Attr{}, slog.String("a", "b"))

	m := decodeJSON(t, &buf)
	if m["a"] != "b" {
		t.Errorf("expected a=b, got %v", m["a"])
	}
	if _, ok := m[""]; ok {
		t.Errorf("unexpected empty attribute key in %v", m)
	}
	for k := range m {
		if k != "level" && k != "message" && k != "a" && k != TimestampFieldName {
			t.Errorf("unexpected attribute key %q", k)
		}
	}
}

func TestSlogHandler_Issue787_EmptyGroupOmitted(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(NewSlogHandler(New(&buf)))
	logger.With("a", "b").WithGroup("G").With("c", "d").WithGroup("H").Info("msg")

	m := decodeJSON(t, &buf)
	if m["a"] != "b" {
		t.Errorf("expected a=b, got %v", m["a"])
	}
	g, ok := m["G"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected G to be an object, got %T (%v)", m["G"], m["G"])
	}
	if g["c"] != "d" {
		t.Errorf("expected G.c=d, got %v", g["c"])
	}
	if _, ok := g["H"]; ok {
		t.Errorf("expected empty group H to be omitted from G, got %v", g["H"])
	}
}

func TestSlogHandler_Issue787_PreformattedWithAttrs(t *testing.T) {
	var buf bytes.Buffer
	h := NewSlogHandler(New(&buf))
	child := h.WithAttrs([]slog.Attr{
		slog.String("service", "billing"),
		slog.Int("env_id", 42),
	})

	slog.New(child).Info("event 1")
	m1 := decodeJSON(t, &buf)
	if m1["service"] != "billing" || m1["env_id"] != float64(42) {
		t.Errorf("expected preformatted attrs in event 1, got %v", m1)
	}

	buf.Reset()
	slog.New(child).Info("event 2")
	m2 := decodeJSON(t, &buf)
	if m2["service"] != "billing" || m2["env_id"] != float64(42) {
		t.Errorf("expected preformatted attrs in event 2, got %v", m2)
	}
}

func TestSlogHandler_SlogtestConformance_Issue787_10(t *testing.T) {
	var buf bytes.Buffer
	newHandler := func(t *testing.T) slog.Handler {
		buf.Reset()
		return NewSlogHandler(New(&buf))
	}
	result := func(t *testing.T) map[string]any {
		m := decodeJSON(t, &buf)
		if msg, ok := m[MessageFieldName]; ok {
			m[slog.MessageKey] = msg
		}
		if caller, ok := m[CallerFieldName]; ok {
			m[slog.SourceKey] = caller
		}
		return m
	}

	slogtest.Run(t, newHandler, result)
}
