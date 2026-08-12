package zerolog

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"testing"
	"testing/slogtest"
	"time"
)

func TestNewSlogHandler(t *testing.T) {
	t.Run("no hooks", func(t *testing.T) {
		h := NewSlogHandler(New(io.Discard))
		if h.hasTimestampHook {
			t.Error("hasTimestampHook, got: true, want: false")
		}
		if h.hasCallerHook {
			t.Error("hasCallerHook, got: true, want: false")
		}
		if got, want := len(h.logger.hooks), 0; got != want {
			t.Errorf("hooks len, got: %d, want: %d", got, want)
		}
	})

	t.Run("timestamp hook stripped", func(t *testing.T) {
		h := NewSlogHandler(New(io.Discard).With().Timestamp().Logger())
		if !h.hasTimestampHook {
			t.Error("hasTimestampHook, got: false, want: true")
		}
		if got, want := len(h.logger.hooks), 0; got != want {
			t.Errorf("hooks len, got: %d, want: %d", got, want)
		}
	})

	t.Run("caller hook stripped", func(t *testing.T) {
		h := NewSlogHandler(New(io.Discard).With().Caller().Logger())
		if !h.hasCallerHook {
			t.Error("hasCallerHook, got: false, want: true")
		}
		if got, want := len(h.logger.hooks), 0; got != want {
			t.Errorf("hooks len, got: %d, want: %d", got, want)
		}
	})

	t.Run("other hooks preserved", func(t *testing.T) {
		logger := New(io.Discard).With().Timestamp().Caller().Logger().Hook(HookFunc(func(e *Event, level Level, msg string) {}))
		h := NewSlogHandler(logger)
		if !h.hasTimestampHook {
			t.Error("hasTimestampHook, got: false, want: true")
		}
		if !h.hasCallerHook {
			t.Error("hasCallerHook, got: false, want: true")
		}
		if got, want := len(h.logger.hooks), 1; got != want {
			t.Errorf("hooks len, got: %d, want: %d", got, want)
		}
	})

	t.Run("original logger not mutated", func(t *testing.T) {
		logger := New(io.Discard).With().Timestamp().Caller().Logger()
		before := len(logger.hooks)
		NewSlogHandler(logger)
		if got, want := len(logger.hooks), before; got != want {
			t.Errorf("original logger hooks mutated, got: %d, want: %d", got, want)
		}
	})
}

func TestSlogHandler_Enabled(t *testing.T) {
	origGlobalLevel := GlobalLevel()
	SetGlobalLevel(DebugLevel)
	t.Cleanup(func() {
		SetGlobalLevel(origGlobalLevel)
	})

	tests := []struct {
		name        string
		loggerLevel Level
		level       slog.Level
		want        bool
	}{
		{"InfoLevel/slog.LevelInfo", InfoLevel, slog.LevelInfo, true},
		{"InfoLevel/slog.LevelDebug", InfoLevel, slog.LevelDebug, false},
		{"InfoLevel/slog.LevelWarn", InfoLevel, slog.LevelWarn, true},
		{"InfoLevel/slog.LevelInfo-1", InfoLevel, slog.LevelInfo - 1, false},
		{"InfoLevel/slog.LevelInfo+1", InfoLevel, slog.LevelInfo + 1, true},
		{"Disabled/slog.LevelDebug", Disabled, slog.LevelDebug, false},
		{"Disabled/slog.LevelInfo", Disabled, slog.LevelInfo, false},
		{"Disabled/slog.LevelWarn", Disabled, slog.LevelWarn, false},
		{"Disabled/slog.LevelError", Disabled, slog.LevelError, false},
		{"Global_DebugLevel/TraceLevel/slog.LevelDebug-2", TraceLevel, slog.LevelDebug - 2, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewSlogHandler(New(io.Discard).Level(tt.loggerLevel))
			if got, want := h.Enabled(context.Background(), tt.level), tt.want; got != want {
				t.Errorf("SlogHandler.Enabled, got: %v, want: %v", got, want)
			}
		})
	}
}

