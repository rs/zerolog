package cbor

import (
	"bytes"
	"encoding/hex"
	"math"
	"testing"
	"time"

	"github.com/rs/zerolog/internal"
)

func TestDecodeInteger(t *testing.T) {
	for _, tc := range internal.IntegerTestCases {
		gotv := decodeInteger(getReader(tc.Binary))
		if gotv != int64(tc.Val) {
			t.Errorf("decodeInteger(0x%s)=0x%d, want: 0x%d",
				hex.EncodeToString([]byte(tc.Binary)), gotv, tc.Val)
		}
	}
}

func TestDecodeString(t *testing.T) {
	for _, tt := range encodeStringTests {
		got := decodeUTF8String(getReader(tt.binary))
		if string(got) != "\""+tt.json+"\"" {
			t.Errorf("DecodeString(0x%s)=%s, want:\"%s\"\n", hex.EncodeToString([]byte(tt.binary)), string(got),
				hex.EncodeToString([]byte(tt.json)))
		}
	}
}

func TestDecodeArray(t *testing.T) {
	for _, tc := range internal.IntegerArrayTestCases {
		buf := bytes.NewBuffer([]byte{})
		array2Json(getReader(tc.Binary), buf)
		if buf.String() != tc.Json {
			t.Errorf("array2Json(0x%s)=%s, want: %s", hex.EncodeToString([]byte(tc.Binary)), buf.String(), tc.Json)
		}
	}
	//Unspecified Length Array
	var infiniteArrayTestCases = []struct {
		in  string
		out string
	}{
		{"\x9f\x20\x00\x18\xc8\x14\xff", "[-1,0,200,20]"},
		{"\x9f\x38\xc7\x29\x18\xc8\x19\x01\x90\xff", "[-200,-10,200,400]"},
		{"\x9f\x01\x02\x03\xff", "[1,2,3]"},
		{"\x9f\x01\x02\x03\x04\x05\x06\x07\x08\x09\x0a\x0b\x0c\x0d\x0e\x0f\x10\x11\x12\x13\x14\x15\x16\x17\x18\x18\x18\x19\xff",
			"[1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17,18,19,20,21,22,23,24,25]"},
	}
	for _, tc := range infiniteArrayTestCases {
		buf := bytes.NewBuffer([]byte{})
		array2Json(getReader(tc.in), buf)
		if buf.String() != tc.out {
			t.Errorf("array2Json(0x%s)=%s, want: %s", hex.EncodeToString([]byte(tc.out)), buf.String(), tc.out)
		}
	}
	for _, tc := range internal.BooleanArrayTestCases {
		buf := bytes.NewBuffer([]byte{})
		array2Json(getReader(tc.Binary), buf)
		if buf.String() != tc.Json {
			t.Errorf("array2Json(0x%s)=%s, want: %s", hex.EncodeToString([]byte(tc.Binary)), buf.String(), tc.Json)
		}
	}
	//TODO add cases for arrays of other types
}

var infiniteMapDecodeTestCases = []struct {
	Bin  []byte
	Json string
}{
	{[]byte("\xbf\x64IETF\x20\xff"), "{\"IETF\":-1}"},
	{[]byte("\xbf\x65Array\x84\x20\x00\x18\xc8\x14\xff"), "{\"Array\":[-1,0,200,20]}"},
}

var mapDecodeTestCases = []struct {
	Bin  []byte
	Json string
}{
	{[]byte("\xa2\x64IETF\x20"), "{\"IETF\":-1}"},
	{[]byte("\xa2\x65Array\x84\x20\x00\x18\xc8\x14"), "{\"Array\":[-1,0,200,20]}"},
	{[]byte("\xa6\x61\x61\x01\x61\x62\x02\x61\x63\x03"), "{\"a\":1,\"b\":2,\"c\":3}"},
	{[]byte("\xbf\x61a\x01\x61b\x02\xff"), "{\"a\":1,\"b\":2}"},
}

