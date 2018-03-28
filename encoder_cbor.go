// +build binary_log

package zerolog

// This file contains bindings to do binary encoding.

import (
	"time"

	"github.com/rs/zerolog/internal/cbor"
)

func appendInterface(dst []byte, i interface{}) []byte {
	return cbor.AppendInterface(dst, i)
}

func appendKey(dst []byte, s string) []byte {
	return cbor.AppendKey(dst, s)
}

func appendFloats64(dst []byte, f []float64) []byte {
	return cbor.AppendFloats64(dst, f)
}

func appendFloat64(dst []byte, f float64) []byte {
	return cbor.AppendFloat64(dst, f)
}

func appendFloats32(dst []byte, f []float32) []byte {
	return cbor.AppendFloats32(dst, f)
}

func appendFloat32(dst []byte, f float32) []byte {
	return cbor.AppendFloat32(dst, f)
}

func appendUints64(dst []byte, i []uint64) []byte {
	return cbor.AppendUints64(dst, i)
}

func appendUint64(dst []byte, i uint64) []byte {
	return cbor.AppendUint64(dst, i)
}

func appendUints32(dst []byte, i []uint32) []byte {
	return cbor.AppendUints32(dst, i)
}

func appendUint32(dst []byte, i uint32) []byte {
	return cbor.AppendUint32(dst, i)
}

func appendUints16(dst []byte, i []uint16) []byte {
	return cbor.AppendUints16(dst, i)
}

func appendUint16(dst []byte, i uint16) []byte {
	return cbor.AppendUint16(dst, i)
}

func appendUints8(dst []byte, i []uint8) []byte {
	return cbor.AppendUints8(dst, i)
}

func appendUint8(dst []byte, i uint8) []byte {
	return cbor.AppendUint8(dst, i)
}

func appendUints(dst []byte, i []uint) []byte {
	return cbor.AppendUints(dst, i)
}

func appendUint(dst []byte, i uint) []byte {
	return cbor.AppendUint(dst, i)
}

func appendInts64(dst []byte, i []int64) []byte {
	return cbor.AppendInts64(dst, i)
}

func appendInt64(dst []byte, i int64) []byte {
	return cbor.AppendInt64(dst, i)
}

func appendInts32(dst []byte, i []int32) []byte {
	return cbor.AppendInts32(dst, i)
}

func appendInt32(dst []byte, i int32) []byte {
	return cbor.AppendInt32(dst, i)
}

func appendInts16(dst []byte, i []int16) []byte {
	return cbor.AppendInts16(dst, i)
}

func appendInt16(dst []byte, i int16) []byte {
	return cbor.AppendInt16(dst, i)
}

func appendInts8(dst []byte, i []int8) []byte {
	return cbor.AppendInts8(dst, i)
}

func appendInt8(dst []byte, i int8) []byte {
	return cbor.AppendInt8(dst, i)
}

func appendInts(dst []byte, i []int) []byte {
	return cbor.AppendInts(dst, i)
}

func appendInt(dst []byte, i int) []byte {
	return cbor.AppendInt(dst, i)
}

func appendBools(dst []byte, b []bool) []byte {
	return cbor.AppendBools(dst, b)
}

func appendBool(dst []byte, b bool) []byte {
	return cbor.AppendBool(dst, b)
}

func appendError(dst []byte, e error) []byte {
	return cbor.AppendError(dst, e)
}

func appendErrors(dst []byte, e []error) []byte {
	return cbor.AppendErrors(dst, e)
}

func appendString(dst []byte, s string) []byte {
	return cbor.AppendString(dst, s)
}

func appendStrings(dst []byte, s []string) []byte {
	return cbor.AppendStrings(dst, s)
}

func appendDuration(dst []byte, t time.Duration, d time.Duration, fmt bool) []byte {
	return cbor.AppendDuration(dst, t, d, fmt)
}

func appendDurations(dst []byte, t []time.Duration, d time.Duration, fmt bool) []byte {
	return cbor.AppendDurations(dst, t, d, fmt)
}

func appendTimes(dst []byte, t []time.Time, fmt string) []byte {
	return cbor.AppendTimes(dst, t, fmt)
}

func appendTime(dst []byte, t time.Time, fmt string) []byte {
	return cbor.AppendTime(dst, t, fmt)
}

func appendEndMarker(dst []byte) []byte {
	return cbor.AppendEndMarker(dst)
}

func appendLineBreak(dst []byte) []byte {
	// No line breaks needed in binary format.
	return dst
}

func appendBeginMarker(dst []byte) []byte {
	return cbor.AppendBeginMarker(dst)
}

func appendBytes(dst []byte, b []byte) []byte {
	return cbor.AppendBytes(dst, b)
}

func appendArrayStart(dst []byte) []byte {
	return cbor.AppendArrayStart(dst)
}

func appendArrayEnd(dst []byte) []byte {
	return cbor.AppendArrayEnd(dst)
}

func appendArrayDelim(dst []byte) []byte {
	return cbor.AppendArrayDelim(dst)
}

func appendObjectData(dst []byte, src []byte) []byte {
        // Map begin character is present in the src, which
        // should not be copied when appending to existing data.
	return cbor.AppendObjectData(dst, src[1:])
}

func appendJSON(dst []byte, j []byte) []byte {
	return cbor.AppendEmbeddedJSON(dst, j)
}

func appendNil(dst []byte) []byte {
	return cbor.AppendNull(dst)
}

func appendHex(dst []byte, val []byte) []byte {
	return cbor.AppendHex(dst, val)
}

// decodeIfBinaryToString - converts a binary formatted log msg to a
// JSON formatted String Log message.
func decodeIfBinaryToString(in []byte) string {
	return cbor.DecodeIfBinaryToString(in)
}

func decodeObjectToStr(in []byte) string {
	return cbor.DecodeObjectToStr(in)
}

// decodeIfBinaryToBytes - converts a binary formatted log msg to a
// JSON formatted Bytes Log message.
func decodeIfBinaryToBytes(in []byte) []byte {
	return cbor.DecodeIfBinaryToBytes(in)
}
