package zerolog

import (
	"context"
	"errors"
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
	Ctx        context.Context
	Durations  []time.Duration
	Errs       []error
	Floats     []float64
	Interfaces []struct {
		Pub  string
		Tag  string `json:"tag"`
		priv int
	}
	Ints     []int
	IPAddrs  []net.IP
	IPPfxs   []net.IPNet
	MACAddr  net.HardwareAddr
	Objects  []fixtureObj
	Stringer net.IP
	Strings  []string
	Times    []time.Time
	Type     reflect.Type
}

func makeFieldFixtures() *fieldFixtures {
	bools := []bool{true, false, true, false, true, false, true, false, true, false}
	ints := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	floats := []float64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
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
	objects := []fixtureObj{
		{"a", "z", 1},
		{"b", "y", 2},
		{"c", "x", 3},
		{"d", "w", 4},
		{"e", "v", 5},
		{"f", "u", 6},
		{"g", "t", 7},
		{"h", "s", 8},
		{"i", "r", 9},
		{"j", "q", 10},
	}
	ipAddrV4 := net.IP{192, 168, 0, 1}
	ipAddrV6 := net.IP{0x20, 0x01, 0x0d, 0xb8, 0x85, 0xa3, 0x00, 0x00, 0x00, 0x00, 0x8a, 0x2e, 0x03, 0x70, 0x73, 0x34}
	ipAddrs := []net.IP{ipAddrV4, ipAddrV6, ipAddrV4, ipAddrV6, ipAddrV4, ipAddrV6, ipAddrV4, ipAddrV6, ipAddrV4, ipAddrV6}
	ipPfxV4 := net.IPNet{IP: net.IP{192, 168, 0, 0}, Mask: net.CIDRMask(24, 32)}
	ipPfxV6 := net.IPNet{IP: net.IP{0x20, 0x01, 0x0d, 0xb8, 0x85, 0xa3, 0x00, 0x00, 0x00, 0x00, 0x8a, 0x2e, 0x03, 0x70, 0x73, 0x00}, Mask: net.CIDRMask(64, 128)}
	ipPfxs := []net.IPNet{ipPfxV4, ipPfxV6, ipPfxV4, ipPfxV6, ipPfxV4, ipPfxV6, ipPfxV4, ipPfxV6, ipPfxV4, ipPfxV6}
	macAddr := net.HardwareAddr{0x00, 0x1A, 0x2B, 0x3C, 0x4D, 0x5E}
	errs := []error{errors.New("a"), errors.New("b"), errors.New("c"), errors.New("d"), errors.New("e")}
	ctx := context.Background()
	stringer := net.IP{127, 0, 0, 1}

	return &fieldFixtures{
		Bools:      bools,
		Ctx:        ctx,
		Durations:  durations,
		Errs:       errs,
		Floats:     floats,
		Interfaces: interfaces,
		Ints:       ints,
		IPAddrs:    ipAddrs,
		IPPfxs:     ipPfxs,
		MACAddr:    macAddr,
		Objects:    objects,
		Stringer:   stringer,
		Strings:    strings,
		Times:      times,
		Type:       reflect.TypeOf(12345),
	}
}