func TestDecodeMap(t *testing.T) {
	for _, tc := range mapDecodeTestCases {
		buf := bytes.NewBuffer([]byte{})
		map2Json(getReader(string(tc.Bin)), buf)
		if buf.String() != tc.Json {
			t.Errorf("map2Json(0x%s)=%s, want: %s", hex.EncodeToString(tc.Bin), buf.String(), tc.Json)
		}
	}
	for _, tc := range infiniteMapDecodeTestCases {
		buf := bytes.NewBuffer([]byte{})
		map2Json(getReader(string(tc.Bin)), buf)
		if buf.String() != tc.Json {
			t.Errorf("map2Json(0x%s)=%s, want: %s", hex.EncodeToString(tc.Bin), buf.String(), tc.Json)
		}
	}
}

func TestDecodeBool(t *testing.T) {
	for _, tc := range internal.BooleanTestCases {
		got := decodeSimpleFloat(getReader(tc.Binary))
		if string(got) != tc.Json {
			t.Errorf("decodeSimpleFloat(0x%s)=%s, want:%s", hex.EncodeToString([]byte(tc.Binary)), string(got), tc.Json)
		}
	}
}

func TestDecodeFloat(t *testing.T) {
	for _, tc := range internal.Float32TestCases {
		got, _ := decodeFloat(getReader(tc.Binary))
		if got != float64(tc.Val) && math.IsNaN(got) != math.IsNaN(float64(tc.Val)) {
			t.Errorf("decodeFloat(0x%s)=%f, want:%f", hex.EncodeToString([]byte(tc.Binary)), got, tc.Val)
		}
	}
	for _, tc := range internal.Float64TestCases {
		got, _ := decodeFloat(getReader(tc.Binary))
		if got != tc.Val && math.IsNaN(got) != math.IsNaN(tc.Val) {
			t.Errorf("decodeFloat(0x%s)=%f, want:%f", hex.EncodeToString([]byte(tc.Binary)), got, tc.Val)
		}
	}

	// Test float64 special values with correct CBOR encoding
	float64Tests := []struct {
		name  string
		input string
		want  float64
	}{
		{"float64 NaN", "\xfb\x7f\xf8\x00\x00\x00\x00\x00\x00", math.NaN()},
		{"float64 +Inf", "\xfb\x7f\xf0\x00\x00\x00\x00\x00\x00", math.Inf(0)},
		{"float64 -Inf", "\xfb\xff\xf0\x00\x00\x00\x00\x00\x00", math.Inf(-1)},
		{"float64 1.0", "\xfb\x3f\xf0\x00\x00\x00\x00\x00\x00", 1.0},
	}

	for _, tt := range float64Tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := decodeFloat(getReader(tt.input))
			if math.IsNaN(tt.want) {
				if !math.IsNaN(got) {
					t.Errorf("decodeFloat(%q) = %f, want NaN", tt.input, got)
				}
			} else if math.IsInf(tt.want, 0) {
				if !math.IsInf(got, 0) {
					t.Errorf("decodeFloat(%q) = %f, want +Inf", tt.input, got)
				}
			} else if math.IsInf(tt.want, -1) {
				if !math.IsInf(got, -1) {
					t.Errorf("decodeFloat(%q) = %f, want -Inf", tt.input, got)
				}
			} else if got != tt.want {
				t.Errorf("decodeFloat(%q) = %f, want %f", tt.input, got, tt.want)
			}
		})
	}
}

