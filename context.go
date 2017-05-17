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

// Str adds the field key with val as a string to the logger context.
func (c Context) Str(key, val string) Context {
	return c.append(fStr(key, val))
}

// Err adds the field "error" with err as a string to the logger context.
// To customize the key name, change zerolog.ErrorFieldName.
func (c Context) Err(err error) Context {
	return c.append(fErr(err))
}

// Bool adds the field key with val as a Boolean to the logger context.
func (c Context) Bool(key string, b bool) Context {
	return c.append(fBool(key, b))
}

// Int adds the field key with i as a int to the logger context.
func (c Context) Int(key string, i int) Context {
	return c.append(fInt(key, i))
}

// Int8 adds the field key with i as a int8 to the logger context.
func (c Context) Int8(key string, i int8) Context {
	return c.append(fInt8(key, i))
}

// Int16 adds the field key with i as a int16 to the logger context.
func (c Context) Int16(key string, i int16) Context {
	return c.append(fInt16(key, i))
}

// Int32 adds the field key with i as a int32 to the logger context.
func (c Context) Int32(key string, i int32) Context {
	return c.append(fInt32(key, i))
}

// Int64 adds the field key with i as a int64 to the logger context.
func (c Context) Int64(key string, i int64) Context {
	return c.append(fInt64(key, i))
}

// Uint adds the field key with i as a uint to the logger context.
func (c Context) Uint(key string, i uint) Context {
	return c.append(fUint(key, i))
}

// Uint8 adds the field key with i as a uint8 to the logger context.
func (c Context) Uint8(key string, i uint8) Context {
	return c.append(fUint8(key, i))
}

// Uint16 adds the field key with i as a uint16 to the logger context.
func (c Context) Uint16(key string, i uint16) Context {
	return c.append(fUint16(key, i))
}

// Uint32 adds the field key with i as a uint32 to the logger context.
func (c Context) Uint32(key string, i uint32) Context {
	return c.append(fUint32(key, i))
}

// Uint64 adds the field key with i as a uint64 to the logger context.
func (c Context) Uint64(key string, i uint64) Context {
	return c.append(fUint64(key, i))
}

// Float32 adds the field key with f as a float32 to the logger context.
func (c Context) Float32(key string, f float32) Context {
	return c.append(fFloat32(key, f))
}

// Float64 adds the field key with f as a float64 to the logger context.
func (c Context) Float64(key string, f float64) Context {
	return c.append(fFloat64(key, f))
}

// Timestamp adds the current local time as UNIX timestamp to the logger context with the "time" key.
// To customize the key name, change zerolog.TimestampFieldName.
func (c Context) Timestamp() Context {
	return c.append(fTimestamp())
}

// Time adds the field key with t formated as string using zerolog.TimeFieldFormat.
func (c Context) Time(key string, t time.Time) Context {
	return c.append(fTime(key, t))
}
