package zerolog

import (
	"fmt"
	"net"
	"os"
	"runtime"
	"sync"
	"time"
)

var eventPool = &sync.Pool{
	New: func() interface{} {
		return &Event{
			Buf: make([]byte, 0, 500),
		}
	},
}

// Event represents a log event. It is instanced by one of the level method of
// Logger and finalized by the Msg or Msgf method.
type Event struct {
	Buf       []byte
	w         LevelWriter
	level     Level
	done      func(msg string)
	stack     bool   // enable error stack trace
	ch        []Hook // hooks from context
	skipFrame int    // The number of additional frames to skip when printing the caller.
}

func putEvent(e *Event) {
	// Proper usage of a sync.Pool requires each entry to have approximately
	// the same memory cost. To obtain this property when the stored type
	// contains a variably-sized buffer, we add a hard limit on the maximum buffer
	// to place back in the pool.
	//
	// See https://golang.org/issue/23199
	const maxSize = 1 << 16 // 64KiB
	if cap(e.Buf) > maxSize {
		return
	}
	eventPool.Put(e)
}

// LogObjectMarshaler provides a strongly-typed and encoding-agnostic interface
// to be implemented by types used with Event/Context's Object methods.
type LogObjectMarshaler interface {
	MarshalZerologObject(e *Event)
}

// LogArrayMarshaler provides a strongly-typed and encoding-agnostic interface
// to be implemented by types used with Event/Context's Array methods.
type LogArrayMarshaler interface {
	MarshalZerologArray(a *Array)
}

func newEvent(w LevelWriter, level Level) *Event {
	e := eventPool.Get().(*Event)
	e.Buf = e.Buf[:0]
	e.ch = nil
	e.Buf = enc.AppendBeginMarker(e.Buf)
	e.w = w
	e.level = level
	e.stack = false
	e.skipFrame = 0
	return e
}

func (e *Event) write() (err error) {
	if e == nil {
		return nil
	}
	if e.level != Disabled {
		e.Buf = enc.AppendEndMarker(e.Buf)
		e.Buf = enc.AppendLineBreak(e.Buf)
		if e.w != nil {
			_, err = e.w.WriteLevel(e.level, e.Buf)
		}
	}
	putEvent(e)
	return
}

// Enabled return false if the *Event is going to be filtered out by
// log level or sampling.
func (e *Event) Enabled() bool {
	return e != nil && e.level != Disabled
}

// Discard disables the event so Msg(f) won't print it.
func (e *Event) Discard() *Event {
	if e == nil {
		return e
	}
	e.level = Disabled
	return nil
}

// Msg sends the *Event with msg added as the message field if not empty.
//
// NOTICE: once this method is called, the *Event should be disposed.
// Calling Msg twice can have unexpected result.
func (e *Event) Msg(msg string) {
	if e == nil {
		return
	}
	e.msg(msg)
}

// Send is equivalent to calling Msg("").
//
// NOTICE: once this method is called, the *Event should be disposed.
func (e *Event) Send() {
	if e == nil {
		return
	}
	e.msg("")
}

// Msgf sends the event with formatted msg added as the message field if not empty.
//
// NOTICE: once this method is called, the *Event should be disposed.
// Calling Msgf twice can have unexpected result.
func (e *Event) Msgf(format string, v ...interface{}) {
	if e == nil {
		return
	}
	e.msg(fmt.Sprintf(format, v...))
}

func (e *Event) MsgFunc(createMsg func() string) {
	if e == nil {
		return
	}
	e.msg(createMsg())
}

func (e *Event) msg(msg string) {
	for _, hook := range e.ch {
		hook.Run(e, e.level, msg)
	}
	if msg != "" {
		e.Buf = enc.AppendString(enc.AppendKey(e.Buf, MessageFieldName), msg)
	}
	if e.done != nil {
		defer e.done(msg)
	}
	if err := e.write(); err != nil {
		if ErrorHandler != nil {
			ErrorHandler(err)
		} else {
			fmt.Fprintf(os.Stderr, "zerolog: could not write event: %v\n", err)
		}
	}
}