func TestDecodeTimestamp(t *testing.T) {
	decodeTimeZone, _ = time.LoadLocation("UTC")
	for _, tc := range internal.TimeIntegerTestcases {
		tm := decodeTagData(getReader(tc.Binary))
		if string(tm) != "\""+tc.RfcStr+"\"" {
			t.Errorf("decodeFloat(0x%s)=%s, want:%s", hex.EncodeToString([]byte(tc.Binary)), tm, tc.RfcStr)
		}
	}
	for _, tc := range internal.TimeFloatTestcases {
		tm := decodeTagData(getReader(tc.Out))
		//Since we convert to float and back - it may be slightly off - so
		//we cannot check for exact equality instead, we'll check it is
		//very close to each other Less than a Microsecond (lets not yet do nanosec)

		got, _ := time.Parse(string(tm), string(tm))
		want, _ := time.Parse(tc.RfcStr, tc.RfcStr)
		if got.Sub(want) > time.Microsecond {
			t.Errorf("decodeFloat(0x%s)=%s, want:%s", hex.EncodeToString([]byte(tc.Out)), tm, tc.RfcStr)
		}
	}

	// Test with decodeTimeZone = nil to cover the else branches
	oldTimeZone := decodeTimeZone
	decodeTimeZone = nil
	defer func() { decodeTimeZone = oldTimeZone }()

	for _, tc := range internal.TimeIntegerTestcases {
		tm := decodeTagData(getReader(tc.Binary))
		if string(tm) != "\""+tc.RfcStr+"\"" {
			t.Errorf("decodeFloat(0x%s)=%s, want:%s", hex.EncodeToString([]byte(tc.Binary)), tm, tc.RfcStr)
		}
	}
	for _, tc := range internal.TimeFloatTestcases {
		tm := decodeTagData(getReader(tc.Out))
		got, _ := time.Parse(string(tm), string(tm))
		want, _ := time.Parse(tc.RfcStr, tc.RfcStr)
		if got.Sub(want) > time.Microsecond {
			t.Errorf("decodeFloat(0x%s)=%s, want:%s", hex.EncodeToString([]byte(tc.Out)), tm, tc.RfcStr)
		}
	}
}

func TestDecodeNetworkAddr(t *testing.T) {
	for _, tc := range internal.IpAddrTestCases {
		d1 := decodeTagData(getReader(tc.Binary))
		if string(d1) != tc.Text {
			t.Errorf("decodeNetworkAddr(0x%s)=%s, want:%s", hex.EncodeToString([]byte(tc.Binary)), d1, tc.Text)
		}
	}
}

func TestDecodeMACAddr(t *testing.T) {
	for _, tc := range internal.MacAddrTestCases {
		d1 := decodeTagData(getReader(tc.Binary))
		if string(d1) != tc.Text {
			t.Errorf("decodeNetworkAddr(0x%s)=%s, want:%s", hex.EncodeToString([]byte(tc.Binary)), d1, tc.Text)
		}
	}
}

func TestDecodeIPPrefix(t *testing.T) {
	for _, tc := range internal.IPPrefixTestCases {
		d1 := decodeTagData(getReader(tc.Binary))
		if string(d1) != tc.Text {
			t.Errorf("decodeIPPrefix(0x%s)=%s, want:%s", hex.EncodeToString([]byte(tc.Binary)), d1, tc.Text)
		}
	}
}

var compositeCborTestCases = []struct {
	Binary []byte
	Json   string
}{
	{[]byte("\xbf\x64IETF\x20\x65Array\x9f\x20\x00\x18\xc8\x14\xff\xff"), "{\"IETF\":-1,\"Array\":[-1,0,200,20]}\n"},
	{[]byte("\xbf\x64IETF\x64YES!\x65Array\x9f\x20\x00\x18\xc8\x14\xff\xff"), "{\"IETF\":\"YES!\",\"Array\":[-1,0,200,20]}\n"},
	{[]byte("\xbf\x61a\x01\x61b\x02\x61c\x03\xff"), "{\"a\":1,\"b\":2,\"c\":3}\n"},
	{[]byte("\xc1\x1a\x51\x0f\x30\xd8"), "\"2013-02-04T03:54:00Z\"\n"},
}

func TestDecodeCbor2Json(t *testing.T) {
	for _, tc := range compositeCborTestCases {
		buf := bytes.NewBuffer([]byte{})
		err := Cbor2JsonManyObjects(getReader(string(tc.Binary)), buf)
		if buf.String() != tc.Json || err != nil {
			t.Errorf("cbor2JsonManyObjects(0x%s)=%s, want: %s, err:%s", hex.EncodeToString(tc.Binary), buf.String(), tc.Json, err.Error())
		}
	}
}

