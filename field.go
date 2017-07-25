package zerolog

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"time"
)

func appendFields(dst []byte, fields map[string]interface{}) []byte {
	keys := make([]string, 0, len(fields))
	for key, _ := range fields {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		switch val := fields[key].(type) {
		case string:
			dst = appendString(dst, key, val)
		case []byte:
			dst = appendBytes(dst, key, val)
		case error:
			dst = appendErrorKey(dst, key, val)
		case []error:
			dst = appendErrorsKey(dst, key, val)
		case bool:
			dst = appendBool(dst, key, val)
		case int:
			dst = appendInt(dst, key, val)
		case int8:
			dst = appendInt8(dst, key, val)
		case int16:
			dst = appendInt16(dst, key, val)
		case int32:
			dst = appendInt32(dst, key, val)
		case int64:
			dst = appendInt64(dst, key, val)
		case uint:
			dst = appendUint(dst, key, val)
		case uint8:
			dst = appendUint8(dst, key, val)
		case uint16:
			dst = appendUint16(dst, key, val)
		case uint32:
			dst = appendUint32(dst, key, val)
		case uint64:
			dst = appendUint64(dst, key, val)
		case float32:
			dst = appendFloat32(dst, key, val)
		case float64:
			dst = appendFloat64(dst, key, val)
		case time.Time:
			dst = appendTime(dst, key, val)
		case time.Duration:
			dst = appendDuration(dst, key, val)
		case []string:
			dst = appendStrings(dst, key, val)
		case []bool:
			dst = appendBools(dst, key, val)
		case []int:
			dst = appendInts(dst, key, val)
		case []int8:
			dst = appendInts8(dst, key, val)
		case []int16:
			dst = appendInts16(dst, key, val)
		case []int32:
			dst = appendInts32(dst, key, val)
		case []int64:
			dst = appendInts64(dst, key, val)
		case []uint:
			dst = appendUints(dst, key, val)
		// case []uint8:
		// 	dst = appendUints8(dst, key, val)
		case []uint16:
			dst = appendUints16(dst, key, val)
		case []uint32:
			dst = appendUints32(dst, key, val)
		case []uint64:
			dst = appendUints64(dst, key, val)
		case []float32:
			dst = appendFloats32(dst, key, val)
		case []float64:
			dst = appendFloats64(dst, key, val)
		case []time.Time:
			dst = appendTimes(dst, key, val)
		case []time.Duration:
			dst = appendDurations(dst, key, val)
		case nil:
			dst = append(appendKey(dst, key), "null"...)
		default:
			dst = appendInterface(dst, key, val)
		}
	}
	return dst
}

func appendKey(dst []byte, key string) []byte {
	if len(dst) > 1 {
		dst = append(dst, ',')
	}
	dst = appendJSONString(dst, key)
	return append(dst, ':')
}

func appendString(dst []byte, key, val string) []byte {
	return appendJSONString(appendKey(dst, key), val)
}

func appendStrings(dst []byte, key string, vals []string) []byte {
	if len(vals) == 0 {
		return append(appendKey(dst, key), '[', ']')
	}
	dst = append(appendKey(dst, key), '[')
	dst = appendJSONString(dst, vals[0])
	if len(vals) > 1 {
		for _, val := range vals[1:] {
			dst = appendJSONString(append(dst, ','), val)
		}
	}
	dst = append(dst, ']')
	return dst
}

func appendBytes(dst []byte, key string, val []byte) []byte {
	return appendJSONBytes(appendKey(dst, key), val)
}

func appendErrorKey(dst []byte, key string, err error) []byte {
	if err == nil {
		return dst
	}
	return appendJSONString(appendKey(dst, key), err.Error())
}

func appendErrorsKey(dst []byte, key string, errs []error) []byte {
	if len(errs) == 0 {
		return append(appendKey(dst, key), '[', ']')
	}
	dst = append(appendKey(dst, key), '[')
	if errs[0] != nil {
		dst = appendJSONString(dst, errs[0].Error())
	} else {
		dst = append(dst, "null"...)
	}
	if len(errs) > 1 {
		for _, err := range errs[1:] {
			if err == nil {
				dst = append(dst, ",null"...)
				continue
			}
			dst = appendJSONString(append(dst, ','), err.Error())
		}
	}
	dst = append(dst, ']')
	return dst
}

