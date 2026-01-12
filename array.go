package zerolog

import (
	"context"
	"net"
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
	buf   []byte
	stack bool            // enable error stack trace
	ctx   context.Context // Optional Go context
	ch    []Hook          // hooks
}

func putArray(a *Array) {
	// prevent any subsequent use of the Array contextual state and truncate the buffer
	a.stack = false
	a.ctx = nil
	a.ch = nil
	a.buf = a.buf[:0]

	// Proper usage of a sync.Pool requires each entry to have approximately
	// the same memory cost. To obtain this property when the stored type
	// contains a variably-sized buffer, we add a hard limit on the maximum buffer
	// to place back in the pool.
	//
	// See https://golang.org/issue/23199
	const maxSize = 1 << 16 // 64KiB
	if cap(a.buf) <= maxSize {
		arrayPool.Put(a)
	}
}

// Arr creates an array to be added to an Event or Context.
// WARNING: This function is deprecated because it does not preserve
// the stack, hooks, and context from the parent event.
// Deprecated: Use Event.CreateArray or Context.CreateArray instead.
func Arr() *Array {
	a := arrayPool.Get().(*Array)
	a.buf = a.buf[:0]
	a.stack = false
	a.ctx = nil
	a.ch = nil
	return a
}

// MarshalZerologArray method here is no-op - since data is
// already in the needed format.
func (*Array) MarshalZerologArray(*Array) {
	// untestable: there's no code to be covered
}

func (a *Array) write(dst []byte) []byte {
	dst = enc.AppendArrayStart(dst)
	if len(a.buf) > 0 {
		dst = append(dst, a.buf...)
	}
	dst = enc.AppendArrayEnd(dst)
	putArray(a)
	return dst
}

// Object marshals an object that implement the LogObjectMarshaler
// interface and appends it to the array.
func (a *Array) Object(obj LogObjectMarshaler) *Array {
	a.buf = appendObject(enc.AppendArrayDelim(a.buf), obj, a.stack, a.ctx, a.ch)
	return a
}

// Str appends the val as a string to the array.
func (a *Array) Str(val string) *Array {
	a.buf = enc.AppendString(enc.AppendArrayDelim(a.buf), val)
	return a
}

// Bytes appends the val as a string to the array.
func (a *Array) Bytes(val []byte) *Array {
	a.buf = enc.AppendBytes(enc.AppendArrayDelim(a.buf), val)
	return a
}

// Hex appends the val as a hex string to the array.
func (a *Array) Hex(val []byte) *Array {
	a.buf = enc.AppendHex(enc.AppendArrayDelim(a.buf), val)
	return a
}

// RawJSON adds already encoded JSON to the array.
func (a *Array) RawJSON(val []byte) *Array {
	a.buf = appendJSON(enc.AppendArrayDelim(a.buf), val)
	return a
}

// Err serializes and appends the err to the array.
func (a *Array) Err(err error) *Array {
	switch m := ErrorMarshalFunc(err).(type) {
	case nil:
		a.buf = enc.AppendNil(enc.AppendArrayDelim(a.buf))
	case LogObjectMarshaler:
		a = a.Object(m)
	case error:
		if !isNilValue(m) {
			a.buf = enc.AppendString(enc.AppendArrayDelim(a.buf), m.Error())
		}
	case string:
		a.buf = enc.AppendString(enc.AppendArrayDelim(a.buf), m)
	default:
		a.buf = enc.AppendInterface(enc.AppendArrayDelim(a.buf), m)
	}

	return a
}

// Errs serializes and appends errors to the array.
func (a *Array) Errs(errs []error) *Array {
	for _, err := range errs {
		switch m := ErrorMarshalFunc(err).(type) {
		case nil:
			a = a.Interface(nil)
		case LogObjectMarshaler:
			a = a.Object(m)
		case error:
			if !isNilValue(m) {
				a = a.Str(m.Error())
			}
		case string:
			a = a.Str(m)
		default:
			a = a.Interface(m)
		}
	}
	return a
}

// Bool appends the val as a bool to the array.
func (a *Array) Bool(b bool) *Array {
	a.buf = enc.AppendBool(enc.AppendArrayDelim(a.buf), b)
	return a
}

// Int appends i as a int to the array.
func (a *Array) Int(i int) *Array {
	a.buf = enc.AppendInt(enc.AppendArrayDelim(a.buf), i)
	return a
}

