package zerolog

import (
	"bytes"
	"strconv"
	"time"
)

var now = time.Now

type fieldMode uint8

const (
	zeroFieldMode fieldMode = iota
	rawFieldMode
	quotedFieldMode
	precomputedFieldMode
	timestampFieldMode
)

// field define a logger field.
type field struct {
	key  string
	mode fieldMode
	val  string
	json []byte
}

func (f field) writeJSON(buf *bytes.Buffer) {
	switch f.mode {
	case zeroFieldMode:
		return
	case precomputedFieldMode:
		buf.Write(f.json)
		return
	case timestampFieldMode:
		writeJSONString(buf, TimestampFieldName)
		buf.WriteByte(':')
		buf.WriteString(strconv.FormatInt(now().Unix(), 10))
		return
	}
	writeJSONString(buf, f.key)
	buf.WriteByte(':')
	switch f.mode {
	case quotedFieldMode:
		writeJSONString(buf, f.val)
	case rawFieldMode:
		buf.WriteString(f.val)
	default:
		panic("unknown field mode")
	}
}

func (f field) compileJSON() field {
	switch f.mode {
	case zeroFieldMode, precomputedFieldMode, timestampFieldMode:
		return f
	}
	buf := &bytes.Buffer{}
	f.writeJSON(buf)
	cf := field{
		mode: precomputedFieldMode,
		json: buf.Bytes(),
	}
	return cf
}

func fStr(key, val string) field {
	return field{key, quotedFieldMode, val, nil}
}

func fErr(err error) field {
	return field{ErrorFieldName, quotedFieldMode, err.Error(), nil}
}

func fBool(key string, b bool) field {
	if b {
		return field{key, rawFieldMode, "true", nil}
	}
	return field{key, rawFieldMode, "false", nil}
}

func fInt(key string, i int) field {
	return field{key, rawFieldMode, strconv.FormatInt(int64(i), 10), nil}
}

func fInt8(key string, i int8) field {
	return field{key, rawFieldMode, strconv.FormatInt(int64(i), 10), nil}
}

func fInt16(key string, i int16) field {
	return field{key, rawFieldMode, strconv.FormatInt(int64(i), 10), nil}
}

func fInt32(key string, i int32) field {
	return field{key, rawFieldMode, strconv.FormatInt(int64(i), 10), nil}
}

func fInt64(key string, i int64) field {
	return field{key, rawFieldMode, strconv.FormatInt(i, 10), nil}
}

func fUint(key string, i uint) field {
	return field{key, rawFieldMode, strconv.FormatUint(uint64(i), 10), nil}
}

func fUint8(key string, i uint8) field {
	return field{key, rawFieldMode, strconv.FormatUint(uint64(i), 10), nil}
}

func fUint16(key string, i uint16) field {
	return field{key, rawFieldMode, strconv.FormatUint(uint64(i), 10), nil}
}

func fUint32(key string, i uint32) field {
	return field{key, rawFieldMode, strconv.FormatUint(uint64(i), 10), nil}
}

func fUint64(key string, i uint64) field {
	return field{key, rawFieldMode, strconv.FormatUint(i, 10), nil}
}

func fFloat32(key string, f float32) field {
	return field{key, rawFieldMode, strconv.FormatFloat(float64(f), 'f', -1, 32), nil}
}

func fFloat64(key string, f float64) field {
	return field{key, rawFieldMode, strconv.FormatFloat(f, 'f', -1, 64), nil}
}

func fTimestamp() field {
	return field{mode: timestampFieldMode}
}

func fTime(key string, t time.Time) field {
	return field{key, quotedFieldMode, t.Format(TimeFieldFormat), nil}
}

func fRaw(key string, raw string) field {
	return field{key, rawFieldMode, raw, nil}
}