// Fields is a helper function to use a map or slice to set fields using type assertion.
// Only map[string]interface{} and []interface{} are accepted. []interface{} must
// alternate string keys and arbitrary values, and extraneous ones are ignored.
func (e *Event) Fields(fields interface{}) *Event {
	if e == nil {
		return e
	}
	e.Buf = appendFields(e.Buf, fields)
	return e
}

// Dict adds the field key with a dict to the event context.
// Use zerolog.Dict() to create the dictionary.
func (e *Event) Dict(key string, dict *Event) *Event {
	if e == nil {
		return e
	}
	dict.Buf = enc.AppendEndMarker(dict.Buf)
	e.Buf = append(enc.AppendKey(e.Buf, key), dict.Buf...)
	putEvent(dict)
	return e
}

// Dict creates an Event to be used with the *Event.Dict method.
// Call usual field methods like Str, Int etc to add fields to this
// event and give it as argument the *Event.Dict method.
func Dict() *Event {
	return newEvent(nil, 0)
}

// Array adds the field key with an array to the event context.
// Use zerolog.Arr() to create the array or pass a type that
// implement the LogArrayMarshaler interface.
func (e *Event) Array(key string, arr LogArrayMarshaler) *Event {
	if e == nil {
		return e
	}
	e.Buf = enc.AppendKey(e.Buf, key)
	var a *Array
	if aa, ok := arr.(*Array); ok {
		a = aa
	} else {
		a = Arr()
		arr.MarshalZerologArray(a)
	}
	e.Buf = a.write(e.Buf)
	return e
}

func (e *Event) appendObject(obj LogObjectMarshaler) {
	e.Buf = enc.AppendBeginMarker(e.Buf)
	obj.MarshalZerologObject(e)
	e.Buf = enc.AppendEndMarker(e.Buf)
}

// Object marshals an object that implement the LogObjectMarshaler interface.
func (e *Event) Object(key string, obj LogObjectMarshaler) *Event {
	if e == nil {
		return e
	}
	e.Buf = enc.AppendKey(e.Buf, key)
	if obj == nil {
		e.Buf = enc.AppendNil(e.Buf)

		return e
	}

	e.appendObject(obj)
	return e
}

// Func allows an anonymous func to run only if the event is enabled.
func (e *Event) Func(f func(e *Event)) *Event {
	if e != nil && e.Enabled() {
		f(e)
	}
	return e
}

// EmbedObject marshals an object that implement the LogObjectMarshaler interface.
func (e *Event) EmbedObject(obj LogObjectMarshaler) *Event {
	if e == nil {
		return e
	}
	if obj == nil {
		return e
	}
	obj.MarshalZerologObject(e)
	return e
}

// Str adds the field key with val as a string to the *Event context.
func (e *Event) Str(key, val string) *Event {
	if e == nil {
		return e
	}
	e.Buf = enc.AppendString(enc.AppendKey(e.Buf, key), val)
	return e
}

// Strs adds the field key with vals as a []string to the *Event context.
func (e *Event) Strs(key string, vals []string) *Event {
	if e == nil {
		return e
	}
	e.Buf = enc.AppendStrings(enc.AppendKey(e.Buf, key), vals)
	return e
}

// Stringer adds the field key with val.String() (or null if val is nil)
// to the *Event context.
func (e *Event) Stringer(key string, val fmt.Stringer) *Event {
	if e == nil {
		return e
	}
	e.Buf = enc.AppendStringer(enc.AppendKey(e.Buf, key), val)
	return e
}

// Stringers adds the field key with vals where each individual val
// is used as val.String() (or null if val is empty) to the *Event
// context.
func (e *Event) Stringers(key string, vals []fmt.Stringer) *Event {
	if e == nil {
		return e
	}
	e.Buf = enc.AppendStringers(enc.AppendKey(e.Buf, key), vals)
	return e
}

