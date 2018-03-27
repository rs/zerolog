package zerolog

import (
	"sort"
	"time"
)

func appendFields(dst []byte, fields map[string]interface{}) []byte {
	keys := make([]string, 0, len(fields))
	for key := range fields {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		dst = appendKey(dst, key)
		switch val := fields[key].(type) {
		case string:
			dst = appendString(dst, val)
		case []byte:
			dst = appendBytes(dst, val)
		case error:
			dst = appendError(dst, val)
		case []error:
			dst = appendErrors(dst, val)
		case bool:
			dst = appendBool(dst, val)
		case int:
			dst = appendInt(dst, val)
		case int8:
			dst = appendInt8(dst, val)
		case int16:
			dst = appendInt16(dst, val)
		case int32:
			dst = appendInt32(dst, val)
		case int64:
			dst = appendInt64(dst, val)
		case uint:
			dst = appendUint(dst, val)
		case uint8:
			dst = appendUint8(dst, val)
		case uint16:
			dst = appendUint16(dst, val)
		case uint32:
			dst = appendUint32(dst, val)
		case uint64:
			dst = appendUint64(dst, val)
		case float32:
			dst = appendFloat32(dst, val)
		case float64:
			dst = appendFloat64(dst, val)
		case time.Time:
			dst = appendTime(dst, val, TimeFieldFormat)
		case time.Duration:
			dst = appendDuration(dst, val, DurationFieldUnit, DurationFieldInteger)
		case *string:
			dst = appendString(dst, *val)
		case *bool:
			dst = appendBool(dst, *val)
		case *int:
			dst = appendInt(dst, *val)
		case *int8:
			dst = appendInt8(dst, *val)
		case *int16:
			dst = appendInt16(dst, *val)
		case *int32:
			dst = appendInt32(dst, *val)
		case *int64:
			dst = appendInt64(dst, *val)
		case *uint:
			dst = appendUint(dst, *val)
		case *uint8:
			dst = appendUint8(dst, *val)
		case *uint16:
			dst = appendUint16(dst, *val)
		case *uint32:
			dst = appendUint32(dst, *val)
		case *uint64:
			dst = appendUint64(dst, *val)
		case *float32:
			dst = appendFloat32(dst, *val)
		case *float64:
			dst = appendFloat64(dst, *val)
		case *time.Time:
			dst = appendTime(dst, *val, TimeFieldFormat)
		case *time.Duration:
			dst = appendDuration(dst, *val, DurationFieldUnit, DurationFieldInteger)
		case []string:
			dst = appendStrings(dst, val)
		case []bool:
			dst = appendBools(dst, val)
		case []int:
			dst = appendInts(dst, val)
		case []int8:
			dst = appendInts8(dst, val)
		case []int16:
			dst = appendInts16(dst, val)
		case []int32:
			dst = appendInts32(dst, val)
		case []int64:
			dst = appendInts64(dst, val)
		case []uint:
			dst = appendUints(dst, val)
		// case []uint8:
		// 	dst = appendUints8(dst, val)
		case []uint16:
			dst = appendUints16(dst, val)
		case []uint32:
			dst = appendUints32(dst, val)
		case []uint64:
			dst = appendUints64(dst, val)
		case []float32:
			dst = appendFloats32(dst, val)
		case []float64:
			dst = appendFloats64(dst, val)
		case []time.Time:
			dst = appendTimes(dst, val, TimeFieldFormat)
		case []time.Duration:
			dst = appendDurations(dst, val, DurationFieldUnit, DurationFieldInteger)
		case nil:
			dst = appendNil(dst)
		default:
			dst = appendInterface(dst, val)
		}
	}
	return dst
}
