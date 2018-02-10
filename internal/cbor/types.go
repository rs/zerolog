package cbor

import (
	"encoding/json"
	"fmt"
	"math"
)

func AppendNull(dst []byte) []byte {
	return append(dst, byte(majorTypeSimpleAndFloat|additionalTypeNull))
}

func AppendBeginMarker(dst []byte) []byte {
	return append(dst, byte(majorTypeMap|additionalTypeInfiniteCount))
}

func AppendEndMarker(dst []byte, terminal bool) []byte {
	return append(dst, byte(majorTypeSimpleAndFloat|additionalTypeBreak))
}

func AppendBool(dst []byte, val bool) []byte {
	b := additionalTypeBoolFalse
	if val {
		b = additionalTypeBoolTrue
	}
	return append(dst, byte(majorTypeSimpleAndFloat|b))
}

func AppendBools(dst []byte, vals []bool) []byte {
	major := majorTypeArray
	l := len(vals)
	if l == 0 {
		return AppendArrayEnd(AppendArrayStart(dst))
	}
	if l <= additionalMax {
		lb := byte(l)
		dst = append(dst, byte(major|lb))
	} else {
		dst = appendCborTypePrefix(dst, major, uint64(l))
	}
	for _, v := range vals {
		dst = AppendBool(dst, v)
	}
	return dst
}

func AppendInt(dst []byte, val int) []byte {
	major := majorTypeUnsignedInt
	contentVal := val
	if val < 0 {
		major = majorTypeNegativeInt
		contentVal = -val - 1
	}
	if contentVal <= additionalMax {
		lb := byte(contentVal)
		dst = append(dst, byte(major|lb))
	} else {
		dst = appendCborTypePrefix(dst, major, uint64(contentVal))
	}
	return dst
}

func AppendInts(dst []byte, vals []int) []byte {
	major := majorTypeArray
	l := len(vals)
	if l == 0 {
		return AppendArrayEnd(AppendArrayStart(dst))
	}
	if l <= additionalMax {
		lb := byte(l)
		dst = append(dst, byte(major|lb))
	} else {
		dst = appendCborTypePrefix(dst, major, uint64(l))
	}
	for _, v := range vals {
		dst = AppendInt(dst, v)
	}
	return dst
}

func AppendInt8(dst []byte, val int8) []byte {
	return AppendInt(dst, int(val))
}

func AppendInts8(dst []byte, vals []int8) []byte {
	major := majorTypeArray
	l := len(vals)
	if l == 0 {
		return AppendArrayEnd(AppendArrayStart(dst))
	}
	if l <= additionalMax {
		lb := byte(l)
		dst = append(dst, byte(major|lb))
	} else {
		dst = appendCborTypePrefix(dst, major, uint64(l))
	}
	for _, v := range vals {
		dst = AppendInt(dst, int(v))
	}
	return dst
}

func AppendInt16(dst []byte, val int16) []byte {
	return AppendInt(dst, int(val))
}

func AppendInts16(dst []byte, vals []int16) []byte {
	major := majorTypeArray
	l := len(vals)
	if l == 0 {
		return AppendArrayEnd(AppendArrayStart(dst))
	}
	if l <= additionalMax {
		lb := byte(l)
		dst = append(dst, byte(major|lb))
	} else {
		dst = appendCborTypePrefix(dst, major, uint64(l))
	}
	for _, v := range vals {
		dst = AppendInt(dst, int(v))
	}
	return dst
}

func AppendInt32(dst []byte, val int32) []byte {
	return AppendInt(dst, int(val))
}

func AppendInts32(dst []byte, vals []int32) []byte {
	major := majorTypeArray
	l := len(vals)
	if l == 0 {
		return AppendArrayEnd(AppendArrayStart(dst))
	}
	if l <= additionalMax {
		lb := byte(l)
		dst = append(dst, byte(major|lb))
	} else {
		dst = appendCborTypePrefix(dst, major, uint64(l))
	}
	for _, v := range vals {
		dst = AppendInt(dst, int(v))
	}
	return dst
}

