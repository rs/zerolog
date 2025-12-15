package json

import (
	"bytes"
	"fmt"
	"math"
	"testing"

	"github.com/rs/zerolog/internal"
)

func TestAppendInts8(t *testing.T) {
	doOne := func(vals []int8) {
		want := make([]byte, 0)
		want = append(want, '[')
		for i, val := range vals {
			want = append(want, []byte(fmt.Sprintf("%d", int8(val)))...)
			if i < len(vals)-1 {
				want = append(want, ',')
			}
		}
		want = append(want, ']')

		got := enc.AppendInts8([]byte{}, vals)
		if !bytes.Equal(got, want) {
			t.Errorf("AppendInts8(%v)\ngot:  %s\nwant: %s",
				vals,
				string(got),
				string(want))
		}
	}

	array := make([]int8, 0)
	for _, tc := range internal.IntegerTestCases {
		if (tc.Val < math.MinInt8) || (tc.Val > math.MaxInt8) {
			continue
		}
		array = append(array, int8(tc.Val))
	}

	doOne(array)
	doOne(array[:1]) // single element
	doOne(array[:0]) // edge case of zero length
}
func TestAppendUints8(t *testing.T) {
	doOne := func(vals []uint8) {
		want := make([]byte, 0)
		want = append(want, '[')
		for i, val := range vals {
			want = append(want, []byte(fmt.Sprintf("%v", uint8(val)))...)
			if i < len(vals)-1 {
				want = append(want, ',')
			}
		}
		want = append(want, ']')

		got := enc.AppendUints8([]byte{}, vals)
		if !bytes.Equal(got, want) {
			t.Errorf("AppendUints8(%v)\ngot:  %s\nwant: %s",
				vals,
				string(got),
				string(want))
		}
	}

	array := make([]uint8, 0)
	for _, tc := range internal.UnsignedIntegerTestCases {
		if tc.Val > math.MaxUint8 {
			continue
		}
		array = append(array, uint8(tc.Val))
	}

	doOne(array)
	doOne(array[:1]) // single element
	doOne(array[:0]) // edge case of zero length
}
func TestAppendInts16(t *testing.T) {
	doOne := func(vals []int16) {
		want := make([]byte, 0)
		want = append(want, '[')
		for i, val := range vals {
			want = append(want, []byte(fmt.Sprintf("%d", int16(val)))...)
			if i < len(vals)-1 {
				want = append(want, ',')
			}
		}
		want = append(want, ']')

		got := enc.AppendInts16([]byte{}, vals)
		if !bytes.Equal(got, want) {
			t.Errorf("AppendInts16(%v)\ngot:  %s\nwant: %s",
				vals,
				string(got),
				string(want))
		}
	}

	array := make([]int16, 0)
	for _, tc := range internal.IntegerTestCases {
		if (tc.Val < math.MinInt16) || (tc.Val > math.MaxInt16) {
			continue
		}
		array = append(array, int16(tc.Val))
	}

	doOne(array)
	doOne(array[:1]) // single element
	doOne(array[:0]) // edge case of zero length
}
func TestAppendUints16(t *testing.T) {
	doOne := func(vals []uint16) {
		want := make([]byte, 0)
		want = append(want, '[')
		for i, val := range vals {
			want = append(want, []byte(fmt.Sprintf("%d", uint16(val)))...)
			if i < len(vals)-1 {
				want = append(want, ',')
			}
		}
		want = append(want, ']')

		got := enc.AppendUints16([]byte{}, vals)
		if !bytes.Equal(got, want) {
			t.Errorf("AppendUints16(%v)\ngot:  %s\nwant: %s",
				vals,
				string(got),
				string(want))
		}
	}

	array := make([]uint16, 0)
	for _, tc := range internal.UnsignedIntegerTestCases {
		if tc.Val > math.MaxUint16 {
			continue
		}
		array = append(array, uint16(tc.Val))
	}

	doOne(array)
	doOne(array[:1]) // single element
	doOne(array[:0]) // edge case of zero length
}
func TestAppendInts32(t *testing.T) {
	doOne := func(vals []int32) {
		want := make([]byte, 0)
		want = append(want, '[')
		for i, val := range vals {
			want = append(want, []byte(fmt.Sprintf("%d", int32(val)))...)
			if i < len(vals)-1 {
				want = append(want, ',')
			}
		}
		want = append(want, ']')

		got := enc.AppendInts32([]byte{}, vals)
		if !bytes.Equal(got, want) {
			t.Errorf("AppendInts32(%v)\ngot:  %s\nwant: %s",
				vals,
				string(got),
				string(want))
		}
	}

	array := make([]int32, 0)
	for _, tc := range internal.IntegerTestCases {
		if (tc.Val < math.MinInt32) || (tc.Val > math.MaxInt32) {
			continue
		}
		array = append(array, int32(tc.Val))
	}

	doOne(array)
	doOne(array[:1]) // single element
	doOne(array[:0]) // edge case of zero length
}
func TestAppendUints32(t *testing.T) {
	doOne := func(vals []uint32) {
		want := make([]byte, 0)
		want = append(want, '[')
		for i, val := range vals {
			want = append(want, []byte(fmt.Sprintf("%d", uint32(val)))...)
			if i < len(vals)-1 {
				want = append(want, ',')
			}
		}
		want = append(want, ']')

		got := enc.AppendUints32([]byte{}, vals)
		if !bytes.Equal(got, want) {
			t.Errorf("AppendUints32(%v)\ngot:  %s\nwant: %s",
				vals,
				string(got),
				string(want))
		}
	}

	array := make([]uint32, 0)
	for _, tc := range internal.UnsignedIntegerTestCases {
		if tc.Val > math.MaxUint32 {
			continue
		}
		array = append(array, uint32(tc.Val))
	}

	doOne(array)
	doOne(array[:1]) // single element
	doOne(array[:0]) // edge case of zero length
}

