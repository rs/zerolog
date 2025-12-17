package zerolog

import (
	"context"
	"errors"
	"fmt"
	"net"
	"reflect"
	"time"
)

type fixtureObj struct {
	Pub  string
	Tag  string `json:"tag"`
	priv int
}

func (o fixtureObj) MarshalZerologObject(e *Event) {
	e.Str("Pub", o.Pub).
		Str("Tag", o.Tag).
		Int("priv", o.priv)
}

type fieldFixtures struct {
	Bools      []bool
	Bytes      []byte
	Ctx        context.Context
	Durations  []time.Duration
	Errs       []error
	Floats32   []float32
	Floats64   []float64
	Interfaces []struct {
		Pub  string
		Tag  string `json:"tag"`
		priv int
	}
	Ints      []int
	Ints8     []int8
	Ints16    []int16
	Ints32    []int32
	Ints64    []int64
	Uints     []uint
	Uints8    []uint8
	Uints16   []uint16
	Uints32   []uint32
	Uints64   []uint64
	IPAddrs   []net.IP
	IPPfxs    []net.IPNet
	MACAddr   net.HardwareAddr
	Objects   []LogObjectMarshaler
	RawCBOR   []byte
	RawJSONs  [][]byte
	Stringers []fmt.Stringer
	Strings   []string
	Times     []time.Time
	Type      reflect.Type
}

func makeFieldFixtures() *fieldFixtures {
	bools := []bool{true, false, true, false, true, false, true, false, true, false}
	bytes := []byte(`abcdef`)
	ints := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	ints8 := []int8{-8, 8}
	ints16 := []int16{-16, 16}
	ints32 := []int32{-32, 32}
	ints64 := []int64{-64, 64}
	uints := []uint{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	uints8 := []uint8{8, uint8(^uint8(0))}
	uints16 := []uint16{16, uint16(^uint16(0))}
	uints32 := []uint32{32, uint32(^uint32(0))}
	uints64 := []uint64{64, uint64(^uint64(0))}
	floats32 := []float32{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	floats64 := []float64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	strings := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j"}
	durations := []time.Duration{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	times := []time.Time{
		time.Unix(0, 0),
		time.Unix(1, 0),
		time.Unix(2, 0),
		time.Unix(3, 0),
		time.Unix(4, 0),
		time.Unix(5, 0),
		time.Unix(6, 0),
		time.Unix(7, 0),
		time.Unix(8, 0),
		time.Unix(9, 0),
	}
	interfaces := []struct {
		Pub  string
		Tag  string `json:"tag"`
		priv int
	}{
		{"A", "j", -5},
		{"B", "i", -4},
		{"C", "h", -3},
		{"D", "g", -2},
		{"E", "f", -1},
		{"F", "e", 0},
		{"G", "d", 1},
		{"H", "c", 2},
		{"I", "b", 3},
		{"J", "a", 4},
	}
	objects := []LogObjectMarshaler{
		fixtureObj{"a", "z", 1},
		fixtureObj{"b", "y", 2},
		fixtureObj{"c", "x", 3},
		fixtureObj{"d", "w", 4},
		fixtureObj{"e", "v", 5},
		fixtureObj{"f", "u", 6},
		fixtureObj{"g", "t", 7},
		fixtureObj{"h", "s", 8},
		fixtureObj{"i", "r", 9},
		fixtureObj{"j", "q", 10},
	}
	ipAddrV4 := net.IP{192, 168, 0, 1}
	ipAddrV6 := net.IP{0x20, 0x01, 0x0d, 0xb8, 0x85, 0xa3, 0x00, 0x00, 0x00, 0x00, 0x8a, 0x2e, 0x03, 0x70, 0x73, 0x34}
	ipAddrs := []net.IP{ipAddrV4, ipAddrV6, ipAddrV4, ipAddrV6, ipAddrV4, ipAddrV6, ipAddrV4, ipAddrV6, ipAddrV4, ipAddrV6}
	ipPfxV4 := net.IPNet{IP: net.IP{192, 168, 0, 0}, Mask: net.CIDRMask(24, 32)}
	ipPfxV6 := net.IPNet{IP: net.IP{0x20, 0x01, 0x0d, 0xb8, 0x85, 0xa3, 0x00, 0x00, 0x00, 0x00, 0x8a, 0x2e, 0x03, 0x70, 0x73, 0x00}, Mask: net.CIDRMask(64, 128)}
	ipPfxs := []net.IPNet{ipPfxV4, ipPfxV6, ipPfxV4, ipPfxV6, ipPfxV4, ipPfxV6, ipPfxV4, ipPfxV6, ipPfxV4, ipPfxV6}
	macAddr := net.HardwareAddr{0x00, 0x1A, 0x2B, 0x3C, 0x4D, 0x5E}
	errs := []error{errors.New("a"), errors.New("b"), errors.New("c"), errors.New("d"), errors.New("e"), nil, loggableError{fmt.Errorf("oops")}}
	ctx := context.Background()
	stringers := []fmt.Stringer{ipAddrs[0], durations[0]}
	rawJSONs := [][]byte{[]byte(`{"some":"json"}`), []byte(`{"longer":[1111,2222,3333,4444,5555]}`)}
	rawCBOR := []byte{0xA1, 0x64, 0x73, 0x6F, 0x6D, 0x65, 0x64, 0x61, 0x74, 0x61} // {"some":"data"}

	return &fieldFixtures{
		Bools:      bools,
		Bytes:      bytes,
		Ctx:        ctx,
		Durations:  durations,
		Errs:       errs,
		Floats32:   floats32,
		Floats64:   floats64,
		Interfaces: interfaces,
		Ints:       ints,
		Ints8:      ints8,
		Ints16:     ints16,
		Ints32:     ints32,
		Ints64:     ints64,
		Uints:      uints,
		Uints8:     uints8,
		Uints16:    uints16,
		Uints32:    uints32,
		Uints64:    uints64,
		IPAddrs:    ipAddrs,
		IPPfxs:     ipPfxs,
		MACAddr:    macAddr,
		Objects:    objects,
		RawCBOR:    rawCBOR,
		RawJSONs:   rawJSONs,
		Stringers:  stringers,
		Strings:    strings,
		Times:      times,
		Type:       reflect.TypeOf(12345),
	}
}
