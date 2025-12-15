package cbor

import (
	"bytes"
	"encoding/hex"
	"errors"
	"math"
	"net"
	"testing"

	"github.com/rs/zerolog/internal"
)

var enc = Encoder{}

func TestAppendNil(t *testing.T) {
	s := enc.AppendNil([]byte{})
	got := string(s)
	want := "\xf6"
	if got != want {
		t.Errorf("AppendNil() = 0x%s, want: 0x%s", hex.EncodeToString(s),
			hex.EncodeToString([]byte(want)))
	}
}
func TestAppendBeginMarker(t *testing.T) {
	s := enc.AppendBeginMarker([]byte{})
	got := string(s)
	want := "\xbf"
	if got != want {
		t.Errorf("AppendBeginMarker() = 0x%s, want: 0x%s", hex.EncodeToString(s),
			hex.EncodeToString([]byte(want)))
	}
}
func TestAppendEndMarker(t *testing.T) {
	s := enc.AppendEndMarker([]byte{})
	got := string(s)
	want := "\xff"
	if got != want {
		t.Errorf("AppendEndMarker() = 0x%s, want: 0x%s", hex.EncodeToString(s),
			hex.EncodeToString([]byte(want)))
	}
}
func TestAppendObjectData(t *testing.T) {
	data := []byte{0xbf, 0x02, 0x03}
	s := enc.AppendObjectData([]byte{}, data)
	got := string(s)
	want := "\x02\x03" // the begin marker is not copied
	if got != want {
		t.Errorf("AppendObjectData() = 0x%s, want: 0x%s", hex.EncodeToString(s),
			hex.EncodeToString([]byte(want)))
	}
}
func TestAppendArrayDelim(t *testing.T) {
	s := enc.AppendArrayDelim([]byte{})
	got := string(s)
	want := ""
	if got != want {
		t.Errorf("AppendArrayDelim() = 0x%s, want: 0x%s", hex.EncodeToString(s),
			hex.EncodeToString([]byte(want)))
	}
}
func TestAppendLineBreak(t *testing.T) {
	s := enc.AppendLineBreak([]byte{})
	got := string(s)
	want := ""
	if got != want {
		t.Errorf("AppendLineBreak() = 0x%s, want: 0x%s", hex.EncodeToString(s),
			hex.EncodeToString([]byte(want)))
	}
}

func TestAppendInterface(t *testing.T) {
	oldJSONMarshalFunc := JSONMarshalFunc
	defer func() {
		JSONMarshalFunc = oldJSONMarshalFunc
	}()

	JSONMarshalFunc = func(v interface{}) ([]byte, error) {
		return internal.InterfaceMarshalFunc(v)
	}

	var i int = 17
	got := enc.AppendInterface([]byte{}, i)
	want := make([]byte, 0)
	want = append(want, 0xd9, 0x01) // start an array
	want = append(want, 0x06, 0x42) // of type interface, two characters
	want = append(want, 0x31, 0x37) // with literal int 17
	if !bytes.Equal(got, want) {
		t.Errorf("AppendInterface\ngot:  0x%s\nwant: 0x%s",
			hex.EncodeToString(got),
			hex.EncodeToString(want))
	}

	JSONMarshalFunc = func(v interface{}) ([]byte, error) {
		return nil, errors.New("test")
	}

	got = enc.AppendInterface([]byte{}, nil)
	want = make([]byte, 0)
	want = append(want, 0x76)                                // string
	want = append(want, []byte("marshaling error: test")...) // of type interface, two characters
	if !bytes.Equal(got, want) {
		t.Errorf("AppendInterface\ngot:  0x%s\nwant: 0x%s",
			hex.EncodeToString(got),
			hex.EncodeToString(want))
	}
}

func TestAppendType(t *testing.T) {
	s := enc.AppendType([]byte{}, "")
	got := string(s)
	want := "\x66string"
	if got != want {
		t.Errorf("AppendType() = 0x%s, want: 0x%s", hex.EncodeToString(s),
			hex.EncodeToString([]byte(want)))
	}

	s = enc.AppendType([]byte{}, nil)
	got = string(s)
	want = "\x65<nil>"
	if got != want {
		t.Errorf("AppendType() = 0x%s, want: 0x%s", hex.EncodeToString(s),
			hex.EncodeToString([]byte(want)))
	}

	var n *int = nil
	s = enc.AppendType([]byte{}, n)
	got = string(s)
	want = "\x64*int"
	if got != want {
		t.Errorf("AppendType() = 0x%s, want: 0x%s", hex.EncodeToString(s),
			hex.EncodeToString([]byte(want)))
	}
}

func TestAppendBool(t *testing.T) {
	for _, tc := range internal.BooleanTestCases {
		s := enc.AppendBool([]byte{}, tc.Val)
		got := string(s)
		if got != tc.Binary {
			t.Errorf("AppendBool(%s)=0x%s, want: 0x%s",
				tc.Json, hex.EncodeToString(s),
				hex.EncodeToString([]byte(tc.Binary)))
		}
	}
}