func TestAppendInt64(t *testing.T) {
	doOne := func(vals []int64) {
		want := make([]byte, 0)
		want = append(want, '[')
		for i, val := range vals {
			want = append(want, []byte(fmt.Sprintf("%d", int64(val)))...)
			if i < len(vals)-1 {
				want = append(want, ',')
			}
		}
		want = append(want, ']')

		got := enc.AppendInts64([]byte{}, vals)
		if !bytes.Equal(got, want) {
			t.Errorf("AppendInts64(%v)\ngot:  %s\nwant: %s",
				vals,
				string(got),
				string(want))
		}
	}

	array := make([]int64, 0)
	for _, tc := range internal.IntegerTestCases {
		array = append(array, int64(tc.Val))
	}

	doOne(array)
	doOne(array[:1]) // single element
	doOne(array[:0]) // edge case of zero length
}
func TestAppendUints64(t *testing.T) {
	doOne := func(vals []uint64) {
		want := make([]byte, 0)
		want = append(want, '[')
		for i, val := range vals {
			want = append(want, []byte(fmt.Sprintf("%d", uint64(val)))...)
			if i < len(vals)-1 {
				want = append(want, ',')
			}
		}
		want = append(want, ']')

		got := enc.AppendUints64([]byte{}, vals)
		if !bytes.Equal(got, want) {
			t.Errorf("AppendUints64(%v)\ngot:  %s\nwant: %s",
				vals,
				string(got),
				string(want))
		}
	}

	array := make([]uint64, 0)
	for _, tc := range internal.UnsignedIntegerTestCases {
		array = append(array, uint64(tc.Val))
	}

	doOne(array)
	doOne(array[:1]) // single element
	doOne(array[:0]) // edge case of zero length
}

func TestAppendInt(t *testing.T) {
	for _, tc := range internal.IntegerTestCases {
		want := []byte(fmt.Sprintf("%d", tc.Val))
		got := enc.AppendInt([]byte{}, tc.Val)
		if !bytes.Equal(got, want) {
			t.Errorf("AppendInt(0x%x)\ngot:  %s\nwant: %s",
				tc.Val,
				string(got),
				string(want))
		}
	}
}
func TestAppendUint(t *testing.T) {
	for _, tc := range internal.UnsignedIntegerTestCases {
		want := []byte(fmt.Sprintf("%d", tc.Val))
		got := enc.AppendUint([]byte{}, tc.Val)
		if !bytes.Equal(got, want) {
			t.Errorf("AppendUint(0x%x)\ngot:  %s\nwant: %s",
				tc.Val,
				string(got),
				string(want))
		}
	}
}

func TestAppendInts(t *testing.T) {
	doOne := func(vals []int) {
		want := make([]byte, 0)
		want = append(want, '[')
		for i, val := range vals {
			want = append(want, []byte(fmt.Sprintf("%d", int(val)))...)
			if i < len(vals)-1 {
				want = append(want, ',')
			}
		}
		want = append(want, ']')

		got := enc.AppendInts([]byte{}, vals)
		if !bytes.Equal(got, want) {
			t.Errorf("AppendInts(%v)\ngot:  %s\nwant: %s",
				vals,
				string(got),
				string(want))
		}
	}

	array := make([]int, 0)
	for _, tc := range internal.IntegerTestCases {
		array = append(array, int(tc.Val))
	}

	doOne(array)
	doOne(array[:1]) // single element
	doOne(array[:0]) // edge case of zero length
}
func TestAppendUints(t *testing.T) {
	doOne := func(vals []uint) {
		want := make([]byte, 0)
		want = append(want, '[')
		for i, val := range vals {
			want = append(want, []byte(fmt.Sprintf("%d", uint(val)))...)
			if i < len(vals)-1 {
				want = append(want, ',')
			}
		}
		want = append(want, ']')

		got := enc.AppendUints([]byte{}, vals)
		if !bytes.Equal(got, want) {
			t.Errorf("AppendUints(%v)\ngot:  %s\nwant: %s",
				vals,
				string(got),
				string(want))
		}
	}

	array := make([]uint, 0)
	for _, tc := range internal.UnsignedIntegerTestCases {
		array = append(array, uint(tc.Val))
	}

	doOne(array)
	doOne(array[:1]) // single element
	doOne(array[:0]) // edge case of zero length
}
