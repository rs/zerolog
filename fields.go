package zerolog

import (
	"net"
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
		dst = enc.AppendKey(dst, key)
		val := fields[key]
		if val, ok := val.(LogObjectMarshaler); ok {
			e := newEvent(nil, 0)
			e.buf = e.buf[:0]
			e.appendObject(val)
			dst = append(dst, e.buf...)
			eventPool.Put(e)
			continue
		}
		switch val := val.(type) {
		case string:
			dst = enc.AppendString(dst, val)
		case []byte:
			dst = enc.AppendBytes(dst, val)
		case error:
			dst = enc.AppendError(dst, val)
		case []error:
			dst = enc.AppendErrors(dst, val)
		case bool:
			dst = enc.AppendBool(dst, val)
		case int:
			dst = enc.AppendInt(dst, val)
		case int8:
			dst = enc.AppendInt8(dst, val)
		case int16:
			dst = enc.AppendInt16(dst, val)
		case int32:
			dst = enc.AppendInt32(dst, val)
		case int64:
			dst = enc.AppendInt64(dst, val)
		case uint:
			dst = enc.AppendUint(dst, val)
		case uint8:
			dst = enc.AppendUint8(dst, val)
		case uint16:
			dst = enc.AppendUint16(dst, val)
		case uint32:
			dst = enc.AppendUint32(dst, val)
		case uint64:
			dst = enc.AppendUint64(dst, val)
		case float32:
			dst = enc.AppendFloat32(dst, val)
		case float64:
			dst = enc.AppendFloat64(dst, val)
		case time.Time:
			dst = enc.AppendTime(dst, val, TimeFieldFormat)
		case time.Duration:
			dst = enc.AppendDuration(dst, val, DurationFieldUnit, DurationFieldInteger)
		case *string:
			dst = enc.AppendString(dst, *val)
		case *bool:
			dst = enc.AppendBool(dst, *val)
		case *int:
			dst = enc.AppendInt(dst, *val)
		case *int8:
			dst = enc.AppendInt8(dst, *val)
		case *int16:
			dst = enc.AppendInt16(dst, *val)
		case *int32:
			dst = enc.AppendInt32(dst, *val)
		case *int64:
			dst = enc.AppendInt64(dst, *val)
		case *uint:
			dst = enc.AppendUint(dst, *val)
		case *uint8:
			dst = enc.AppendUint8(dst, *val)
		case *uint16:
			dst = enc.AppendUint16(dst, *val)
		case *uint32:
			dst = enc.AppendUint32(dst, *val)
		case *uint64:
			dst = enc.AppendUint64(dst, *val)
		case *float32:
			dst = enc.AppendFloat32(dst, *val)
		case *float64:
			dst = enc.AppendFloat64(dst, *val)
		case *time.Time:
			dst = enc.AppendTime(dst, *val, TimeFieldFormat)
		case *time.Duration:
			dst = enc.AppendDuration(dst, *val, DurationFieldUnit, DurationFieldInteger)
		case []string:
			dst = enc.AppendStrings(dst, val)
		case []bool:
			dst = enc.AppendBools(dst, val)
		case []int:
			dst = enc.AppendInts(dst, val)
		case []int8:
			dst = enc.AppendInts8(dst, val)
		case []int16:
			dst = enc.AppendInts16(dst, val)
		case []int32:
			dst = enc.AppendInts32(dst, val)
		case []int64:
			dst = enc.AppendInts64(dst, val)
		case []uint:
			dst = enc.AppendUints(dst, val)
		// case []uint8:
		// 	dst = enc.AppendUints8(dst, val)
		case []uint16:
			dst = enc.AppendUints16(dst, val)
		case []uint32:
			dst = enc.AppendUints32(dst, val)
		case []uint64:
			dst = enc.AppendUints64(dst, val)
		case []float32:
			dst = enc.AppendFloats32(dst, val)
		case []float64:
			dst = enc.AppendFloats64(dst, val)
		case []time.Time:
			dst = enc.AppendTimes(dst, val, TimeFieldFormat)
		case []time.Duration:
			dst = enc.AppendDurations(dst, val, DurationFieldUnit, DurationFieldInteger)
		case nil:
			dst = enc.AppendNil(dst)
		case net.IP:
			dst = enc.AppendIPAddr(dst, val)
		case net.IPNet:
			dst = enc.AppendIPPrefix(dst, val)
		case net.HardwareAddr:
			dst = enc.AppendMACAddr(dst, val)
		default:
			dst = enc.AppendInterface(dst, val)
		}
	}
	return dst
}
