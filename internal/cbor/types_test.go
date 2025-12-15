package cbor

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"math"
	"net"
	"testing"
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

// inline copy from globals.go of InterfaceMarshalFunc used in tests to avoid import cycle
func interfaceMarshalFunc(v interface{}) ([]byte, error) {
	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	encoder.SetEscapeHTML(false)
	err := encoder.Encode(v)
	if err != nil {
		return nil, err
	}
	b := buf.Bytes()
	if len(b) > 0 {
		// Remove trailing \n which is added by Encode.
		return b[:len(b)-1], nil
	}
	return b, nil
}

func TestAppendInterface(t *testing.T) {
	oldJSONMarshalFunc := JSONMarshalFunc
	defer func() {
		JSONMarshalFunc = oldJSONMarshalFunc
	}()

	JSONMarshalFunc = func(v interface{}) ([]byte, error) {
		return interfaceMarshalFunc(v)
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

var booleanTestCases = []struct {
	val    bool
	binary string
	json   string
}{
	{true, "\xf5", "true"},
	{false, "\xf4", "false"},
}

func TestAppendBool(t *testing.T) {
	for _, tc := range booleanTestCases {
		s := enc.AppendBool([]byte{}, tc.val)
		got := string(s)
		if got != tc.binary {
			t.Errorf("AppendBool(%s)=0x%s, want: 0x%s",
				tc.json, hex.EncodeToString(s),
				hex.EncodeToString([]byte(tc.binary)))
		}
	}
}

var booleanArrayTestCases = []struct {
	val    []bool
	binary string
	json   string
}{
	{[]bool{}, "\x9f\xff", "[]"},
	{[]bool{false}, "\x81\xf4", "[false]"},
	{[]bool{true, false, true}, "\x83\xf5\xf4\xf5", "[true,false,true]"},
	{[]bool{true, false, false, true, false, true}, "\x86\xf5\xf4\xf4\xf5\xf4\xf5", "[true,false,false,true,false,true]"},
}

func TestAppendBoolArray(t *testing.T) {
	for _, tc := range booleanArrayTestCases {
		s := enc.AppendBools([]byte{}, tc.val)
		got := string(s)
		if got != tc.binary {
			t.Errorf("AppendBools(%s)=0x%s, want: 0x%s",
				tc.json, hex.EncodeToString(s),
				hex.EncodeToString([]byte(tc.binary)))
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

var integerTestCases = []struct {
	val    int
	binary string
}{
	// Value included in the type.
	{0, "\x00"},
	{1, "\x01"},
	{2, "\x02"},
	{3, "\x03"},
	{8, "\x08"},
	{9, "\x09"},
	{10, "\x0a"},
	{22, "\x16"},
	{23, "\x17"},

	// Value in 1 byte.
	{24, "\x18\x18"},
	{25, "\x18\x19"},
	{26, "\x18\x1a"},
	{127, "\x18\x7f"},
	{254, "\x18\xfe"},
	{255, "\x18\xff"},

	// Value in 2 bytes.
	{256, "\x19\x01\x00"},
	{257, "\x19\x01\x01"},
	{1000, "\x19\x03\xe8"},
	{0xFFFF, "\x19\xff\xff"},

	// Value in 4 bytes.
	{0x10000, "\x1a\x00\x01\x00\x00"},
	{0x7FFFFFFE, "\x1a\x7f\xff\xff\xfe"},
	{1000000, "\x1a\x00\x0f\x42\x40"},

	// Negative number test cases.
	// Value included in the type.
	{-1, "\x20"},
	{-2, "\x21"},
	{-3, "\x22"},
	{-10, "\x29"},
	{-21, "\x34"},
	{-22, "\x35"},
	{-23, "\x36"},
	{-24, "\x37"},

	// Value in 1 byte.
	{-25, "\x38\x18"},
	{-26, "\x38\x19"},
	{-100, "\x38\x63"},
	{-128, "\x38\x7f"},
	{-254, "\x38\xfd"},
	{-255, "\x38\xfe"},
	{-256, "\x38\xff"},

	// Value in 2 bytes.
	{-257, "\x39\x01\x00"},
	{-258, "\x39\x01\x01"},
	{-1000, "\x39\x03\xe7"},

	// Value in 4 bytes.
	{-0x10001, "\x3a\x00\x01\x00\x00"},
	{-0x7FFFFFFE, "\x3a\x7f\xff\xff\xfd"},
	{-1000000, "\x3a\x00\x0f\x42\x3f"},

	//Constants
	{math.MaxInt8, "\x18\x7f"},
	{math.MinInt8, "\x38\x7f"},
	{math.MaxInt16, "\x19\x7f\xff"},
	{math.MinInt16, "\x39\x7f\xff"},
	{math.MaxInt32, "\x1a\x7f\xff\xff\xff"},
	{math.MinInt32, "\x3a\x7f\xff\xff\xff"},
	{math.MaxInt64, "\x1b\x7f\xff\xff\xff\xff\xff\xff\xff"},
	{math.MinInt64, "\x3b\x7f\xff\xff\xff\xff\xff\xff\xff"},
}

func TestAppendInt8(t *testing.T) {
	for _, tc := range integerTestCases {
		if (tc.val < math.MinInt8) || (tc.val > math.MaxInt8) {
			continue
		}
		s := enc.AppendInt8([]byte{}, int8(tc.val))
		got := string(s)
		if got != tc.binary {
			t.Errorf("AppendInt8(0x%x)=0x%s, want: 0x%s",
				tc.val, hex.EncodeToString(s),
				hex.EncodeToString([]byte(tc.binary)))
		}
	}
}
func TestAppendInts8(t *testing.T) {
	array := make([]int8, 0)
	want := make([]byte, 0)
	want = append(want, 0x98) // start array
	want = append(want, 0x1b) // for signed 8-bit elements
	for _, tc := range integerTestCases {
		if (tc.val < math.MinInt8) || (tc.val > math.MaxInt8) {
			continue
		}
		array = append(array, int8(tc.val))
		want = append(want, tc.binary...)
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
	for _, tc := range integerTestCases {
		if (tc.val < math.MinInt16) || (tc.val > math.MaxInt16) {
			continue
		}
		s := enc.AppendInt16([]byte{}, int16(tc.val))
		got := string(s)
		if got != tc.binary {
			t.Errorf("AppendInt16(0x%x)=0x%s, want: 0x%s",
				tc.val, hex.EncodeToString(s),
				hex.EncodeToString([]byte(tc.binary)))
		}
	}
}
func TestAppendInts16(t *testing.T) {
	array := make([]int16, 0)
	want := make([]byte, 0)
	want = append(want, 0x98) // start array
	want = append(want, 0x28) // for signed 16-bit elements
	for _, tc := range integerTestCases {
		if (tc.val < math.MinInt16) || (tc.val > math.MaxInt16) {
			continue
		}
		array = append(array, int16(tc.val))
		want = append(want, tc.binary...)
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
	for _, tc := range integerTestCases {
		if (tc.val < math.MinInt32) || (tc.val > math.MaxInt32) {
			continue
		}
		s := enc.AppendInt32([]byte{}, int32(tc.val))
		got := string(s)
		if got != tc.binary {
			t.Errorf("AppendInt32(0x%x)=0x%s, want: 0x%s",
				tc.val, hex.EncodeToString(s),
				hex.EncodeToString([]byte(tc.binary)))
		}
	}
}
func TestAppendInts32(t *testing.T) {
	array := make([]int32, 0)
	want := make([]byte, 0)
	want = append(want, 0x98) // start array
	want = append(want, 0x31) // for signed 32-bit elements
	for _, tc := range integerTestCases {
		if (tc.val < math.MinInt32) || (tc.val > math.MaxInt32) {
			continue
		}
		array = append(array, int32(tc.val))
		want = append(want, tc.binary...)
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
	for _, tc := range integerTestCases {
		s := enc.AppendInt64([]byte{}, int64(tc.val))
		got := string(s)
		if got != tc.binary {
			t.Errorf("AppendInt64(0x%x)=0x%s, want: 0x%s",
				tc.val, hex.EncodeToString(s),
				hex.EncodeToString([]byte(tc.binary)))
		}
	}
}
func TestAppendInts64(t *testing.T) {
	array := make([]int64, 0)
	want := make([]byte, 0)
	want = append(want, 0x98) // start array
	want = append(want, 0x33) // for signed 64-bit elements
	for _, tc := range integerTestCases {
		array = append(array, int64(tc.val))
		want = append(want, tc.binary...)
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
	for _, tc := range integerTestCases {
		s := enc.AppendInt([]byte{}, tc.val)
		got := string(s)
		if got != tc.binary {
			t.Errorf("AppendInt(0x%x)=0x%s, want: 0x%s",
				tc.val, hex.EncodeToString(s),
				hex.EncodeToString([]byte(tc.binary)))
		}
	}
}
func TestAppendInts(t *testing.T) {
	array := make([]int, 0)
	want := make([]byte, 0)
	want = append(want, 0x98) // start array
	want = append(want, 0x33) // for signed int elements
	for _, tc := range integerTestCases {
		array = append(array, int(tc.val))
		want = append(want, tc.binary...)
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

type unsignedIntTestCase struct {
	val       uint
	binary    string
	bigbinary string
}

var additionalUnsignedIntegerTestCases = []unsignedIntTestCase{
	{0x7FFFFFFF, "\x18\xff", "\x1a\x7f\xff\xff\xff"},
	{0x80000000, "\x19\xff\xff", "\x1a\x80\x00\x00\x00"},
	{1000000, "\x1b\x80\x00\x00\x00\x00\x00\x00\x00", "\x1a\x00\x0f\x42\x40"},

	//Constants
	{math.MaxUint8, "\x18\xff", "\x18\xff"},
	{math.MaxUint16, "\x19\xff\xff", "\x19\xff\xff"},
	{math.MaxUint32, "\x1a\xff\xff\xff\xff", "\x1a\xff\xff\xff\xff"},
	{math.MaxUint64, "\x1b\xff\xff\xff\xff\xff\xff\xff\xff", "\x1b\xff\xff\xff\xff\xff\xff\xff\xff"},
}

func UnsignedIntegerTestCases() []unsignedIntTestCase {
	size := len(integerTestCases) + len(additionalUnsignedIntegerTestCases)
	cases := make([]unsignedIntTestCase, 0, size)
	cases = append(cases, additionalUnsignedIntegerTestCases...)
	for _, itc := range integerTestCases {
		if itc.val < 0 {
			continue
		}
		cases = append(cases, unsignedIntTestCase{val: uint(itc.val), binary: itc.binary, bigbinary: itc.binary})
	}
	return cases
}

var unsignedIntegerTestCases = UnsignedIntegerTestCases()

func TestAppendUint8(t *testing.T) {
	for _, tc := range unsignedIntegerTestCases {
		if tc.val > math.MaxUint8 {
			continue
		}
		s := enc.AppendUint8([]byte{}, uint8(tc.val))
		got := string(s)
		if got != tc.binary {
			t.Errorf("AppendUint8(0x%x)=0x%s, want: 0x%s",
				tc.val, hex.EncodeToString(s),
				hex.EncodeToString([]byte(tc.binary)))
		}
	}
}
func TestAppendUints8(t *testing.T) {
	array := make([]uint8, 0)
	want := make([]byte, 0)
	want = append(want, 0x91) // start array for unsigned 8-bit elements
	for _, tc := range unsignedIntegerTestCases {
		if tc.val > math.MaxUint8 {
			continue
		}
		array = append(array, uint8(tc.val))
		want = append(want, tc.binary...)
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
	for _, tc := range unsignedIntegerTestCases {
		if tc.val > math.MaxUint16 {
			continue
		}
		s := enc.AppendUint16([]byte{}, uint16(tc.val))
		got := string(s)
		if got != tc.binary {
			t.Errorf("AppendUint16(0x%x)=0x%s, want: 0x%s",
				tc.val, hex.EncodeToString(s),
				hex.EncodeToString([]byte(tc.binary)))
		}

	}
}
func TestAppendUints16(t *testing.T) {
	array := make([]uint16, 0)
	want := make([]byte, 0)
	want = append(want, 0x97) // start array for unsigned 16-bit elements
	for _, tc := range unsignedIntegerTestCases {
		if tc.val > math.MaxUint16 {
			continue
		}
		array = append(array, uint16(tc.val))
		want = append(want, tc.binary...)
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
	for _, tc := range unsignedIntegerTestCases {
		if tc.val > math.MaxUint32 {
			continue
		}
		s := enc.AppendUint32([]byte{}, uint32(tc.val))
		got := string(s)
		want := tc.bigbinary
		if got != want {
			t.Errorf("AppendUint32(0x%x)=0x%s, want: 0x%s",
				tc.val, hex.EncodeToString(s),
				hex.EncodeToString([]byte(want)))
		}
	}
}
func TestAppendUints32(t *testing.T) {
	array := make([]uint32, 0)
	want := make([]byte, 0)
	want = append(want, 0x98) // start array
	want = append(want, 0x1f) // for unsigned  32-bit elements
	for _, tc := range unsignedIntegerTestCases {
		if tc.val > math.MaxUint32 {
			continue
		}
		array = append(array, uint32(tc.val))
		want = append(want, tc.bigbinary...)
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
	for _, tc := range unsignedIntegerTestCases {
		s := enc.AppendUint64([]byte{}, uint64(tc.val))
		got := string(s)
		want := tc.bigbinary
		if got != want {
			t.Errorf("AppendUint64(0x%x)=0x%s, want: 0x%s",
				tc.val, hex.EncodeToString(s),
				hex.EncodeToString([]byte(want)))
		}
	}
}
func TestAppendUints64(t *testing.T) {
	array := make([]uint64, 0)
	want := make([]byte, 0)
	want = append(want, 0x98) // start array
	want = append(want, 0x21) // for unsigned 64-bit elements
	for _, tc := range unsignedIntegerTestCases {
		array = append(array, uint64(tc.val))
		want = append(want, tc.bigbinary...)
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
	for _, tc := range unsignedIntegerTestCases {
		s := enc.AppendUint([]byte{}, tc.val)
		got := string(s)
		want := tc.bigbinary
		if tc.val == math.MaxUint64 {
			want = "\x20" // this is special case for uint max value when using AppendUint
		}
		if got != want {
			t.Errorf("AppendUint(0x%x)=0x%s, want: 0x%s",
				tc.val, hex.EncodeToString(s),
				hex.EncodeToString([]byte(want)))
		}
	}
}
func TestAppendUints(t *testing.T) {
	array := make([]uint, 0)
	want := make([]byte, 0)
	want = append(want, 0x98) // start array
	want = append(want, 0x21) // for unsigned int elements
	for _, tc := range unsignedIntegerTestCases {
		array = append(array, uint(tc.val))
		expected := tc.bigbinary
		if tc.val == math.MaxUint64 {
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

var integerArrayTestCases = []struct {
	val    []int
	binary string
	json   string
}{
	{[]int{}, "\x9f\xff", "[]"},
	{[]int{32768}, "\x81\x19\x80\x00", "[32768]"},
	{[]int{-1, 0, 200, 20}, "\x84\x20\x00\x18\xc8\x14", "[-1,0,200,20]"},
	{[]int{-200, -10, 200, 400}, "\x84\x38\xc7\x29\x18\xc8\x19\x01\x90", "[-200,-10,200,400]"},
	{[]int{1, 2, 3}, "\x83\x01\x02\x03", "[1,2,3]"},
	{[]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25},
		"\x98\x19\x01\x02\x03\x04\x05\x06\x07\x08\x09\x0a\x0b\x0c\x0d\x0e\x0f\x10\x11\x12\x13\x14\x15\x16\x17\x18\x18\x18\x19",
		"[1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17,18,19,20,21,22,23,24,25]"},
}

func TestAppendIntArray(t *testing.T) {
	for _, tc := range integerArrayTestCases {
		s := enc.AppendInts([]byte{}, tc.val)
		got := string(s)
		if got != tc.binary {
			t.Errorf("AppendInts(%s)=0x%s, want: 0x%s",
				tc.json, hex.EncodeToString(s),
				hex.EncodeToString([]byte(tc.binary)))
		}
	}
}

var float32TestCases = []struct {
	val    float32
	binary string
}{
	{0.0, "\xfa\x00\x00\x00\x00"},
	{-0.0, "\xfa\x00\x00\x00\x00"},
	{1.0, "\xfa\x3f\x80\x00\x00"},
	{1.5, "\xfa\x3f\xc0\x00\x00"},
	{65504.0, "\xfa\x47\x7f\xe0\x00"},
	{-4.0, "\xfa\xc0\x80\x00\x00"},
	{0.00006103515625, "\xfa\x38\x80\x00\x00"},
	{float32(math.Inf(0)), "\xfa\x7f\x80\x00\x00"},
	{float32(math.Inf(-1)), "\xfa\xff\x80\x00\x00"},
	{float32(math.NaN()), "\xfa\x7f\xc0\x00\x00"},
	{math.SmallestNonzeroFloat32, "\xfa\x00\x00\x00\x01"},
	{math.MaxFloat32, "\xfa\x7f\x7f\xff\xff"},
}

func TestAppendFloat32(t *testing.T) {
	for _, tc := range float32TestCases {
		s := enc.AppendFloat32([]byte{}, tc.val, -1)
		got := string(s)
		want := tc.binary
		if got != want {
			t.Errorf("AppendFloat32(0x%x)=0x%s, want: 0x%s",
				tc.val, hex.EncodeToString(s),
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

var float64TestCases = []struct {
	val    float64
	binary string
}{
	{0.0, "\xfa\x00\x00\x00\x00"},
	{-0.0, "\xfa\x00\x00\x00\x00"},
	{1.0, "\xfa\x3f\x80\x00\x00"},
	{1.5, "\xfa\x3f\xc0\x00\x00"},
	{65504.0, "\xfa\x47\x7f\xe0\x00"},
	{-4.0, "\xfa\xc0\x80\x00\x00"},
	{0.00006103515625, "\xfa\x38\x80\x00\x00"},
	{math.Inf(0), "\xfa\x7f\x80\x00\x00\x00\x00\x00\x00"},
	{math.Inf(-1), "\xfa\xff\x80\x00\x00\x00\x00\x00\x00"},
	{math.NaN(), "\xfb\x7f\xf8\x00\x00\x00\x00\x00\x00"},
	{math.SmallestNonzeroFloat64, "\xfa\x00\x00\x00\x00\x00\x00\x00\x01"},
	{math.MaxFloat64, "\xfa\x7f\x7f\xff\xff"},
}

func TestAppendFloat64(t *testing.T) {
	for _, tc := range float64TestCases {
		s := enc.AppendFloat64([]byte{}, tc.val, -1)
		got := string(s)
		if got != tc.binary && ((got == "NaN") != math.IsNaN(tc.val)) {
			t.Errorf("AppendFloat64(%f)=0x%s, want: 0x%s",
				tc.val, hex.EncodeToString(s),
				hex.EncodeToString([]byte(tc.binary)))
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

var ipAddrTestCases = []struct {
	ipaddr net.IP
	text   string // ASCII representation of ipaddr
	binary string // CBOR representation of ipaddr
}{
	{net.IP{10, 0, 0, 1}, "\"10.0.0.1\"", "\xd9\x01\x04\x44\x0a\x00\x00\x01"},
	{net.IP{0x20, 0x01, 0x0d, 0xb8, 0x85, 0xa3, 0x0, 0x0, 0x0, 0x0, 0x8a, 0x2e, 0x03, 0x70, 0x73, 0x34},
		"\"2001:db8:85a3::8a2e:370:7334\"",
		"\xd9\x01\x04\x50\x20\x01\x0d\xb8\x85\xa3\x00\x00\x00\x00\x8a\x2e\x03\x70\x73\x34"},
}

func TestAppendNetworkAddr(t *testing.T) {
	for _, tc := range ipAddrTestCases {
		s := enc.AppendIPAddr([]byte{}, tc.ipaddr)
		got := string(s)
		if got != tc.binary {
			t.Errorf("AppendIPAddr(%s)=0x%s, want: 0x%s",
				tc.ipaddr, hex.EncodeToString(s),
				hex.EncodeToString([]byte(tc.binary)))
		}
	}
}

var IPAddrArrayTestCases = []struct {
	val    []net.IP
	binary string
	json   string
}{
	{[]net.IP{}, "\x9f\xff", "[]"},
	{[]net.IP{{127, 0, 0, 0}}, "\x81\xd9\x01\x04\x44\x7f\x00\x00\x00", "[127.0.0.0]"},
	{[]net.IP{{0, 0, 0, 0}, {192, 168, 0, 100}}, "\x82\xd9\x01\x04\x44\x00\x00\x00\x00\xd9\x01\x04\x44\xc0\xa8\x00\x64", "[0.0.0.0,192.168.0.100]"},
}

func TestAppendIPAddrArray(t *testing.T) {
	for _, tc := range IPAddrArrayTestCases {
		s := enc.AppendIPAddrs([]byte{}, tc.val)
		got := string(s)
		if got != tc.binary {
			t.Errorf("AppendIPAddr(%s)=0x%s, want: 0x%s",
				tc.json, hex.EncodeToString(s),
				hex.EncodeToString([]byte(tc.binary)))
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

var macAddrTestCases = []struct {
	macaddr net.HardwareAddr
	text    string // ASCII representation of macaddr
	binary  string // CBOR representation of macaddr
}{
	{net.HardwareAddr{0x12, 0x34, 0x56, 0x78, 0x90, 0xab}, "\"12:34:56:78:90:ab\"", "\xd9\x01\x04\x46\x12\x34\x56\x78\x90\xab"},
	{net.HardwareAddr{0x20, 0x01, 0x0d, 0xb8, 0x85, 0xa3}, "\"20:01:0d:b8:85:a3\"", "\xd9\x01\x04\x46\x20\x01\x0d\xb8\x85\xa3"},
}

func TestAppendMACAddr(t *testing.T) {
	for _, tc := range macAddrTestCases {
		s := enc.AppendMACAddr([]byte{}, tc.macaddr)
		got := string(s)
		if got != tc.binary {
			t.Errorf("AppendMACAddr(%s)=0x%s, want: 0x%s",
				tc.macaddr.String(), hex.EncodeToString(s),
				hex.EncodeToString([]byte(tc.binary)))
		}
	}
}

var IPPrefixTestCases = []struct {
	pfx    net.IPNet
	text   string // ASCII representation of pfx
	binary string // CBOR representation of pfx
}{
	{net.IPNet{IP: net.IP{0, 0, 0, 0}, Mask: net.CIDRMask(0, 32)}, "\"0.0.0.0/0\"", "\xd9\x01\x05\xa1\x44\x00\x00\x00\x00\x00"},
	{net.IPNet{IP: net.IP{192, 168, 0, 100}, Mask: net.CIDRMask(24, 32)}, "\"192.168.0.100/24\"",
		"\xd9\x01\x05\xa1\x44\xc0\xa8\x00\x64\x18\x18"},
}

func TestAppendIPPrefix(t *testing.T) {
	for _, tc := range IPPrefixTestCases {
		s := enc.AppendIPPrefix([]byte{}, tc.pfx)
		got := string(s)
		if got != tc.binary {
			t.Errorf("AppendIPPrefix(%s)=0x%s, want: 0x%s",
				tc.pfx.String(), hex.EncodeToString(s),
				hex.EncodeToString([]byte(tc.binary)))
		}
	}
}

var IPPrefixArrayTestCases = []struct {
	val    []net.IPNet
	binary string
	json   string
}{
	{[]net.IPNet{}, "\x9f\xff", "[]"},
	{[]net.IPNet{{IP: net.IP{127, 0, 0, 0}, Mask: net.CIDRMask(24, 32)}}, "\x81\xd9\x01\x05\xa1\x44\x7f\x00\x00\x00\x18\x18", "[127.0.0.0/24]"},
	{[]net.IPNet{{IP: net.IP{0, 0, 0, 0}, Mask: net.CIDRMask(0, 32)}, {IP: net.IP{192, 168, 0, 100}, Mask: net.CIDRMask(24, 32)}}, "\x82\xd9\x01\x05\xa1\x44\x00\x00\x00\x00\x00\xd9\x01\x05\xa1\x44\xc0\xa8\x00\x64\x18\x18", "[0.0.0.0/0,192.168.0.100/24]"},
}

func TestAppendIPPrefixArray(t *testing.T) {
	for _, tc := range IPPrefixArrayTestCases {
		s := enc.AppendIPPrefixes([]byte{}, tc.val)
		got := string(s)
		if got != tc.binary {
			t.Errorf("AppendIPPrefix(%s)=0x%s, want: 0x%s",
				tc.json, hex.EncodeToString(s),
				hex.EncodeToString([]byte(tc.binary)))
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