func TestAppendBoolArray(t *testing.T) {
	for _, tc := range internal.BooleanArrayTestCases {
		s := enc.AppendBools([]byte{}, tc.Val)
		got := string(s)
		if got != tc.Binary {
			t.Errorf("AppendBools(%s)=0x%s, want: 0x%s",
				tc.Json, hex.EncodeToString(s),
				hex.EncodeToString([]byte(tc.Binary)))
		}
	}

	// now empty array case
	array := make([]bool, 0)
	want := make([]byte, 0)
	want = append(want, 0x9f) // start and end array
	want = append(want, 0xff) // for empty array
	got := enc.AppendBools([]byte{}, array)
	if !bytes.Equal(got, want) {
		t.Errorf("AppendBools(%v)\ngot:  0x%s\nwant: 0x%s",
			array, hex.EncodeToString(got),
			hex.EncodeToString(want))
	}

	// now a large array case
	array = make([]bool, 24)
	want = make([]byte, 0)
	want = append(want, 0x98) // start a large array
	want = append(want, 0x18) // for 24 elements
	for i := 0; i < 24; i++ {
		array[i] = bool(i%2 == 1)
		want = append(want, 0xf4|byte(i&0x01)) // 0xf4 is false, 0xf5 is true
	}
	got = enc.AppendBools([]byte{}, array)
	if !bytes.Equal(got, want) {
		t.Errorf("AppendBools(%v)\ngot:  0x%s\nwant: 0x%s",
			array, hex.EncodeToString(got),
			hex.EncodeToString(want))
	}
}

func TestAppendInt8(t *testing.T) {
	for _, tc := range internal.IntegerTestCases {
		if (tc.Val < math.MinInt8) || (tc.Val > math.MaxInt8) {
			continue
		}
		s := enc.AppendInt8([]byte{}, int8(tc.Val))
		got := string(s)
		if got != tc.Binary {
			t.Errorf("AppendInt8(0x%x)=0x%s, want: 0x%s",
				tc.Val, hex.EncodeToString(s),
				hex.EncodeToString([]byte(tc.Binary)))
		}
	}
}

func TestAppendInts8(t *testing.T) {
	array := make([]int8, 0)
	want := make([]byte, 0)
	want = append(want, 0x98) // start array
	want = append(want, 0x1b) // for signed 8-bit elements
	for _, tc := range internal.IntegerTestCases {
		if (tc.Val < math.MinInt8) || (tc.Val > math.MaxInt8) {
			continue
		}
		array = append(array, int8(tc.Val))
		want = append(want, tc.Binary...)
	}

	got := enc.AppendInts8([]byte{}, array)
	if !bytes.Equal(got, want) {
		t.Errorf("AppendInts8(%v)\ngot:  0x%s\nwant: 0x%s",
			array, hex.EncodeToString(got),
			hex.EncodeToString(want))
	}

	// now empty array case
	array = make([]int8, 0)
	want = make([]byte, 0)
	want = append(want, 0x9f) // start and end array
	want = append(want, 0xff) // for empty array
	got = enc.AppendInts8([]byte{}, array)
	if !bytes.Equal(got, want) {
		t.Errorf("AppendInts8(%v)\ngot:  0x%s\nwant: 0x%s",
			array, hex.EncodeToString(got),
			hex.EncodeToString(want))
	}

	// now small array case
	array = make([]int8, 21)
	want = make([]byte, 0)
	want = append(want, 0x95) // start a small array
	for i := 0; i < 21; i++ {
		array[i] = int8(i)
		want = append(want, byte(i))
	}
	got = enc.AppendInts8([]byte{}, array)
	if !bytes.Equal(got, want) {
		t.Errorf("AppendInts8(%v)\ngot:  0x%s\nwant: 0x%s",
			array, hex.EncodeToString(got),
			hex.EncodeToString(want))
	}
}

func TestAppendInt16(t *testing.T) {
	for _, tc := range internal.IntegerTestCases {
		if (tc.Val < math.MinInt16) || (tc.Val > math.MaxInt16) {
			continue
		}
		s := enc.AppendInt16([]byte{}, int16(tc.Val))
		got := string(s)
		if got != tc.Binary {
			t.Errorf("AppendInt16(0x%x)=0x%s, want: 0x%s",
				tc.Val, hex.EncodeToString(s),
				hex.EncodeToString([]byte(tc.Binary)))
		}
	}
}

