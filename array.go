package zerolog

import (
	"sync"
	"time"

	"github.com/rs/zerolog/internal/cbor"
	"github.com/rs/zerolog/internal/json"
)

var arrayPool = &sync.Pool{
	New: func() interface{} {
		return &Array{
			buf: make([]byte, 0, 500),
		}
	},
}

type Array struct {
	buf      []byte
	isBinary bool
}

// Arr creates an array to be added to an Event or Context.
func Arr() *Array {
	a := arrayPool.Get().(*Array)
	a.buf = a.buf[:0]
	a.isBinary = false
	return a
}

func ArrBinary() *Array {
	a := arrayPool.Get().(*Array)
	a.buf = a.buf[:0]
	a.isBinary = true
	return a
}

func (*Array) MarshalZerologArray(*Array) {
}

func (a *Array) write(dst []byte) []byte {
	if a.isBinary {
		dst = cbor.AppendArrayStart(dst)
	} else {
		dst = json.AppendArrayStart(dst)
	}
	if len(a.buf) > 0 {
		dst = append(append(dst, a.buf...))
	}
	if a.isBinary {
		dst = cbor.AppendArrayEnd(dst)
	} else {
		dst = json.AppendArrayEnd(dst)
	}
	arrayPool.Put(a)
	return dst
}

// Object marshals an object that implement the LogObjectMarshaler
// interface and append it to the array.
func (a *Array) Object(obj LogObjectMarshaler) *Array {
	if a.isBinary {
		a.buf = cbor.AppendArrayDelim(a.buf)
	} else {
		a.buf = json.AppendArrayDelim(a.buf)
	}
	e := Dict()
	e.isBinary = a.isBinary
	obj.MarshalZerologObject(e)
	if a.isBinary {
		e.buf = cbor.AppendEndMarker(e.buf, false)
	} else {
		e.buf = json.AppendEndMarker(e.buf, false)
	}
	a.buf = append(a.buf, e.buf...)
	eventPool.Put(e)
	return a
}

// Str append the val as a string to the array.
func (a *Array) Str(val string) *Array {
	if a.isBinary {
		a.buf = cbor.AppendArrayDelim(a.buf)
		a.buf = cbor.AppendString(a.buf, val)
	} else {
		a.buf = json.AppendArrayDelim(a.buf)
		a.buf = json.AppendString(a.buf, val)
	}
	return a
}

// Bytes append the val as a string to the array.
func (a *Array) Bytes(val []byte) *Array {
	if a.isBinary {
		a.buf = cbor.AppendArrayDelim(a.buf)
		a.buf = cbor.AppendBytes(a.buf, val)
	} else {
		a.buf = json.AppendArrayDelim(a.buf)
		a.buf = json.AppendBytes(a.buf, val)
	}
	return a
}

// Err append the err as a string to the array.
func (a *Array) Err(err error) *Array {
	if a.isBinary {
		a.buf = cbor.AppendArrayDelim(a.buf)
		a.buf = cbor.AppendError(a.buf, err)
	} else {
		a.buf = json.AppendArrayDelim(a.buf)
		a.buf = json.AppendError(a.buf, err)
	}
	return a
}

// Bool append the val as a bool to the array.
func (a *Array) Bool(b bool) *Array {
	if a.isBinary {
		a.buf = cbor.AppendArrayDelim(a.buf)
		a.buf = cbor.AppendBool(a.buf, b)
	} else {
		a.buf = json.AppendArrayDelim(a.buf)
		a.buf = json.AppendBool(a.buf, b)
	}
	return a
}

// Int append i as a int to the array.
func (a *Array) Int(i int) *Array {
	if a.isBinary {
		a.buf = cbor.AppendArrayDelim(a.buf)
		a.buf = cbor.AppendInt(a.buf, i)
	} else {
		a.buf = json.AppendArrayDelim(a.buf)
		a.buf = json.AppendInt(a.buf, i)
	}
	return a
}

// Int8 append i as a int8 to the array.
func (a *Array) Int8(i int8) *Array {
	if a.isBinary {
		a.buf = cbor.AppendArrayDelim(a.buf)
		a.buf = cbor.AppendInt8(a.buf, i)
	} else {
		a.buf = json.AppendArrayDelim(a.buf)
		a.buf = json.AppendInt8(a.buf, i)
	}
	return a
}

// Int16 append i as a int16 to the array.
func (a *Array) Int16(i int16) *Array {
	if a.isBinary {
		a.buf = cbor.AppendArrayDelim(a.buf)
		a.buf = cbor.AppendInt16(a.buf, i)
	} else {
		a.buf = json.AppendArrayDelim(a.buf)
		a.buf = json.AppendInt16(a.buf, i)
	}
	return a
}

