package json

import (
	"encoding/json"
	"fmt"
	"math"
	"strconv"
)

// AppendBool converts the input bool to a string and
// appends the encoded string to the input byte slice.
func AppendBool(dst []byte, val bool) []byte {
	return strconv.AppendBool(dst, val)
}

// AppendBools encodes the input bools to json and
// appends the encoded string list to the input byte slice.
func AppendBools(dst []byte, vals []bool) []byte {
	if len(vals) == 0 {
		return append(dst, '[', ']')
	}
	dst = append(dst, '[')
	dst = strconv.AppendBool(dst, vals[0])
	if len(vals) > 1 {
		for _, val := range vals[1:] {
			dst = strconv.AppendBool(append(dst, ','), val)
		}
	}
	dst = append(dst, ']')
	return dst
}

// AppendInt converts the input int to a string and
// appends the encoded string to the input byte slice.
func AppendInt(dst []byte, val int) []byte {
	return strconv.AppendInt(dst, int64(val), 10)
}

// AppendInts encodes the input ints to json and
// appends the encoded string list to the input byte slice.
func AppendInts(dst []byte, vals []int) []byte {
	if len(vals) == 0 {
		return append(dst, '[', ']')
	}
	dst = append(dst, '[')
	dst = strconv.AppendInt(dst, int64(vals[0]), 10)
	if len(vals) > 1 {
		for _, val := range vals[1:] {
			dst = strconv.AppendInt(append(dst, ','), int64(val), 10)
		}
	}
	dst = append(dst, ']')
	return dst
}

// AppendInt8 converts the input []int8 to a string and
// appends the encoded string to the input byte slice.
func AppendInt8(dst []byte, val int8) []byte {
	return strconv.AppendInt(dst, int64(val), 10)
}

// AppendInts8 encodes the input int8s to json and
// appends the encoded string list to the input byte slice.
func AppendInts8(dst []byte, vals []int8) []byte {
	if len(vals) == 0 {
		return append(dst, '[', ']')
	}
	dst = append(dst, '[')
	dst = strconv.AppendInt(dst, int64(vals[0]), 10)
	if len(vals) > 1 {
		for _, val := range vals[1:] {
			dst = strconv.AppendInt(append(dst, ','), int64(val), 10)
		}
	}
	dst = append(dst, ']')
	return dst
}

// AppendInt16 converts the input int16 to a string and
// appends the encoded string to the input byte slice.
func AppendInt16(dst []byte, val int16) []byte {
	return strconv.AppendInt(dst, int64(val), 10)
}

// AppendInts16 encodes the input int16s to json and
// appends the encoded string list to the input byte slice.
func AppendInts16(dst []byte, vals []int16) []byte {
	if len(vals) == 0 {
		return append(dst, '[', ']')
	}
	dst = append(dst, '[')
	dst = strconv.AppendInt(dst, int64(vals[0]), 10)
	if len(vals) > 1 {
		for _, val := range vals[1:] {
			dst = strconv.AppendInt(append(dst, ','), int64(val), 10)
		}
	}
	dst = append(dst, ']')
	return dst
}

// AppendInt32 converts the input int32 to a string and
// appends the encoded string to the input byte slice.
func AppendInt32(dst []byte, val int32) []byte {
	return strconv.AppendInt(dst, int64(val), 10)
}

// AppendInts32 encodes the input int32s to json and
// appends the encoded string list to the input byte slice.
func AppendInts32(dst []byte, vals []int32) []byte {
	if len(vals) == 0 {
		return append(dst, '[', ']')
	}
	dst = append(dst, '[')
	dst = strconv.AppendInt(dst, int64(vals[0]), 10)
	if len(vals) > 1 {
		for _, val := range vals[1:] {
			dst = strconv.AppendInt(append(dst, ','), int64(val), 10)
		}
	}
	dst = append(dst, ']')
	return dst
}

// AppendInt64 converts the input int64 to a string and
// appends the encoded string to the input byte slice.
func AppendInt64(dst []byte, val int64) []byte {
	return strconv.AppendInt(dst, val, 10)
}

// AppendInts64 encodes the input int64s to json and
// appends the encoded string list to the input byte slice.
func AppendInts64(dst []byte, vals []int64) []byte {
	if len(vals) == 0 {
		return append(dst, '[', ']')
	}
	dst = append(dst, '[')
	dst = strconv.AppendInt(dst, vals[0], 10)
	if len(vals) > 1 {
		for _, val := range vals[1:] {
			dst = strconv.AppendInt(append(dst, ','), val, 10)
		}
	}
	dst = append(dst, ']')
	return dst
}

// AppendUint converts the input uint to a string and
// appends the encoded string to the input byte slice.
func AppendUint(dst []byte, val uint) []byte {
	return strconv.AppendUint(dst, uint64(val), 10)
}

// AppendUints encodes the input uints to json and
// appends the encoded string list to the input byte slice.
func AppendUints(dst []byte, vals []uint) []byte {
	if len(vals) == 0 {
		return append(dst, '[', ']')
	}
	dst = append(dst, '[')
	dst = strconv.AppendUint(dst, uint64(vals[0]), 10)
	if len(vals) > 1 {
		for _, val := range vals[1:] {
			dst = strconv.AppendUint(append(dst, ','), uint64(val), 10)
		}
	}
	dst = append(dst, ']')
	return dst
}

// AppendUint8 converts the input uint8 to a string and
// appends the encoded string to the input byte slice.
func AppendUint8(dst []byte, val uint8) []byte {
	return strconv.AppendUint(dst, uint64(val), 10)
}