func TestAppendInts16(t *testing.T) {
	array := make([]int16, 0)
	want := make([]byte, 0)
	want = append(want, 0x98) // start array
	want = append(want, 0x28) // for signed 16-bit elements
	for _, tc := range internal.IntegerTestCases {
		if (tc.Val < math.MinInt16) || (tc.Val > math.MaxInt16) {
			continue
		}
		array = append(array, int16(tc.Val))
		want = append(want, tc.Binary...)
	}

	got := enc.AppendInts16([]byte{}, array)
	if !bytes.Equal(got, want) {
		t.Errorf("AppendInts16(%v)\ngot:  0x%s\nwant: 0x%s",
			array, hex.EncodeToString(got),
			hex.EncodeToString(want))
	}

	// now empty array case
	array = make([]int16, 0)
	want = make([]byte, 0)
	want = append(want, 0x9f) // start and end array
	want = append(want, 0xff) // for empty array
	got = enc.AppendInts16([]byte{}, array)
	if !bytes.Equal(got, want) {
		t.Errorf("AppendInts16(%v)\ngot:  0x%s\nwant: 0x%s",
			array, hex.EncodeToString(got),
			hex.EncodeToString(want))
	}

	// now a small array case
	array = make([]int16, 21)
	want = make([]byte, 0)
	want = append(want, 0x95) // start a smaller array
	for i := 0; i < 21; i++ {
		array[i] = int16(i)
		want = append(want, byte(i))
	}
	got = enc.AppendInts16([]byte{}, array)
	if !bytes.Equal(got, want) {
		t.Errorf("AppendInts16(%v)\ngot:  0x%s\nwant: 0x%s",
			array, hex.EncodeToString(got),
			hex.EncodeToString(want))
	}
}

func TestAppendInt32(t *testing.T) {
	for _, tc := range internal.IntegerTestCases {
		if (tc.Val < math.MinInt32) || (tc.Val > math.MaxInt32) {
			continue
		}
		s := enc.AppendInt32([]byte{}, int32(tc.Val))
		got := string(s)
		if got != tc.Binary {
			t.Errorf("AppendInt32(0x%x)=0x%s, want: 0x%s",
				tc.Val, hex.EncodeToString(s),
				hex.EncodeToString([]byte(tc.Binary)))
		}
	}
}

func TestAppendInts32(t *testing.T) {
	array := make([]int32, 0)
	want := make([]byte, 0)
	want = append(want, 0x98) // start array
	want = append(want, 0x31) // for signed 32-bit elements
	for _, tc := range internal.IntegerTestCases {
		if (tc.Val < math.MinInt32) || (tc.Val > math.MaxInt32) {
			continue
		}
		array = append(array, int32(tc.Val))
		want = append(want, tc.Binary...)
	}

	got := enc.AppendInts32([]byte{}, array)
	if !bytes.Equal(got, want) {
		t.Errorf("AppendInts32(%v)\ngot:  0x%s\nwant: 0x%s",
			array, hex.EncodeToString(got),
			hex.EncodeToString(want))
	}

	// now empty array case
	array = make([]int32, 0)
	want = make([]byte, 0)
	want = append(want, 0x9f) // start and end array
	want = append(want, 0xff) // for empty array
	got = enc.AppendInts32([]byte{}, array)
	if !bytes.Equal(got, want) {
		t.Errorf("AppendInts32(%v)\ngot:  0x%s\nwant: 0x%s",
			array, hex.EncodeToString(got),
			hex.EncodeToString(want))
	}

	// now a small array case
	array = make([]int32, 21)
	want = make([]byte, 0)
	want = append(want, 0x95) // start a smaller array
	for i := 0; i < 21; i++ {
		array[i] = int32(i)
		want = append(want, byte(i))
	}
	got = enc.AppendInts32([]byte{}, array)
	if !bytes.Equal(got, want) {
		t.Errorf("AppendInts32(%v)\ngot:  0x%s\nwant: 0x%s",
			array, hex.EncodeToString(got),
			hex.EncodeToString(want))
	}
}

func TestAppendInt64(t *testing.T) {
	for _, tc := range internal.IntegerTestCases {
		s := enc.AppendInt64([]byte{}, int64(tc.Val))
		got := string(s)
		if got != tc.Binary {
			t.Errorf("AppendInt64(0x%x)=0x%s, want: 0x%s",
				tc.Val, hex.EncodeToString(s),
				hex.EncodeToString([]byte(tc.Binary)))
		}
	}
}

func TestAppendInts64(t *testing.T) {
	array := make([]int64, 0)
	want := make([]byte, 0)
	want = append(want, 0x98) // start array
	want = append(want, 0x33) // for signed 64-bit elements
	for _, tc := range internal.IntegerTestCases {
		array = append(array, int64(tc.Val))
		want = append(want, tc.Binary...)
	}

	got := enc.AppendInts64([]byte{}, array)
	if !bytes.Equal(got, want) {
		t.Errorf("AppendInts64(%v)\ngot:  0x%s\nwant: 0x%s",
			array, hex.EncodeToString(got),
			hex.EncodeToString(want))
	}

	// now empty array case
	array = make([]int64, 0)
	want = make([]byte, 0)
	want = append(want, 0x9f) // start and end array
	want = append(want, 0xff) // for empty array
	got = enc.AppendInts64([]byte{}, array)
	if !bytes.Equal(got, want) {
		t.Errorf("AppendInts64(%v)\ngot:  0x%s\nwant: 0x%s",
			array, hex.EncodeToString(got),
			hex.EncodeToString(want))
	}

	// now a small array case
	array = make([]int64, 21)
	want = make([]byte, 0)
	want = append(want, 0x95) // start a smaller array
	for i := 0; i < 21; i++ {
		array[i] = int64(i)
		want = append(want, byte(i))
	}
	got = enc.AppendInts64([]byte{}, array)
	if !bytes.Equal(got, want) {
		t.Errorf("AppendInts64(%v)\ngot:  0x%s\nwant: 0x%s",
			array, hex.EncodeToString(got),
			hex.EncodeToString(want))
	}
}