func TestSlogHandler_EnabledNilWriter(t *testing.T) {
	h := NewSlogHandler(Logger{})
	if got, want := h.Enabled(context.Background(), slog.LevelError), false; got != want {
		t.Errorf("SlogHandler.Enabled, got: %v, want: %v", got, want)
	}
}

func TestSlogHandler_LevelMapping(t *testing.T) {
	tests := []struct {
		name  string
		level slog.Level
		want  string
	}{
		{"below debug", slog.LevelDebug - 2, `{"level":"trace","message":"msg"}` + "\n"},
		{"debug", slog.LevelDebug, `{"level":"debug","message":"msg"}` + "\n"},
		{"between debug and info", slog.LevelInfo - 2, `{"level":"debug","message":"msg"}` + "\n"},
		{"info", slog.LevelInfo, `{"level":"info","message":"msg"}` + "\n"},
		{"between info and warn", slog.LevelWarn - 2, `{"level":"info","message":"msg"}` + "\n"},
		{"warn", slog.LevelWarn, `{"level":"warn","message":"msg"}` + "\n"},
		{"between warn and error", slog.LevelError - 2, `{"level":"warn","message":"msg"}` + "\n"},
		{"error", slog.LevelError, `{"level":"error","message":"msg"}` + "\n"},
		{"above error", slog.LevelError + 2, `{"level":"error","message":"msg"}` + "\n"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := &bytes.Buffer{}
			logger := slog.New(NewSlogHandler(New(out)))
			logger.Log(context.Background(), tt.level, "msg")
			if got, want := decodeIfBinaryToString(out.Bytes()), tt.want; got != want {
				t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
			}
		})
	}
}

func TestSlogHandler_Timestamp(t *testing.T) {
	t.Run("no timestamp hook omits timestamp field", func(t *testing.T) {
		out := &bytes.Buffer{}
		logger := slog.New(NewSlogHandler(New(out)))
		logger.Info("msg")
		if got, want := decodeIfBinaryToString(out.Bytes()), `{"level":"info","message":"msg"}`+"\n"; got != want {
			t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
		}
	})

	t.Run("zero time omits timestamp field", func(t *testing.T) {
		out := &bytes.Buffer{}
		h := NewSlogHandler(New(out).With().Timestamp().Logger())
		record := slog.NewRecord(time.Time{}, slog.LevelInfo, "msg", 0)
		if err := h.Handle(context.Background(), record); err != nil {
			t.Fatal(err)
		}
		if got, want := decodeIfBinaryToString(out.Bytes()), `{"level":"info","message":"msg"}`+"\n"; got != want {
			t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
		}
	})

	t.Run("timestamp field logged from record.Time", func(t *testing.T) {
		out := &bytes.Buffer{}
		h := NewSlogHandler(New(out).With().Timestamp().Logger())
		record := slog.NewRecord(time.Date(2001, time.February, 3, 4, 5, 6, 7, time.UTC), slog.LevelInfo, "msg", 0)
		if err := h.Handle(context.Background(), record); err != nil {
			t.Fatal(err)
		}
		if got, want := decodeIfBinaryToString(out.Bytes()), `{"level":"info","time":"2001-02-03T04:05:06Z","message":"msg"}`+"\n"; got != want {
			t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
		}
	})

	t.Run("custom TimestampFieldName", func(t *testing.T) {
		origTimestampFieldName := TimestampFieldName
		TimestampFieldName = "t"
		t.Cleanup(func() { TimestampFieldName = origTimestampFieldName })
		out := &bytes.Buffer{}
		h := NewSlogHandler(New(out).With().Timestamp().Logger())
		record := slog.NewRecord(time.Date(2001, time.February, 3, 4, 5, 6, 7, time.UTC), slog.LevelInfo, "msg", 0)
		if err := h.Handle(context.Background(), record); err != nil {
			t.Fatal(err)
		}
		if got, want := decodeIfBinaryToString(out.Bytes()), `{"level":"info","t":"2001-02-03T04:05:06Z","message":"msg"}`+"\n"; got != want {
			t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
		}
	})
}

