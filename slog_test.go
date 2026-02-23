package zerolog_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"log/slog"
	"testing"
	"time"

	"github.com/rs/zerolog"
)

func newSlogLogger(buf *bytes.Buffer) *slog.Logger {
	zl := zerolog.New(buf)
	return slog.New(zerolog.NewSlogHandler(zl))
}

func decodeJSON(t *testing.T, buf *bytes.Buffer) map[string]interface{} {
	t.Helper()
	var m map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &m); err != nil {
		t.Fatalf("failed to decode JSON %q: %v", buf.String(), err)
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
	zl := zerolog.New(&buf).Level(zerolog.DebugLevel)
	logger := slog.New(zerolog.NewSlogHandler(zl))

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
	zl := zerolog.New(&buf)
	handler := zerolog.NewSlogHandler(zl)

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
	zl := zerolog.New(&buf)
	handler := zerolog.NewSlogHandler(zl)

	// WithAttrs with empty slice should return same handler
	child := handler.WithAttrs(nil)
	if child != handler {
		t.Error("expected WithAttrs(nil) to return same handler")
	}
}

func TestSlogHandler_WithGroup(t *testing.T) {
	var buf bytes.Buffer
	zl := zerolog.New(&buf)
	handler := zerolog.NewSlogHandler(zl)

	child := handler.WithGroup("request")
	logger := slog.New(child)

	logger.Info("handled", "method", "GET", "status", 200)

	m := decodeJSON(t, &buf)
	if m["request.method"] != "GET" {
		t.Errorf("expected request.method=GET, got %v", m["request.method"])
	}
	if m["request.status"] != float64(200) {
		t.Errorf("expected request.status=200, got %v", m["request.status"])
	}
}

func TestSlogHandler_WithGroupEmpty(t *testing.T) {
	var buf bytes.Buffer
	zl := zerolog.New(&buf)
	handler := zerolog.NewSlogHandler(zl)

	// WithGroup with empty name should return same handler
	child := handler.WithGroup("")
	if child != handler {
		t.Error("expected WithGroup('') to return same handler")
	}
}

func TestSlogHandler_WithNestedGroups(t *testing.T) {
	var buf bytes.Buffer
	zl := zerolog.New(&buf)
	handler := zerolog.NewSlogHandler(zl)

	child := handler.WithGroup("http").WithGroup("request")
	logger := slog.New(child)

	logger.Info("handled", "method", "POST")

	m := decodeJSON(t, &buf)
	if m["http.request.method"] != "POST" {
		t.Errorf("expected http.request.method=POST, got %v", m["http.request.method"])
	}
}

func TestSlogHandler_WithGroupAndAttrs(t *testing.T) {
	var buf bytes.Buffer
	zl := zerolog.New(&buf)
	handler := zerolog.NewSlogHandler(zl)

	child := handler.WithGroup("server").WithAttrs([]slog.Attr{
		slog.String("host", "localhost"),
	})
	logger := slog.New(child)

	logger.Info("started", "port", 8080)

	m := decodeJSON(t, &buf)
	if m["server.host"] != "localhost" {
		t.Errorf("expected server.host=localhost, got %v", m["server.host"])
	}
	if m["server.port"] != float64(8080) {
		t.Errorf("expected server.port=8080, got %v", m["server.port"])
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
	if m["user.name"] != "alice" {
		t.Errorf("expected user.name=alice, got %v", m["user.name"])
	}
	if m["user.age"] != float64(30) {
		t.Errorf("expected user.age=30, got %v", m["user.age"])
	}
}

func TestSlogHandler_LevelFiltering(t *testing.T) {
	var buf bytes.Buffer
	zl := zerolog.New(&buf).Level(zerolog.WarnLevel)
	handler := zerolog.NewSlogHandler(zl)

	// Debug should be filtered
	if handler.Enabled(nil, slog.LevelDebug) {
		t.Error("expected debug to be filtered at warn level")
	}
	// Info should be filtered
	if handler.Enabled(nil, slog.LevelInfo) {
		t.Error("expected info to be filtered at warn level")
	}
	// Warn should pass
	if !handler.Enabled(nil, slog.LevelWarn) {
		t.Error("expected warn to be enabled at warn level")
	}
	// Error should pass
	if !handler.Enabled(nil, slog.LevelError) {
		t.Error("expected error to be enabled at warn level")
	}
}

func TestSlogHandler_FilteredMessageNotWritten(t *testing.T) {
	var buf bytes.Buffer
	zl := zerolog.New(&buf).Level(zerolog.ErrorLevel)
	logger := slog.New(zerolog.NewSlogHandler(zl))

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
	// LogValuer resolves to a group
	if m["addr.host"] != "example.com" {
		t.Errorf("expected addr.host=example.com, got %v", m["addr.host"])
	}
	if m["addr.port"] != float64(443) {
		t.Errorf("expected addr.port=443, got %v", m["addr.port"])
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
	var buf1, buf2 bytes.Buffer
	zl1 := zerolog.New(&buf1)
	zl2 := zerolog.New(&buf2)

	handler := zerolog.NewSlogHandler(zl1)
	child1 := handler.WithAttrs([]slog.Attr{slog.String("from", "child1")})
	_ = zerolog.NewSlogHandler(zl2).WithAttrs([]slog.Attr{slog.String("from", "child2")})

	slog.New(child1).Info("test")

	m := decodeJSON(t, &buf1)
	if m["from"] != "child1" {
		t.Errorf("expected from=child1, got %v", m["from"])
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
		zl := zerolog.New(&buf).Level(zerolog.TraceLevel)
		logger := slog.New(zerolog.NewSlogHandler(zl))

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
}

func TestSlogHandler_WithContext(t *testing.T) {
	var buf bytes.Buffer
	zl := zerolog.New(&buf).With().Str("service", "api").Logger()
	logger := slog.New(zerolog.NewSlogHandler(zl))

	logger.Info("request")

	m := decodeJSON(t, &buf)
	if m["service"] != "api" {
		t.Errorf("expected service=api, got %v", m["service"])
	}
	if m["message"] != "request" {
		t.Errorf("expected message 'request', got %v", m["message"])
	}
}
