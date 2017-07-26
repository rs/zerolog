package zerolog

import (
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
		Float32(11).
		Float64(12).
		Str("a").
		Time(time.Time{}).
		Dur(0)
	want := `[true,1,2,3,4,5,6,7,8,9,10,11,12,"a","0001-01-01T00:00:00Z",0]`
	if got := string(a.write([]byte{})); got != want {
		t.Errorf("Array.write()\ngot:  %s\nwant: %s", got, want)
	}
}