// Bytes adds the field key with val as a string to the *Event context.
//
// Runes outside of normal ASCII ranges will be hex-encoded in the resulting
// JSON.
func (e *Event) Bytes(key string, val []byte) *Event {
	if e == nil {
		return e
	}
	e.Buf = enc.AppendBytes(enc.AppendKey(e.Buf, key), val)
	return e
}

// Hex adds the field key with val as a hex string to the *Event context.
func (e *Event) Hex(key string, val []byte) *Event {
	if e == nil {
		return e
	}
	e.Buf = enc.AppendHex(enc.AppendKey(e.Buf, key), val)
	return e
}

// RawJSON adds already encoded JSON to the log line under key.
//
// No sanity check is performed on b; it must not contain carriage returns and
// be valid JSON.
func (e *Event) RawJSON(key string, b []byte) *Event {
	if e == nil {
		return e
	}
	e.Buf = appendJSON(enc.AppendKey(e.Buf, key), b)
	return e
}

// AnErr adds the field key with serialized err to the *Event context.
// If err is nil, no field is added.
func (e *Event) AnErr(key string, err error) *Event {
	if e == nil {
		return e
	}
	switch m := ErrorMarshalFunc(err).(type) {
	case nil:
		return e
	case LogObjectMarshaler:
		return e.Object(key, m)
	case error:
		if m == nil || isNilValue(m) {
			return e
		} else {
			return e.Str(key, m.Error())
		}
	case string:
		return e.Str(key, m)
	default:
		return e.Interface(key, m)
	}
}

// Errs adds the field key with errs as an array of serialized errors to the
// *Event context.
func (e *Event) Errs(key string, errs []error) *Event {
	if e == nil {
		return e
	}
	arr := Arr()
	for _, err := range errs {
		switch m := ErrorMarshalFunc(err).(type) {
		case LogObjectMarshaler:
			arr = arr.Object(m)
		case error:
			arr = arr.Err(m)
		case string:
			arr = arr.Str(m)
		default:
			arr = arr.Interface(m)
		}
	}

	return e.Array(key, arr)
}

// Err adds the field "error" with serialized err to the *Event context.
// If err is nil, no field is added.
//
// To customize the key name, change zerolog.ErrorFieldName.
//
// If Stack() has been called before and zerolog.ErrorStackMarshaler is defined,
// the err is passed to ErrorStackMarshaler and the result is appended to the
// zerolog.ErrorStackFieldName.
func (e *Event) Err(err error) *Event {
	if e == nil {
		return e
	}
	if e.stack && ErrorStackMarshaler != nil {
		switch m := ErrorStackMarshaler(err).(type) {
		case nil:
		case LogObjectMarshaler:
			e.Object(ErrorStackFieldName, m)
		case error:
			if m != nil && !isNilValue(m) {
				e.Str(ErrorStackFieldName, m.Error())
			}
		case string:
			e.Str(ErrorStackFieldName, m)
		default:
			e.Interface(ErrorStackFieldName, m)
		}
	}
	return e.AnErr(ErrorFieldName, err)
}

// Stack enables stack trace printing for the error passed to Err().
//
// ErrorStackMarshaler must be set for this method to do something.
func (e *Event) Stack() *Event {
	if e != nil {
		e.stack = true
	}
	return e
}

// Bool adds the field key with val as a bool to the *Event context.
func (e *Event) Bool(key string, b bool) *Event {
	if e == nil {
		return e
	}
	e.Buf = enc.AppendBool(enc.AppendKey(e.Buf, key), b)
	return e
}

// Bools adds the field key with val as a []bool to the *Event context.
func (e *Event) Bools(key string, b []bool) *Event {
	if e == nil {
		return e
	}
	e.Buf = enc.AppendBools(enc.AppendKey(e.Buf, key), b)
	return e
}