func appendError(dst []byte, err error) []byte {
	return appendErrorKey(dst, ErrorFieldName, err)
}

func appendBool(dst []byte, key string, val bool) []byte {
	return strconv.AppendBool(appendKey(dst, key), val)
}

func appendBools(dst []byte, key string, vals []bool) []byte {
	if len(vals) == 0 {
		return append(appendKey(dst, key), '[', ']')
	}
	dst = append(appendKey(dst, key), '[')
	dst = strconv.AppendBool(dst, vals[0])
	if len(vals) > 1 {
		for _, val := range vals[1:] {
			dst = strconv.AppendBool(append(dst, ','), val)
		}
	}
	dst = append(dst, ']')
	return dst
}

func appendInt(dst []byte, key string, val int) []byte {
	return strconv.AppendInt(appendKey(dst, key), int64(val), 10)
}

func appendInts(dst []byte, key string, vals []int) []byte {
	if len(vals) == 0 {
		return append(appendKey(dst, key), '[', ']')
	}
	dst = append(appendKey(dst, key), '[')
	dst = strconv.AppendInt(dst, int64(vals[0]), 10)
	if len(vals) > 1 {
		for _, val := range vals[1:] {
			dst = strconv.AppendInt(append(dst, ','), int64(val), 10)
		}
	}
	dst = append(dst, ']')
	return dst
}

func appendInt8(dst []byte, key string, val int8) []byte {
	return strconv.AppendInt(appendKey(dst, key), int64(val), 10)
}

func appendInts8(dst []byte, key string, vals []int8) []byte {
	if len(vals) == 0 {
		return append(appendKey(dst, key), '[', ']')
	}
	dst = append(appendKey(dst, key), '[')
	dst = strconv.AppendInt(dst, int64(vals[0]), 10)
	if len(vals) > 1 {
		for _, val := range vals[1:] {
			dst = strconv.AppendInt(append(dst, ','), int64(val), 10)
		}
	}
	dst = append(dst, ']')
	return dst
}

func appendInt16(dst []byte, key string, val int16) []byte {
	return strconv.AppendInt(appendKey(dst, key), int64(val), 10)
}

func appendInts16(dst []byte, key string, vals []int16) []byte {
	if len(vals) == 0 {
		return append(appendKey(dst, key), '[', ']')
	}
	dst = append(appendKey(dst, key), '[')
	dst = strconv.AppendInt(dst, int64(vals[0]), 10)
	if len(vals) > 1 {
		for _, val := range vals[1:] {
			dst = strconv.AppendInt(append(dst, ','), int64(val), 10)
		}
	}
	dst = append(dst, ']')
	return dst
}

func appendInt32(dst []byte, key string, val int32) []byte {
	return strconv.AppendInt(appendKey(dst, key), int64(val), 10)
}

func appendInts32(dst []byte, key string, vals []int32) []byte {
	if len(vals) == 0 {
		return append(appendKey(dst, key), '[', ']')
	}
	dst = append(appendKey(dst, key), '[')
	dst = strconv.AppendInt(dst, int64(vals[0]), 10)
	if len(vals) > 1 {
		for _, val := range vals[1:] {
			dst = strconv.AppendInt(append(dst, ','), int64(val), 10)
		}
	}
	dst = append(dst, ']')
	return dst
}

func appendInt64(dst []byte, key string, val int64) []byte {
	return strconv.AppendInt(appendKey(dst, key), val, 10)
}

func appendInts64(dst []byte, key string, vals []int64) []byte {
	if len(vals) == 0 {
		return append(appendKey(dst, key), '[', ']')
	}
	dst = append(appendKey(dst, key), '[')
	dst = strconv.AppendInt(dst, vals[0], 10)
	if len(vals) > 1 {
		for _, val := range vals[1:] {
			dst = strconv.AppendInt(append(dst, ','), val, 10)
		}
	}
	dst = append(dst, ']')
	return dst
}

