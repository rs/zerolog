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

// Dict adds the field key with the dict to the logger context.
func (c Context) Dict(key string, dict *Event) Context {
	dict.buf = append(dict.buf, '}')
	c.l.context = append(appendKey(c.l.context, key), dict.buf...)
	eventPool.Put(dict)
	return c
}

// Str adds the field key with val as a string to the logger context.
func (c Context) Str(key, val string) Context {
	c.l.context = appendString(c.l.context, key, val)
	return c
}

// Err adds the field "error" with err as a string to the logger context.
// To customize the key name, change zerolog.ErrorFieldName.
func (c Context) Err(err error) Context {
	c.l.context = appendError(c.l.context, err)
	return c
}

// Bool adds the field key with val as a Boolean to the logger context.
func (c Context) Bool(key string, b bool) Context {
	c.l.context = appendBool(c.l.context, key, b)
	return c
}

// Int adds the field key with i as a int to the logger context.
func (c Context) Int(key string, i int) Context {
	c.l.context = appendInt(c.l.context, key, i)
	return c
}

// Int8 adds the field key with i as a int8 to the logger context.
func (c Context) Int8(key string, i int8) Context {
	c.l.context = appendInt8(c.l.context, key, i)
	return c
}

// Int16 adds the field key with i as a int16 to the logger context.
func (c Context) Int16(key string, i int16) Context {
	c.l.context = appendInt16(c.l.context, key, i)
	return c
}

// Int32 adds the field key with i as a int32 to the logger context.
func (c Context) Int32(key string, i int32) Context {
	c.l.context = appendInt32(c.l.context, key, i)
	return c
}

// Int64 adds the field key with i as a int64 to the logger context.
func (c Context) Int64(key string, i int64) Context {
	c.l.context = appendInt64(c.l.context, key, i)
	return c
}

// Uint adds the field key with i as a uint to the logger context.
func (c Context) Uint(key string, i uint) Context {
	c.l.context = appendUint(c.l.context, key, i)
	return c
}

// Uint8 adds the field key with i as a uint8 to the logger context.
func (c Context) Uint8(key string, i uint8) Context {
	c.l.context = appendUint8(c.l.context, key, i)
	return c
}

// Uint16 adds the field key with i as a uint16 to the logger context.
func (c Context) Uint16(key string, i uint16) Context {
	c.l.context = appendUint16(c.l.context, key, i)
	return c
}

// Uint32 adds the field key with i as a uint32 to the logger context.
func (c Context) Uint32(key string, i uint32) Context {
	c.l.context = appendUint32(c.l.context, key, i)
	return c
}

// Uint64 adds the field key with i as a uint64 to the logger context.
func (c Context) Uint64(key string, i uint64) Context {
	c.l.context = appendUint64(c.l.context, key, i)
	return c
}

// Float32 adds the field key with f as a float32 to the logger context.
func (c Context) Float32(key string, f float32) Context {
	c.l.context = appendFloat32(c.l.context, key, f)
	return c
}

// Float64 adds the field key with f as a float64 to the logger context.
func (c Context) Float64(key string, f float64) Context {
	c.l.context = appendFloat64(c.l.context, key, f)
	return c
}

// Timestamp adds the current local time as UNIX timestamp to the logger context with the "time" key.
// To customize the key name, change zerolog.TimestampFieldName.
func (c Context) Timestamp() Context {
	if len(c.l.context) > 0 {
		c.l.context[0] = 1
	} else {
		c.l.context = append(c.l.context, 1)
	}
	return c
}

// Time adds the field key with t formated as string using zerolog.TimeFieldFormat.
func (c Context) Time(key string, t time.Time) Context {
	c.l.context = appendTime(c.l.context, key, t)
	return c
}

// Object adds the field key with obj marshaled using reflection.
func (c Context) Object(key string, obj interface{}) Context {
	c.l.context = appendObject(c.l.context, key, obj)
	return c
}
