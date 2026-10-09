package zerolog

import (
	"bytes"
	"context"
	"log/slog"
	"runtime"
	"slices"
	"time"
)

type groupInfo struct {
	name  string
	attrs []slog.Attr
}

// SlogHandler implements the slog.Handler interface using a zerolog.Logger
// as the underlying log backend. This allows code that uses the standard
// library's slog package to route log output through zerolog.
type SlogHandler struct {
	logger    Logger
	hasCaller bool
	groups    []groupInfo
}

// NewSlogHandler creates a new slog.Handler that writes log records to the
// given zerolog.Logger. The handler maps slog levels to zerolog levels and
// converts slog attributes to zerolog fields.
func NewSlogHandler(logger Logger) *SlogHandler {
	hasCaller := false
	var cleanHooks []Hook
	for _, h := range logger.hooks {
		if _, ok := h.(timestampHook); ok {
			// Built-in timestamp hook is dropped; SlogHandler emits record.Time directly.
			continue
		}
		if _, ok := h.(callerHook); ok {
			// Built-in caller hook is dropped; SlogHandler emits record.PC directly.
			hasCaller = true
			continue
		}
		if sh, ok := h.(slogHook); ok {
			cleanHooks = append(cleanHooks, sh)
			continue
		}
		cleanHooks = append(cleanHooks, slogHook{hook: h})
	}
	logger.hooks = cleanHooks
	return &SlogHandler{
		logger:    logger,
		hasCaller: hasCaller,
	}
}

// slogHook wraps existing hooks to prevent duplicate timestamp or caller fields
// from being written by custom hooks when slog is managing them.
type slogHook struct {
	hook Hook
}

func (sh slogHook) Run(e *Event, level Level, msg string) {
	start := len(e.buf)
	sh.hook.Run(e, level, msg)
	if len(e.buf) > start {
		e.buf = stripFieldFromBuf(e.buf, start, TimestampFieldName)
		e.buf = stripFieldFromBuf(e.buf, start, CallerFieldName)
	}
}

func stripFieldFromBuf(buf []byte, start int, fieldName string) []byte {
	if fieldName == "" || len(buf) <= start {
		return buf
	}
	keyWithoutComma := enc.AppendKey([]byte{'{'}, fieldName)[1:]
	keyWithComma := append([]byte{','}, keyWithoutComma...)

	for {
		if len(buf) <= start {
			break
		}
		added := buf[start:]
		idx := bytes.Index(added, keyWithComma)
		hasLeadingComma := true
		keyLen := len(keyWithComma)
		if idx < 0 {
			idx = bytes.Index(added, keyWithoutComma)
			hasLeadingComma = false
			keyLen = len(keyWithoutComma)
		}
		if idx < 0 {
			break
		}

		fieldStart := start + idx
		valStart := fieldStart + keyLen
		valEnd := valStart

		if valStart < len(buf) {
			if buf[valStart] == '"' {
				valEnd++
				for valEnd < len(buf) {
					if buf[valEnd] == '\\' {
						valEnd += 2
						continue
					}
					if buf[valEnd] == '"' {
						valEnd++
						break
					}
					valEnd++
				}
			} else {
				for valEnd < len(buf) && buf[valEnd] != ',' && buf[valEnd] != '}' {
					valEnd++
				}
			}
		}

		if hasLeadingComma {
			buf = append(buf[:fieldStart], buf[valEnd:]...)
		} else {
			if valEnd < len(buf) && buf[valEnd] == ',' {
				valEnd++
			}
			buf = append(buf[:fieldStart], buf[valEnd:]...)
		}
	}
	return buf
}

// Enabled reports whether the handler handles records at the given level.
// It mirrors Logger.should's level and writer checks (without sampling).
func (h *SlogHandler) Enabled(_ context.Context, level slog.Level) bool {
	if h.logger.w == nil {
		return false
	}
	zl := slogToZerologLevel(level)
	if zl < GlobalLevel() {
		return false
	}
	return zl >= h.logger.level
}