var negativeCborTestCases = []struct {
	Binary []byte
	errStr string
}{
	{[]byte("\xb9\x64IETF\x20\x65Array\x9f\x20\x00\x18\xc8\x14"), "Tried to Read 18 Bytes.. But hit end of file"},
	{[]byte("\xbf\x64IETF\x20\x65Array\x9f\x20\x00\x18\xc8\x14"), "EOF"},
	{[]byte("\xbf\x14IETF\x20\x65Array\x9f\x20\x00\x18\xc8\x14"), "Tried to Read 40736 Bytes.. But hit end of file"},
	{[]byte("\xbf\x64IETF"), "EOF"},
	{[]byte("\xbf\x64IETF\x20\x65Array\x9f\x20\x00\x18\xc8\xff\xff\xff"), "Invalid Additional Type: 31 in decodeSimpleFloat"},
	{[]byte("\xbf\x64IETF\x20\x65Array"), "EOF"},
	{[]byte("\xbf\x64"), "Tried to Read 4 Bytes.. But hit end of file"},
}

func TestDecodeNegativeCbor2Json(t *testing.T) {
	for _, tc := range negativeCborTestCases {
		buf := bytes.NewBuffer([]byte{})
		err := Cbor2JsonManyObjects(getReader(string(tc.Binary)), buf)
		if err == nil || err.Error() != tc.errStr {
			t.Errorf("Expected error got:%s, want:%s", err, tc.errStr)
		}
	}
}

func TestBinaryFmt(t *testing.T) {
	tests := []struct {
		input []byte
		want  bool
	}{
		{[]byte{}, false},
		{[]byte{0x00}, false},
		{[]byte{0x7F}, false},
		{[]byte{0x80}, true},
		{[]byte{0xFF}, true},
		{[]byte{0x00, 0x80}, false}, // Only checks first byte
	}

	for _, tt := range tests {
		got := binaryFmt(tt.input)
		if got != tt.want {
			t.Errorf("binaryFmt(%v) = %v, want %v", tt.input, got, tt.want)
		}
	}
}

