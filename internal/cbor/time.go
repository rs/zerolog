package cbor

import (
	"time"
)

func appendIntegerTimestamp(dst []byte, t time.Time) []byte {
	major := majorTypeTags
	minor := additionalTypeTimestamp
	dst = append(dst, byte(major|minor))
	secs := t.Unix()
	var val uint64
	if secs < 0 {
		major = majorTypeNegativeInt
		val = uint64(-secs - 1)
	} else {
		major = majorTypeUnsignedInt
		val = uint64(secs)
	}
	dst = appendCborTypePrefix(dst, major, uint64(val))
	return dst
}

func appendFloatTimestamp(dst []byte, t time.Time) []byte {
	major := majorTypeTags
	minor := additionalTypeTimestamp
	dst = append(dst, byte(major|minor))
	secs := t.Unix()
	nanos := t.Nanosecond()
	var val float64
	val = float64(secs)*1.0 + float64(nanos)*1E-9
	return AppendFloat64(dst, val)
}

// AppendTime encodes and adds a timestamp to the dst byte array.
func AppendTime(dst []byte, t time.Time, unused string) []byte {
	utc := t.UTC()
	if utc.Nanosecond() == 0 {
		return appendIntegerTimestamp(dst, utc)
	}
	return appendFloatTimestamp(dst, utc)
}

// AppendTimes encodes and adds an array of timestamps to the dst byte array.
func AppendTimes(dst []byte, vals []time.Time, unused string) []byte {
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

	for _, t := range vals {
		dst = AppendTime(dst, t, unused)
	}
	return dst
}

// AppendDuration encodes and adds a duration to the dst byte array.
// useInt field indicates whether to store the duration as seconds (integer) or
// as seconds+nanoseconds (float).
func AppendDuration(dst []byte, d time.Duration, unit time.Duration, useInt bool) []byte {
	if useInt {
		return AppendInt64(dst, int64(d/unit))
	}
	return AppendFloat64(dst, float64(d)/float64(unit))
}

// AppendDurations encodes and adds an array of durations to the dst byte array.
// useInt field indicates whether to store the duration as seconds (integer) or
// as seconds+nanoseconds (float).
func AppendDurations(dst []byte, vals []time.Duration, unit time.Duration, useInt bool) []byte {
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
	for _, d := range vals {
		dst = AppendDuration(dst, d, unit, useInt)
	}
	return dst
}
