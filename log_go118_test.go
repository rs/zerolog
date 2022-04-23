//go:build go1.18
// +build go1.18

package zerolog

import (
	"bytes"
	"fmt"
	"net/netip"
	"runtime"
	"testing"
)

func TestWith_netip(t *testing.T) {
	out := &bytes.Buffer{}
	ctx := New(out).With().
		NetipAddr("addr", netip.AddrFrom4([4]byte{127, 0, 0, 1})).
		NetipAddrPort("addrPort", netip.AddrPortFrom(netip.AddrFrom4([4]byte{127, 0, 0, 1}), 8080)).
		NetipPrefix("prefix", netip.PrefixFrom(netip.AddrFrom4([4]byte{192, 168, 0, 1}), 24))
	_, file, line, _ := runtime.Caller(0)
	caller := fmt.Sprintf("%s:%d", file, line+3)
	log := ctx.Caller().Logger()
	log.Log().Send()
	if got, want := decodeIfBinaryToString(out.Bytes()), `{"addr":"127.0.0.1","addrPort":"127.0.0.1:8080","prefix":"192.168.0.1/24","caller":"`+caller+`"}`+"\n"; got != want {
		t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
	}
}

func TestFieldsMap_netip(t *testing.T) {
	out := &bytes.Buffer{}
	log := New(out)
	log.Log().Fields(map[string]interface{}{
		"addr":     netip.AddrFrom4([4]byte{10, 0, 0, 1}),
		"addrPort": netip.AddrPortFrom(netip.AddrFrom4([4]byte{172, 20, 0, 1}), 8080),
		"prefix":   netip.PrefixFrom(netip.AddrFrom4([4]byte{192, 168, 0, 1}), 24),
	}).Send()
	if got, want := decodeIfBinaryToString(out.Bytes()), `{"addr":"10.0.0.1","addrPort":"172.20.0.1:8080","prefix":"192.168.0.1/24"}`+"\n"; got != want {
		t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
	}
}

func TestFieldsSlice_netip(t *testing.T) {
	out := &bytes.Buffer{}
	log := New(out)
	log.Log().Fields([]interface{}{
		"addr", netip.AddrFrom4([4]byte{10, 0, 0, 1}),
		"addrPort", netip.AddrPortFrom(netip.AddrFrom4([4]byte{172, 20, 0, 1}), 8080),
		"prefix", netip.PrefixFrom(netip.AddrFrom4([4]byte{192, 168, 0, 1}), 24),
	}).Send()
	if got, want := decodeIfBinaryToString(out.Bytes()), `{"addr":"10.0.0.1","addrPort":"172.20.0.1:8080","prefix":"192.168.0.1/24"}`+"\n"; got != want {
		t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
	}
}
