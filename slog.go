package zerolog

import (
	"context"
	"log/slog"
	"net"
	"runtime"
	"time"
)

// SlogHandler implements slog.Handler using a zerolog.Logger as the backend.
type SlogHandler struct {
	logger Logger

	unopenedGroups []string // groups from WithGroup that haven't been opened in the output buffer
	nOpenGroups    int      // number of groups currently open in the output buffer

	hasTimestampHook bool
	hasCallerHook    bool
}

// NewSlogHandler creates a slog.Handler using the given zerolog.Logger as the backend.
//
// If the logger was configured with Timestamp() or Caller(), those hooks are stripped and their values sourced from slog.Record instead.
// This is required for two reasons:
//
//   - A zero slog.Record.Time must produce no time field (slog.Handler contract).
//
//   - slog.Record.PC resolves the correct call site.
func NewSlogHandler(logger Logger) *SlogHandler {
	h := &SlogHandler{}

	hooks := make([]Hook, 0, len(logger.hooks))
	for _, hook := range logger.hooks {
		if _, ok := hook.(timestampHook); ok {
			h.hasTimestampHook = true
			continue
		}
		if _, ok := hook.(callerHook); ok {
			h.hasCallerHook = true
			continue
		}
		hooks = append(hooks, hook)
	}
	logger.hooks = hooks

	h.logger = logger
	return h
}

// Enabled reports whether the handler handles records at the given level.
func (h *SlogHandler) Enabled(_ context.Context, level slog.Level) bool {
	// It is equivalent to Logger.should, except that it does not invoke the configured Sampler.
	// slog.Logger calls Handler.Enabled before creating a slog.Record, while zerolog calls Logger.should (which performs sampling) when creating an Event.
	// Calling Logger.should here would cause the sampler to run twice for the same log entry, potentially dropping logs incorrectly.
	zl := slogToZerologLevel(level)

	if h.logger.disabled() {
		return false
	}
	if zl < h.logger.level || zl < GlobalLevel() {
		return false
	}
	return true
}

// Handle converts the slog.Record to a zerolog event and writes it.
func (h *SlogHandler) Handle(ctx context.Context, record slog.Record) error {
	e := h.logger.WithLevel(slogToZerologLevel(record.Level))
	if e == nil {
		return nil
	}

	if ctx != nil {
		e = e.Ctx(ctx)
	}

	nOpenGroups := h.nOpenGroups

	if record.NumAttrs() > 0 {
		e.buf = h.appendUnopenedGroups(e.buf)
		nOpenGroups += len(h.unopenedGroups)

		record.Attrs(func(a slog.Attr) bool {
			e = appendSlogAttrToEvent(e, a)
			return true
		})
	}

	// Close all opened groups
	for range nOpenGroups {
		e.buf = enc.AppendEndMarker(e.buf)
	}

	// Add timestamp using slog.Record.Time
	if h.hasTimestampHook && !record.Time.IsZero() {
		e = e.Time(TimestampFieldName, record.Time)
	}

	// Add caller using slog.Record.PC
	if h.hasCallerHook && record.PC != 0 {
		f, _ := runtime.CallersFrames([]uintptr{record.PC}).Next()
		if f.PC != 0 {
			e.buf = enc.AppendString(enc.AppendKey(e.buf, CallerFieldName), CallerMarshalFunc(f.PC, f.File, f.Line))
		}
	}

	e.Msg(record.Message)
	return nil
}

// WithAttrs returns a new Handler with the given attributes pre-attached.
// These attributes will be included in every subsequent log record.
func (h *SlogHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	if len(attrs) == 0 {
		return h
	}

	h2 := *h

	ctx := h.logger.With()

	ctx.l.context = h.appendUnopenedGroups(ctx.l.context)
	h2.nOpenGroups += len(h2.unopenedGroups)
	h2.unopenedGroups = nil

	for _, attr := range attrs {
		ctx = appendSlogAttrToContext(ctx, attr)
	}

	h2.logger = ctx.Logger()
	return &h2
}

func (h *SlogHandler) appendUnopenedGroups(buf []byte) []byte {
	for _, group := range h.unopenedGroups {
		buf = enc.AppendBeginMarker(enc.AppendKey(buf, group))
	}
	return buf
}

// WithGroup returns a new Handler with the given group name. All subsequent attributes are nested under name in the output.
func (h *SlogHandler) WithGroup(name string) slog.Handler {
	if name == "" {
		return h
	}

	h2 := *h
	h2.unopenedGroups = make([]string, len(h.unopenedGroups)+1)
	copy(h2.unopenedGroups, h.unopenedGroups)
	h2.unopenedGroups[len(h2.unopenedGroups)-1] = name
	return &h2
}

// slogToZerologLevel maps slog levels to zerolog levels.
//
// slog levels:  Debug=-4, Info=0, Warn=4, Error=8
// zerolog levels: Trace=-1, Debug=0, Info=1, Warn=2, Error=3, Fatal=4, Panic=5
func slogToZerologLevel(level slog.Level) Level {
	switch {
	case level < slog.LevelDebug:
		return TraceLevel
	case level < slog.LevelInfo:
		return DebugLevel
	case level < slog.LevelWarn:
		return InfoLevel
	case level < slog.LevelError:
		return WarnLevel
	default:
		return ErrorLevel
	}
}

