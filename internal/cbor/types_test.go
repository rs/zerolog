package cbor

import (
	"encoding/hex"
	"testing"
)

func TestAppendNull(t *testing.T) {
	s := AppendNull([]byte{})
	got := string(s)
	want := "\xf6"
	if got != want {
		t.Errorf("appendNull() = 0x%s, want: 0x%s", hex.EncodeToString(s),
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
		s := AppendBool([]byte{}, tc.val)
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
	{[]bool{true, false, true}, "\x83\xf5\xf4\xf5", "[true,false,true]"},
	{[]bool{true, false, false, true, false, true}, "\x86\xf5\xf4\xf4\xf5\xf4\xf5", "[true,false,false,true,false,true]"},
}

func TestAppendBoolArray(t *testing.T) {
	for _, tc := range booleanArrayTestCases {
		s := AppendBools([]byte{}, tc.val)
		got := string(s)
		if got != tc.binary {
			t.Errorf("AppendBools(%s)=0x%s, want: 0x%s",
				tc.json, hex.EncodeToString(s),
				hex.EncodeToString([]byte(tc.binary)))
		}
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
	{100, "\x18\x64"},
	{254, "\x18\xfe"},
	{255, "\x18\xff"},
	// Value in 2 bytes.
	{256, "\x19\x01\x00"},
	{257, "\x19\x01\x01"},
	{1000, "\x19\x03\xe8"},
	{0xFFFF, "\x19\xff\xff"},
	// Value in 4 bytes.
	{0x10000, "\x1a\x00\x01\x00\x00"},
	{0xFFFFFFFE, "\x1a\xff\xff\xff\xfe"},
	{1000000, "\x1a\x00\x0f\x42\x40"},
	// Value in 8 bytes.
	{0xabcd100000000, "\x1b\x00\x0a\xbc\xd1\x00\x00\x00\x00"},
	{1000000000000, "\x1b\x00\x00\x00\xe8\xd4\xa5\x10\x00"},
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
	{-254, "\x38\xfd"},
	{-255, "\x38\xfe"},
	{-256, "\x38\xff"},
	// Value in 2 bytes.
	{-257, "\x39\x01\x00"},
	{-258, "\x39\x01\x01"},
	{-1000, "\x39\x03\xe7"},
	// Value in 4 bytes.
	{-0x10001, "\x3a\x00\x01\x00\x00"},
	{-0xFFFFFFFE, "\x3a\xff\xff\xff\xfd"},
	{-1000000, "\x3a\x00\x0f\x42\x3f"},
	// Value in 8 bytes.
	{-0xabcd100000001, "\x3b\x00\x0a\xbc\xd1\x00\x00\x00\x00"},
	{-1000000000001, "\x3b\x00\x00\x00\xe8\xd4\xa5\x10\x00"},
}

func TestAppendInt(t *testing.T) {
	for _, tc := range integerTestCases {
		s := AppendInt([]byte{}, tc.val)
		got := string(s)
		if got != tc.binary {
			t.Errorf("AppendInt(0x%x)=0x%s, want: 0x%s",
				tc.val, hex.EncodeToString(s),
				hex.EncodeToString([]byte(tc.binary)))
		}
	}
}

var integerArrayTestCases = []struct {
	val    []int
	binary string
	json   string
}{
	{[]int{-1, 0, 200, 20}, "\x84\x20\x00\x18\xc8\x14", "[-1,0,200,20]"},
	{[]int{-200, -10, 200, 400}, "\x84\x38\xc7\x29\x18\xc8\x19\x01\x90", "[-200,-10,200,400]"},
	{[]int{1, 2, 3}, "\x83\x01\x02\x03", "[1,2,3]"},
	{[]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25},
		"\x98\x19\x01\x02\x03\x04\x05\x06\x07\x08\x09\x0a\x0b\x0c\x0d\x0e\x0f\x10\x11\x12\x13\x14\x15\x16\x17\x18\x18\x18\x19",
		"[1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17,18,19,20,21,22,23,24,25]"},
}

func TestAppendIntArray(t *testing.T) {
	for _, tc := range integerArrayTestCases {
		s := AppendInts([]byte{}, tc.val)
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
}

func TestAppendFloat32(t *testing.T) {
	for _, tc := range float32TestCases {
		s := AppendFloat32([]byte{}, tc.val)
		got := string(s)
		if got != tc.binary {
			t.Errorf("AppendFloat32(%f)=0x%s, want: 0x%s",
				tc.val, hex.EncodeToString(s),
				hex.EncodeToString([]byte(tc.binary)))
		}
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
					_ = AppendInt(buf, int(str.val))
				case 1:
					_ = AppendUint8(buf, uint8(str.val))
				case 2:
					_ = AppendUint16(buf, uint16(str.val))
				case 4:
					_ = AppendUint32(buf, uint32(str.val))
				case 8:
					_ = AppendUint64(buf, uint64(str.val))
				case 21:
					_ = AppendInt8(buf, int8(str.val))
				case 22:
					_ = AppendInt16(buf, int16(str.val))
				case 23:
					_ = AppendInt32(buf, int32(str.val))
				case 24:
					_ = AppendInt64(buf, int64(str.val))
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
					_ = AppendFloat32(buf, float32(str.val))
				case 8:
					_ = AppendFloat64(buf, str.val)
				}
			}
		})
	}
}
