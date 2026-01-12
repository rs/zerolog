package zerolog

import (
	"context"
	"encoding/json"
	"io"
	"net"
	"reflect"
	"sort"
	"time"
)

func isNilValue(e error) bool {
	switch reflect.TypeOf(e).Kind() {
	case reflect.Ptr:
		return reflect.ValueOf(e).IsNil()
	default:
		return false
	}
}

func appendFields(dst []byte, fields interface{}, stack bool, ctx context.Context, hooks []Hook) []byte {
	switch fields := fields.(type) {
	case []interface{}:
		if n := len(fields); n&0x1 == 1 { // odd number
			fields = fields[:n-1]
		}
		dst = appendFieldList(dst, fields, stack, ctx, hooks)
	case map[string]interface{}:
		keys := make([]string, 0, len(fields))
		for key := range fields {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		kv := make([]interface{}, 2)
		for _, key := range keys {
			kv[0], kv[1] = key, fields[key]
			dst = appendFieldList(dst, kv, stack, ctx, hooks)
		}
	}
	return dst
}

func appendObject(dst []byte, obj LogObjectMarshaler, stack bool, ctx context.Context, hooks []Hook) []byte {
	e := newEvent(LevelWriterAdapter{io.Discard}, DebugLevel, stack, ctx, hooks)
	e.buf = e.buf[:0] // discard the beginning marker added by newEvent
	e.appendObject(obj)
	dst = append(dst, e.buf...)
	putEvent(e)
	return dst
}

func appendFieldList(dst []byte, kvList []interface{}, stack bool, ctx context.Context, hooks []Hook) []byte {
	for i, n := 0, len(kvList); i < n; i += 2 {
		key, val := kvList[i], kvList[i+1]
		if key, ok := key.(string); ok {
			dst = enc.AppendKey(dst, key)
		} else {
			continue
		}
		switch val := val.(type) {
		case string:
			dst = enc.AppendString(dst, val)
		case []byte:
			dst = enc.AppendBytes(dst, val)
		case error:
			switch m := ErrorMarshalFunc(val).(type) {
			case nil:
				dst = enc.AppendNil(dst)
			case LogObjectMarshaler:
				dst = appendObject(dst, m, stack, ctx, hooks)
			case error:
				if !isNilValue(m) {
					dst = enc.AppendString(dst, m.Error())
				}
			case string:
				dst = enc.AppendString(dst, m)
			default:
				dst = enc.AppendInterface(dst, m)
			}

			if stack && ErrorStackMarshaler != nil {
				switch m := ErrorStackMarshaler(val).(type) {
				case nil:
					return dst // do nothing with nil errors
				case LogObjectMarshaler:
					dst = enc.AppendKey(dst, ErrorStackFieldName)
					dst = appendObject(dst, m, stack, ctx, hooks)
				case error:
					dst = enc.AppendKey(dst, ErrorStackFieldName)
					dst = enc.AppendString(dst, m.Error())
				case string:
					dst = enc.AppendKey(dst, ErrorStackFieldName)
					dst = enc.AppendString(dst, m)
				default:
					dst = enc.AppendKey(dst, ErrorStackFieldName)
					dst = enc.AppendInterface(dst, m)
				}
			}
		case []error:
			dst = enc.AppendArrayStart(dst)
			for i, err := range val {
				switch m := ErrorMarshalFunc(err).(type) {
				case nil:
					dst = enc.AppendNil(dst)
				case LogObjectMarshaler:
					dst = appendObject(dst, m, stack, ctx, hooks)
				case error:
					if !isNilValue(m) {
						dst = enc.AppendString(dst, m.Error())
					}
				case string:
					dst = enc.AppendString(dst, m)
				default:
					dst = enc.AppendInterface(dst, m)
				}

				if i < (len(val) - 1) {
					dst = enc.AppendArrayDelim(dst)
				}
			}
			dst = enc.AppendArrayEnd(dst)
		case []LogObjectMarshaler:
			dst = enc.AppendArrayStart(dst)
			for i, obj := range val {
				dst = appendObject(dst, obj, stack, ctx, hooks)
				if i < (len(val) - 1) {
					dst = enc.AppendArrayDelim(dst)
				}
			}
			dst = enc.AppendArrayEnd(dst)
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
			dst = enc.AppendFloat32(dst, val, FloatingPointPrecision)
		case float64:
			dst = enc.AppendFloat64(dst, val, FloatingPointPrecision)
		case time.Time:
			dst = enc.AppendTime(dst, val, TimeFieldFormat)
		case time.Duration:
			dst = enc.AppendDuration(dst, val, DurationFieldUnit, DurationFieldFormat, DurationFieldInteger, FloatingPointPrecision)
		case *string:
			if val != nil {
				dst = enc.AppendString(dst, *val)
			} else {
				dst = enc.AppendNil(dst)
			}
		case *bool:
			if val != nil {
				dst = enc.AppendBool(dst, *val)
			} else {
				dst = enc.AppendNil(dst)
			}
		case *int:
			if val != nil {
				dst = enc.AppendInt(dst, *val)
			} else {
				dst = enc.AppendNil(dst)
			}
		case *int8:
			if val != nil {
				dst = enc.AppendInt8(dst, *val)
			} else {
				dst = enc.AppendNil(dst)
			}
		case *int16:
			if val != nil {
				dst = enc.AppendInt16(dst, *val)
			} else {
				dst = enc.AppendNil(dst)
			}
		case *int32:
			if val != nil {
				dst = enc.AppendInt32(dst, *val)
			} else {
				dst = enc.AppendNil(dst)
			}
		case *int64:
			if val != nil {
				dst = enc.AppendInt64(dst, *val)
			} else {
				dst = enc.AppendNil(dst)
			}
		case *uint:
			if val != nil {
				dst = enc.AppendUint(dst, *val)
			} else {
				dst = enc.AppendNil(dst)
			}
		case *uint8:
			if val != nil {
				dst = enc.AppendUint8(dst, *val)
			} else {
				dst = enc.AppendNil(dst)
			}
		case *uint16:
			if val != nil {
				dst = enc.AppendUint16(dst, *val)
			} else {
				dst = enc.AppendNil(dst)
			}
		case *uint32:
			if val != nil {
				dst = enc.AppendUint32(dst, *val)
			} else {
				dst = enc.AppendNil(dst)
			}
		case *uint64:
			if val != nil {
				dst = enc.AppendUint64(dst, *val)
			} else {
				dst = enc.AppendNil(dst)
			}
		case *float32:
			if val != nil {
				dst = enc.AppendFloat32(dst, *val, FloatingPointPrecision)
			} else {
				dst = enc.AppendNil(dst)
			}
		case *float64:
			if val != nil {
				dst = enc.AppendFloat64(dst, *val, FloatingPointPrecision)
			} else {
				dst = enc.AppendNil(dst)
			}
		case *time.Time:
			if val != nil {
				dst = enc.AppendTime(dst, *val, TimeFieldFormat)
			} else {
				dst = enc.AppendNil(dst)
			}
		case *time.Duration:
			if val != nil {
				dst = enc.AppendDuration(dst, *val, DurationFieldUnit, DurationFieldFormat, DurationFieldInteger, FloatingPointPrecision)
			} else {
				dst = enc.AppendNil(dst)
			}
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
		// case []uint8: is handled as []byte above
		case []uint16:
			dst = enc.AppendUints16(dst, val)
		case []uint32:
			dst = enc.AppendUints32(dst, val)
		case []uint64:
			dst = enc.AppendUints64(dst, val)
		case []float32:
			dst = enc.AppendFloats32(dst, val, FloatingPointPrecision)
		case []float64:
			dst = enc.AppendFloats64(dst, val, FloatingPointPrecision)
		case []time.Time:
			dst = enc.AppendTimes(dst, val, TimeFieldFormat)
		case []time.Duration:
			dst = enc.AppendDurations(dst, val, DurationFieldUnit, DurationFieldFormat, DurationFieldInteger, FloatingPointPrecision)
		case nil:
			dst = enc.AppendNil(dst)
		case net.IP:
			dst = enc.AppendIPAddr(dst, val)
		case []net.IP:
			dst = enc.AppendIPAddrs(dst, val)
		case net.IPNet:
			dst = enc.AppendIPPrefix(dst, val)
		case []net.IPNet:
			dst = enc.AppendIPPrefixes(dst, val)
		case net.HardwareAddr:
			dst = enc.AppendMACAddr(dst, val)
		case json.RawMessage:
			dst = appendJSON(dst, val)
		default:
			if lom, ok := val.(LogObjectMarshaler); ok {
				dst = appendObject(dst, lom, stack, ctx, hooks)
			} else {
				dst = enc.AppendInterface(dst, val)
			}
		}
	}
	return dst
}
