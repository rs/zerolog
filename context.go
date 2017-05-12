package zerolog

import "time"

// Context configures a new sub-logger with contextual fields.
type Context struct {
	l Logger
}

// Logger returns the logger with the context previously set.
func (c Context) Logger() Logger {
	return c.l
}

func (c Context) append(f field) Context {
	return Context{
		l: Logger{
			parent:  c.l,
			w:       c.l.w,
			field:   f.compileJSON(),
			level:   c.l.level,
			sample:  c.l.sample,
			counter: c.l.counter,
		},
	}
}

func (c Context) Str(key, val string) Context {
	return c.append(fStr(key, val))
}

func (c Context) Err(err error) Context {
	return c.append(fErr(err))
}

func (c Context) Bool(key string, b bool) Context {
	return c.append(fBool(key, b))
}

func (c Context) Int(key string, i int) Context {
	return c.append(fInt(key, i))
}

func (c Context) Int8(key string, i int8) Context {
	return c.append(fInt8(key, i))
}

func (c Context) Int16(key string, i int16) Context {
	return c.append(fInt16(key, i))
}

func (c Context) Int32(key string, i int32) Context {
	return c.append(fInt32(key, i))
}

func (c Context) Int64(key string, i int64) Context {
	return c.append(fInt64(key, i))
}

func (c Context) Uint(key string, i uint) Context {
	return c.append(fUint(key, i))
}

func (c Context) Uint8(key string, i uint8) Context {
	return c.append(fUint8(key, i))
}

func (c Context) Uint16(key string, i uint16) Context {
	return c.append(fUint16(key, i))
}

func (c Context) Uint32(key string, i uint32) Context {
	return c.append(fUint32(key, i))
}

func (c Context) Uint64(key string, i uint64) Context {
	return c.append(fUint64(key, i))
}

func (c Context) Float32(key string, f float32) Context {
	return c.append(fFloat32(key, f))
}

func (c Context) Float64(key string, f float64) Context {
	return c.append(fFloat64(key, f))
}

func (c Context) Timestamp() Context {
	return c.append(fTimestamp())
}

func (c Context) Time(key string, t time.Time) Context {
	return c.append(fTime(key, t))
}