// Int adds the field key with i as a int to the *Event context.
func (e *Event) Int(key string, i int) *Event {
	if e == nil {
		return e
	}
	e.Buf = enc.AppendInt(enc.AppendKey(e.Buf, key), i)
	return e
}

// Ints adds the field key with i as a []int to the *Event context.
func (e *Event) Ints(key string, i []int) *Event {
	if e == nil {
		return e
	}
	e.Buf = enc.AppendInts(enc.AppendKey(e.Buf, key), i)
	return e
}

// Int8 adds the field key with i as a int8 to the *Event context.
func (e *Event) Int8(key string, i int8) *Event {
	if e == nil {
		return e
	}
	e.Buf = enc.AppendInt8(enc.AppendKey(e.Buf, key), i)
	return e
}

// Ints8 adds the field key with i as a []int8 to the *Event context.
func (e *Event) Ints8(key string, i []int8) *Event {
	if e == nil {
		return e
	}
	e.Buf = enc.AppendInts8(enc.AppendKey(e.Buf, key), i)
	return e
}

// Int16 adds the field key with i as a int16 to the *Event context.
func (e *Event) Int16(key string, i int16) *Event {
	if e == nil {
		return e
	}
	e.Buf = enc.AppendInt16(enc.AppendKey(e.Buf, key), i)
	return e
}

// Ints16 adds the field key with i as a []int16 to the *Event context.
func (e *Event) Ints16(key string, i []int16) *Event {
	if e == nil {
		return e
	}
	e.Buf = enc.AppendInts16(enc.AppendKey(e.Buf, key), i)
	return e
}

// Int32 adds the field key with i as a int32 to the *Event context.
func (e *Event) Int32(key string, i int32) *Event {
	if e == nil {
		return e
	}
	e.Buf = enc.AppendInt32(enc.AppendKey(e.Buf, key), i)
	return e
}

// Ints32 adds the field key with i as a []int32 to the *Event context.
func (e *Event) Ints32(key string, i []int32) *Event {
	if e == nil {
		return e
	}
	e.Buf = enc.AppendInts32(enc.AppendKey(e.Buf, key), i)
	return e
}

// Int64 adds the field key with i as a int64 to the *Event context.
func (e *Event) Int64(key string, i int64) *Event {
	if e == nil {
		return e
	}
	e.Buf = enc.AppendInt64(enc.AppendKey(e.Buf, key), i)
	return e
}

// Ints64 adds the field key with i as a []int64 to the *Event context.
func (e *Event) Ints64(key string, i []int64) *Event {
	if e == nil {
		return e
	}
	e.Buf = enc.AppendInts64(enc.AppendKey(e.Buf, key), i)
	return e
}

// Uint adds the field key with i as a uint to the *Event context.
func (e *Event) Uint(key string, i uint) *Event {
	if e == nil {
		return e
	}
	e.Buf = enc.AppendUint(enc.AppendKey(e.Buf, key), i)
	return e
}

// Uints adds the field key with i as a []int to the *Event context.
func (e *Event) Uints(key string, i []uint) *Event {
	if e == nil {
		return e
	}
	e.Buf = enc.AppendUints(enc.AppendKey(e.Buf, key), i)
	return e
}

// Uint8 adds the field key with i as a uint8 to the *Event context.
func (e *Event) Uint8(key string, i uint8) *Event {
	if e == nil {
		return e
	}
	e.Buf = enc.AppendUint8(enc.AppendKey(e.Buf, key), i)
	return e
}

// Uints8 adds the field key with i as a []int8 to the *Event context.
func (e *Event) Uints8(key string, i []uint8) *Event {
	if e == nil {
		return e
	}
	e.Buf = enc.AppendUints8(enc.AppendKey(e.Buf, key), i)
	return e
}

