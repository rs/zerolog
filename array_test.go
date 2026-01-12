package zerolog

import (
	"fmt"
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
		RawJSON([]byte(`{"some":"json"}`)).
		RawJSON([]byte(`{"longer":[1111,2222,3333,4444,5555]}`)).
		Time(time.Time{}).
		IPAddr(net.IP{192, 168, 0, 10}).
		IPPrefix(net.IPNet{IP: net.IP{127, 0, 0, 0}, Mask: net.CIDRMask(24, 32)}).
		MACAddr(net.HardwareAddr{0x01, 0x23, 0x45, 0x67, 0x89, 0xab}).
		Interface(struct {
			Pub  string
			Tag  string `json:"tag"`
			priv int
		}{"A", "j", -5}).
		Interface(logObjectMarshalerImpl{
			name: "ZOT",
			age:  35,
		}).
		Dur(0).
		Dict(Dict().
			Str("bar", "baz").
			Int("n", 1),
		).
		Err(nil).
		Err(fmt.Errorf("failure")).
		Err(loggableError{fmt.Errorf("oops")}).
		Object(logObjectMarshalerImpl{
			name: "ZIT",
			age:  22,
		}).
		Type(3.14)

	want := `[true,1,2,3,4,5,6,7,8,9,10,11.98122,12.987654321,"a","b","1f",{"some":"json"},{"longer":[1111,2222,3333,4444,5555]},"0001-01-01T00:00:00Z","192.168.0.10","127.0.0.0/24","01:23:45:67:89:ab",{"Pub":"A","tag":"j"},{"name":"zot","age":-35},0,{"bar":"baz","n":1},null,"failure",{"l":"OOPS"},{"name":"zit","age":-22},"float64"]`
	if got := decodeObjectToStr(a.write([]byte{})); got != want {
		t.Errorf("Array.write()\ngot:  %s\nwant: %s", got, want)
	}
}

func TestArray_MarshalZerologArray(t *testing.T) {
	a := Arr()
	a.MarshalZerologArray(nil) // no-op method, should not panic
}