func TestAppendInt(t *testing.T) {
	for _, tc := range internal.IntegerTestCases {
		s := enc.AppendInt([]byte{}, tc.Val)
		got := string(s)
		if got != tc.Binary {
			t.Errorf("AppendInt(0x%x)=0x%s, want: 0x%s",
				tc.Val, hex.EncodeToString(s),
				hex.EncodeToString([]byte(tc.Binary)))
		}
	}
}

func TestAppendInts(t *testing.T) {
	array := make([]int, 0)
	want := make([]byte, 0)
	want = append(want, 0x98) // start array
	want = append(want, 0x33) // for signed int elements
	for _, tc := range internal.IntegerTestCases {
		array = append(array, int(tc.Val))
		want = append(want, tc.Binary...)
	}

	got := enc.AppendInts([]byte{}, array)
	if !bytes.Equal(got, want) {
		t.Errorf("AppendInts(%v)\ngot:  0x%s\nwant: 0x%s",
			array, hex.EncodeToString(got),
			hex.EncodeToString(want))
	}

	// now empty array case
	array = make([]int, 0)
	want = make([]byte, 0)
	want = append(want, 0x9f) // start and end array
	want = append(want, 0xff) // for empty array
	got = enc.AppendInts([]byte{}, array)
	if !bytes.Equal(got, want) {
		t.Errorf("AppendInts(%v)\ngot:  0x%s\nwant: 0x%s",
			array, hex.EncodeToString(got),
			hex.EncodeToString(want))
	}

	// now a small array case
	array = make([]int, 21)
	want = make([]byte, 0)
	want = append(want, 0x95) // start a smaller array
	for i := 0; i < 21; i++ {
		array[i] = int(i)
		want = append(want, byte(i))
	}
	got = enc.AppendInts([]byte{}, array)
	if !bytes.Equal(got, want) {
		t.Errorf("AppendInts(%v)\ngot:  0x%s\nwant: 0x%s",
			array, hex.EncodeToString(got),
			hex.EncodeToString(want))
	}
}

func TestAppendUint8(t *testing.T) {
	for _, tc := range internal.UnsignedIntegerTestCases {
		if tc.Val > math.MaxUint8 {
			continue
		}
		s := enc.AppendUint8([]byte{}, uint8(tc.Val))
		got := string(s)
		if got != tc.Binary {
			t.Errorf("AppendUint8(0x%x)=0x%s, want: 0x%s",
				tc.Val, hex.EncodeToString(s),
				hex.EncodeToString([]byte(tc.Binary)))
		}
	}
}

func TestAppendUints8(t *testing.T) {
	array := make([]uint8, 0)
	want := make([]byte, 0)
	want = append(want, 0x91) // start array for unsigned 8-bit elements
	for _, tc := range internal.UnsignedIntegerTestCases {
		if tc.Val > math.MaxUint8 {
			continue
		}
		array = append(array, uint8(tc.Val))
		want = append(want, tc.Binary...)
	}

	got := enc.AppendUints8([]byte{}, array)
	if !bytes.Equal(got, want) {
		t.Errorf("AppendUints8(%v)=0x%s, want: 0x%s",
			array, hex.EncodeToString(got),
			hex.EncodeToString(want))
	}

	// now empty array case
	array = make([]uint8, 0)
	want = make([]byte, 0)
	want = append(want, 0x9f) // start and end array
	want = append(want, 0xff) // for empty array
	got = enc.AppendUints8([]byte{}, array)
	if !bytes.Equal(got, want) {
		t.Errorf("AppendUints8(%v)\ngot:  0x%s\nwant: 0x%s",
			array, hex.EncodeToString(got),
			hex.EncodeToString(want))
	}

	// now a large array case
	array = make([]uint8, 24)
	want = make([]byte, 0)
	want = append(want, 0x98) // start a large array
	want = append(want, 0x18) // for 24 elements
	for i := 0; i < 24; i++ {
		array[i] = uint8(i)
		want = append(want, byte(i))
	}
	got = enc.AppendUints8([]byte{}, array)
	if !bytes.Equal(got, want) {
		t.Errorf("AppendUints8(%v)\ngot:  0x%s\nwant: 0x%s",
			array, hex.EncodeToString(got),
			hex.EncodeToString(want))
	}
}

func TestAppendUint16(t *testing.T) {
	for _, tc := range internal.UnsignedIntegerTestCases {
		if tc.Val > math.MaxUint16 {
			continue
		}
		s := enc.AppendUint16([]byte{}, uint16(tc.Val))
		got := string(s)
		if got != tc.Binary {
			t.Errorf("AppendUint16(0x%x)=0x%s, want: 0x%s",
				tc.Val, hex.EncodeToString(s),
				hex.EncodeToString([]byte(tc.Binary)))
		}

	}
}

