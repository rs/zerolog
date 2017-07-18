package zerolog

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"sync"
	"time"

	pkgErrors "github.com/pkg/errors"
)

var eventPool = &sync.Pool{
	New: func() interface{} {
		return &Event{
			buf: make([]byte, 0, 500),
		}
	},
}

// Event represents a log event. It is instanced by one of the level method of
// Logger and finalized by the Msg or Msgf method.
type Event struct {
	buf        []byte
	errKey     string
	err        error
	w          LevelWriter
	level      Level
	enabled    bool
	stackTrace bool
	done       func(msg string)
}

func newEvent(w LevelWriter, level Level, enabled bool) *Event {
	if !enabled {
		return &Event{}
	}
	e := eventPool.Get().(*Event)
	e.buf = e.buf[:1]
	e.buf[0] = '{'
	e.err = nil
	e.errKey = ""
	e.w = w
	e.level = level
	e.enabled = true
	e.stackTrace = false
	return e
}

func (e *Event) write() (err error) {
	if !e.enabled {
		return nil
	}
	e.buf = append(e.buf, '}', '\n')
	_, err = e.w.WriteLevel(e.level, e.buf)
	eventPool.Put(e)
	return
}

// Enabled return false if the *Event is going to be filtered out by
// log level or sampling.
func (e *Event) Enabled() bool {
	return e.enabled
}

// Error sends the *Event with msg added as the message field if not empty.  If
// the current error is not wrapped, automatically wrap the error using
// github.com/pkg/errors.Wrap
//
// NOTICE: once this method is called, the *Event should be disposed.  Calling
// Error or Msg twice can have unexpected result.
func (e *Event) Error(msg string) (err error) {
	if !e.enabled {
		return
	}

	// If an error hasn't been created, create one now
	if e.err == nil {
		e.err = errors.New(msg)
	}

	e.Msg(msg)
	return err
}

// Errorf sends the *Event with msg added as the message field if not empty.  If
// the current error is not wrapped, automatically wrap the error using
// github.com/pkg/errors.Wrap
//
// NOTICE: once this method is called, the *Event should be disposed.  Calling
// Error or Msg twice can have unexpected result.
func (e *Event) Errorf(format string, v ...interface{}) (err error) {
	if !e.enabled {
		return
	}

	msg := fmt.Sprintf(format, v...)

	// If an error hasn't been created, create one now
	if e.err == nil {
		e.err = errors.New(msg)
	}

	e.Msg(msg)
	return err
}

// Msg sends the *Event with msg added as the message field if not empty.
//
// NOTICE: once this method is called, the *Event should be disposed.  Calling
// Error or Msg twice can have unexpected result.
func (e *Event) Msg(msg string) {
	if !e.enabled {
		return
	}
	if msg != "" {
		e.buf = appendString(e.buf, MessageFieldName, msg)
	}
	if e.err != nil {
		if e.errKey == "" {
			e.errKey = ErrorFieldName
		}

		type causer interface {
			Cause() error
		}
		if cause, ok := e.err.(causer); ok {
			err := cause.(error)
			e.buf = appendJSONString(appendKey(e.buf, e.errKey), err.Error())
		} else if !ok && e.stackTrace {
			e.err = pkgErrors.WithStack(e.err)
		} else {
			e.buf = appendJSONString(appendKey(e.buf, e.errKey), e.err.Error())
		}

		if e.stackTrace {
			type withStackTracer interface {
				StackTrace() pkgErrors.StackTrace
			}
			if stack, ok := e.err.(withStackTracer); ok {
				e.buf = appendJSONStack(appendKey(e.buf, StackFieldName), stack.StackTrace())
			}
		}
	}
	if e.done != nil {
		defer e.done(msg)
	}
	if err := e.write(); err != nil {
		fmt.Fprintf(os.Stderr, "zerolog: could not write event: %v", err)
	}
}