func TestSlogHandler_Caller(t *testing.T) {
	t.Run("no caller hook omits caller field", func(t *testing.T) {
		out := &bytes.Buffer{}
		logger := slog.New(NewSlogHandler(New(out)))
		logger.Info("msg")
		if got, want := decodeIfBinaryToString(out.Bytes()), `{"level":"info","message":"msg"}`+"\n"; got != want {
			t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
		}
	})

	t.Run("zero PC omits caller field", func(t *testing.T) {
		out := &bytes.Buffer{}
		h := NewSlogHandler(New(out).With().Caller().Logger())
		record := slog.NewRecord(time.Time{}, slog.LevelInfo, "msg", 0)
		if err := h.Handle(context.Background(), record); err != nil {
			t.Fatal(err)
		}
		if got, want := decodeIfBinaryToString(out.Bytes()), `{"level":"info","message":"msg"}`+"\n"; got != want {
			t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
		}
	})

	t.Run("caller field resolves from record.PC", func(t *testing.T) {
		out := &bytes.Buffer{}
		logger := slog.New(NewSlogHandler(New(out).With().Caller().Logger()))
		pc, file, line, _ := runtime.Caller(0)
		caller := CallerMarshalFunc(pc, file, line+2)
		logger.Info("msg")
		if got, want := decodeIfBinaryToString(out.Bytes()), `{"level":"info","caller":"`+caller+`","message":"msg"}`+"\n"; got != want {
			t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
		}
	})

	t.Run("custom CallerMarshalFunc", func(t *testing.T) {
		origCallerMarshalFunc := CallerMarshalFunc
		CallerMarshalFunc = func(pc uintptr, file string, line int) string {
			parts := strings.Split(file, "/")
			if len(parts) > 1 {
				return strings.Join(parts[len(parts)-2:], "/") + ":" + strconv.Itoa(line)
			}

			return runtime.FuncForPC(pc).Name() + ":" + file + ":" + strconv.Itoa(line)
		}
		t.Cleanup(func() { CallerMarshalFunc = origCallerMarshalFunc })
		out := &bytes.Buffer{}
		logger := slog.New(NewSlogHandler(New(out).With().Caller().Logger()))
		pc, file, line, _ := runtime.Caller(0)
		caller := CallerMarshalFunc(pc, file, line+2)
		logger.Info("msg")
		if got, want := decodeIfBinaryToString(out.Bytes()), `{"level":"info","caller":"`+caller+`","message":"msg"}`+"\n"; got != want {
			t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
		}
	})

	t.Run("custom CallerFieldName", func(t *testing.T) {
		origCallerFieldName := CallerFieldName
		CallerFieldName = "source"
		t.Cleanup(func() { CallerFieldName = origCallerFieldName })
		out := &bytes.Buffer{}
		logger := slog.New(NewSlogHandler(New(out).With().Caller().Logger()))
		pc, file, line, _ := runtime.Caller(0)
		caller := CallerMarshalFunc(pc, file, line+2)
		logger.Info("msg")
		if got, want := decodeIfBinaryToString(out.Bytes()), `{"level":"info","source":"`+caller+`","message":"msg"}`+"\n"; got != want {
			t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
		}
	})
}

func TestSlogHandler_Message(t *testing.T) {
	t.Run("empty message omits message field", func(t *testing.T) {
		out := &bytes.Buffer{}
		logger := slog.New(NewSlogHandler(New(out)))
		logger.Info("")
		if got, want := decodeIfBinaryToString(out.Bytes()), `{"level":"info"}`+"\n"; got != want {
			t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
		}
	})

	t.Run("logs message field", func(t *testing.T) {
		out := &bytes.Buffer{}
		logger := slog.New(NewSlogHandler(New(out)))
		logger.Info("msg")
		if got, want := decodeIfBinaryToString(out.Bytes()), `{"level":"info","message":"msg"}`+"\n"; got != want {
			t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
		}
	})

	t.Run("custom MessageFieldName", func(t *testing.T) {
		origMessageFieldName := MessageFieldName
		MessageFieldName = "msg"
		t.Cleanup(func() { MessageFieldName = origMessageFieldName })
		out := &bytes.Buffer{}
		logger := slog.New(NewSlogHandler(New(out)))
		logger.Info("msg")
		if got, want := decodeIfBinaryToString(out.Bytes()), `{"level":"info","msg":"msg"}`+"\n"; got != want {
			t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
		}
	})
}

