package cbor

import (
	"bytes"
	"encoding/hex"
	"testing"
	"time"
)

func TestDecodeInteger(t *testing.T) {
	for _, tc := range integerTestCases {
		gotv, gotc, err := decodeInteger([]byte(tc.binary))
		if gotv != int64(tc.val) || int(gotc) != len(tc.binary) || err != nil {
			t.Errorf("decodeInteger(0x%s)=0x%d, want: 0x%d",
				hex.EncodeToString([]byte(tc.binary)), gotv, tc.val)
		}
	}
}

func TestDecodeString(t *testing.T) {
	for _, tt := range encodeStringTests {
		got, _, err := decodeUTF8String([]byte(tt.binary))
		if err != nil {
			t.Errorf("Got Error for the case: %s", hex.EncodeToString([]byte(tt.binary)))
		}
		if string(got) != "\""+tt.json+"\"" {
			t.Errorf("DecodeString(0x%s)=%s, want:\"%s\"\n", hex.EncodeToString([]byte(tt.binary)), string(got),
				hex.EncodeToString([]byte(tt.json)))
		}
	}
}

func TestDecodeArray(t *testing.T) {
	for _, tc := range integerArrayTestCases {
		buf := bytes.NewBuffer([]byte{})
		_, err := array2Json([]byte(tc.binary), buf)
		if err != nil {
			panic(err)
		}
		if buf.String() != tc.json {
			t.Errorf("array2Json(0x%s)=%s, want: %s", hex.EncodeToString([]byte(tc.binary)), buf.String(), tc.json)
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
		_, err := array2Json([]byte(tc.in), buf)
		if err != nil {
			panic(err)
		}
		if buf.String() != tc.out {
			t.Errorf("array2Json(0x%s)=%s, want: %s", hex.EncodeToString([]byte(tc.out)), buf.String(), tc.out)
		}
	}
	for _, tc := range booleanArrayTestCases {
		buf := bytes.NewBuffer([]byte{})
		_, err := array2Json([]byte(tc.binary), buf)
		if err != nil {
			t.Errorf("array2Json(0x%s) errored out: %s", hex.EncodeToString([]byte(tc.binary)), err.Error())
		}
		if buf.String() != tc.json {
			t.Errorf("array2Json(0x%s)=%s, want: %s", hex.EncodeToString([]byte(tc.binary)), buf.String(), tc.json)
		}
	}
	//TODO add cases for arrays of other types
}

var infiniteMapDecodeTestCases = []struct {
	bin  []byte
	json string
}{
	{[]byte("\xbf\x64IETF\x20\xff"), "{\"IETF\":-1}"},
	{[]byte("\xbf\x65Array\x84\x20\x00\x18\xc8\x14\xff"), "{\"Array\":[-1,0,200,20]}"},
}

var mapDecodeTestCases = []struct {
	bin  []byte
	json string
}{
	{[]byte("\xa2\x64IETF\x20"), "{\"IETF\":-1}"},
	{[]byte("\xa2\x65Array\x84\x20\x00\x18\xc8\x14"), "{\"Array\":[-1,0,200,20]}"},
}

func TestDecodeMap(t *testing.T) {
	for _, tc := range mapDecodeTestCases {
		buf := bytes.NewBuffer([]byte{})
		_, err := map2Json(tc.bin, buf)
		if err != nil {
			t.Errorf("map2Json(0x%s) returned error", err)
		}
		if buf.String() != tc.json {
			t.Errorf("map2Json(0x%s)=%s, want: %s", hex.EncodeToString(tc.bin), buf.String(), tc.json)
		}
	}
	for _, tc := range infiniteMapDecodeTestCases {
		buf := bytes.NewBuffer([]byte{})
		_, err := map2Json(tc.bin, buf)
		if err != nil {
			t.Errorf("map2Json(0x%s) returned error", err)
		}
		if buf.String() != tc.json {
			t.Errorf("map2Json(0x%s)=%s, want: %s", hex.EncodeToString(tc.bin), buf.String(), tc.json)
		}
	}
}

func TestDecodeBool(t *testing.T) {
	for _, tc := range booleanTestCases {
		got, _, err := decodeSimpleFloat([]byte(tc.binary))
		if err != nil {
			t.Errorf("decodeSimpleFloat(0x%s) errored %s", hex.EncodeToString([]byte(tc.binary)), err.Error())
		}
		if string(got) != tc.json {
			t.Errorf("decodeSimpleFloat(0x%s)=%s, want:%s", hex.EncodeToString([]byte(tc.binary)), string(got), tc.json)
		}
	}
}

func TestDecodeFloat(t *testing.T) {
	for _, tc := range float32TestCases {
		got, _, err := decodeFloat([]byte(tc.binary))
		if err != nil {
			t.Errorf("decodeFloat(0x%s) returned error: %s", hex.EncodeToString([]byte(tc.binary)), err.Error())
		}
		if got != float64(tc.val) {
			t.Errorf("decodeFloat(0x%s)=%f, want:%f", hex.EncodeToString([]byte(tc.binary)), got, tc.val)
		}
	}
}

func TestDecodeTimestamp(t *testing.T) {
	decodeTimeZone, _ = time.LoadLocation("UTC")
	for _, tc := range timeIntegerTestcases {
		tm, _, err := decodeTagData([]byte(tc.binary))
		if err != nil {
			t.Errorf("decodeTagData(0x%s) returned error: %s", hex.EncodeToString([]byte(tc.binary)), err.Error())
		}
		if string(tm) != "\""+tc.rfcStr+"\"" {
			t.Errorf("decodeFloat(0x%s)=%s, want:%s", hex.EncodeToString([]byte(tc.binary)), tm, tc.rfcStr)
		}
	}
	for _, tc := range timeFloatTestcases {
		tm, _, err := decodeTagData([]byte(tc.out))
		if err != nil {
			t.Errorf("decodeTagData(0x%s) returned error: %s", hex.EncodeToString([]byte(tc.out)), err.Error())
		}
		//Since we convert to float and back - it may be slightly off - so
		//we cannot check for exact equality instead, we'll check it is
		//very close to each other Less than a Microsecond (lets not yet do nanosec)

		got, _ := time.Parse(string(tm), string(tm))
		want, _ := time.Parse(tc.rfcStr, tc.rfcStr)
		if got.Sub(want) > time.Microsecond {
			t.Errorf("decodeFloat(0x%s)=%s, want:%s", hex.EncodeToString([]byte(tc.out)), tm, tc.rfcStr)
		}
	}
}

var compositeCborTestCases = []struct {
	binary []byte
	json   string
}{
	{[]byte("\xbf\x64IETF\x20\x65Array\x9f\x20\x00\x18\xc8\x14\xff\xff"), "{\"IETF\":-1,\"Array\":[-1,0,200,20]}\n"},
	{[]byte("\xbf\x64IETF\x64YES!\x65Array\x9f\x20\x00\x18\xc8\x14\xff\xff"), "{\"IETF\":\"YES!\",\"Array\":[-1,0,200,20]}\n"},
}

func TestDecodeCbor2Json(t *testing.T) {
	for _, tc := range compositeCborTestCases {
		buf := bytes.NewBuffer([]byte{})
		_, err := Cbor2JsonManyObjects(tc.binary, buf)
		if err != nil {
			t.Errorf("cbor2JsonManyObjects(0x%s) returned error", err)
		}
		if buf.String() != tc.json {
			t.Errorf("cbor2JsonManyObjects(0x%s)=%s, want: %s", hex.EncodeToString(tc.binary), buf.String(), tc.json)
		}
	}
}