// Msgf sends the *Event with formated msg added as the message field if not empty.
//
// NOTICE: once this methid is called, the *Event should be disposed.
// Calling Msg twice can have unexpected result.
func (e *Event) Msgf(format string, v ...interface{}) {
	if !e.enabled {
		return
	}
	msg := fmt.Sprintf(format, v...)
	if msg != "" {
		e.buf = appendString(e.buf, MessageFieldName, msg)
	}
	if e.err != nil {
		if e.errKey == "" {
			e.errKey = ErrorFieldName
		}

		type causer interface {
			Cause() error
		}
		if cause, ok := e.err.(causer); ok {
			err := cause.(error)
			e.buf = appendJSONString(appendKey(e.buf, e.errKey), err.Error())
		} else if !ok && e.stackTrace {
			e.err = pkgErrors.WithStack(e.err)
		} else {
			e.buf = appendJSONString(appendKey(e.buf, e.errKey), e.err.Error())
		}

		if e.stackTrace {
			type withStackTracer interface {
				StackTrace() pkgErrors.StackTrace
			}
			if stack, ok := e.err.(withStackTracer); ok {
				e.buf = appendJSONStack(appendKey(e.buf, StackFieldName), stack.StackTrace())
			}
		}
	}
	if e.done != nil {
		defer e.done(msg)
	}
	if err := e.write(); err != nil {
		fmt.Fprintf(os.Stderr, "zerolog: could not write event: %v", err)
	}
}

// Fields is a helper function to use a map to set fields using type assertion.
func (e *Event) Fields(fields map[string]interface{}) *Event {
	if !e.enabled {
		return e
	}
	e.buf = appendFields(e.buf, fields)
	return e
}

// Dict adds the field key with a dict to the event context.
// Use zerolog.Dict() to create the dictionary.
func (e *Event) Dict(key string, dict *Event) *Event {
	if !e.enabled {
		return e
	}
	e.buf = append(append(appendKey(e.buf, key), dict.buf...), '}')
	eventPool.Put(dict)
	return e
}

// Dict creates an Event to be used with the *Event.Dict method.
// Call usual field methods like Str, Int etc to add fields to this
// event and give it as argument the *Event.Dict method.
func Dict() *Event {
	return newEvent(levelWriterAdapter{ioutil.Discard}, 0, true)
}

// StackTrace enables the dumping of a wrapped error's stacktrace if an error is
// not nil.
func (e *Event) StackTrace() *Event {
	if !e.enabled {
		return e
	}
	e.stackTrace = true
	return e
}

// Str adds the field key with val as a string to the *Event context.
func (e *Event) Str(key, val string) *Event {
	if !e.enabled {
		return e
	}
	e.buf = appendString(e.buf, key, val)
	return e
}

// Bytes adds the field key with val as a []byte to the *Event context.
func (e *Event) Bytes(key string, val []byte) *Event {
	if !e.enabled {
		return e
	}
	e.buf = appendBytes(e.buf, key, val)
	return e
}

// AnErr adds the field key with err as a string to the *Event context.  If err
// is nil, no field is added.  If an error has already been set, the existing
// error is not overwritten.
func (e *Event) AnErr(key string, err error) *Event {
	if !e.enabled {
		return e
	}
	if e.err == nil {
		e.errKey = key
		e.err = err
	}
	return e
}

// Err adds the field "error" with err as a string to the *Event context.
// If err is nil, no field is added.
// To customize the key name, change zerolog.ErrorFieldName.
func (e *Event) Err(err error) *Event {
	if !e.enabled {
		return e
	}
	if e.err == nil {
		e.err = err
		e.errKey = ""
	}
	return e
}

// Bool adds the field key with val as a Boolean to the *Event context.
func (e *Event) Bool(key string, b bool) *Event {
	if !e.enabled {
		return e
	}
	e.buf = appendBool(e.buf, key, b)
	return e
}

// Int adds the field key with i as a int to the *Event context.
func (e *Event) Int(key string, i int) *Event {
	if !e.enabled {
		return e
	}
	e.buf = appendInt(e.buf, key, i)
	return e
}

// Int8 adds the field key with i as a int8 to the *Event context.
func (e *Event) Int8(key string, i int8) *Event {
	if !e.enabled {
		return e
	}
	e.buf = appendInt8(e.buf, key, i)
	return e
}