type testCtxKey struct{}

func TestSlogHandler_PropagatesContext(t *testing.T) {
	ctx := context.WithValue(context.Background(), testCtxKey{}, "v")

	hook := HookFunc(func(e *Event, level Level, message string) {
		ctx := e.GetCtx()
		if ctx == nil {
			return
		}
		if v, ok := ctx.Value(testCtxKey{}).(string); ok {
			e.Str("k", v)
		}
	})

	out := &bytes.Buffer{}
	logger := slog.New(NewSlogHandler(New(out).Hook(hook)))
	logger.InfoContext(ctx, "msg")
	if got, want := decodeIfBinaryToString(out.Bytes()),
		`{"level":"info","k":"v","message":"msg"}`+"\n"; got != want {
		t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
	}
}

type testLogValuer struct {
	v any
}

func (t testLogValuer) LogValue() slog.Value { return slog.AnyValue(t.v) }

type testAny struct {
	Name string `json:"name"`
}

func TestSlogHandler_AttrsTypes(t *testing.T) {
	tests := []struct {
		name string
		args []any
		want string
	}{
		{"zero attr", []any{slog.Attr{}}, `{"level":"info","message":"msg"}` + "\n"},
		{"empty key", []any{"", "v"}, `{"level":"info","":"v","message":"msg"}` + "\n"},
		{"empty group", []any{slog.Group("g")}, `{"level":"info","message":"msg"}` + "\n"},
		{"basic types", []any{
			"string", "hello",
			"int64", -1,
			"uint64", uint64(1),
			"float64", 3.14,
			"bool", true,
			"time", time.Date(2001, time.February, 3, 4, 5, 6, 7, time.UTC),
			"duration", time.Second,
		}, `{"level":"info","string":"hello","int64":-1,"uint64":1,"float64":3.14,"bool":true,"time":"2001-02-03T04:05:06Z","duration":1000,"message":"msg"}` + "\n"},
		{"group", []any{slog.Group("g", "k", "v")}, `{"level":"info","g":{"k":"v"},"message":"msg"}` + "\n"},
		{"empty group name", []any{slog.Group("", "k", "v")}, `{"level":"info","k":"v","message":"msg"}` + "\n"},
		{"nested group", []any{slog.Group("g1", slog.Group("g2", "k", "v"))}, `{"level":"info","g1":{"g2":{"k":"v"}},"message":"msg"}` + "\n"},
		{"slog.LogValuer", []any{"k", &testLogValuer{"v"}}, `{"level":"info","k":"v","message":"msg"}` + "\n"},
		{"slog.LogValuer in group", []any{slog.Group("g", "k", &testLogValuer{"v"})}, `{"level":"info","g":{"k":"v"},"message":"msg"}` + "\n"},
		{"any", []any{"user", testAny{Name: "john"}}, `{"level":"info","user":{"name":"john"},"message":"msg"}` + "\n"},
		{"other types", []any{
			"strs", []string{"one", "two"},
			"ints", []int{-1, 2},
			"ints8", []int8{-1, 2},
			"ints16", []int16{-1, 2},
			"ints32", []int32{-1, 2},
			"ints64", []int64{-1, 2},
			"uints", []uint{1, 2},
			"uints16", []uint16{1, 2},
			"uints32", []uint32{1, 2},
			"uints64", []uint64{1, 2},
			"floats32", []float32{3.14, 1.414},
			"floats64", []float64{3.14, 1.414},
			"bools", []bool{true, false},
			"times", []time.Time{time.Date(2001, time.February, 3, 4, 5, 6, 7, time.UTC), time.Date(2002, time.February, 3, 4, 5, 6, 7, time.UTC)},
			"durs", []time.Duration{time.Second, time.Minute},
			"error", errors.New("test error"),
			"errs", []error{errors.New("first error"), errors.New("second error")},
			"bytes", []byte("foo"),
			"ipAddr", net.IP{192, 168, 0, 10},
			"ipAddrs", []net.IP{{192, 168, 0, 10}, {127, 0, 0, 0}},
			"ipPrefix", net.IPNet{IP: net.IP{192, 168, 0, 10}, Mask: net.CIDRMask(24, 32)},
			"ipPrefixes", []net.IPNet{{IP: net.IP{192, 168, 0, 10}, Mask: net.CIDRMask(24, 32)}, {IP: net.IP{127, 0, 0, 0}, Mask: net.CIDRMask(24, 32)}},
			"macAddr", net.HardwareAddr{0x01, 0x23, 0x45, 0x67, 0x89, 0xab},
		}, `{"level":"info","strs":["one","two"],"ints":[-1,2],"ints8":[-1,2],"ints16":[-1,2],"ints32":[-1,2],"ints64":[-1,2],"uints":[1,2],"uints16":[1,2],"uints32":[1,2],"uints64":[1,2],"floats32":[3.14,1.414],"floats64":[3.14,1.414],"bools":[true,false],"times":["2001-02-03T04:05:06Z","2002-02-03T04:05:06Z"],"durs":[1000,60000],"error":"test error","errs":["first error","second error"],"bytes":"foo","ipAddr":"192.168.0.10","ipAddrs":["192.168.0.10","127.0.0.0"],"ipPrefix":"192.168.0.10/24","ipPrefixes":["192.168.0.10/24","127.0.0.0/24"],"macAddr":"01:23:45:67:89:ab","message":"msg"}` + "\n"},
	}

	for _, tt := range tests {
		t.Run("RecordAttrs/"+tt.name, func(t *testing.T) {
			out := &bytes.Buffer{}
			logger := slog.New(NewSlogHandler(New(out)))
			logger.Info("msg", tt.args...)
			if got, want := decodeIfBinaryToString(out.Bytes()), tt.want; got != want {
				t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
			}
		})

		t.Run("WithAttrs/"+tt.name, func(t *testing.T) {
			out := &bytes.Buffer{}
			logger := slog.New(NewSlogHandler(New(out))).With(tt.args...)
			logger.Info("msg")
			if got, want := decodeIfBinaryToString(out.Bytes()), tt.want; got != want {
				t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
			}
		})
	}
}