// Int32 append i as a int32 to the array.
func (a *Array) Int32(i int32) *Array {
	if a.isBinary {
		a.buf = cbor.AppendArrayDelim(a.buf)
		a.buf = cbor.AppendInt32(a.buf, i)
	} else {
		a.buf = json.AppendArrayDelim(a.buf)
		a.buf = json.AppendInt32(a.buf, i)
	}
	return a
}

// Int64 append i as a int64 to the array.
func (a *Array) Int64(i int64) *Array {
	if a.isBinary {
		a.buf = cbor.AppendArrayDelim(a.buf)
		a.buf = cbor.AppendInt64(a.buf, i)
	} else {
		a.buf = json.AppendArrayDelim(a.buf)
		a.buf = json.AppendInt64(a.buf, i)
	}
	return a
}

// Uint append i as a uint to the array.
func (a *Array) Uint(i uint) *Array {
	if a.isBinary {
		a.buf = cbor.AppendArrayDelim(a.buf)
		a.buf = cbor.AppendUint(a.buf, i)
	} else {
		a.buf = json.AppendArrayDelim(a.buf)
		a.buf = json.AppendUint(a.buf, i)
	}
	return a
}

// Uint8 append i as a uint8 to the array.
func (a *Array) Uint8(i uint8) *Array {
	if a.isBinary {
		a.buf = cbor.AppendArrayDelim(a.buf)
		a.buf = cbor.AppendUint8(a.buf, i)
	} else {
		a.buf = json.AppendArrayDelim(a.buf)
		a.buf = json.AppendUint8(a.buf, i)
	}
	return a
}

// Uint16 append i as a uint16 to the array.
func (a *Array) Uint16(i uint16) *Array {
	if a.isBinary {
		a.buf = cbor.AppendArrayDelim(a.buf)
		a.buf = cbor.AppendUint16(a.buf, i)
	} else {
		a.buf = json.AppendArrayDelim(a.buf)
		a.buf = json.AppendUint16(a.buf, i)
	}
	return a
}

// Uint32 append i as a uint32 to the array.
func (a *Array) Uint32(i uint32) *Array {
	if a.isBinary {
		a.buf = cbor.AppendArrayDelim(a.buf)
		a.buf = cbor.AppendUint32(a.buf, i)
	} else {
		a.buf = json.AppendArrayDelim(a.buf)
		a.buf = json.AppendUint32(a.buf, i)
	}
	return a
}

// Uint64 append i as a uint64 to the array.
func (a *Array) Uint64(i uint64) *Array {
	if a.isBinary {
		a.buf = cbor.AppendArrayDelim(a.buf)
		a.buf = cbor.AppendUint64(a.buf, i)
	} else {
		a.buf = json.AppendArrayDelim(a.buf)
		a.buf = json.AppendUint64(a.buf, i)
	}
	return a
}

// Float32 append f as a float32 to the array.
func (a *Array) Float32(f float32) *Array {
	if a.isBinary {
		a.buf = cbor.AppendArrayDelim(a.buf)
		a.buf = cbor.AppendFloat32(a.buf, f)
	} else {
		a.buf = json.AppendArrayDelim(a.buf)
		a.buf = json.AppendFloat32(a.buf, f)
	}
	return a
}

// Float64 append f as a float64 to the array.
func (a *Array) Float64(f float64) *Array {
	if a.isBinary {
		a.buf = cbor.AppendArrayDelim(a.buf)
		a.buf = cbor.AppendFloat64(a.buf, f)
	} else {
		a.buf = json.AppendArrayDelim(a.buf)
		a.buf = json.AppendFloat64(a.buf, f)
	}
	return a
}

// Time append t formated as string using zerolog.TimeFieldFormat.
func (a *Array) Time(t time.Time) *Array {
	if a.isBinary {
		a.buf = cbor.AppendArrayDelim(a.buf)
		a.buf = cbor.AppendTime(a.buf, t, TimeFieldFormat)
	} else {
		a.buf = json.AppendArrayDelim(a.buf)
		a.buf = json.AppendTime(a.buf, t, TimeFieldFormat)
	}
	return a
}

// Dur append d to the array.
func (a *Array) Dur(d time.Duration) *Array {
	if a.isBinary {
		a.buf = cbor.AppendArrayDelim(a.buf)
		a.buf = cbor.AppendDuration(a.buf, d, DurationFieldUnit, DurationFieldInteger)
	} else {
		a.buf = json.AppendArrayDelim(a.buf)
		a.buf = json.AppendDuration(a.buf, d, DurationFieldUnit, DurationFieldInteger)
	}
	return a
}

// Interface append i marshaled using reflection.
func (a *Array) Interface(i interface{}) *Array {
	if obj, ok := i.(LogObjectMarshaler); ok {
		return a.Object(obj)
	}
	if a.isBinary {
		a.buf = cbor.AppendArrayDelim(a.buf)
		a.buf = cbor.AppendInterface(a.buf, i)
	} else {
		a.buf = json.AppendArrayDelim(a.buf)
		a.buf = json.AppendInterface(a.buf, i)
	}
	return a
}