// Uint16 adds the field key with i as a uint16 to the *Event context.
func (e *Event) Uint16(key string, i uint16) *Event {
	if e == nil {
		return e
	}
	e.Buf = enc.AppendUint16(enc.AppendKey(e.Buf, key), i)
	return e
}

// Uints16 adds the field key with i as a []int16 to the *Event context.
func (e *Event) Uints16(key string, i []uint16) *Event {
	if e == nil {
		return e
	}
	e.Buf = enc.AppendUints16(enc.AppendKey(e.Buf, key), i)
	return e
}

// Uint32 adds the field key with i as a uint32 to the *Event context.
func (e *Event) Uint32(key string, i uint32) *Event {
	if e == nil {
		return e
	}
	e.Buf = enc.AppendUint32(enc.AppendKey(e.Buf, key), i)
	return e
}

// Uints32 adds the field key with i as a []int32 to the *Event context.
func (e *Event) Uints32(key string, i []uint32) *Event {
	if e == nil {
		return e
	}
	e.Buf = enc.AppendUints32(enc.AppendKey(e.Buf, key), i)
	return e
}

// Uint64 adds the field key with i as a uint64 to the *Event context.
func (e *Event) Uint64(key string, i uint64) *Event {
	if e == nil {
		return e
	}
	e.Buf = enc.AppendUint64(enc.AppendKey(e.Buf, key), i)
	return e
}

// Uints64 adds the field key with i as a []int64 to the *Event context.
func (e *Event) Uints64(key string, i []uint64) *Event {
	if e == nil {
		return e
	}
	e.Buf = enc.AppendUints64(enc.AppendKey(e.Buf, key), i)
	return e
}

// Float32 adds the field key with f as a float32 to the *Event context.
func (e *Event) Float32(key string, f float32) *Event {
	if e == nil {
		return e
	}
	e.Buf = enc.AppendFloat32(enc.AppendKey(e.Buf, key), f)
	return e
}

// Floats32 adds the field key with f as a []float32 to the *Event context.
func (e *Event) Floats32(key string, f []float32) *Event {
	if e == nil {
		return e
	}
	e.Buf = enc.AppendFloats32(enc.AppendKey(e.Buf, key), f)
	return e
}

// Float64 adds the field key with f as a float64 to the *Event context.
func (e *Event) Float64(key string, f float64) *Event {
	if e == nil {
		return e
	}
	e.Buf = enc.AppendFloat64(enc.AppendKey(e.Buf, key), f)
	return e
}

// Floats64 adds the field key with f as a []float64 to the *Event context.
func (e *Event) Floats64(key string, f []float64) *Event {
	if e == nil {
		return e
	}
	e.Buf = enc.AppendFloats64(enc.AppendKey(e.Buf, key), f)
	return e
}

// Timestamp adds the current local time as UNIX timestamp to the *Event context with the "time" key.
// To customize the key name, change zerolog.TimestampFieldName.
//
// NOTE: It won't dedupe the "time" key if the *Event (or *Context) has one
// already.
func (e *Event) Timestamp() *Event {
	if e == nil {
		return e
	}
	e.Buf = enc.AppendTime(enc.AppendKey(e.Buf, TimestampFieldName), TimestampFunc(), TimeFieldFormat)
	return e
}

// Time adds the field key with t formatted as string using zerolog.TimeFieldFormat.
func (e *Event) Time(key string, t time.Time) *Event {
	if e == nil {
		return e
	}
	e.Buf = enc.AppendTime(enc.AppendKey(e.Buf, key), t, TimeFieldFormat)
	return e
}

// Times adds the field key with t formatted as string using zerolog.TimeFieldFormat.
func (e *Event) Times(key string, t []time.Time) *Event {
	if e == nil {
		return e
	}
	e.Buf = enc.AppendTimes(enc.AppendKey(e.Buf, key), t, TimeFieldFormat)
	return e
}