func TestSlogHandler_WithGroup(t *testing.T) {
	t.Run("empty name returns same handler", func(t *testing.T) {
		h := NewSlogHandler(New(io.Discard))
		if got := h.WithGroup(""); got != h {
			t.Error("empty group name: expected same handler")
		}
	})

	t.Run("original handler not mutated", func(t *testing.T) {
		out := &bytes.Buffer{}
		logger := slog.New(NewSlogHandler(New(out)))
		logger.WithGroup("g")
		logger.Info("msg", "k", "v")
		if got, want := decodeIfBinaryToString(out.Bytes()), `{"level":"info","k":"v","message":"msg"}`+"\n"; got != want {
			t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
		}
	})

	t.Run("no attrs produces no group key", func(t *testing.T) {
		out := &bytes.Buffer{}
		logger := slog.New(NewSlogHandler(New(out))).WithGroup("g")
		logger.Info("msg")
		if got, want := decodeIfBinaryToString(out.Bytes()), `{"level":"info","message":"msg"}`+"\n"; got != want {
			t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
		}
	})

	t.Run("no attrs produces no group key (multiple nested groups)", func(t *testing.T) {
		out := &bytes.Buffer{}
		logger := slog.New(NewSlogHandler(New(out))).WithGroup("g1").WithGroup("g2")
		logger.Info("msg")
		if got, want := decodeIfBinaryToString(out.Bytes()), `{"level":"info","message":"msg"}`+"\n"; got != want {
			t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
		}
	})

	t.Run("record attrs nested under group", func(t *testing.T) {
		out := &bytes.Buffer{}
		logger := slog.New(NewSlogHandler(New(out))).WithGroup("g")
		logger.Info("msg", "k", "v")
		if got, want := decodeIfBinaryToString(out.Bytes()), `{"level":"info","g":{"k":"v"},"message":"msg"}`+"\n"; got != want {
			t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
		}
	})

	t.Run("sibling groups", func(t *testing.T) {
		out := &bytes.Buffer{}
		base := slog.New(NewSlogHandler(New(out)))
		base.WithGroup("g1").Info("msg", "k", "v")
		base.WithGroup("g2").Info("msg", "k", "v")
		if got, want := decodeIfBinaryToString(out.Bytes()), `{"level":"info","g1":{"k":"v"},"message":"msg"}`+"\n"+`{"level":"info","g2":{"k":"v"},"message":"msg"}`+"\n"; got != want {
			t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
		}
	})

	t.Run("nested groups", func(t *testing.T) {
		out := &bytes.Buffer{}
		logger := slog.New(NewSlogHandler(New(out))).WithGroup("g1").WithGroup("g2")
		logger.Info("msg", "k", "v")
		if got, want := decodeIfBinaryToString(out.Bytes()), `{"level":"info","g1":{"g2":{"k":"v"}},"message":"msg"}`+"\n"; got != want {
			t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
		}
	})
}

