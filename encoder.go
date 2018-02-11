// +build !zerolog_binary

package zerolog

import (
	"time"

	"github.com/rs/zerolog/internal/json"
)

func appendKey(dst []byte, key string) []byte {
	return json.AppendKey(dst, key)
}

func appendError(dst []byte, err error) []byte {
	return json.AppendError(dst, err)
}

func appendErrors(dst []byte, errs []error) []byte {
	return json.AppendErrors(dst, errs)
}

func appendStrings(dst []byte, vals []string) []byte {
	return json.AppendStrings(dst, vals)
}

func appendString(dst []byte, s string) []byte {
	return json.AppendString(dst, s)
}

func appendBytes(dst, b []byte) []byte {
	return json.AppendBytes(dst, b)
}

func appendTime(dst []byte, t time.Time, format string) []byte {
	return json.AppendTime(dst, t, format)
}

func appendTimes(dst []byte, vals []time.Time, format string) []byte {
	return json.AppendTimes(dst, vals, format)
}

func appendDuration(dst []byte, d time.Duration, unit time.Duration, useInt bool) []byte {
	return json.AppendDuration(dst, d, unit, useInt)
}

func appendDurations(dst []byte, vals []time.Duration, unit time.Duration, useInt bool) []byte {
	return json.AppendDurations(dst, vals, unit, useInt)
}

func appendBool(dst []byte, val bool) []byte {
	return json.AppendBool(dst, val)
}

func appendBools(dst []byte, vals []bool) []byte {
	return json.AppendBools(dst, vals)
}

func appendInt(dst []byte, val int) []byte {
	return json.AppendInt(dst, val)
}

func appendInts(dst []byte, vals []int) []byte {
	return json.AppendInts(dst, vals)
}

func appendInt8(dst []byte, val int8) []byte {
	return json.AppendInt8(dst, val)
}

func appendInts8(dst []byte, vals []int8) []byte {
	return json.AppendInts8(dst, vals)
}

func appendInt16(dst []byte, val int16) []byte {
	return json.AppendInt16(dst, val)
}

func appendInts16(dst []byte, vals []int16) []byte {
	return json.AppendInts16(dst, vals)
}

func appendInt32(dst []byte, val int32) []byte {
	return json.AppendInt32(dst, val)
}

func appendInts32(dst []byte, vals []int32) []byte {
	return json.AppendInts32(dst, vals)
}

func appendInt64(dst []byte, val int64) []byte {
	return json.AppendInt64(dst, val)
}

func appendInts64(dst []byte, vals []int64) []byte {
	return json.AppendInts64(dst, vals)
}

func appendUint(dst []byte, val uint) []byte {
	return json.AppendUint(dst, val)
}

func appendUints(dst []byte, vals []uint) []byte {
	return json.AppendUints(dst, vals)
}

func appendUint8(dst []byte, val uint8) []byte {
	return json.AppendUint8(dst, val)
}

func appendUints8(dst []byte, vals []uint8) []byte {
	return json.AppendUints8(dst, vals)
}

func appendUint16(dst []byte, val uint16) []byte {
	return json.AppendUint16(dst, val)
}

func appendUints16(dst []byte, vals []uint16) []byte {
	return json.AppendUints16(dst, vals)
}

func appendUint32(dst []byte, val uint32) []byte {
	return json.AppendUint32(dst, val)
}

func appendUints32(dst []byte, vals []uint32) []byte {
	return json.AppendUints32(dst, vals)
}

func appendUint64(dst []byte, val uint64) []byte {
	return json.AppendUint64(dst, val)
}

func appendUints64(dst []byte, vals []uint64) []byte {
	return json.AppendUints64(dst, vals)
}

func appendFloat(dst []byte, val float64, bitSize int) []byte {
	return json.AppendFloat(dst, val, bitSize)
}

func appendFloat32(dst []byte, val float32) []byte {
	return json.AppendFloat32(dst, val)
}

func appendFloats32(dst []byte, vals []float32) []byte {
	return json.AppendFloats32(dst, vals)
}

func appendFloat64(dst []byte, val float64) []byte {
	return json.AppendFloat64(dst, val)
}

func appendFloats64(dst []byte, vals []float64) []byte {
	return json.AppendFloats64(dst, vals)
}

func appendInterface(dst []byte, i interface{}) []byte {
	return json.AppendInterface(dst, i)
}
