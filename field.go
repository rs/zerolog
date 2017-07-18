package zerolog

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"time"

	pkgErrors "github.com/pkg/errors"
)

func appendFields(dst []byte, fields map[string]interface{}) []byte {
	keys := make([]string, 0, len(fields))
	for key := range fields {
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

func appendBytes(dst []byte, key string, val []byte) []byte {
	return appendJSONBytes(appendKey(dst, key), val)
}

func appendErrorKey(dst []byte, key string, err error) []byte {
	if err == nil {
		return dst
	}
	return appendJSONString(appendKey(dst, key), err.Error())
}

func appendJSONStack(dst []byte, stackTrace pkgErrors.StackTrace) []byte {
	// Even though this is a "slow path" operation, make this as performant as
	// possible given the interface defined by github.com/pkg/errors.
	buf := &bytes.Buffer{}
	buf.Grow(500)
	dst = append(dst, '[')
	for i, frame := range stackTrace {
		dst = append(dst, '{')

		dst = appendJSONString(dst, StackSourceFileName)
		dst = append(dst, ':')
		fmt.Fprintf(buf, "%s", frame)
		dst = appendJSONBytes(dst, buf.Bytes())
		buf.Reset()

		dst = appendKey(dst, StackSourceLineName)
		fmt.Fprintf(buf, "%d", frame)
		dst = appendJSONBytes(dst, buf.Bytes())
		buf.Reset()

		dst = appendKey(dst, StackSourceFunctionName)
		fmt.Fprintf(buf, "%n", frame)
		dst = appendJSONBytes(dst, buf.Bytes())
		buf.Reset()

		dst = append(dst, '}')
		if i < len(stackTrace)-1 {
			dst = append(dst, ',')
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

func appendInt(dst []byte, key string, val int) []byte {
	return strconv.AppendInt(appendKey(dst, key), int64(val), 10)
}

func appendInt8(dst []byte, key string, val int8) []byte {
	return strconv.AppendInt(appendKey(dst, key), int64(val), 10)
}

func appendInt16(dst []byte, key string, val int16) []byte {
	return strconv.AppendInt(appendKey(dst, key), int64(val), 10)
}

func appendInt32(dst []byte, key string, val int32) []byte {
	return strconv.AppendInt(appendKey(dst, key), int64(val), 10)
}

func appendInt64(dst []byte, key string, val int64) []byte {
	return strconv.AppendInt(appendKey(dst, key), int64(val), 10)
}

func appendUint(dst []byte, key string, val uint) []byte {
	return strconv.AppendUint(appendKey(dst, key), uint64(val), 10)
}

func appendUint8(dst []byte, key string, val uint8) []byte {
	return strconv.AppendUint(appendKey(dst, key), uint64(val), 10)
}

func appendUint16(dst []byte, key string, val uint16) []byte {
	return strconv.AppendUint(appendKey(dst, key), uint64(val), 10)
}

func appendUint32(dst []byte, key string, val uint32) []byte {
	return strconv.AppendUint(appendKey(dst, key), uint64(val), 10)
}

func appendUint64(dst []byte, key string, val uint64) []byte {
	return strconv.AppendUint(appendKey(dst, key), uint64(val), 10)
}

func appendFloat32(dst []byte, key string, val float32) []byte {
	return strconv.AppendFloat(appendKey(dst, key), float64(val), 'f', -1, 32)
}

func appendFloat64(dst []byte, key string, val float64) []byte {
	return strconv.AppendFloat(appendKey(dst, key), float64(val), 'f', -1, 32)
}

func appendTime(dst []byte, key string, t time.Time) []byte {
	if TimeFieldFormat == "" {
		return appendInt64(dst, key, t.Unix())
	}
	return append(t.AppendFormat(append(appendKey(dst, key), '"'), TimeFieldFormat), '"')
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

func appendInterface(dst []byte, key string, i interface{}) []byte {
	marshaled, err := json.Marshal(i)
	if err != nil {
		return appendString(dst, key, fmt.Sprintf("marshaling error: %v", err))
	}
	return append(appendKey(dst, key), marshaled...)
}