func TestSlogHandler_WithAttrs(t *testing.T) {
	t.Run("empty attrs returns same handler", func(t *testing.T) {
		h := NewSlogHandler(New(io.Discard))
		if got := h.WithAttrs(nil); got != h {
			t.Error("nil attrs: expected same handler")
		}
		if got := h.WithAttrs([]slog.Attr{}); got != h {
			t.Error("empty attrs: expected same handler")
		}
	})

	t.Run("original handler not mutated", func(t *testing.T) {
		out := &bytes.Buffer{}
		logger := slog.New(NewSlogHandler(New(out)))
		logger.With("k", "v")
		logger.Info("msg")
		if got, want := decodeIfBinaryToString(out.Bytes()), `{"level":"info","message":"msg"}`+"\n"; got != want {
			t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
		}
	})

	t.Run("attrs applied to every record", func(t *testing.T) {
		out := &bytes.Buffer{}
		logger := slog.New(NewSlogHandler(New(out))).With("k", "v")
		logger.Info("first")
		logger.Info("second")
		if got, want := decodeIfBinaryToString(out.Bytes()), `{"level":"info","k":"v","message":"first"}`+"\n"+`{"level":"info","k":"v","message":"second"}`+"\n"; got != want {
			t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
		}
	})

	t.Run("attrs accumulate across multiple calls", func(t *testing.T) {
		out := &bytes.Buffer{}
		logger := slog.New(NewSlogHandler(New(out))).With("k1", "v1").With("k2", "v2")
		logger.Info("msg")
		if got, want := decodeIfBinaryToString(out.Bytes()), `{"level":"info","k1":"v1","k2":"v2","message":"msg"}`+"\n"; got != want {
			t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
		}
	})
}

