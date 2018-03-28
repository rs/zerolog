// +build !binary_log

package zerolog

// encoder_json.go file contains bindings to generate
// JSON encoded byte stream.

import (
	"strconv"
	"time"

	"github.com/rs/zerolog/internal/json"
)

func appendInterface(dst []byte, i interface{}) []byte {
	return json.AppendInterface(dst, i)
}

func appendKey(dst []byte, s string) []byte {
	return json.AppendKey(dst, s)
}

func appendFloats64(dst []byte, f []float64) []byte {
	return json.AppendFloats64(dst, f)
}

func appendFloat64(dst []byte, f float64) []byte {
	return json.AppendFloat64(dst, f)
}

func appendFloats32(dst []byte, f []float32) []byte {
	return json.AppendFloats32(dst, f)
}

func appendFloat32(dst []byte, f float32) []byte {
	return json.AppendFloat32(dst, f)
}

func appendUints64(dst []byte, i []uint64) []byte {
	return json.AppendUints64(dst, i)
}

func appendUint64(dst []byte, i uint64) []byte {
	return strconv.AppendUint(dst, uint64(i), 10)
}

func appendUints32(dst []byte, i []uint32) []byte {
	return json.AppendUints32(dst, i)
}

func appendUint32(dst []byte, i uint32) []byte {
	return strconv.AppendUint(dst, uint64(i), 10)
}

func appendUints16(dst []byte, i []uint16) []byte {
	return json.AppendUints16(dst, i)
}

func appendUint16(dst []byte, i uint16) []byte {
	return strconv.AppendUint(dst, uint64(i), 10)
}

func appendUints8(dst []byte, i []uint8) []byte {
	return json.AppendUints8(dst, i)
}

func appendUint8(dst []byte, i uint8) []byte {
	return strconv.AppendUint(dst, uint64(i), 10)
}

func appendUints(dst []byte, i []uint) []byte {
	return json.AppendUints(dst, i)
}

func appendUint(dst []byte, i uint) []byte {
	return strconv.AppendUint(dst, uint64(i), 10)
}

func appendInts64(dst []byte, i []int64) []byte {
	return json.AppendInts64(dst, i)
}

func appendInt64(dst []byte, i int64) []byte {
	return strconv.AppendInt(dst, int64(i), 10)
}

func appendInts32(dst []byte, i []int32) []byte {
	return json.AppendInts32(dst, i)
}

func appendInt32(dst []byte, i int32) []byte {
	return strconv.AppendInt(dst, int64(i), 10)
}

func appendInts16(dst []byte, i []int16) []byte {
	return json.AppendInts16(dst, i)
}

func appendInt16(dst []byte, i int16) []byte {
	return strconv.AppendInt(dst, int64(i), 10)
}

func appendInts8(dst []byte, i []int8) []byte {
	return json.AppendInts8(dst, i)
}

func appendInt8(dst []byte, i int8) []byte {
	return strconv.AppendInt(dst, int64(i), 10)
}

func appendInts(dst []byte, i []int) []byte {
	return json.AppendInts(dst, i)
}

func appendInt(dst []byte, i int) []byte {
	return strconv.AppendInt(dst, int64(i), 10)
}

func appendBools(dst []byte, b []bool) []byte {
	return json.AppendBools(dst, b)
}

func appendBool(dst []byte, b bool) []byte {
	return strconv.AppendBool(dst, b)
}

func appendError(dst []byte, e error) []byte {
	return json.AppendError(dst, e)
}

func appendErrors(dst []byte, e []error) []byte {
	return json.AppendErrors(dst, e)
}

func appendString(dst []byte, s string) []byte {
	return json.AppendString(dst, s)
}

func appendStrings(dst []byte, s []string) []byte {
	return json.AppendStrings(dst, s)
}

func appendDuration(dst []byte, t time.Duration, d time.Duration, fmt bool) []byte {
	return json.AppendDuration(dst, t, d, fmt)
}

func appendDurations(dst []byte, t []time.Duration, d time.Duration, fmt bool) []byte {
	return json.AppendDurations(dst, t, d, fmt)
}

func appendTimes(dst []byte, t []time.Time, fmt string) []byte {
	return json.AppendTimes(dst, t, fmt)
}

func appendTime(dst []byte, t time.Time, fmt string) []byte {
	return json.AppendTime(dst, t, fmt)
}

func appendEndMarker(dst []byte) []byte {
	return append(dst, '}')
}

func appendLineBreak(dst []byte) []byte {
	return append(dst, '\n')
}

func appendBeginMarker(dst []byte) []byte {
	return append(dst, '{')
}

func appendBytes(dst []byte, b []byte) []byte {
	return json.AppendBytes(dst, b)
}

func appendArrayStart(dst []byte) []byte {
	return append(dst, '[')
}

func appendArrayEnd(dst []byte) []byte {
	return append(dst, ']')
}

func appendArrayDelim(dst []byte) []byte {
	if len(dst) > 0 {
		return append(dst, ',')
	}
	return dst
}

func appendObjectData(dst []byte, src []byte) []byte {
	return json.AppendObjectData(dst, src)
}

func appendJSON(dst []byte, j []byte) []byte {
	return append(dst, j...)
}

func appendNil(dst []byte) []byte {
	return append(dst, "null"...)
}

func decodeIfBinaryToString(in []byte) string {
	return string(in)
}

func decodeObjectToStr(in []byte) string {
	return string(in)
}

func decodeIfBinaryToBytes(in []byte) []byte {
	return in
}

func appendHex(in []byte, val []byte) []byte {
    return json.AppendHex(in, val)
}