// Dur adds the field key with duration d stored as zerolog.DurationFieldUnit.
// If zerolog.DurationFieldInteger is true, durations are rendered as integer
// instead of float.
func (e *Event) Dur(key string, d time.Duration) *Event {
	if e == nil {
		return e
	}
	e.Buf = enc.AppendDuration(enc.AppendKey(e.Buf, key), d, DurationFieldUnit, DurationFieldInteger)
	return e
}

// Durs adds the field key with duration d stored as zerolog.DurationFieldUnit.
// If zerolog.DurationFieldInteger is true, durations are rendered as integer
// instead of float.
func (e *Event) Durs(key string, d []time.Duration) *Event {
	if e == nil {
		return e
	}
	e.Buf = enc.AppendDurations(enc.AppendKey(e.Buf, key), d, DurationFieldUnit, DurationFieldInteger)
	return e
}

// TimeDiff adds the field key with positive duration between time t and start.
// If time t is not greater than start, duration will be 0.
// Duration format follows the same principle as Dur().
func (e *Event) TimeDiff(key string, t time.Time, start time.Time) *Event {
	if e == nil {
		return e
	}
	var d time.Duration
	if t.After(start) {
		d = t.Sub(start)
	}
	e.Buf = enc.AppendDuration(enc.AppendKey(e.Buf, key), d, DurationFieldUnit, DurationFieldInteger)
	return e
}

// Interface adds the field key with i marshaled using reflection.
func (e *Event) Interface(key string, i interface{}) *Event {
	if e == nil {
		return e
	}
	if obj, ok := i.(LogObjectMarshaler); ok {
		return e.Object(key, obj)
	}
	e.Buf = enc.AppendInterface(enc.AppendKey(e.Buf, key), i)
	return e
}

// Type adds the field key with val's type using reflection.
func (e *Event) Type(key string, val interface{}) *Event {
	if e == nil {
		return e
	}
	e.Buf = enc.AppendType(enc.AppendKey(e.Buf, key), val)
	return e
}

// CallerSkipFrame instructs any future Caller calls to skip the specified number of frames.
// This includes those added via hooks from the context.
func (e *Event) CallerSkipFrame(skip int) *Event {
	if e == nil {
		return e
	}
	e.skipFrame += skip
	return e
}

// Caller adds the file:line of the caller with the zerolog.CallerFieldName key.
// The argument skip is the number of stack frames to ascend
// Skip If not passed, use the global variable CallerSkipFrameCount
func (e *Event) Caller(skip ...int) *Event {
	sk := CallerSkipFrameCount
	if len(skip) > 0 {
		sk = skip[0] + CallerSkipFrameCount
	}
	return e.caller(sk)
}

func (e *Event) caller(skip int) *Event {
	if e == nil {
		return e
	}
	pc, file, line, ok := runtime.Caller(skip + e.skipFrame)
	if !ok {
		return e
	}
	e.Buf = enc.AppendString(enc.AppendKey(e.Buf, CallerFieldName), CallerMarshalFunc(pc, file, line))
	return e
}

// IPAddr adds IPv4 or IPv6 Address to the event
func (e *Event) IPAddr(key string, ip net.IP) *Event {
	if e == nil {
		return e
	}
	e.Buf = enc.AppendIPAddr(enc.AppendKey(e.Buf, key), ip)
	return e
}

// IPPrefix adds IPv4 or IPv6 Prefix (address and mask) to the event
func (e *Event) IPPrefix(key string, pfx net.IPNet) *Event {
	if e == nil {
		return e
	}
	e.Buf = enc.AppendIPPrefix(enc.AppendKey(e.Buf, key), pfx)
	return e
}

// MACAddr adds MAC address to the event
func (e *Event) MACAddr(key string, ha net.HardwareAddr) *Event {
	if e == nil {
		return e
	}
	e.Buf = enc.AppendMACAddr(enc.AppendKey(e.Buf, key), ha)
	return e
}