func TestDecodeIfBinaryToString(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
		want  string
	}{
		{
			name:  "non-binary input",
			input: []byte(`{"key":"value"}`),
			want:  `{"key":"value"}`,
		},
		{
			name:  "binary input - simple object",
			input: []byte("\xbf\x64IETF\x20\xff"), // {"IETF": -1} in indefinite length CBOR
			want:  "{\"IETF\":-1}\n",
		},
		{
			name:  "binary input - multiple objects",
			input: []byte("\xbf\x64IETF\x20\xff\xbf\x65Array\x84\x20\x00\x18\xc8\x14\xff"), // Two objects
			want:  "{\"IETF\":-1}\n{\"Array\":[-1,0,200,20]}\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DecodeIfBinaryToString(tt.input)
			if got != tt.want {
				t.Errorf("DecodeIfBinaryToString() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestDecodeObjectToStr(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
		want  string
	}{
		{
			name:  "non-binary input",
			input: []byte(`{"key":"value"}`),
			want:  `{"key":"value"}`,
		},
		{
			name:  "binary input - simple object",
			input: []byte("\xbf\x64IETF\x20\xff"), // {"IETF": -1} in indefinite length CBOR
			want:  "{\"IETF\":-1}",
		},
		{
			name:  "binary input - array",
			input: []byte("\x84\x20\x00\x18\xc8\x14"), // [-1, 0, 200, 20]
			want:  "[-1,0,200,20]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DecodeObjectToStr(tt.input)
			if got != tt.want {
				t.Errorf("DecodeObjectToStr() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestDecodeIfBinaryToBytes(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
		want  []byte
	}{
		{
			name:  "non-binary input",
			input: []byte(`{"key":"value"}`),
			want:  []byte(`{"key":"value"}`),
		},
		{
			name:  "binary input - simple object",
			input: []byte("\xbf\x64IETF\x20\xff"), // {"IETF": -1} in indefinite length CBOR
			want:  []byte("{\"IETF\":-1}\n"),
		},
		{
			name:  "binary input - multiple objects",
			input: []byte("\xbf\x64IETF\x20\xff\xbf\x65Array\x84\x20\x00\x18\xc8\x14\xff"), // Two objects
			want:  []byte("{\"IETF\":-1}\n{\"Array\":[-1,0,200,20]}\n"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DecodeIfBinaryToBytes(tt.input)
			if !bytes.Equal(got, tt.want) {
				t.Errorf("DecodeIfBinaryToBytes() = %q, want %q", string(got), string(tt.want))
			}
		})
	}
}

func TestDecodeEmbeddedCBOR(t *testing.T) {
	// Test embedded CBOR tag: 0xD8 0x3F (tag 63) followed by byte string
	// 0xD8 = major type 6 (tags) + additional type 24 (uint8 follows)
	// 0x3F = 63 (additionalTypeEmbeddedCBOR)
	// 0x43 = major type 2 (byte string) + length 3
	// 0x01 0x02 0x03 = the embedded CBOR data

	embeddedCBOR := []byte("\xd8\x3f\x43\x01\x02\x03")
	expected := "\"data:application/cbor;base64,AQID\""

	got := decodeTagData(getReader(string(embeddedCBOR)))
	if string(got) != expected {
		t.Errorf("decodeTagData(embedded CBOR) = %q, want %q", string(got), expected)
	}
}

func TestDecodeEmbeddedJSON(t *testing.T) {
	t.Run("valid embedded JSON", func(t *testing.T) {
		// Test embedded JSON tag: 0xD9 0x01 0x06 (tag 262) followed by byte string.
		// 0xD9 = major type 6 (tags) + additional type 25 (uint16 follows)
		// 0x01 0x06 = 262 (additionalTypeEmbeddedJSON)
		// 0x47 = major type 2 (byte string) + length 7
		// {"a":1} = embedded JSON payload (no surrounding quotes expected)
		embeddedJSON := []byte("\xd9\x01\x06\x47{\"a\":1}")
		expected := "{\"a\":1}"

		got := decodeTagData(getReader(string(embeddedJSON)))
		if string(got) != expected {
			t.Errorf("decodeTagData(embedded JSON) = %q, want %q", string(got), expected)
		}
	})

	t.Run("unsupported embedded type panics", func(t *testing.T) {
		// Same embedded JSON tag, but followed by a UTF-8 string instead of a byte string.
		// This should hit the "Unsupported embedded Type" panic branch.
		bad := []byte("\xd9\x01\x06\x61x")

		defer func() {
			if r := recover(); r == nil {
				t.Fatalf("expected panic, got none")
			}
		}()

		_ = decodeTagData(getReader(string(bad)))
	})
}

func TestDecodeHexString(t *testing.T) {
	// Test hex string tag: 0xD9 0x01 0x07 (tag 263) followed by byte string
	// 0xD9 = major type 6 (tags) + additional type 25 (uint16 follows)
	// 0x01 0x07 = 263 (additionalTypeTagHexString)
	// 0x43 = major type 2 (byte string) + length 3
	// 0x01 0x02 0x03 = the byte data to hex encode

	hexString := []byte("\xd9\x01\x07\x43\x01\x02\x03")
	expected := "\"010203\""

	got := decodeTagData(getReader(string(hexString)))
	if string(got) != expected {
		t.Errorf("decodeTagData(hex string) = %q, want %q", string(got), expected)
	}
}

func TestDecodeSimpleFloat(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Boolean and null cases (already covered)
		{"true", "\xf5", "true"},
		{"false", "\xf4", "false"},
		{"null", "\xf6", "null"},

		// Float32 cases
		{"float32 1.0", "\xfa\x3f\x80\x00\x00", "1"},
		{"float32 1.5", "\xfa\x3f\xc0\x00\x00", "1.5"},
		{"float32 +Inf", "\xfa\x7f\x80\x00\x00", "\"+Inf\""},
		{"float32 -Inf", "\xfa\xff\x80\x00\x00", "\"-Inf\""},
		{"float32 NaN", "\xfa\x7f\xc0\x00\x00", "\"NaN\""},

		// Float64 cases
		{"float64 1.0", "\xfb\x3f\xf0\x00\x00\x00\x00\x00\x00", "1"},
		{"float64 +Inf", "\xfb\x7f\xf0\x00\x00\x00\x00\x00\x00", "\"+Inf\""},
		{"float64 -Inf", "\xfb\xff\xf0\x00\x00\x00\x00\x00\x00", "\"-Inf\""},
		{"float64 NaN", "\xfb\x7f\xf8\x00\x00\x00\x00\x00\x00", "\"NaN\""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := decodeSimpleFloat(getReader(tt.input))
			if string(got) != tt.want {
				t.Errorf("decodeSimpleFloat(%q) = %q, want %q", tt.input, string(got), tt.want)
			}
		})
	}
}