// AppendUints8 encodes the input uint8s to json and
// appends the encoded string list to the input byte slice.
func AppendUints8(dst []byte, vals []uint8) []byte {
	if len(vals) == 0 {
		return append(dst, '[', ']')
	}
	dst = append(dst, '[')
	dst = strconv.AppendUint(dst, uint64(vals[0]), 10)
	if len(vals) > 1 {
		for _, val := range vals[1:] {
			dst = strconv.AppendUint(append(dst, ','), uint64(val), 10)
		}
	}
	dst = append(dst, ']')
	return dst
}

// AppendUint16 converts the input uint16 to a string and
// appends the encoded string to the input byte slice.
func AppendUint16(dst []byte, val uint16) []byte {
	return strconv.AppendUint(dst, uint64(val), 10)
}

// AppendUints16 encodes the input uint16s to json and
// appends the encoded string list to the input byte slice.
func AppendUints16(dst []byte, vals []uint16) []byte {
	if len(vals) == 0 {
		return append(dst, '[', ']')
	}
	dst = append(dst, '[')
	dst = strconv.AppendUint(dst, uint64(vals[0]), 10)
	if len(vals) > 1 {
		for _, val := range vals[1:] {
			dst = strconv.AppendUint(append(dst, ','), uint64(val), 10)
		}
	}
	dst = append(dst, ']')
	return dst
}

// AppendUint32 converts the input uint32 to a string and
// appends the encoded string to the input byte slice.
func AppendUint32(dst []byte, val uint32) []byte {
	return strconv.AppendUint(dst, uint64(val), 10)
}

// AppendUints32 encodes the input uint32s to json and
// appends the encoded string list to the input byte slice.
func AppendUints32(dst []byte, vals []uint32) []byte {
	if len(vals) == 0 {
		return append(dst, '[', ']')
	}
	dst = append(dst, '[')
	dst = strconv.AppendUint(dst, uint64(vals[0]), 10)
	if len(vals) > 1 {
		for _, val := range vals[1:] {
			dst = strconv.AppendUint(append(dst, ','), uint64(val), 10)
		}
	}
	dst = append(dst, ']')
	return dst
}

// AppendUint64 converts the input uint64 to a string and
// appends the encoded string to the input byte slice.
func AppendUint64(dst []byte, val uint64) []byte {
	return strconv.AppendUint(dst, uint64(val), 10)
}

// AppendUints64 encodes the input uint64s to json and
// appends the encoded string list to the input byte slice.
func AppendUints64(dst []byte, vals []uint64) []byte {
	if len(vals) == 0 {
		return append(dst, '[', ']')
	}
	dst = append(dst, '[')
	dst = strconv.AppendUint(dst, vals[0], 10)
	if len(vals) > 1 {
		for _, val := range vals[1:] {
			dst = strconv.AppendUint(append(dst, ','), val, 10)
		}
	}
	dst = append(dst, ']')
	return dst
}

// AppendFloat converts the input float to a string and
// appends the encoded string to the input byte slice.
func AppendFloat(dst []byte, val float64, bitSize int) []byte {
	// JSON does not permit NaN or Infinity. A typical JSON encoder would fail
	// with an error, but a logging library wants the data to get thru so we
	// make a tradeoff and store those types as string.
	switch {
	case math.IsNaN(val):
		return append(dst, `"NaN"`...)
	case math.IsInf(val, 1):
		return append(dst, `"+Inf"`...)
	case math.IsInf(val, -1):
		return append(dst, `"-Inf"`...)
	}
	return strconv.AppendFloat(dst, val, 'f', -1, bitSize)
}

// AppendFloat32 converts the input float32 to a string and
// appends the encoded string to the input byte slice.
func AppendFloat32(dst []byte, val float32) []byte {
	return AppendFloat(dst, float64(val), 32)
}

// AppendFloats32 encodes the input float32s to json and
// appends the encoded string list to the input byte slice.
func AppendFloats32(dst []byte, vals []float32) []byte {
	if len(vals) == 0 {
		return append(dst, '[', ']')
	}
	dst = append(dst, '[')
	dst = AppendFloat(dst, float64(vals[0]), 32)
	if len(vals) > 1 {
		for _, val := range vals[1:] {
			dst = AppendFloat(append(dst, ','), float64(val), 32)
		}
	}
	dst = append(dst, ']')
	return dst
}

// AppendFloat64 converts the input float64 to a string and
// appends the encoded string to the input byte slice.
func AppendFloat64(dst []byte, val float64) []byte {
	return AppendFloat(dst, val, 64)
}

// AppendFloats64 encodes the input float64s to json and
// appends the encoded string list to the input byte slice.
func AppendFloats64(dst []byte, vals []float64) []byte {
	if len(vals) == 0 {
		return append(dst, '[', ']')
	}
	dst = append(dst, '[')
	dst = AppendFloat(dst, vals[0], 32)
	if len(vals) > 1 {
		for _, val := range vals[1:] {
			dst = AppendFloat(append(dst, ','), val, 64)
		}
	}
	dst = append(dst, ']')
	return dst
}

// AppendInterface marshals the input interface to a string and
// appends the encoded string to the input byte slice.
func AppendInterface(dst []byte, i interface{}) []byte {
	marshaled, err := json.Marshal(i)
	if err != nil {
		return AppendString(dst, fmt.Sprintf("marshaling error: %v", err))
	}
	return append(dst, marshaled...)
}

func AppendObjectData(dst []byte, o []byte) []byte {
	// Two conditions we want to put a ',' between existing content and
	// new content:
	// 1. new content starts with '{' - which shd be dropped   OR
	// 2. existing content has already other fields
	if o[0] == '{' {
		o[0] = ','
	} else if len(dst) > 1 {
		dst = append(dst, ',')
	}
	return append(dst, o...)
}