// Handle handles the Record. It converts the slog.Record into a zerolog event
// and writes it using the underlying zerolog.Logger.
func (h *SlogHandler) Handle(ctx context.Context, record slog.Record) error {
	zlevel := slogToZerologLevel(record.Level)
	event := h.logger.WithLevel(zlevel)
	if event == nil {
		return nil
	}

	// Propagate slog context to the zerolog event so that hooks
	// relying on Event.GetCtx() (e.g. tracing) can access it.
	if ctx != nil {
		event = event.Ctx(ctx)
	}

	// Add timestamp from the slog record.
	// Contract: if record.Time is zero, ignore the time.
	if !record.Time.IsZero() {
		event.Time(TimestampFieldName, record.Time)
	}

	// Add caller from slog record PC if caller is enabled on the logger.
	// Contract: if record.PC is zero, ignore it.
	if h.hasCaller && record.PC != 0 {
		fs := runtime.CallersFrames([]uintptr{record.PC})
		f, _ := fs.Next()
		if f.File != "" {
			event.Str(CallerFieldName, CallerMarshalFunc(record.PC, f.File, f.Line))
		}
	}

	// Format groups and record attributes.
	if len(h.groups) == 0 {
		record.Attrs(func(a slog.Attr) bool {
			appendSlogAttr(event, a)
			return true
		})
	} else {
		h.appendGroups(event, record)
	}

	event.Msg(record.Message)
	return nil
}

func (h *SlogHandler) appendGroups(event *Event, record slog.Record) {
	appendGroupLevel(event, h.groups, 0, record)
}

func appendGroupLevel(parent *Event, groups []groupInfo, level int, record slog.Record) bool {
	dict := parent.CreateDict()
	appended := false

	for _, a := range groups[level].attrs {
		if appendSlogAttr(dict, a) {
			appended = true
		}
	}

	if level == len(groups)-1 {
		record.Attrs(func(a slog.Attr) bool {
			if appendSlogAttr(dict, a) {
				appended = true
			}
			return true
		})
	} else {
		if appendGroupLevel(dict, groups, level+1, record) {
			appended = true
		}
	}

	if !appended {
		putEvent(dict)
		return false
	}

	parent.Dict(groups[level].name, dict)
	return true
}

// WithAttrs returns a new Handler with the given attributes pre-attached.
// These attributes will be included in every subsequent log record.
func (h *SlogHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	if len(attrs) == 0 {
		return h
	}
	h2 := h.clone()
	if len(h2.groups) == 0 {
		ctx := h2.logger.With()
		for _, a := range attrs {
			ctx, _ = appendSlogAttrToContext(ctx, a)
		}
		h2.logger = ctx.Logger()
		return h2
	}
	h2.groups[len(h2.groups)-1].attrs = append(h2.groups[len(h2.groups)-1].attrs, attrs...)
	return h2
}

// WithGroup returns a new Handler with the given group name. All subsequent
// attributes will be nested under this group name in the output.
func (h *SlogHandler) WithGroup(name string) slog.Handler {
	if name == "" {
		return h
	}
	h2 := h.clone()
	h2.groups = append(h2.groups, groupInfo{name: name})
	return h2
}

