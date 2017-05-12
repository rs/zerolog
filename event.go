package zerolog

import (
	"bytes"
	"fmt"
	"sync"
	"time"
)

var pool = &sync.Pool{
	New: func() interface{} {
		return bytes.NewBuffer(make([]byte, 500))
	},
}

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

func (e Event) Str(key, val string) Event {
	if !e.enabled {
		return e
	}
	return e.append(fStr(key, val))
}

func (e Event) Err(err error) Event {
	if !e.enabled {
		return e
	}
	return e.append(fErr(err))
}

func (e Event) Bool(key string, b bool) Event {
	if !e.enabled {
		return e
	}
	return e.append(fBool(key, b))
}

func (e Event) Int(key string, i int) Event {
	if !e.enabled {
		return e
	}
	return e.append(fInt(key, i))
}

func (e Event) Int8(key string, i int8) Event {
	if !e.enabled {
		return e
	}
	return e.append(fInt8(key, i))
}

func (e Event) Int16(key string, i int16) Event {
	if !e.enabled {
		return e
	}
	return e.append(fInt16(key, i))
}

func (e Event) Int32(key string, i int32) Event {
	if !e.enabled {
		return e
	}
	return e.append(fInt32(key, i))
}

func (e Event) Int64(key string, i int64) Event {
	if !e.enabled {
		return e
	}
	return e.append(fInt64(key, i))
}

func (e Event) Uint(key string, i uint) Event {
	if !e.enabled {
		return e
	}
	return e.append(fUint(key, i))
}

func (e Event) Uint8(key string, i uint8) Event {
	if !e.enabled {
		return e
	}
	return e.append(fUint8(key, i))
}

func (e Event) Uint16(key string, i uint16) Event {
	if !e.enabled {
		return e
	}
	return e.append(fUint16(key, i))
}

func (e Event) Uint32(key string, i uint32) Event {
	if !e.enabled {
		return e
	}
	return e.append(fUint32(key, i))
}

func (e Event) Uint64(key string, i uint64) Event {
	if !e.enabled {
		return e
	}
	return e.append(fUint64(key, i))
}

func (e Event) Float32(key string, f float32) Event {
	if !e.enabled {
		return e
	}
	return e.append(fFloat32(key, f))
}

func (e Event) Float64(key string, f float64) Event {
	if !e.enabled {
		return e
	}
	return e.append(fFloat64(key, f))
}

func (e Event) Timestamp() Event {
	if !e.enabled {
		return e
	}
	return e.append(fTimestamp())
}

func (e Event) Time(key string, t time.Time) Event {
	if !e.enabled {
		return e
	}
	return e.append(fTime(key, t))
}
