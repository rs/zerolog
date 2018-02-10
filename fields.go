package zerolog

import (
	"sort"
	"time"

	"github.com/rs/zerolog/internal/cbor"
	"github.com/rs/zerolog/internal/json"
)

func appendFieldsText(dst []byte, fields map[string]interface{}) []byte {
	keys := make([]string, 0, len(fields))
	for key := range fields {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		dst = json.AppendKey(dst, key)
		switch val := fields[key].(type) {
		case string:
			dst = json.AppendString(dst, val)
		case []byte:
			dst = json.AppendBytes(dst, val)
		case error:
			dst = json.AppendError(dst, val)
		case []error:
			dst = json.AppendErrors(dst, val)
		case bool:
			dst = json.AppendBool(dst, val)
		case int:
			dst = json.AppendInt(dst, val)
		case int8:
			dst = json.AppendInt8(dst, val)
		case int16:
			dst = json.AppendInt16(dst, val)
		case int32:
			dst = json.AppendInt32(dst, val)
		case int64:
			dst = json.AppendInt64(dst, val)
		case uint:
			dst = json.AppendUint(dst, val)
		case uint8:
			dst = json.AppendUint8(dst, val)
		case uint16:
			dst = json.AppendUint16(dst, val)
		case uint32:
			dst = json.AppendUint32(dst, val)
		case uint64:
			dst = json.AppendUint64(dst, val)
		case float32:
			dst = json.AppendFloat32(dst, val)
		case float64:
			dst = json.AppendFloat64(dst, val)
		case time.Time:
			dst = json.AppendTime(dst, val, TimeFieldFormat)
		case time.Duration:
			dst = json.AppendDuration(dst, val, DurationFieldUnit, DurationFieldInteger)
		case []string:
			dst = json.AppendStrings(dst, val)
		case []bool:
			dst = json.AppendBools(dst, val)
		case []int:
			dst = json.AppendInts(dst, val)
		case []int8:
			dst = json.AppendInts8(dst, val)
		case []int16:
			dst = json.AppendInts16(dst, val)
		case []int32:
			dst = json.AppendInts32(dst, val)
		case []int64:
			dst = json.AppendInts64(dst, val)
		case []uint:
			dst = json.AppendUints(dst, val)
		// case []uint8:
		// 	dst = appendUints8(dst, val)
		case []uint16:
			dst = json.AppendUints16(dst, val)
		case []uint32:
			dst = json.AppendUints32(dst, val)
		case []uint64:
			dst = json.AppendUints64(dst, val)
		case []float32:
			dst = json.AppendFloats32(dst, val)
		case []float64:
			dst = json.AppendFloats64(dst, val)
		case []time.Time:
			dst = json.AppendTimes(dst, val, TimeFieldFormat)
		case []time.Duration:
			dst = json.AppendDurations(dst, val, DurationFieldUnit, DurationFieldInteger)
		case nil:
			dst = append(dst, "null"...)
		default:
			dst = json.AppendInterface(dst, val)
		}
	}
	return dst
}

func appendFieldsBinary(dst []byte, fields map[string]interface{}) []byte {
	keys := make([]string, 0, len(fields))
	for key := range fields {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		dst = cbor.AppendKey(dst, key)
		switch val := fields[key].(type) {
		case string:
			dst = cbor.AppendString(dst, val)
		case []byte:
			dst = cbor.AppendBytes(dst, val)
		case error:
			dst = cbor.AppendError(dst, val)
		case []error:
			dst = cbor.AppendErrors(dst, val)
		case bool:
			dst = cbor.AppendBool(dst, val)
		case int:
			dst = cbor.AppendInt(dst, val)
		case int8:
			dst = cbor.AppendInt8(dst, val)
		case int16:
			dst = cbor.AppendInt16(dst, val)
		case int32:
			dst = cbor.AppendInt32(dst, val)
		case int64:
			dst = cbor.AppendInt64(dst, val)
		case uint:
			dst = cbor.AppendUint(dst, val)
		case uint8:
			dst = cbor.AppendUint8(dst, val)
		case uint16:
			dst = cbor.AppendUint16(dst, val)
		case uint32:
			dst = cbor.AppendUint32(dst, val)
		case uint64:
			dst = cbor.AppendUint64(dst, val)
		case float32:
			dst = cbor.AppendFloat32(dst, val)
		case float64:
			dst = cbor.AppendFloat64(dst, val)
		case time.Time:
			dst = cbor.AppendTime(dst, val, TimeFieldFormat)
		case time.Duration:
			dst = cbor.AppendDuration(dst, val, DurationFieldUnit, DurationFieldInteger)
		case []string:
			dst = cbor.AppendStrings(dst, val)
		case []bool:
			dst = cbor.AppendBools(dst, val)
		case []int:
			dst = cbor.AppendInts(dst, val)
		case []int8:
			dst = cbor.AppendInts8(dst, val)
		case []int16:
			dst = cbor.AppendInts16(dst, val)
		case []int32:
			dst = cbor.AppendInts32(dst, val)
		case []int64:
			dst = cbor.AppendInts64(dst, val)
		case []uint:
			dst = cbor.AppendUints(dst, val)
		// case []uint8:
		// 	dst = appendUints8(dst, val)
		case []uint16:
			dst = cbor.AppendUints16(dst, val)
		case []uint32:
			dst = cbor.AppendUints32(dst, val)
		case []uint64:
			dst = cbor.AppendUints64(dst, val)
		case []float32:
			dst = cbor.AppendFloats32(dst, val)
		case []float64:
			dst = cbor.AppendFloats64(dst, val)
		case []time.Time:
			dst = cbor.AppendTimes(dst, val, TimeFieldFormat)
		case []time.Duration:
			dst = cbor.AppendDurations(dst, val, DurationFieldUnit, DurationFieldInteger)
		case nil:
			dst = append(dst, "null"...)
		default:
			dst = cbor.AppendInterface(dst, val)
		}
	}
	return dst
}