func appendUint(dst []byte, key string, val uint) []byte {
	return strconv.AppendUint(appendKey(dst, key), uint64(val), 10)
}

func appendUints(dst []byte, key string, vals []uint) []byte {
	if len(vals) == 0 {
		return append(appendKey(dst, key), '[', ']')
	}
	dst = append(appendKey(dst, key), '[')
	dst = strconv.AppendUint(dst, uint64(vals[0]), 10)
	if len(vals) > 1 {
		for _, val := range vals[1:] {
			dst = strconv.AppendUint(append(dst, ','), uint64(val), 10)
		}
	}
	dst = append(dst, ']')
	return dst
}

func appendUint8(dst []byte, key string, val uint8) []byte {
	return strconv.AppendUint(appendKey(dst, key), uint64(val), 10)
}

func appendUints8(dst []byte, key string, vals []uint8) []byte {
	if len(vals) == 0 {
		return append(appendKey(dst, key), '[', ']')
	}
	dst = append(appendKey(dst, key), '[')
	dst = strconv.AppendUint(dst, uint64(vals[0]), 10)
	if len(vals) > 1 {
		for _, val := range vals[1:] {
			dst = strconv.AppendUint(append(dst, ','), uint64(val), 10)
		}
	}
	dst = append(dst, ']')
	return dst
}

func appendUint16(dst []byte, key string, val uint16) []byte {
	return strconv.AppendUint(appendKey(dst, key), uint64(val), 10)
}

func appendUints16(dst []byte, key string, vals []uint16) []byte {
	if len(vals) == 0 {
		return append(appendKey(dst, key), '[', ']')
	}
	dst = append(appendKey(dst, key), '[')
	dst = strconv.AppendUint(dst, uint64(vals[0]), 10)
	if len(vals) > 1 {
		for _, val := range vals[1:] {
			dst = strconv.AppendUint(append(dst, ','), uint64(val), 10)
		}
	}
	dst = append(dst, ']')
	return dst
}

func appendUint32(dst []byte, key string, val uint32) []byte {
	return strconv.AppendUint(appendKey(dst, key), uint64(val), 10)
}

func appendUints32(dst []byte, key string, vals []uint32) []byte {
	if len(vals) == 0 {
		return append(appendKey(dst, key), '[', ']')
	}
	dst = append(appendKey(dst, key), '[')
	dst = strconv.AppendUint(dst, uint64(vals[0]), 10)
	if len(vals) > 1 {
		for _, val := range vals[1:] {
			dst = strconv.AppendUint(append(dst, ','), uint64(val), 10)
		}
	}
	dst = append(dst, ']')
	return dst
}

func appendUint64(dst []byte, key string, val uint64) []byte {
	return strconv.AppendUint(appendKey(dst, key), uint64(val), 10)
}

func appendUints64(dst []byte, key string, vals []uint64) []byte {
	if len(vals) == 0 {
		return append(appendKey(dst, key), '[', ']')
	}
	dst = append(appendKey(dst, key), '[')
	dst = strconv.AppendUint(dst, vals[0], 10)
	if len(vals) > 1 {
		for _, val := range vals[1:] {
			dst = strconv.AppendUint(append(dst, ','), val, 10)
		}
	}
	dst = append(dst, ']')
	return dst
}

func appendFloat32(dst []byte, key string, val float32) []byte {
	return strconv.AppendFloat(appendKey(dst, key), float64(val), 'f', -1, 32)
}

func appendFloats32(dst []byte, key string, vals []float32) []byte {
	if len(vals) == 0 {
		return append(appendKey(dst, key), '[', ']')
	}
	dst = append(appendKey(dst, key), '[')
	dst = strconv.AppendFloat(dst, float64(vals[0]), 'f', -1, 32)
	if len(vals) > 1 {
		for _, val := range vals[1:] {
			dst = strconv.AppendFloat(append(dst, ','), float64(val), 'f', -1, 32)
		}
	}
	dst = append(dst, ']')
	return dst
}

