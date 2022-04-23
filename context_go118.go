//go:build go1.18
// +build go1.18

package zerolog

import "net/netip"

// NetipAddr adds IPv4 or IPv6 Address to the context.
func (c Context) NetipAddr(key string, addr netip.Addr) Context {
	c.l.context = enc.AppendNetipAddr(enc.AppendKey(c.l.context, key), addr)
	return c
}

// NetipAddrPort adds IPv4 or IPv6 Address and port to the context.
func (c Context) NetipAddrPort(key string, addrPort netip.AddrPort) Context {
	c.l.context = enc.AppendNetipAddrPort(enc.AppendKey(c.l.context, key), addrPort)
	return c
}

// NetipPrefix adds IPv4 or IPv6 Prefix (address and mask) to the context.
func (c Context) NetipPrefix(key string, pfx netip.Prefix) Context {
	c.l.context = enc.AppendNetipPrefix(enc.AppendKey(c.l.context, key), pfx)
	return c
}