// Int16 adds the field key with i as a int16 to the *Event context.
func (e *Event) Int16(key string, i int16) *Event {
	if !e.enabled {
		return e
	}
	e.buf = appendInt16(e.buf, key, i)
	return e
}

// Int32 adds the field key with i as a int32 to the *Event context.
func (e *Event) Int32(key string, i int32) *Event {
	if !e.enabled {
		return e
	}
	e.buf = appendInt32(e.buf, key, i)
	return e
}

// Int64 adds the field key with i as a int64 to the *Event context.
func (e *Event) Int64(key string, i int64) *Event {
	if !e.enabled {
		return e
	}
	e.buf = appendInt64(e.buf, key, i)
	return e
}

// Uint adds the field key with i as a uint to the *Event context.
func (e *Event) Uint(key string, i uint) *Event {
	if !e.enabled {
		return e
	}
	e.buf = appendUint(e.buf, key, i)
	return e
}

// Uint8 adds the field key with i as a uint8 to the *Event context.
func (e *Event) Uint8(key string, i uint8) *Event {
	if !e.enabled {
		return e
	}
	e.buf = appendUint8(e.buf, key, i)
	return e
}

// Uint16 adds the field key with i as a uint16 to the *Event context.
func (e *Event) Uint16(key string, i uint16) *Event {
	if !e.enabled {
		return e
	}
	e.buf = appendUint16(e.buf, key, i)
	return e
}

// Uint32 adds the field key with i as a uint32 to the *Event context.
func (e *Event) Uint32(key string, i uint32) *Event {
	if !e.enabled {
		return e
	}
	e.buf = appendUint32(e.buf, key, i)
	return e
}

// Uint64 adds the field key with i as a uint64 to the *Event context.
func (e *Event) Uint64(key string, i uint64) *Event {
	if !e.enabled {
		return e
	}
	e.buf = appendUint64(e.buf, key, i)
	return e
}

// Float32 adds the field key with f as a float32 to the *Event context.
func (e *Event) Float32(key string, f float32) *Event {
	if !e.enabled {
		return e
	}
	e.buf = appendFloat32(e.buf, key, f)
	return e
}

// Float64 adds the field key with f as a float64 to the *Event context.
func (e *Event) Float64(key string, f float64) *Event {
	if !e.enabled {
		return e
	}
	e.buf = appendFloat64(e.buf, key, f)
	return e
}

// Timestamp adds the current local time as UNIX timestamp to the *Event context with the "time" key.
// To customize the key name, change zerolog.TimestampFieldName.
func (e *Event) Timestamp() *Event {
	if !e.enabled {
		return e
	}
	e.buf = appendTimestamp(e.buf)
	return e
}

// Time adds the field key with t formated as string using zerolog.TimeFieldFormat.
func (e *Event) Time(key string, t time.Time) *Event {
	if !e.enabled {
		return e
	}
	e.buf = appendTime(e.buf, key, t)
	return e
}

// Dur adds the field key with duration d stored as zerolog.DurationFieldUnit.
// If zerolog.DurationFieldInteger is true, durations are rendered as integer
// instead of float.
func (e *Event) Dur(key string, d time.Duration) *Event {
	if !e.enabled {
		return e
	}
	e.buf = appendDuration(e.buf, key, d)
	return e
}

// TimeDiff adds the field key with positive duration between time t and start.
// If time t is not greater than start, duration will be 0.
// Duration format follows the same principle as Dur().
func (e *Event) TimeDiff(key string, t time.Time, start time.Time) *Event {
	if !e.enabled {
		return e
	}
	var d time.Duration
	if t.After(start) {
		d = t.Sub(start)
	}
	e.buf = appendDuration(e.buf, key, d)
	return e
}

// Interface adds the field key with i marshaled using reflection.
func (e *Event) Interface(key string, i interface{}) *Event {
	if !e.enabled {
		return e
	}
	e.buf = appendInterface(e.buf, key, i)
	return e
}