// Int8 appends i as a int8 to the array.
func (a *Array) Int8(i int8) *Array {
	a.buf = enc.AppendInt8(enc.AppendArrayDelim(a.buf), i)
	return a
}

// Int16 appends i as a int16 to the array.
func (a *Array) Int16(i int16) *Array {
	a.buf = enc.AppendInt16(enc.AppendArrayDelim(a.buf), i)
	return a
}

// Int32 appends i as a int32 to the array.
func (a *Array) Int32(i int32) *Array {
	a.buf = enc.AppendInt32(enc.AppendArrayDelim(a.buf), i)
	return a
}

// Int64 appends i as a int64 to the array.
func (a *Array) Int64(i int64) *Array {
	a.buf = enc.AppendInt64(enc.AppendArrayDelim(a.buf), i)
	return a
}

// Uint appends i as a uint to the array.
func (a *Array) Uint(i uint) *Array {
	a.buf = enc.AppendUint(enc.AppendArrayDelim(a.buf), i)
	return a
}

// Uint8 appends i as a uint8 to the array.
func (a *Array) Uint8(i uint8) *Array {
	a.buf = enc.AppendUint8(enc.AppendArrayDelim(a.buf), i)
	return a
}

// Uint16 appends i as a uint16 to the array.
func (a *Array) Uint16(i uint16) *Array {
	a.buf = enc.AppendUint16(enc.AppendArrayDelim(a.buf), i)
	return a
}

// Uint32 appends i as a uint32 to the array.
func (a *Array) Uint32(i uint32) *Array {
	a.buf = enc.AppendUint32(enc.AppendArrayDelim(a.buf), i)
	return a
}

// Uint64 appends i as a uint64 to the array.
func (a *Array) Uint64(i uint64) *Array {
	a.buf = enc.AppendUint64(enc.AppendArrayDelim(a.buf), i)
	return a
}

// Float32 appends f as a float32 to the array.
func (a *Array) Float32(f float32) *Array {
	a.buf = enc.AppendFloat32(enc.AppendArrayDelim(a.buf), f, FloatingPointPrecision)
	return a
}

// Float64 appends f as a float64 to the array.
func (a *Array) Float64(f float64) *Array {
	a.buf = enc.AppendFloat64(enc.AppendArrayDelim(a.buf), f, FloatingPointPrecision)
	return a
}

// Time appends t formatted as string using zerolog.TimeFieldFormat.
func (a *Array) Time(t time.Time) *Array {
	a.buf = enc.AppendTime(enc.AppendArrayDelim(a.buf), t, TimeFieldFormat)
	return a
}

// Dur appends d to the array.
func (a *Array) Dur(d time.Duration) *Array {
	a.buf = enc.AppendDuration(enc.AppendArrayDelim(a.buf), d, DurationFieldUnit, DurationFieldFormat, DurationFieldInteger, FloatingPointPrecision)
	return a
}

// Interface appends i marshaled using reflection.
func (a *Array) Interface(i interface{}) *Array {
	if obj, ok := i.(LogObjectMarshaler); ok {
		return a.Object(obj)
	}
	a.buf = enc.AppendInterface(enc.AppendArrayDelim(a.buf), i)
	return a
}

// IPAddr adds a net.IP IPv4 or IPv6 address to the array
func (a *Array) IPAddr(ip net.IP) *Array {
	a.buf = enc.AppendIPAddr(enc.AppendArrayDelim(a.buf), ip)
	return a
}

// IPPrefix adds a net.IPNet IPv4 or IPv6 Prefix (IP + mask) to the array
func (a *Array) IPPrefix(pfx net.IPNet) *Array {
	a.buf = enc.AppendIPPrefix(enc.AppendArrayDelim(a.buf), pfx)
	return a
}

// MACAddr adds a net.HardwareAddr MAC (Ethernet) address to the array
func (a *Array) MACAddr(ha net.HardwareAddr) *Array {
	a.buf = enc.AppendMACAddr(enc.AppendArrayDelim(a.buf), ha)
	return a
}

// Dict adds the dict Event to the array
func (a *Array) Dict(dict *Event) *Array {
	dict.buf = enc.AppendEndMarker(dict.buf)
	a.buf = append(enc.AppendArrayDelim(a.buf), dict.buf...)
	putEvent(dict)
	return a
}

// Type adds the val's type using reflection to the array.
func (a *Array) Type(val interface{}) *Array {
	a.buf = enc.AppendType(enc.AppendArrayDelim(a.buf), val)
	return a
}