func AppendInt64(dst []byte, val int64) []byte {
	major := majorTypeUnsignedInt
	contentVal := val
	if val < 0 {
		major = majorTypeNegativeInt
		contentVal = -val - 1
	}
	if contentVal <= additionalMax {
		lb := byte(contentVal)
		dst = append(dst, byte(major|lb))
	} else {
		dst = appendCborTypePrefix(dst, major, uint64(contentVal))
	}
	return dst
}

func AppendInts64(dst []byte, vals []int64) []byte {
	major := majorTypeArray
	l := len(vals)
	if l == 0 {
		return AppendArrayEnd(AppendArrayStart(dst))
	}
	if l <= additionalMax {
		lb := byte(l)
		dst = append(dst, byte(major|lb))
	} else {
		dst = appendCborTypePrefix(dst, major, uint64(l))
	}
	for _, v := range vals {
		dst = AppendInt64(dst, v)
	}
	return dst
}

func AppendUint(dst []byte, val uint) []byte {
	return AppendInt64(dst, int64(val))
}

func AppendUints(dst []byte, vals []uint) []byte {
	major := majorTypeArray
	l := len(vals)
	if l == 0 {
		return AppendArrayEnd(AppendArrayStart(dst))
	}
	if l <= additionalMax {
		lb := byte(l)
		dst = append(dst, byte(major|lb))
	} else {
		dst = appendCborTypePrefix(dst, major, uint64(l))
	}
	for _, v := range vals {
		dst = AppendUint(dst, v)
	}
	return dst
}

func AppendUint8(dst []byte, val uint8) []byte {
	return AppendUint(dst, uint(val))
}

func AppendUints8(dst []byte, vals []uint8) []byte {
	major := majorTypeArray
	l := len(vals)
	if l == 0 {
		return AppendArrayEnd(AppendArrayStart(dst))
	}
	if l <= additionalMax {
		lb := byte(l)
		dst = append(dst, byte(major|lb))
	} else {
		dst = appendCborTypePrefix(dst, major, uint64(l))
	}
	for _, v := range vals {
		dst = AppendUint8(dst, v)
	}
	return dst
}

func AppendUint16(dst []byte, val uint16) []byte {
	return AppendUint(dst, uint(val))
}

func AppendUints16(dst []byte, vals []uint16) []byte {
	major := majorTypeArray
	l := len(vals)
	if l == 0 {
		return AppendArrayEnd(AppendArrayStart(dst))
	}
	if l <= additionalMax {
		lb := byte(l)
		dst = append(dst, byte(major|lb))
	} else {
		dst = appendCborTypePrefix(dst, major, uint64(l))
	}
	for _, v := range vals {
		dst = AppendUint16(dst, v)
	}
	return dst
}

func AppendUint32(dst []byte, val uint32) []byte {
	return AppendUint(dst, uint(val))
}

func AppendUints32(dst []byte, vals []uint32) []byte {
	major := majorTypeArray
	l := len(vals)
	if l == 0 {
		return AppendArrayEnd(AppendArrayStart(dst))
	}
	if l <= additionalMax {
		lb := byte(l)
		dst = append(dst, byte(major|lb))
	} else {
		dst = appendCborTypePrefix(dst, major, uint64(l))
	}
	for _, v := range vals {
		dst = AppendUint32(dst, v)
	}
	return dst
}

func AppendUint64(dst []byte, val uint64) []byte {
	major := majorTypeUnsignedInt
	contentVal := val
	if contentVal <= additionalMax {
		lb := byte(contentVal)
		dst = append(dst, byte(major|lb))
	} else {
		dst = appendCborTypePrefix(dst, major, uint64(contentVal))
	}
	return dst
}

