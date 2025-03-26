package zerolog

import (
	"encoding/base64"
	"net"
	"testing"
	"time"
)

func TestArray(t *testing.T) {
	a := Arr().
		Bool(true).
		Int(1).
		Int8(2).
		Int16(3).
		Int32(4).
		Int64(5).
		Uint(6).
		Uint8(7).
		Uint16(8).
		Uint32(9).
		Uint64(10).
		Float32(11.98122).
		Float64(12.987654321).
		Str("a").
		Bytes([]byte("b")).
		Hex([]byte{0x1f}).
		Base64([]byte{0x12, 0xef, 0x29, 0x30, 0xff}).
		Base64Custom(base64.RawURLEncoding, []byte{0xcc, 0xbb, 0xaa, 0xff}).
		RawJSON([]byte(`{"some":"json"}`)).
		Time(time.Time{}).
		IPAddr(net.IP{192, 168, 0, 10}).
		Dur(0).
		Dict(Dict().
			Str("bar", "baz").
			Int("n", 1),
		)
	want := `[true,1,2,3,4,5,6,7,8,9,10,11.98122,12.987654321,"a","b","1f","Eu8pMP8=","zLuq_w",{"some":"json"},"0001-01-01T00:00:00Z","192.168.0.10",0,{"bar":"baz","n":1}]`
	if got := decodeObjectToStr(a.write([]byte{})); got != want {
		t.Errorf("Array.write()\ngot:  %s\nwant: %s", got, want)
	}
}