func TestAppendUints16(t *testing.T) {
	array := make([]uint16, 0)
	want := make([]byte, 0)
	want = append(want, 0x97) // start array for unsigned 16-bit elements
	for _, tc := range internal.UnsignedIntegerTestCases {
		if tc.Val > math.MaxUint16 {
			continue
		}
		array = append(array, uint16(tc.Val))
		want = append(want, tc.Binary...)
	}

	got := enc.AppendUints16([]byte{}, array)
	if !bytes.Equal(got, want) {
		t.Errorf("AppendUints16(%v)=0x%s, want: 0x%s",
			array, hex.EncodeToString(got),
			hex.EncodeToString(want))
	}

	// now empty array case
	array = make([]uint16, 0)
	want = make([]byte, 0)
	want = append(want, 0x9f) // start and end array
	want = append(want, 0xff) // for empty array
	got = enc.AppendUints16([]byte{}, array)
	if !bytes.Equal(got, want) {
		t.Errorf("AppendUints8(%v)\ngot:  0x%s\nwant: 0x%s",
			array, hex.EncodeToString(got),
			hex.EncodeToString(want))
	}

	// now a large array case
	array = make([]uint16, 24)
	want = make([]byte, 0)
	want = append(want, 0x98) // start a larger array
	want = append(want, 0x18) // for 24 elements
	for i := 0; i < 24; i++ {
		array[i] = uint16(i)
		want = append(want, byte(i))
	}
	got = enc.AppendUints16([]byte{}, array)
	if !bytes.Equal(got, want) {
		t.Errorf("AppendUints16(%v)\ngot:  0x%s\nwant: 0x%s",
			array, hex.EncodeToString(got),
			hex.EncodeToString(want))
	}
}

func TestAppendUint32(t *testing.T) {
	for _, tc := range internal.UnsignedIntegerTestCases {
		if tc.Val > math.MaxUint32 {
			continue
		}
		s := enc.AppendUint32([]byte{}, uint32(tc.Val))
		got := string(s)
		want := tc.Bigbinary
		if got != want {
			t.Errorf("AppendUint32(0x%x)=0x%s, want: 0x%s",
				tc.Val, hex.EncodeToString(s),
				hex.EncodeToString([]byte(want)))
		}
	}
}

func TestAppendUints32(t *testing.T) {
	array := make([]uint32, 0)
	want := make([]byte, 0)
	want = append(want, 0x98) // start array
	want = append(want, 0x1f) // for unsigned  32-bit elements
	for _, tc := range internal.UnsignedIntegerTestCases {
		if tc.Val > math.MaxUint32 {
			continue
		}
		array = append(array, uint32(tc.Val))
		want = append(want, tc.Bigbinary...)
	}

	got := enc.AppendUints32([]byte{}, array)
	if !bytes.Equal(got, want) {
		t.Errorf("AppendUints32(%v)\ngot:  0x%s\nwant: 0x%s",
			array, hex.EncodeToString(got),
			hex.EncodeToString(want))
	}

	// now empty array case
	array = make([]uint32, 0)
	want = make([]byte, 0)
	want = append(want, 0x9f) // start and end array
	want = append(want, 0xff) // for empty array
	got = enc.AppendUints32([]byte{}, array)
	if !bytes.Equal(got, want) {
		t.Errorf("AppendUints32(%v)\ngot:  0x%s\nwant: 0x%s",
			array, hex.EncodeToString(got),
			hex.EncodeToString(want))
	}

	// now a small array case
	array = make([]uint32, 21)
	want = make([]byte, 0)
	want = append(want, 0x95) // start a smaller array
	for i := 0; i < 21; i++ {
		array[i] = uint32(i)
		want = append(want, byte(i))
	}
	got = enc.AppendUints32([]byte{}, array)
	if !bytes.Equal(got, want) {
		t.Errorf("AppendUints32(%v)\ngot:  0x%s\nwant: 0x%s",
			array, hex.EncodeToString(got),
			hex.EncodeToString(want))
	}
}

func TestAppendUint64(t *testing.T) {
	for _, tc := range internal.UnsignedIntegerTestCases {
		s := enc.AppendUint64([]byte{}, uint64(tc.Val))
		got := string(s)
		want := tc.Bigbinary
		if got != want {
			t.Errorf("AppendUint64(0x%x)=0x%s, want: 0x%s",
				tc.Val, hex.EncodeToString(s),
				hex.EncodeToString([]byte(want)))
		}
	}
}

