package zerolog

import (
	"context"
	"log/slog"
	"time"
)

// SlogHandler implements the slog.Handler interface using a zerolog.Logger
// as the underlying log backend. This allows code that uses the standard
// library's slog package to route log output through zerolog.
type SlogHandler struct {
	logger Logger
	prefix string // group prefix for nested groups
	attrs  []slog.Attr
}

// NewSlogHandler creates a new slog.Handler that writes log records to the
// given zerolog.Logger. The handler maps slog levels to zerolog levels and
// converts slog attributes to zerolog fields.
func NewSlogHandler(logger Logger) *SlogHandler {
	return &SlogHandler{logger: logger}
}

// Enabled reports whether the handler handles records at the given level.
func (h *SlogHandler) Enabled(_ context.Context, level slog.Level) bool {
	return h.logger.GetLevel() <= slogToZerologLevel(level)
}

// Handle handles the Record. It converts the slog.Record into a zerolog event
// and writes it using the underlying zerolog.Logger.
func (h *SlogHandler) Handle(_ context.Context, record slog.Record) error {
	zlevel := slogToZerologLevel(record.Level)
	event := h.logger.WithLevel(zlevel)
	if event == nil {
		return nil
	}

	// Add pre-attached attrs from WithAttrs
	for _, a := range h.attrs {
		event = appendSlogAttr(event, a, h.prefix)
	}

	// Add attrs from the record itself
	record.Attrs(func(a slog.Attr) bool {
		event = appendSlogAttr(event, a, h.prefix)
		return true
	})

	// Add timestamp
	if !record.Time.IsZero() {
		event.Time(TimestampFieldName, record.Time)
	}

	event.Msg(record.Message)
	return nil
}

// WithAttrs returns a new Handler with the given attributes pre-attached.
// These attributes will be included in every subsequent log record.
func (h *SlogHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	if len(attrs) == 0 {
		return h
	}
	h2 := h.clone()
	h2.attrs = append(h2.attrs, attrs...)
	return h2
}

// WithGroup returns a new Handler with the given group name. All subsequent
// attributes will be nested under this group name in the output.
func (h *SlogHandler) WithGroup(name string) slog.Handler {
	if name == "" {
		return h
	}
	h2 := h.clone()
	if h2.prefix != "" {
		h2.prefix = h2.prefix + "." + name
	} else {
		h2.prefix = name
	}
	return h2
}

func (h *SlogHandler) clone() *SlogHandler {
	h2 := &SlogHandler{
		logger: h.logger,
		prefix: h.prefix,
	}
	if len(h.attrs) > 0 {
		h2.attrs = make([]slog.Attr, len(h.attrs))
		copy(h2.attrs, h.attrs)
	}
	return h2
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

// zerologToSlogLevel maps zerolog levels to slog levels.
func zerologToSlogLevel(level Level) slog.Level {
	switch level {
	case TraceLevel:
		return slog.LevelDebug - 4
	case DebugLevel:
		return slog.LevelDebug
	case InfoLevel:
		return slog.LevelInfo
	case WarnLevel:
		return slog.LevelWarn
	case ErrorLevel:
		return slog.LevelError
	case FatalLevel:
		return slog.LevelError + 4
	case PanicLevel:
		return slog.LevelError + 8
	default:
		return slog.LevelInfo
	}
}

// appendSlogAttr appends a single slog.Attr to the zerolog event, handling
// type-specific encoding to avoid reflection where possible.
func appendSlogAttr(event *Event, attr slog.Attr, prefix string) *Event {
	if event == nil {
		return event
	}

	// Resolve the attribute to handle LogValuer types
	attr.Value = attr.Value.Resolve()

	key := attr.Key
	if key == "" && attr.Value.Kind() != slog.KindGroup {
		return event
	}
	if prefix != "" && key != "" {
		key = prefix + "." + key
	}

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
		if len(attrs) == 0 {
			return event
		}
		groupPrefix := prefix
		if key != "" {
			groupPrefix = key
		}
		for _, ga := range attrs {
			event = appendSlogAttr(event, ga, groupPrefix)
		}
	case slog.KindAny:
		v := val.Any()
		switch cv := v.(type) {
		case error:
			event = event.AnErr(key, cv)
		case time.Duration:
			event = event.Dur(key, cv)
		case time.Time:
			event = event.Time(key, cv)
		case []byte:
			event = event.Bytes(key, cv)
		default:
			event = event.Interface(key, v)
		}
	default:
		event = event.Interface(key, val.Any())
	}

	return event
}

// Verify at compile time that SlogHandler satisfies the slog.Handler interface.
var _ slog.Handler = (*SlogHandler)(nil)
