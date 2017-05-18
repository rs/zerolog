package zerolog

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"sync"
	"time"
)

var pool = &sync.Pool{
	New: func() interface{} {
		return bytes.NewBuffer(make([]byte, 500))
	},
}

// Event represents a log event. It is instancied by one of the level method of
// Logger and finalized by the Msg or Msgf method.
type Event struct {
	buf     *bytes.Buffer
	w       LevelWriter
	level   Level
	enabled bool
	done    func(msg string)
}

func newEvent(w LevelWriter, level Level, enabled bool) Event {
	if !enabled {
		return Event{}
	}
	buf := pool.Get().(*bytes.Buffer)
	buf.Reset()
	buf.WriteByte('{')
	return Event{
		buf:     buf,
		w:       w,
		level:   level,
		enabled: true,
	}
}

func (e Event) write() (n int, err error) {
	if !e.enabled {
		return 0, nil
	}
	e.buf.WriteByte('}')
	e.buf.WriteByte('\n')
	n, err = e.w.WriteLevel(e.level, e.buf.Bytes())
	pool.Put(e.buf)
	return
}

func (e Event) append(f field) Event {
	if !e.enabled {
		return e
	}
	if e.buf.Len() > 1 {
		e.buf.WriteByte(',')
	}
	f.writeJSON(e.buf)
	return e
}

// Enabled return false if the event is going to be filtered out by
// log level or sampling.
func (e Event) Enabled() bool {
	return e.enabled
}

// Msg sends the event with msg added as the message field if not empty.
//
// NOTICE: once this methid is called, the Event should be disposed.
// Calling Msg twice can have unexpected result.
func (e Event) Msg(msg string) (n int, err error) {
	if !e.enabled {
		return 0, nil
	}
	if msg != "" {
		e.append(fStr(MessageFieldName, msg))
	}
	if e.done != nil {
		defer e.done(msg)
	}
	return e.write()
}

// Msgf sends the event with formated msg added as the message field if not empty.
//
// NOTICE: once this methid is called, the Event should be disposed.
// Calling Msg twice can have unexpected result.
func (e Event) Msgf(format string, v ...interface{}) (n int, err error) {
	if !e.enabled {
		return 0, nil
	}
	msg := fmt.Sprintf(format, v...)
	if msg != "" {
		e.append(fStr(MessageFieldName, msg))
	}
	if e.done != nil {
		defer e.done(msg)
	}
	return e.write()
}

// Dict adds the field key with a dict to the event context.
// Use zerolog.Dict() to create the dictionary.
func (e Event) Dict(key string, dict Event) Event {
	if !e.enabled {
		return e
	}
	if e.buf.Len() > 1 {
		e.buf.WriteByte(',')
	}
	io.Copy(e.buf, dict.buf)
	e.buf.WriteByte('}')
	return e
}

// Dict creates an Event to be used with the event.Dict method.
// Call usual field methods like Str, Int etc to add fields to this
// event and give it as argument the event.Dict method.
func Dict() Event {
	return newEvent(levelWriterAdapter{ioutil.Discard}, 0, true)
}

// Str adds the field key with val as a string to the event context.
func (e Event) Str(key, val string) Event {
	if !e.enabled {
		return e
	}
	return e.append(fStr(key, val))
}

// Err adds the field "error" with err as a string to the event context.
// To customize the key name, change zerolog.ErrorFieldName.
func (e Event) Err(err error) Event {
	if !e.enabled {
		return e
	}
	return e.append(fErr(err))
}

// Bool adds the field key with val as a Boolean to the event context.
func (e Event) Bool(key string, b bool) Event {
	if !e.enabled {
		return e
	}
	return e.append(fBool(key, b))
}

// Int adds the field key with i as a int to the event context.
func (e Event) Int(key string, i int) Event {
	if !e.enabled {
		return e
	}
	return e.append(fInt(key, i))
}

// Int8 adds the field key with i as a int8 to the event context.
func (e Event) Int8(key string, i int8) Event {
	if !e.enabled {
		return e
	}
	return e.append(fInt8(key, i))
}

// Int16 adds the field key with i as a int16 to the event context.
func (e Event) Int16(key string, i int16) Event {
	if !e.enabled {
		return e
	}
	return e.append(fInt16(key, i))
}

// Int32 adds the field key with i as a int32 to the event context.
func (e Event) Int32(key string, i int32) Event {
	if !e.enabled {
		return e
	}
	return e.append(fInt32(key, i))
}

// Int64 adds the field key with i as a int64 to the event context.
func (e Event) Int64(key string, i int64) Event {
	if !e.enabled {
		return e
	}
	return e.append(fInt64(key, i))
}

// Uint adds the field key with i as a uint to the event context.
func (e Event) Uint(key string, i uint) Event {
	if !e.enabled {
		return e
	}
	return e.append(fUint(key, i))
}

// Uint8 adds the field key with i as a uint8 to the event context.
func (e Event) Uint8(key string, i uint8) Event {
	if !e.enabled {
		return e
	}
	return e.append(fUint8(key, i))
}

// Uint16 adds the field key with i as a uint16 to the event context.
func (e Event) Uint16(key string, i uint16) Event {
	if !e.enabled {
		return e
	}
	return e.append(fUint16(key, i))
}

// Uint32 adds the field key with i as a uint32 to the event context.
func (e Event) Uint32(key string, i uint32) Event {
	if !e.enabled {
		return e
	}
	return e.append(fUint32(key, i))
}

// Uint64 adds the field key with i as a uint64 to the event context.
func (e Event) Uint64(key string, i uint64) Event {
	if !e.enabled {
		return e
	}
	return e.append(fUint64(key, i))
}

// Float32 adds the field key with f as a float32 to the event context.
func (e Event) Float32(key string, f float32) Event {
	if !e.enabled {
		return e
	}
	return e.append(fFloat32(key, f))
}

// Float64 adds the field key with f as a float64 to the event context.
func (e Event) Float64(key string, f float64) Event {
	if !e.enabled {
		return e
	}
	return e.append(fFloat64(key, f))
}

// Timestamp adds the current local time as UNIX timestamp to the event context with the "time" key.
// To customize the key name, change zerolog.TimestampFieldName.
func (e Event) Timestamp() Event {
	if !e.enabled {
		return e
	}
	return e.append(fTimestamp())
}

// Time adds the field key with t formated as string using zerolog.TimeFieldFormat.
func (e Event) Time(key string, t time.Time) Event {
	if !e.enabled {
		return e
	}
	return e.append(fTime(key, t))
}