func TestSlogHandler_WithGroupAndAttrs(t *testing.T) {
	t.Run("WithGroup/WithAttrs", func(t *testing.T) {
		out := &bytes.Buffer{}
		logger := slog.New(NewSlogHandler(New(out))).WithGroup("g").With("a", "b")
		logger.Info("msg", "k", "v")
		if got, want := decodeIfBinaryToString(out.Bytes()), `{"level":"info","g":{"a":"b","k":"v"},"message":"msg"}`+"\n"; got != want {
			t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
		}
	})

	t.Run("WithAttrs/WithGroup", func(t *testing.T) {
		out := &bytes.Buffer{}
		logger := slog.New(NewSlogHandler(New(out))).With("a", "b").WithGroup("g")
		logger.Info("msg", "k", "v")
		if got, want := decodeIfBinaryToString(out.Bytes()), `{"level":"info","a":"b","g":{"k":"v"},"message":"msg"}`+"\n"; got != want {
			t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
		}
	})

	t.Run("WithAttrs/WithGroup/WithAttrs", func(t *testing.T) {
		out := &bytes.Buffer{}
		logger := slog.New(NewSlogHandler(New(out))).With("a", "b").WithGroup("g").With("c", "d")
		logger.Info("msg", "k", "v")
		if got, want := decodeIfBinaryToString(out.Bytes()), `{"level":"info","a":"b","g":{"c":"d","k":"v"},"message":"msg"}`+"\n"; got != want {
			t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
		}
	})

	t.Run("WithGroup/WithAttrs/WithGroup", func(t *testing.T) {
		out := &bytes.Buffer{}
		logger := slog.New(NewSlogHandler(New(out))).WithGroup("g1").With("a", "b").WithGroup("g2")
		logger.Info("msg", "k", "v")
		if got, want := decodeIfBinaryToString(out.Bytes()), `{"level":"info","g1":{"a":"b","g2":{"k":"v"}},"message":"msg"}`+"\n"; got != want {
			t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
		}
	})
}

func TestSlogHandler_WithSampler(t *testing.T) {
	out := &bytes.Buffer{}
	sampler := &BasicSampler{N: 2}
	logger := slog.New(NewSlogHandler(New(out).Sample(sampler)))
	logger.Info("msg", "i", 1)
	logger.Info("msg", "i", 2)
	logger.Info("msg", "i", 3)
	logger.Info("msg", "i", 4)
	if got, want := decodeIfBinaryToString(out.Bytes()), `{"level":"info","i":1,"message":"msg"}`+"\n"+`{"level":"info","i":3,"message":"msg"}`+"\n"; got != want {
		t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
	}
}

func TestSlogHandler_WithZerologContext(t *testing.T) {
	tests := []struct {
		name     string
		group    string
		withArgs []any
		logArgs  []any
		want     string
	}{
		{"Empty/Empty/Empty", "", nil, nil, `{"level":"info","k":"v","message":"msg"}` + "\n"},
		{"Empty/Empty/Attrs", "", nil, []any{"k1", "v1"}, `{"level":"info","k":"v","k1":"v1","message":"msg"}` + "\n"},
		{"Empty/WithAttrs/Empty", "", []any{"k1", "v1"}, nil, `{"level":"info","k":"v","k1":"v1","message":"msg"}` + "\n"},
		{"Empty/WithAttrs/Attrs", "", []any{"k1", "v1"}, []any{"k2", "v2"}, `{"level":"info","k":"v","k1":"v1","k2":"v2","message":"msg"}` + "\n"},
		{"Group/Empty/Empty", "g", nil, nil, `{"level":"info","k":"v","message":"msg"}` + "\n"},
		{"Group/Empty/Attrs", "g", nil, []any{"k1", "v1"}, `{"level":"info","k":"v","g":{"k1":"v1"},"message":"msg"}` + "\n"},
		{"Group/WithAttrs/Empty", "g", []any{"k1", "v1"}, nil, `{"level":"info","k":"v","g":{"k1":"v1"},"message":"msg"}` + "\n"},
		{"Group/WithAttrs/Attrs", "g", []any{"k1", "v1"}, []any{"k2", "v2"}, `{"level":"info","k":"v","g":{"k1":"v1","k2":"v2"},"message":"msg"}` + "\n"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := &bytes.Buffer{}
			logger := slog.New(NewSlogHandler(New(out).With().Str("k", "v").Logger())).WithGroup(tt.group).With(tt.withArgs...)
			logger.Info("msg", tt.logArgs...)
			if got, want := decodeIfBinaryToString(out.Bytes()), tt.want; got != want {
				t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
			}
		})
	}
}

