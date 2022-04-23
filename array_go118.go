//go:build go1.18
// +build go1.18

package zerolog

import (
	"net/netip"
)

// NetipAddr adds IPv4 or IPv6 address to the array.
func (a *Array) NetipAddr(addr netip.Addr) *Array {
	a.buf = enc.AppendNetipAddr(enc.AppendArrayDelim(a.buf), addr)
	return a
}

// NetipAddrPort adds IPv4 or IPv6 address and port to the array.
func (a *Array) NetipAddrPort(addrPort netip.AddrPort) *Array {
	a.buf = enc.AppendNetipAddrPort(enc.AppendArrayDelim(a.buf), addrPort)
	return a
}

// NetipPrefix adds IPv4 or IPv6 address prefixto the array.
func (a *Array) NetipPrefix(pfx netip.Prefix) *Array {
	a.buf = enc.AppendNetipPrefix(enc.AppendArrayDelim(a.buf), pfx)
	return a
}