func appendFloat64(dst []byte, key string, val float64) []byte {
	return strconv.AppendFloat(appendKey(dst, key), val, 'f', -1, 32)
}

func appendFloats64(dst []byte, key string, vals []float64) []byte {
	if len(vals) == 0 {
		return append(appendKey(dst, key), '[', ']')
	}
	dst = append(appendKey(dst, key), '[')
	dst = strconv.AppendFloat(dst, vals[0], 'f', -1, 32)
	if len(vals) > 1 {
		for _, val := range vals[1:] {
			dst = strconv.AppendFloat(append(dst, ','), val, 'f', -1, 32)
		}
	}
	dst = append(dst, ']')
	return dst
}

func appendTime(dst []byte, key string, t time.Time) []byte {
	if TimeFieldFormat == "" {
		return appendInt64(dst, key, t.Unix())
	}
	return append(t.AppendFormat(append(appendKey(dst, key), '"'), TimeFieldFormat), '"')
}

func appendTimes(dst []byte, key string, vals []time.Time) []byte {
	if TimeFieldFormat == "" {
		return appendUnixTimes(dst, key, vals)
	}
	if len(vals) == 0 {
		return append(appendKey(dst, key), '[', ']')
	}
	dst = append(appendKey(dst, key), '[')
	dst = append(vals[0].AppendFormat(append(dst, '"'), TimeFieldFormat), '"')
	if len(vals) > 1 {
		for _, t := range vals[1:] {
			dst = append(t.AppendFormat(append(dst, ',', '"'), TimeFieldFormat), '"')
		}
	}
	dst = append(dst, ']')
	return dst
}

func appendUnixTimes(dst []byte, key string, vals []time.Time) []byte {
	if len(vals) == 0 {
		return append(appendKey(dst, key), '[', ']')
	}
	dst = append(appendKey(dst, key), '[')
	dst = strconv.AppendInt(dst, vals[0].Unix(), 10)
	if len(vals) > 1 {
		for _, t := range vals[1:] {
			dst = strconv.AppendInt(dst, t.Unix(), 10)
		}
	}
	dst = append(dst, ']')
	return dst
}

func appendTimestamp(dst []byte) []byte {
	return appendTime(dst, TimestampFieldName, TimestampFunc())
}

func appendDuration(dst []byte, key string, d time.Duration) []byte {
	if DurationFieldInteger {
		return appendInt64(dst, key, int64(d/DurationFieldUnit))
	}
	return appendFloat64(dst, key, float64(d)/float64(DurationFieldUnit))
}

func appendDurations(dst []byte, key string, vals []time.Duration) []byte {
	if DurationFieldInteger {
		return appendIntDurations(dst, key, vals)
	}
	if len(vals) == 0 {
		return append(appendKey(dst, key), '[', ']')
	}
	dst = append(appendKey(dst, key), '[')
	dst = strconv.AppendFloat(dst, float64(vals[0])/float64(DurationFieldUnit), 'f', -1, 32)
	if len(vals) > 1 {
		for _, d := range vals[1:] {
			dst = strconv.AppendFloat(append(dst, ','), float64(d)/float64(DurationFieldUnit), 'f', -1, 32)
		}
	}
	dst = append(dst, ']')
	return dst
}

func appendIntDurations(dst []byte, key string, vals []time.Duration) []byte {
	if len(vals) == 0 {
		return append(appendKey(dst, key), '[', ']')
	}
	dst = append(appendKey(dst, key), '[')
	dst = strconv.AppendInt(dst, int64(vals[0]/DurationFieldUnit), 10)
	if len(vals) > 1 {
		for _, d := range vals[1:] {
			dst = strconv.AppendInt(append(dst, ','), int64(d/DurationFieldUnit), 10)
		}
	}
	dst = append(dst, ']')
	return dst
}

func appendInterface(dst []byte, key string, i interface{}) []byte {
	marshaled, err := json.Marshal(i)
	if err != nil {
		return appendString(dst, key, fmt.Sprintf("marshaling error: %v", err))
	}
	return append(appendKey(dst, key), marshaled...)
}