func TestSlogHandler_WithHook(t *testing.T) {
	called := false
	hook := HookFunc(func(e *Event, level Level, msg string) {
		called = true
	})
	logger := slog.New(NewSlogHandler(New(io.Discard).Hook(hook)))
	logger.Info("msg")
	if !called {
		t.Error("called, got: false, want: true")
	}
}

// Run a few different loggers with concurrent logs
// in an attempt to trip up 'go test -race' and discover any data races.
//
// Adopted from :- https://github.com/uber-go/zap/blob/018b91390e74732e9e40f8d356887b8d06461886/exp/zapslog/handler_test.go#L271
func TestSlogHandler_ConcurrentLogs(t *testing.T) {
	t.Parallel()

	const (
		NumWorkers = 10
		NumLogs    = 100
	)

	tests := []struct {
		name string
		log  func(l *slog.Logger, worker, log int)
	}{
		{
			name: "default",
			log: func(l *slog.Logger, worker, log int) {
				l.Info("default", "worker", worker, "log", log)
			},
		},
		{
			name: "WithGroup",
			log: func(l *slog.Logger, worker, log int) {
				l.WithGroup("g").Info("with group", "worker", worker, "log", log)
			},
		},
		{
			name: "WithAttrs",
			log: func(l *slog.Logger, worker, log int) {
				l.With("a", "b").Info("with attrs", "worker", worker, "log", log)
			},
		},
		{
			name: "WithGroupAndAttrs",
			log: func(l *slog.Logger, worker, log int) {
				l.WithGroup("g").With("a", "b").Info("with group and attrs", "worker", worker, "log", log)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			out := &bytes.Buffer{}
			logger := slog.New(NewSlogHandler(New(SyncWriter(out)).With().Timestamp().Caller().Str("x", "y").Logger()))

			// Use two wait groups to coordinate the workers:
			//
			// - ready: indicates when all workers should start logging.
			// - done: indicates when all workers have finished logging.
			var ready, done sync.WaitGroup
			ready.Add(NumWorkers)
			done.Add(NumWorkers)

			for i := range NumWorkers {
				go func() {
					defer done.Done()

					ready.Done() // I'm ready.
					ready.Wait() // Are others?

					for j := range NumLogs {
						tt.log(logger, i, j)
					}
				}()
			}

			done.Wait()

			if got, want := strings.Count(decodeIfBinaryToString(out.Bytes()), "\n"), NumWorkers*NumLogs; got != want {
				t.Errorf("number of logs, got: %d, want: %d", got, want)
			}
		})
	}
}

func TestSlogHandler_Slogtest(t *testing.T) {
	origMessageFieldName := MessageFieldName
	origCallerFieldName := CallerFieldName
	origLevelFieldName := LevelFieldName
	origTimestampFieldName := TimestampFieldName

	// Set field names as expected by slogtest
	MessageFieldName = slog.MessageKey
	CallerFieldName = slog.SourceKey
	LevelFieldName = slog.LevelKey
	TimestampFieldName = slog.TimeKey

	t.Cleanup(func() {
		MessageFieldName = origMessageFieldName
		CallerFieldName = origCallerFieldName
		LevelFieldName = origLevelFieldName
		TimestampFieldName = origTimestampFieldName
	})

	out := &bytes.Buffer{}
	handler := NewSlogHandler(New(out).With().Timestamp().Caller().Logger())

	if err := slogtest.TestHandler(
		handler,
		func() []map[string]any {
			var entries []map[string]any
			dec := json.NewDecoder(bytes.NewReader(decodeIfBinaryToBytes(out.Bytes())))
			for dec.More() {
				var decoded map[string]any
				if err := dec.Decode(&decoded); err != nil {
					t.Fatal(err)
				}
				entries = append(entries, decoded)
			}
			return entries
		},
	); err != nil {
		t.Fatal(err)
	}
}