// appendSlogAttrToEvent appends a single slog.Attr to the zerolog event, handling type-specific encoding to avoid reflection where possible.
func appendSlogAttrToEvent(event *Event, attr slog.Attr) *Event {
	if event == nil {
		return event
	}

	attr.Value = attr.Value.Resolve()

	if attr.Equal(slog.Attr{}) {
		return event
	}

	key := attr.Key
	val := attr.Value
	switch val.Kind() {
	case slog.KindString:
		event = event.Str(key, val.String())
	case slog.KindInt64:
		event = event.Int64(key, val.Int64())
	case slog.KindUint64:
		event = event.Uint64(key, val.Uint64())
	case slog.KindFloat64:
		event = event.Float64(key, val.Float64())
	case slog.KindBool:
		event = event.Bool(key, val.Bool())
	case slog.KindDuration:
		event = event.Dur(key, val.Duration())
	case slog.KindTime:
		event = event.Time(key, val.Time())
	case slog.KindGroup:
		attrs := val.Group()
		if len(attrs) > 0 {
			if key == "" {
				for _, ga := range attrs {
					event = appendSlogAttrToEvent(event, ga)
				}
			} else {
				dict := event.CreateDict()
				for _, ga := range attrs {
					dict = appendSlogAttrToEvent(dict, ga)
				}
				event = event.Dict(key, dict)
			}
		}
	case slog.KindAny:
		v := val.Any()
		switch cv := v.(type) {
		case []time.Duration:
			event = event.Durs(key, cv)
		case net.IP:
			event = event.IPAddr(key, cv)
		case []net.IP:
			event = event.IPAddrs(key, cv)
		case net.IPNet:
			event = event.IPPrefix(key, cv)
		case []net.IPNet:
			event = event.IPPrefixes(key, cv)
		case net.HardwareAddr:
			event = event.MACAddr(key, cv)
		case []time.Time:
			event = event.Times(key, cv)
		case error:
			event = event.AnErr(key, cv)
		case []bool:
			event = event.Bools(key, cv)
		case []byte:
			event = event.Bytes(key, cv)
		case []error:
			event = event.Errs(key, cv)
		case []float32:
			event = event.Floats32(key, cv)
		case []float64:
			event = event.Floats64(key, cv)
		case []int:
			event = event.Ints(key, cv)
		case []int8:
			event = event.Ints8(key, cv)
		case []int16:
			event = event.Ints16(key, cv)
		case []int32:
			event = event.Ints32(key, cv)
		case []int64:
			event = event.Ints64(key, cv)
		case []string:
			event = event.Strs(key, cv)
		case []uint:
			event = event.Uints(key, cv)
		case []uint16:
			event = event.Uints16(key, cv)
		case []uint32:
			event = event.Uints32(key, cv)
		case []uint64:
			event = event.Uints64(key, cv)
		default:
			event = event.Interface(key, v)
		}
	default:
		event = event.Interface(key, val.Any())
	}

	return event
}

func appendSlogAttrToContext(ctx Context, attr slog.Attr) Context {
	attr.Value = attr.Value.Resolve()

	if attr.Equal(slog.Attr{}) {
		return ctx
	}

	key := attr.Key
	val := attr.Value
	switch val.Kind() {
	case slog.KindString:
		ctx = ctx.Str(key, val.String())
	case slog.KindInt64:
		ctx = ctx.Int64(key, val.Int64())
	case slog.KindUint64:
		ctx = ctx.Uint64(key, val.Uint64())
	case slog.KindFloat64:
		ctx = ctx.Float64(key, val.Float64())
	case slog.KindBool:
		ctx = ctx.Bool(key, val.Bool())
	case slog.KindDuration:
		ctx = ctx.Dur(key, val.Duration())
	case slog.KindTime:
		ctx = ctx.Time(key, val.Time())
	case slog.KindGroup:
		attrs := val.Group()
		if len(attrs) > 0 {
			if key == "" {
				for _, ga := range attrs {
					ctx = appendSlogAttrToContext(ctx, ga)
				}
			} else {
				dict := ctx.CreateDict()
				for _, ga := range attrs {
					dict = appendSlogAttrToEvent(dict, ga)
				}
				ctx = ctx.Dict(key, dict)
			}
		}
	case slog.KindAny:
		v := val.Any()
		switch cv := v.(type) {
		case []time.Duration:
			ctx = ctx.Durs(key, cv)
		case net.IP:
			ctx = ctx.IPAddr(key, cv)
		case []net.IP:
			ctx = ctx.IPAddrs(key, cv)
		case net.IPNet:
			ctx = ctx.IPPrefix(key, cv)
		case []net.IPNet:
			ctx = ctx.IPPrefixes(key, cv)
		case net.HardwareAddr:
			ctx = ctx.MACAddr(key, cv)
		case []time.Time:
			ctx = ctx.Times(key, cv)
		case error:
			ctx = ctx.AnErr(key, cv)
		case []bool:
			ctx = ctx.Bools(key, cv)
		case []byte:
			ctx = ctx.Bytes(key, cv)
		case []error:
			ctx = ctx.Errs(key, cv)
		case []float32:
			ctx = ctx.Floats32(key, cv)
		case []float64:
			ctx = ctx.Floats64(key, cv)
		case []int:
			ctx = ctx.Ints(key, cv)
		case []int8:
			ctx = ctx.Ints8(key, cv)
		case []int16:
			ctx = ctx.Ints16(key, cv)
		case []int32:
			ctx = ctx.Ints32(key, cv)
		case []int64:
			ctx = ctx.Ints64(key, cv)
		case []string:
			ctx = ctx.Strs(key, cv)
		case []uint:
			ctx = ctx.Uints(key, cv)
		case []uint16:
			ctx = ctx.Uints16(key, cv)
		case []uint32:
			ctx = ctx.Uints32(key, cv)
		case []uint64:
			ctx = ctx.Uints64(key, cv)
		default:
			ctx = ctx.Interface(key, v)
		}
	default:
		ctx = ctx.Interface(key, val.Any())
	}

	return ctx
}

// Verify at compile time that SlogHandler satisfies the slog.Handler interface.
var _ slog.Handler = (*SlogHandler)(nil)