func TestAppendUints64(t *testing.T) {
	array := make([]uint64, 0)
	want := make([]byte, 0)
	want = append(want, 0x98) // start array
	want = append(want, 0x21) // for unsigned 64-bit elements
	for _, tc := range internal.UnsignedIntegerTestCases {
		array = append(array, uint64(tc.Val))
		want = append(want, tc.Bigbinary...)
	}

	got := enc.AppendUints64([]byte{}, array)
	if !bytes.Equal(got, want) {
		t.Errorf("AppendUints64(%v)\ngot:  0x%s\nwant: 0x%s",
			array, hex.EncodeToString(got),
			hex.EncodeToString(want))
	}

	// now empty array case
	array = make([]uint64, 0)
	want = make([]byte, 0)
	want = append(want, 0x9f) // start and end array
	want = append(want, 0xff) // for empty array
	got = enc.AppendUints64([]byte{}, array)
	if !bytes.Equal(got, want) {
		t.Errorf("AppendUints64(%v)\ngot:  0x%s\nwant: 0x%s",
			array, hex.EncodeToString(got),
			hex.EncodeToString(want))
	}

	// now a small array case
	array = make([]uint64, 21)
	want = make([]byte, 0)
	want = append(want, 0x95) // start a smaller array
	for i := 0; i < 21; i++ {
		array[i] = uint64(i)
		want = append(want, byte(i))
	}
	got = enc.AppendUints64([]byte{}, array)
	if !bytes.Equal(got, want) {
		t.Errorf("AppendUints64(%v)\ngot:  0x%s\nwant: 0x%s",
			array, hex.EncodeToString(got),
			hex.EncodeToString(want))
	}
}

func TestAppendUint(t *testing.T) {
	for _, tc := range internal.UnsignedIntegerTestCases {
		s := enc.AppendUint([]byte{}, tc.Val)
		got := string(s)
		want := tc.Bigbinary
		if tc.Val == math.MaxUint64 {
			want = "\x20" // this is special case for uint max value when using AppendUint
		}
		if got != want {
			t.Errorf("AppendUint(0x%x)=0x%s, want: 0x%s",
				tc.Val, hex.EncodeToString(s),
				hex.EncodeToString([]byte(want)))
		}
	}
}

func TestAppendUints(t *testing.T) {
	array := make([]uint, 0)
	want := make([]byte, 0)
	want = append(want, 0x98) // start array
	want = append(want, 0x21) // for unsigned int elements
	for _, tc := range internal.UnsignedIntegerTestCases {
		array = append(array, uint(tc.Val))
		expected := tc.Bigbinary
		if tc.Val == math.MaxUint64 {
			expected = "\x20" // this is special case for uint max value when using AppendUint
		}
		want = append(want, expected...)
	}

	got := enc.AppendUints([]byte{}, array)
	if !bytes.Equal(got, want) {
		t.Errorf("AppendUints(%v)\ngot:  0x%s\nwant: 0x%s",
			array, hex.EncodeToString(got),
			hex.EncodeToString(want))
	}

	// now empty array case
	array = make([]uint, 0)
	want = make([]byte, 0)
	want = append(want, 0x9f) // start and end array
	want = append(want, 0xff) // for empty array
	got = enc.AppendUints([]byte{}, array)
	if !bytes.Equal(got, want) {
		t.Errorf("AppendUints(%v)\ngot:  0x%s\nwant: 0x%s",
			array, hex.EncodeToString(got),
			hex.EncodeToString(want))
	}

	// now a small array case
	array = make([]uint, 21)
	want = make([]byte, 0)
	want = append(want, 0x95) // start a smaller array
	for i := 0; i < 21; i++ {
		array[i] = uint(i)
		want = append(want, byte(i))
	}
	got = enc.AppendUints([]byte{}, array)
	if !bytes.Equal(got, want) {
		t.Errorf("AppendUints(%v)\ngot:  0x%s\nwant: 0x%s",
			array, hex.EncodeToString(got),
			hex.EncodeToString(want))
	}
}

func TestAppendIntArray(t *testing.T) {
	for _, tc := range internal.IntegerArrayTestCases {
		s := enc.AppendInts([]byte{}, tc.Val)
		got := string(s)
		if got != tc.Binary {
			t.Errorf("AppendInts(%s)=0x%s, want: 0x%s",
				tc.Json, hex.EncodeToString(s),
				hex.EncodeToString([]byte(tc.Binary)))
		}
	}
}

func TestAppendFloat32(t *testing.T) {
	for _, tc := range internal.Float32TestCases {
		s := enc.AppendFloat32([]byte{}, tc.Val, -1)
		got := string(s)
		want := tc.Binary
		if got != want {
			t.Errorf("AppendFloat32(0x%x)=0x%s, want: 0x%s",
				tc.Val, hex.EncodeToString(s),
				hex.EncodeToString([]byte(want)))
		}
	}
}

func TestAppendFloats32(t *testing.T) {
	array := []float32{1.0, 1.5}
	want := make([]byte, 0)
	want = append(want, 0x82)                         // start array for float elements
	want = append(want, 0xfa, 0x3f, 0x80, 0x00, 0x00) // 32 bit 1.0
	want = append(want, 0xfa, 0x3f, 0xc0, 0x00, 0x00) // 32 bit 1.5

	got := enc.AppendFloats32([]byte{}, array, -1)
	if !bytes.Equal(got, want) {
		t.Errorf("AppendFloats32(%v)=0x%s, want: 0x%s",
			array, hex.EncodeToString(got),
			hex.EncodeToString(want))
	}

	// now empty array case
	array = []float32{}
	want = make([]byte, 0)
	want = append(want, 0x9f, 0xff) // start and end array
	got = enc.AppendFloats32([]byte{}, array, -1)
	if !bytes.Equal(got, want) {
		t.Errorf("AppendFloats32(%v)\ngot:  0x%s\nwant: 0x%s",
			array, hex.EncodeToString(got),
			hex.EncodeToString(want))
	}

	// now a large array case
	array = make([]float32, 24)
	want = make([]byte, 0)
	want = append(want, 0x98) // start a larger array
	want = append(want, 0x18) // for 24 elements
	for i := 0; i < 24; i++ {
		want = append(want, 0xfa, 0x00, 0x00, 0x00, 0x00)
	}
	got = enc.AppendFloats32([]byte{}, array, -1)
	if !bytes.Equal(got, want) {
		t.Errorf("AppendFloats32(%v)\ngot:  0x%s\nwant: 0x%s",
			array, hex.EncodeToString(got),
			hex.EncodeToString(want))
	}
}

