package zerolog

import (
	"sync"
	"time"
)

var arrayPool = &sync.Pool{
	New: func() interface{} {
		return &Array{
			buf: make([]byte, 0, 500),
		}
	},
}

// Array is used to prepopulate an array of items
// which can be re-used to add to log messages.
type Array struct {
	buf []byte
}

// Arr creates an array to be added to an Event or Context.
func Arr() *Array {
	a := arrayPool.Get().(*Array)
	a.buf = a.buf[:0]
	return a
}

// MarshalZerologArray method here is no-op - since data is
// already in the needed format.
func (*Array) MarshalZerologArray(*Array) {
}

func (a *Array) write(dst []byte) []byte {
	dst = appendArrayStart(dst)
	if len(a.buf) > 0 {
		dst = append(append(dst, a.buf...))
	}
	dst = appendArrayEnd(dst)
	arrayPool.Put(a)
	return dst
}

// Object marshals an object that implement the LogObjectMarshaler
// interface and append it to the array.
func (a *Array) Object(obj LogObjectMarshaler) *Array {
	e := Dict()
	obj.MarshalZerologObject(e)
	e.buf = appendEndMarker(e.buf)
	a.buf = append(appendArrayDelim(a.buf), e.buf...)
	eventPool.Put(e)
	return a
}

// Str append the val as a string to the array.
func (a *Array) Str(val string) *Array {
	a.buf = appendString(appendArrayDelim(a.buf), val)
	return a
}

// Bytes append the val as a string to the array.
func (a *Array) Bytes(val []byte) *Array {
	a.buf = appendBytes(appendArrayDelim(a.buf), val)
	return a
}

// Hex append the val as a hex string to the array.
func (a *Array) Hex(val []byte) *Array {
	a.buf = appendHex(appendArrayDelim(a.buf), val)
	return a
}

// Err append the err as a string to the array.
func (a *Array) Err(err error) *Array {
	a.buf = appendError(appendArrayDelim(a.buf), err)
	return a
}

// Bool append the val as a bool to the array.
func (a *Array) Bool(b bool) *Array {
	a.buf = appendBool(appendArrayDelim(a.buf), b)
	return a
}

// Int append i as a int to the array.
func (a *Array) Int(i int) *Array {
	a.buf = appendInt(appendArrayDelim(a.buf), i)
	return a
}

// Int8 append i as a int8 to the array.
func (a *Array) Int8(i int8) *Array {
	a.buf = appendInt8(appendArrayDelim(a.buf), i)
	return a
}

// Int16 append i as a int16 to the array.
func (a *Array) Int16(i int16) *Array {
	a.buf = appendInt16(appendArrayDelim(a.buf), i)
	return a
}

// Int32 append i as a int32 to the array.
func (a *Array) Int32(i int32) *Array {
	a.buf = appendInt32(appendArrayDelim(a.buf), i)
	return a
}

// Int64 append i as a int64 to the array.
func (a *Array) Int64(i int64) *Array {
	a.buf = appendInt64(appendArrayDelim(a.buf), i)
	return a
}

// Uint append i as a uint to the array.
func (a *Array) Uint(i uint) *Array {
	a.buf = appendUint(appendArrayDelim(a.buf), i)
	return a
}

// Uint8 append i as a uint8 to the array.
func (a *Array) Uint8(i uint8) *Array {
	a.buf = appendUint8(appendArrayDelim(a.buf), i)
	return a
}

// Uint16 append i as a uint16 to the array.
func (a *Array) Uint16(i uint16) *Array {
	a.buf = appendUint16(appendArrayDelim(a.buf), i)
	return a
}

// Uint32 append i as a uint32 to the array.
func (a *Array) Uint32(i uint32) *Array {
	a.buf = appendUint32(appendArrayDelim(a.buf), i)
	return a
}

// Uint64 append i as a uint64 to the array.
func (a *Array) Uint64(i uint64) *Array {
	a.buf = appendUint64(appendArrayDelim(a.buf), i)
	return a
}

// Float32 append f as a float32 to the array.
func (a *Array) Float32(f float32) *Array {
	a.buf = appendFloat32(appendArrayDelim(a.buf), f)
	return a
}

// Float64 append f as a float64 to the array.
func (a *Array) Float64(f float64) *Array {
	a.buf = appendFloat64(appendArrayDelim(a.buf), f)
	return a
}

// Time append t formated as string using zerolog.TimeFieldFormat.
func (a *Array) Time(t time.Time) *Array {
	a.buf = appendTime(appendArrayDelim(a.buf), t, TimeFieldFormat)
	return a
}

// Dur append d to the array.
func (a *Array) Dur(d time.Duration) *Array {
	a.buf = appendDuration(appendArrayDelim(a.buf), d, DurationFieldUnit, DurationFieldInteger)
	return a
}

// Interface append i marshaled using reflection.
func (a *Array) Interface(i interface{}) *Array {
	if obj, ok := i.(LogObjectMarshaler); ok {
		return a.Object(obj)
	}
	a.buf = appendInterface(appendArrayDelim(a.buf), i)
	return a
}
