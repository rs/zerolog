//go:build go1.18
// +build go1.18

package json

import (
	"net/netip"
)

// AppendNetipAddr adds IPv4 or IPv6 address to dst.
func (e Encoder) AppendNetipAddr(dst []byte, addr netip.Addr) []byte {
	return e.AppendString(dst, addr.String())
}

// AppendNetipAddrPort adds IPv4 or IPv6 address and port to dst.
func (e Encoder) AppendNetipAddrPort(dst []byte, addrPort netip.AddrPort) []byte {
	return e.AppendString(dst, addrPort.String())
}

// AppendNetipPrefix adds IPv4 or IPv6 Prefix (address & mask) to dst.
func (e Encoder) AppendNetipPrefix(dst []byte, pfx netip.Prefix) []byte {
	return e.AppendString(dst, pfx.String())
}