func TestAppendFloat64(t *testing.T) {
	for _, tc := range internal.Float64TestCases {
		s := enc.AppendFloat64([]byte{}, tc.Val, -1)
		got := string(s)
		if got != tc.Binary && ((got == "NaN") != math.IsNaN(tc.Val)) {
			t.Errorf("AppendFloat64(%f)=0x%s, want: 0x%s",
				tc.Val, hex.EncodeToString(s),
				hex.EncodeToString([]byte(tc.Binary)))
		}
	}
}

func TestAppendFloats64(t *testing.T) {
	array := []float64{1.0, 1.5}
	want := make([]byte, 0)
	want = append(want, 0x82)                                                 // start array for float elements
	want = append(want, 0xfb, 0x3f, 0xf0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00) // 64 bit 1.0
	want = append(want, 0xfb, 0x3f, 0xf8, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00) // 64 bit 1.5

	got := enc.AppendFloats64([]byte{}, array, -1)
	if !bytes.Equal(got, want) {
		t.Errorf("AppendFloats64(%v)=0x%s, want: 0x%s",
			array, hex.EncodeToString(got),
			hex.EncodeToString(want))
	}

	// now empty array case
	array = []float64{}
	want = make([]byte, 0)
	want = append(want, 0x9f, 0xff) // start and end array
	got = enc.AppendFloats64([]byte{}, array, -1)
	if !bytes.Equal(got, want) {
		t.Errorf("AppendFloats64(%v)\ngot:  0x%s\nwant: 0x%s",
			array, hex.EncodeToString(got),
			hex.EncodeToString(want))
	}

	// now a large array case
	array = make([]float64, 24)
	want = make([]byte, 0)
	want = append(want, 0x98) // start a larger array
	want = append(want, 0x18) // for 24 elements
	for i := 0; i < 24; i++ {
		want = append(want, 0xfb, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00)
	}
	got = enc.AppendFloats64([]byte{}, array, -1)
	if !bytes.Equal(got, want) {
		t.Errorf("AppendFloats64(%v)\ngot:  0x%s\nwant: 0x%s",
			array, hex.EncodeToString(got),
			hex.EncodeToString(want))
	}
}

func TestAppendNetworkAddr(t *testing.T) {
	for _, tc := range internal.IpAddrTestCases {
		s := enc.AppendIPAddr([]byte{}, tc.Ipaddr)
		got := string(s)
		if got != tc.Binary {
			t.Errorf("AppendIPAddr(%s)=0x%s, want: 0x%s",
				tc.Ipaddr, hex.EncodeToString(s),
				hex.EncodeToString([]byte(tc.Binary)))
		}
	}
}

func TestAppendIPAddrArray(t *testing.T) {
	for _, tc := range internal.IPAddrArrayTestCases {
		s := enc.AppendIPAddrs([]byte{}, tc.Val)
		got := string(s)
		if got != tc.Binary {
			t.Errorf("AppendIPAddr(%s)=0x%s, want: 0x%s",
				tc.Json, hex.EncodeToString(s),
				hex.EncodeToString([]byte(tc.Binary)))
		}
	}

	// now a large array case
	array := make([]net.IP, 24)
	want := make([]byte, 0)
	want = append(want, 0x98) // start a larger array
	want = append(want, 0x18) // for 24 elements
	for i := 0; i < 24; i++ {
		array[i] = net.IP{0, 0, 0, byte(i)}
		want = append(want, 0xd9, 0x01, 0x04, 0x44, 0x00, 0x00, 0x00, byte(i))
	}
	got := enc.AppendIPAddrs([]byte{}, array)
	if !bytes.Equal(got, want) {
		t.Errorf("AppendIPAddrs(%v)\ngot:  0x%s\nwant: 0x%s",
			array, hex.EncodeToString(got),
			hex.EncodeToString(want))
	}
}

func TestAppendMACAddr(t *testing.T) {
	for _, tc := range internal.MacAddrTestCases {
		s := enc.AppendMACAddr([]byte{}, tc.Macaddr)
		got := string(s)
		if got != tc.Binary {
			t.Errorf("AppendMACAddr(%s)=0x%s, want: 0x%s",
				tc.Macaddr.String(), hex.EncodeToString(s),
				hex.EncodeToString([]byte(tc.Binary)))
		}
	}
}