func (h *SlogHandler) clone() *SlogHandler {
	h2 := &SlogHandler{
		logger:    h.logger,
		hasCaller: h.hasCaller,
	}
	if len(h.groups) > 0 {
		h2.groups = make([]groupInfo, len(h.groups))
		for i, g := range h.groups {
			h2.groups[i] = groupInfo{
				name:  g.name,
				attrs: slices.Clone(g.attrs),
			}
		}
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

// appendSlogAttr appends a single slog.Attr to the zerolog event.
func appendSlogAttr(event *Event, attr slog.Attr) bool {
	if event == nil {
		return false
	}

	attr.Value = attr.Value.Resolve()
	if attr.Equal(slog.Attr{}) {
		return false
	}

	if attr.Value.Kind() == slog.KindGroup {
		attrs := attr.Value.Group()
		if len(attrs) == 0 {
			return false
		}
		if attr.Key == "" {
			appended := false
			for _, ga := range attrs {
				if appendSlogAttr(event, ga) {
					appended = true
				}
			}
			return appended
		}
		sub := event.CreateDict()
		appended := false
		for _, ga := range attrs {
			if appendSlogAttr(sub, ga) {
				appended = true
			}
		}
		if !appended {
			putEvent(sub)
			return false
		}
		event.Dict(attr.Key, sub)
		return true
	}

	key := attr.Key
	val := attr.Value

	switch val.Kind() {
	case slog.KindString:
		event.Str(key, val.String())
	case slog.KindInt64:
		event.Int64(key, val.Int64())
	case slog.KindUint64:
		event.Uint64(key, val.Uint64())
	case slog.KindFloat64:
		event.Float64(key, val.Float64())
	case slog.KindBool:
		event.Bool(key, val.Bool())
	case slog.KindDuration:
		event.Dur(key, val.Duration())
	case slog.KindTime:
		event.Time(key, val.Time())
	case slog.KindAny:
		v := val.Any()
		switch cv := v.(type) {
		case error:
			event.AnErr(key, cv)
		case time.Duration:
			event.Dur(key, cv)
		case time.Time:
			event.Time(key, cv)
		case []byte:
			event.Bytes(key, cv)
		default:
			event.Interface(key, v)
		}
	default:
		event.Interface(key, val.Any())
	}

	return true
}

func appendSlogAttrToContext(ctx Context, attr slog.Attr) (Context, bool) {
	attr.Value = attr.Value.Resolve()
	if attr.Equal(slog.Attr{}) {
		return ctx, false
	}

	if attr.Value.Kind() == slog.KindGroup {
		attrs := attr.Value.Group()
		if len(attrs) == 0 {
			return ctx, false
		}
		if attr.Key == "" {
			appended := false
			for _, ga := range attrs {
				var ok bool
				ctx, ok = appendSlogAttrToContext(ctx, ga)
				if ok {
					appended = true
				}
			}
			return ctx, appended
		}
		sub := ctx.CreateDict()
		appended := false
		for _, ga := range attrs {
			if appendSlogAttr(sub, ga) {
				appended = true
			}
		}
		if !appended {
			putEvent(sub)
			return ctx, false
		}
		ctx = ctx.Dict(attr.Key, sub)
		return ctx, true
	}

	key := attr.Key
	val := attr.Value

	switch val.Kind() {
	case slog.KindString:
		return ctx.Str(key, val.String()), true
	case slog.KindInt64:
		return ctx.Int64(key, val.Int64()), true
	case slog.KindUint64:
		return ctx.Uint64(key, val.Uint64()), true
	case slog.KindFloat64:
		return ctx.Float64(key, val.Float64()), true
	case slog.KindBool:
		return ctx.Bool(key, val.Bool()), true
	case slog.KindDuration:
		return ctx.Dur(key, val.Duration()), true
	case slog.KindTime:
		return ctx.Time(key, val.Time()), true
	case slog.KindAny:
		v := val.Any()
		switch cv := v.(type) {
		case error:
			return ctx.AnErr(key, cv), true
		case time.Duration:
			return ctx.Dur(key, cv), true
		case time.Time:
			return ctx.Time(key, cv), true
		case []byte:
			return ctx.Bytes(key, cv), true
		default:
			return ctx.Interface(key, v), true
		}
	default:
		return ctx.Interface(key, val.Any()), true
	}
}

// Verify at compile time that SlogHandler satisfies the slog.Handler interface.
var _ slog.Handler = (*SlogHandler)(nil)
