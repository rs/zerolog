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