func TestAppendIPPrefix(t *testing.T) {
	for _, tc := range internal.IPPrefixTestCases {
		s := enc.AppendIPPrefix([]byte{}, tc.Pfx)
		got := string(s)
		if got != tc.Binary {
			t.Errorf("AppendIPPrefix(%s)=0x%s, want: 0x%s",
				tc.Pfx.String(), hex.EncodeToString(s),
				hex.EncodeToString([]byte(tc.Binary)))
		}
	}
}

func TestAppendIPPrefixArray(t *testing.T) {
	for _, tc := range internal.IPPrefixArrayTestCases {
		s := enc.AppendIPPrefixes([]byte{}, tc.Val)
		got := string(s)
		if got != tc.Binary {
			t.Errorf("AppendIPPrefix(%s)=0x%s, want: 0x%s",
				tc.Json, hex.EncodeToString(s),
				hex.EncodeToString([]byte(tc.Binary)))
		}
	}

	// now a large array case
	array := make([]net.IPNet, 24)
	want := make([]byte, 0)
	want = append(want, 0x98) // start a larger array
	want = append(want, 0x18) // for 24 elements
	for i := 0; i < 24; i++ {
		array[i] = net.IPNet{IP: net.IP{0, 0, 0, byte(i)}, Mask: net.CIDRMask(24, 32)}
		want = append(want, 0xd9, 0x01, 0x05, 0xa1, 0x44, 0x00, 0x00, 0x00, byte(i), 0x18, 0x18)
	}
	got := enc.AppendIPPrefixes([]byte{}, array)
	if !bytes.Equal(got, want) {
		t.Errorf("AppendIPPrefixes(%v)\ngot:  0x%s\nwant: 0x%s",
			array, hex.EncodeToString(got),
			hex.EncodeToString(want))
	}
}

func TestAppendHex(t *testing.T) {
	array := []byte{0x01, 0x02}
	want := make([]byte, 0)
	want = append(want, 0xd9, 0x01) // start array
	want = append(want, 0x07, 0x42) // array of two elements
	want = append(want, 0x01, 0x02) // 0x01, 0x02
	got := enc.AppendHex([]byte{}, array)
	if !bytes.Equal(got, want) {
		t.Errorf("AppendHex(%v)=0x%s, want: 0x%s",
			array, hex.EncodeToString(got),
			hex.EncodeToString(want))
	}

	// now empty array case
	array = make([]byte, 0)
	want = make([]byte, 0)
	want = append(want, 0xd9, 0x01) // start an array
	want = append(want, 0x07, 0x40) // array of zero elements
	got = enc.AppendHex([]byte{}, array)
	if !bytes.Equal(got, want) {
		t.Errorf("AppendHex(%v)\ngot:  0x%s\nwant: 0x%s",
			array, hex.EncodeToString(got),
			hex.EncodeToString(want))
	}
}

func BenchmarkAppendInt(b *testing.B) {
	type st struct {
		sz  byte
		val int64
	}
	tests := map[string]st{
		"int-Positive": {sz: 0, val: 10000},
		"int-Negative": {sz: 0, val: -10000},
		"uint8":        {sz: 1, val: 100},
		"uint16":       {sz: 2, val: 0xfff},
		"uint32":       {sz: 4, val: 0xffffff},
		"uint64":       {sz: 8, val: 0xffffffffff},
		"int8":         {sz: 21, val: -120},
		"int16":        {sz: 22, val: -1200},
		"int32":        {sz: 23, val: 32000},
		"int64":        {sz: 24, val: 0xffffffffff},
	}
	for name, str := range tests {
		b.Run(name, func(b *testing.B) {
			buf := make([]byte, 0, 100)
			for i := 0; i < b.N; i++ {
				switch str.sz {
				case 0:
					_ = enc.AppendInt(buf, int(str.val))
				case 1:
					_ = enc.AppendUint8(buf, uint8(str.val))
				case 2:
					_ = enc.AppendUint16(buf, uint16(str.val))
				case 4:
					_ = enc.AppendUint32(buf, uint32(str.val))
				case 8:
					_ = enc.AppendUint64(buf, uint64(str.val))
				case 21:
					_ = enc.AppendInt8(buf, int8(str.val))
				case 22:
					_ = enc.AppendInt16(buf, int16(str.val))
				case 23:
					_ = enc.AppendInt32(buf, int32(str.val))
				case 24:
					_ = enc.AppendInt64(buf, int64(str.val))
				}
			}
		})
	}
}

func BenchmarkAppendFloat(b *testing.B) {
	type st struct {
		sz  byte
		val float64
	}
	tests := map[string]st{
		"Float32": {sz: 4, val: 10000.12345},
		"Float64": {sz: 8, val: -10000.54321},
	}
	for name, str := range tests {
		b.Run(name, func(b *testing.B) {
			buf := make([]byte, 0, 100)
			for i := 0; i < b.N; i++ {
				switch str.sz {
				case 4:
					_ = enc.AppendFloat32(buf, float32(str.val), -1)
				case 8:
					_ = enc.AppendFloat64(buf, str.val, -1)
				}
			}
		})
	}
}