func AppendUints64(dst []byte, vals []uint64) []byte {
	major := majorTypeArray
	l := len(vals)
	if l == 0 {
		return AppendArrayEnd(AppendArrayStart(dst))
	}
	if l <= additionalMax {
		lb := byte(l)
		dst = append(dst, byte(major|lb))
	} else {
		dst = appendCborTypePrefix(dst, major, uint64(l))
	}
	for _, v := range vals {
		dst = AppendUint64(dst, v)
	}
	return dst
}

func AppendFloat32(dst []byte, val float32) []byte {
	switch {
	case math.IsNaN(float64(val)):
		return append(dst, "\xfa\x7f\xc0\x00\x00"...)
	case math.IsInf(float64(val), 1):
		return append(dst, "\xfa\x7f\x80\x00\x00"...)
	case math.IsInf(float64(val), -1):
		return append(dst, "\xfa\xff\x80\x00\x00"...)
	}
	major := majorTypeSimpleAndFloat
	subType := additionalTypeFloat32
	n := math.Float32bits(val)
	var buf [4]byte
	for i := uint(0); i < 4; i++ {
		buf[i] = byte(n >> ((3 - i) * 8))
	}
	return append(append(dst, byte(major|subType)), buf[0], buf[1], buf[2], buf[3])
}

func AppendFloats32(dst []byte, vals []float32) []byte {
	major := majorTypeArray
	l := len(vals)
	if l == 0 {
		return AppendArrayEnd(AppendArrayStart(dst))
	}
	if l <= additionalMax {
		lb := byte(l)
		dst = append(dst, byte(major|lb))
	} else {
		dst = appendCborTypePrefix(dst, major, uint64(l))
	}
	for _, v := range vals {
		dst = AppendFloat32(dst, v)
	}
	return dst
}

func AppendFloat64(dst []byte, val float64) []byte {
	switch {
	case math.IsNaN(val):
		return append(dst, "\xfb\x7f\xf8\x00\x00\x00\x00\x00\x00"...)
	case math.IsInf(val, 1):
		return append(dst, "\xfb\x7f\xf0\x00\x00\x00\x00\x00\x00"...)
	case math.IsInf(val, -1):
		return append(dst, "\xfb\xff\xf0\x00\x00\x00\x00\x00\x00"...)
	}
	major := majorTypeSimpleAndFloat
	subType := additionalTypeFloat64
	n := math.Float64bits(val)
	dst = append(dst, byte(major|subType))
	for i := uint(1); i <= 8; i++ {
		b := byte(n >> ((8 - i) * 8))
		dst = append(dst, b)
	}
	return dst
}

func AppendFloats64(dst []byte, vals []float64) []byte {
	major := majorTypeArray
	l := len(vals)
	if l == 0 {
		return AppendArrayEnd(AppendArrayStart(dst))
	}
	if l <= additionalMax {
		lb := byte(l)
		dst = append(dst, byte(major|lb))
	} else {
		dst = appendCborTypePrefix(dst, major, uint64(l))
	}
	for _, v := range vals {
		dst = AppendFloat64(dst, v)
	}
	return dst
}

func AppendInterface(dst []byte, i interface{}) []byte {
	//TODO - gotto use reflect to find out the contents and print accordingly
	//For now we'll use JSON to reflect to a string - until we build
	//a CBOR based reflection
	marshaled, err := json.Marshal(i)
	if err != nil {
		return AppendString(dst, fmt.Sprintf("marshaling error: %v", err))
	}
	return AppendBytes(dst, marshaled)
}

func AppendObjectData(dst []byte, o []byte) []byte {
	//TODO - check is this sufficient or do we need to something more..
	return append(dst, o...)
}

func AppendArrayStart(dst []byte) []byte {
	return append(dst, byte(majorTypeArray|additionalTypeInfiniteCount))
}

func AppendArrayEnd(dst []byte) []byte {
	return append(dst, byte(majorTypeSimpleAndFloat|additionalTypeBreak))
}

func AppendArrayDelim(dst []byte) []byte {
	//No delimiters needed in cbor
	return dst
}
