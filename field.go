package zerolog

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

var now = time.Now

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

func appendError(dst []byte, err error) []byte {
	return appendJSONString(appendKey(dst, ErrorFieldName), err.Error())
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
	return append(t.AppendFormat(append(appendKey(dst, key), '"'), TimeFieldFormat), '"')
}

func appendTimestamp(dst []byte) []byte {
	return appendInt64(dst, TimestampFieldName, now().Unix())
}

func appendObject(dst []byte, key string, obj interface{}) []byte {
	marshaled, err := json.Marshal(obj)
	if err != nil {
		return appendString(dst, key, fmt.Sprintf("marshaling error: %v", err))
	}
	return append(appendKey(dst, key), marshaled...)
}
